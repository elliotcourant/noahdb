package sql

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/commands"
	"github.com/elliotcourant/noahdb/pkg/pgerror"
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
				err = s.executeStatement(cmd.Statement, result, nil)
			case commands.ExecutePortal:
				// Make sure the portal exists, if it doesn't then we want to break early.
				portal, ok := s.portals[cmd.Name]
				if !ok {
					err = pgerror.NewErrorf(pgerror.CodeInvalidCursorNameError,
						"unknown portal [%s]", cmd.Name)
					break
				}

				// At this point the portal exists, but we need to make sure that the query is valid
				if portal.Stmt.Statement == nil {
					result = commands.CreateEmptyQueryResult(s.Backend())
					break
				}

				result = commands.CreateExecutePortalResult(s.Backend(), portal.Stmt.Statement)
				err = s.executeStatement(portal.Stmt.Statement, result, nil)
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
