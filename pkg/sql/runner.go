package sql

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/commands"
	"github.com/elliotcourant/timber"
	"io"
	"reflect"
)

func Run(stx sessionContext, terminateChannel chan bool) error {
	s := newSession(stx)
	for {
		select {
		case <-terminateChannel:
			timber.Debugf("terminating runner")
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
				timber.Debugf("found null command, advancing 1")
				s.StatementBuffer().AdvanceOne()
			}

			result := &commands.CommandResult{}
			switch cmd := c.(type) {
			case commands.ExecuteStatement:
				result = commands.CreateExecuteCommandResult(s.Backend(), cmd.Statement)
				err = s.executeStatement(cmd, result)
			case commands.ExecutePortal:
				fmt.Println("oh")
			case commands.PrepareStatement:
				result = commands.CreatePreparedStatementResult(s.Backend(), cmd.Statement)
				err = s.executePrepare(cmd, result)
			case commands.DescribeStatement:
				result = commands.CreateDescribeStatementResult(s.Backend())
				err = s.executeDescribe(cmd, result)
			case commands.BindStatement:
				result = commands.CreateBindStatementResult(s.Backend())
				err = s.executeBind(cmd, result)
			case commands.DeletePreparedStatement:
			case commands.SendError:
				result = commands.CreateErrorResult(s.Backend(), cmd.Err)
			case commands.Sync:
				result = commands.CreateSyncCommandResult(s.Backend())
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
