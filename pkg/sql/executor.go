package sql

import (
	"database/sql"
	"github.com/elliotcourant/noahdb/pkg/core"
	"github.com/elliotcourant/noahdb/pkg/drivers/rqliter"
	"github.com/elliotcourant/noahdb/pkg/pgproto"
	"github.com/elliotcourant/noahdb/pkg/pgwirebase"
	"github.com/elliotcourant/noahdb/pkg/types"
	"github.com/readystock/golog"
	"reflect"
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
					case *pgproto.ErrorResponse:
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
				var colTypes []*sql.ColumnType = nil
				for rows.Next() {
					if colTypes == nil {
						// colTypes, _ = rows.ColumnTypes()
					}
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

				for i, typ := range colTypes {
					field := pgproto.FieldDescription{
						Name:   typ.Name(),
						Format: int16(pgwirebase.FormatText),
					}
					t := reflect.New(typ.ScanType()).Interface()
					typs[i] = t
					switch t.(type) {
					case *int64:
						field.DataTypeOID = uint32(core.Type_int8)
					}
					rowDescription.Fields[i] = field
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
						case *int64:
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
