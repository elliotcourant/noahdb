package sql

import (
	"github.com/elliotcourant/noahdb/pkg/ast"
	"github.com/elliotcourant/noahdb/pkg/core"
	"github.com/elliotcourant/noahdb/pkg/drivers/rqliter"
	"github.com/elliotcourant/noahdb/pkg/pgproto"
	"github.com/elliotcourant/noahdb/pkg/pgwirebase"
	"github.com/elliotcourant/noahdb/pkg/types"
	"time"
)

type responsePipe struct {
	conn core.PoolConnection
	err  error
}

func (s *session) executeExpandedPlan(plan ExpandedPlan) error {
	if len(plan.OutFormats) == 0 {
		plan.OutFormats = []pgwirebase.FormatCode{
			pgwirebase.FormatText,
		}
	}
	startTimestamp := time.Now()
	defer func() {
		s.log.Verbosef("[%s] execution of statement", time.Since(startTimestamp))
	}()

	switch plan.Target {
	case PlanTarget_STANDARD:
		responses := make(chan *responsePipe, len(plan.Tasks))

		for i, task := range plan.Tasks {
			go func(index int, task ExpandedPlanTask) {
				var response = &responsePipe{}
				defer func() {
					s.log.Verbosef("[%s] dispatch of query to data node shard [%d]", time.Since(startTimestamp), task.DataNodeShardID)
					responses <- response
				}()

				frontend, err := s.GetConnectionForDataNodeShard(task.DataNodeShardID)
				if err != nil {
					s.log.Errorf(
						"could not retrieve connection from pool for data node shard [%d]: %s",
						task.DataNodeShardID, err.Error())
					response.err = err
					return
				}
				s.log.Verbosef("{%d} executing: %s", task.DataNodeShardID, task.Query)

				queryMode := s.GetQueryMode()
				if task.Type != ast.Rows {
					queryMode = QueryModeStandard
				}

				switch queryMode {
				case QueryModeStandard:
					// In the standard query mode we don't need to care about the output format
					// Since we will be writing a row description header anyway.
					if err := frontend.Send(&pgproto.Query{
						String: task.Query,
					}); err != nil {
						s.log.Errorf(
							"could not send query to data node shard [%d]: %s",
							task.DataNodeShardID, err.Error())
						response.err = err
						return
					}
				case QueryModeExtended:
					// When we are in extended query mode we want to send the query in the same
					// extended query mode.
					if err := frontend.Send(&pgproto.Parse{
						Name:  "",
						Query: task.Query,
					}); err != nil {
						s.log.Errorf(
							"could not send query to data node shard [%d]: %s",
							task.DataNodeShardID, err.Error())
						response.err = err
						return
					}

					if err := frontend.Send(&pgproto.Describe{
						ObjectType: 'S',
						Name:       "",
					}); err != nil {
						s.log.Errorf(
							"could not describe query on data node shard [%d]: %s",
							task.DataNodeShardID, err.Error())
						response.err = err
						return
					}

					if err := frontend.Send(&pgproto.Bind{
						DestinationPortal: "",
						PreparedStatement: "",
						ResultFormatCodes: plan.OutFormats,
					}); err != nil {
						s.log.Errorf(
							"could not bind on data node shard [%d]: %s",
							task.DataNodeShardID, err.Error())
						response.err = err
						return
					}

					if err := frontend.Send(&pgproto.Execute{
						Portal:  "",
						MaxRows: 0,
					}); err != nil {
						s.log.Errorf(
							"could not bind on data node shard [%d]: %s",
							task.DataNodeShardID, err.Error())
						response.err = err
						return
					}

					if err := frontend.Send(&pgproto.Sync{}); err != nil {
						s.log.Errorf(
							"could not sync on data node shard [%d]: %s",
							task.DataNodeShardID, err.Error())
						response.err = err
						return
					}
				}

				response.conn = frontend
			}(i, task)
		}

		// If we are committing or rolling back a transaction then clear the transaction state.
		if plan.DistPlanType != DistributedPlanType_NONE {
			s.SetTransactionState(TransactionState_None)
		}

		for i := 0; i < len(plan.Tasks); i++ {
			sentRowDescription := s.GetQueryMode() == QueryModeExtended
			err := func(response *responsePipe) error {
				if response.err != nil {
					return response.err
				}
				frontend := response.conn
				// If we are not in a transaction then we can throw this connection away.
				if s.GetTransactionState() == TransactionState_None {
					defer s.ReleaseConnectionForDataNodeShard(frontend)
				}
				canExit := false
				for {
					message, err := frontend.Receive()
					if err != nil {
						s.log.Errorf("received error from frontend: %s", err.Error())
						return err
					}

					switch message.(type) {
					case *pgproto.RowDescription:
						if sentRowDescription {
							continue
						}
						if err := s.Backend().Send(message); err != nil {
							return err
						}
						canExit = true
						sentRowDescription = true
					case *pgproto.DataRow:
						if err := s.Backend().Send(message); err != nil {
							return err
						}
						canExit = true
					case *pgproto.ErrorResponse:
						if err := s.Backend().Send(message); err != nil {
							return err
						}
						canExit = true
						return nil
					case *pgproto.CommandComplete:
						canExit = true
						return nil
					case *pgproto.ReadyForQuery:
						if canExit {
							return nil
						}
					default:
						s.log.Tracef("received default message [%T]", message)
						// Do nothing
					}
				}
			}(<-responses)
			if err != nil {
				return err
			}
		}
	case PlanTarget_INTERNAL:
		for i, task := range plan.Tasks {
			return func() error {
				s.log.Verbosef("executing task %d on internal data store", i)
				response, err := s.Colony().Query(task.Query)
				rows := rqliter.NewRqlRows(response)
				if err != nil {
					s.log.Errorf("could not execute internal query: %s", err.Error())
					return err
				}
				result := make([][]interface{}, 0)
				columns := rows.Columns()
				for rows.Next() {
					row := make([]interface{}, len(columns))
					for i := 0; i < len(columns); i++ {
						row[i] = new(interface{})
					}
					if err := rows.Scan(row...); err != nil {
						s.log.Errorf("could not scan row: %s", err.Error())
						return err
					}
					result = append(result, row)
				}
				if err := rows.Err(); err != nil {
					s.log.Errorf("could not query internal store: %s", err.Error())
					return err
				}

				rowDescription := pgproto.RowDescription{
					Fields: make([]pgproto.FieldDescription, len(columns)),
				}

				typs := make([]interface{}, len(columns))

				for i := 0; i < len(columns); i++ {
					func() {
						field := pgproto.FieldDescription{
							Name:   columns[i],
							Format: int16(pgwirebase.FormatText),
						}
						defer func() {
							rowDescription.Fields[i] = field
						}()
						if len(result) > 0 {
							if len(result[0])-1 >= i {
								// This is some weird pointer magic to determine the type of the cell
								// basically without the *interface{} cast T would be 3 types at once?
								// which shouldn't be possible to my knowledge but it would show
								// as interface{} | *interface{} | int64 with a select 1 query.
								// So none of the types in the switch case would evaluate properly.
								// This weird magic fixes that.
								t := result[0][i].(*interface{})
								typs[i] = *t
								switch (*t).(type) {
								case int64:
									field.DataTypeOID = uint32(types.Type_int8)
								}
							}
						}
					}()
				}

				if err := s.Backend().Send(&rowDescription); err != nil {
					return err
				}

				for _, row := range result {
					dataRow := pgproto.DataRow{
						Values: make([][]byte, len(columns)),
					}
					for x, col := range row {
						switch typs[x].(type) {
						case int64:
							val := types.Int8{}
							if err := val.Set(col); err != nil {
								return err
							}
							dataRow.Values[x], err = val.EncodeText(nil, nil)
							if err != nil {
								return err
							}
						}
					}
					if err := s.Backend().Send(&dataRow); err != nil {
						return err
					}
				}
				return nil
			}()
		}
	}

	return nil
}
