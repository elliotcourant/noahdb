package sql

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/commands"
	"github.com/elliotcourant/noahdb/pkg/pgproto"
	"github.com/readystock/golog"
	"io"
	"reflect"
)

func Run(stx sessionContext, terminateChannel chan bool) error {
	s := newSession(stx)
	for {
		select {
		case <-terminateChannel:
			golog.Debugf("terminating ingress runner")
			return nil
		default:
			c, _, err := s.StatementBuffer().CurrentCommand()
			if err != nil {
				if err == io.EOF {
					return nil
				}
				return err
			}
			if c == nil {
				golog.Debugf("found null command, advancing 1")
				s.StatementBuffer().AdvanceOne()
			}
			result := &commands.CommandResult{}
			switch cmd := c.(type) {
			case commands.ExecuteStatement:
				result = commands.CreateExecuteCommandResult(s)
				err = s.ExecuteStatement(cmd, result)
				// val := types.Int4{}
				// val.Set(1)
				// bytes, _ := val.EncodeText(types.NewConnInfo(), nil)
				// err = s.Backend().Send(pgproto.BackendMessages{
				// 	&pgproto.RowDescription{
				// 		Fields: []pgproto.FieldDescription{
				// 			{
				// 				Name:                 "",
				// 				TableOID:             0,
				// 				TableAttributeNumber: 0,
				// 				DataTypeOID:          types.Int2OID,
				// 				DataTypeSize:         4,
				// 				TypeModifier:         0,
				// 				Format:               pgproto.TextFormat,
				// 			},
				// 		},
				// 	},
				// 	&pgproto.DataRow{
				// 		Values: [][]byte{
				// 			bytes,
				// 		},
				// 	},
				// 	&pgproto.CommandComplete{
				// 		CommandTag: "SELECT 1",
				// 	},
				// })
				// // err = session.Backend().Send()
				// // err = session.Backend().Send()
			case commands.ExecutePortal:
			case commands.PrepareStatement:
				err = s.ExecutePrepare(cmd, result)
			case commands.DescribeStatement:
			case commands.BindStatement:
			case commands.DeletePreparedStatement:
			case commands.SendError:
				result = commands.CreateErrorResult(s, cmd.Err)
				err = s.Backend().Send(&pgproto.ErrorResponse{
					Message: cmd.Err.Error(),
				})
			case commands.Sync:
				result = commands.CreateSyncCommandResult(s)
			case commands.Flush:
			case commands.CopyIn:
			default:
				panic(fmt.Sprintf("unsupported command type [%s]", reflect.TypeOf(cmd).Name()))
			}

			if err != nil {
				if err = result.CloseWithErr(err); err != nil {
					return err
				}
				if err = s.StatementBuffer().SeekToNextBatch(); err != nil {
					return err
				}
			} else {
				if resultError := result.Err(); resultError != nil {
					if err := result.CloseWithErr(resultError); err != nil {
						return err
					}
				} else {
					if err := result.Close(); err != nil {
						return err
					}
				}
				s.StatementBuffer().AdvanceOne()
			}
		}
	}
}
