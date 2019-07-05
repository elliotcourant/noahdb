package sql

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/commands"
	"github.com/readystock/golog"
	"io"
	"reflect"
)

func Run(stx sessionContext, terminateChannel chan bool) error {
	s := newSession(stx)
	for {
		select {
		case <-terminateChannel:
			golog.Debugf("terminating runner")
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
				result = commands.CreateExecuteCommandResult(s, cmd.Statement)
				err = s.ExecuteStatement(cmd, result)
			case commands.ExecutePortal:
			case commands.PrepareStatement:
				result = commands.CreatePreparedStatementResult(s, cmd.Statement)
				err = s.ExecutePrepare(cmd, result)
			case commands.DescribeStatement:
				result = commands.CreateDescribeStatementResult(s)
				err = s.ExecuteDescribe(cmd, result)
			case commands.BindStatement:
			case commands.DeletePreparedStatement:
			case commands.SendError:
				result = commands.CreateErrorResult(s, cmd.Err)
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
