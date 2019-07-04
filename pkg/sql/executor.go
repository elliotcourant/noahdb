package sql

import (
	"github.com/elliotcourant/noahdb/pkg/core"
	"github.com/elliotcourant/noahdb/pkg/drivers/rqliter"
	"github.com/elliotcourant/noahdb/pkg/pgproto"
	"github.com/elliotcourant/noahdb/pkg/pgwirebase"
	"github.com/elliotcourant/noahdb/pkg/types"
	"github.com/readystock/golog"
	"time"
)

type responsePipe struct {
	conn core.PoolConnection
	err  error
}

func (s *session) executeExpandedPlan(plan ExpandedPlan) error {
	startTimestamp := time.Now()
	defer func() {
		golog.Debugf("execution of statement took %s", time.Since(startTimestamp))
	}()

	switch plan.Target {
	case PlanTarget_STANDARD:
		responses := make(chan *responsePipe, len(plan.Tasks))

		for i, task := range plan.Tasks {
			go func(index int, task ExpandedPlanTask) {
				var response = &responsePipe{}
				defer func() {
					responses <- response
				}()

				frontend, err := s.Colony().Pool().GetConnectionForDataNodeShard(task.DataNodeShardID)
				if err != nil {
					golog.Errorf("could not retrieve connection from pool for data node shard [%d]: %s", task.DataNodeShardID, err.Error())
					response.err = err
					return
				}

				if err := frontend.Send(&pgproto.Query{
					String: task.Query,
				}); err != nil {
					golog.Errorf("could not send query to data node [%d]: %s", task.DataNodeShardID, err.Error())
					response.err = err
					return
				}
				response.conn = frontend
			}(i, task)
		}

		for i := 0; i < len(plan.Tasks); i++ {
			return func(response *responsePipe) error {
				if response.err != nil {
					return response.err
				}
				frontend := response.conn
				defer frontend.Release()
				canExit := false
				for {
					message, err := frontend.Receive()
					if err != nil {
						golog.Errorf("received error from frontend: %s", err.Error())
						return err
					}

					switch message.(type) {
					case *pgproto.RowDescription, *pgproto.DataRow:
						if err := s.Backend().Send(message); err != nil {
							return err
						}
						canExit = true
					case *pgproto.CommandComplete:
						canExit = true
						return nil
					case *pgproto.ErrorResponse:
						if err := s.Backend().Send(message); err != nil {
							return err
						}
						canExit = true
						return nil
					case *pgproto.ReadyForQuery:
						if canExit {
							return nil
						}
					default:
						// Do nothing
					}
				}
			}(<-responses)
		}
	case PlanTarget_INTERNAL:
		for i, task := range plan.Tasks {
			return func() error {
				golog.Verbosef("executing task %d on internal data store", i)
				response, err := s.Colony().Query(task.Query)
				rows := rqliter.NewRqlRows(response)
				if err != nil {
					golog.Errorf("could not execute internal query: %s", err.Error())
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
						golog.Errorf("could not scan row: %s", err.Error())
						return err
					}
					result = append(result, row)
				}
				if err := rows.Err(); err != nil {
					golog.Errorf("could not query internal store: %s", err.Error())
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
									field.DataTypeOID = uint32(core.Type_int8)
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
