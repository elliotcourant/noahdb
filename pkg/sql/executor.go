package sql

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/pgproto"
	"github.com/readystock/golog"
	"net"
	"time"
)

func (s *session) executeExpandedPlan(plan ExpandedPlan) error {
	startTimestamp := time.Now()
	defer func() {
		golog.Debugf("execution of statement took %s", time.Since(startTimestamp))
	}()

	responses := make(chan *pgproto.Frontend, len(plan.Tasks))

	for i, task := range plan.Tasks {
		go func(index int, task ExpandedPlanTask) {
			frontend := new(pgproto.Frontend)
			frontend = nil
			defer func() {
				responses <- frontend
			}()

			golog.Verbosef("preparing [%s] for data node [%d] shard [%d]", task.Query, task.DataNode.DataNodeID, task.Shard.ShardID)

			addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", task.DataNode.Address, task.DataNode.Port))
			if err != nil {
				golog.Errorf("could not resolve address for data node [%d]: %s", task.DataNode.DataNodeID, err.Error())
				return
			}

			conn, err := net.DialTCP("tcp", nil, addr)
			if err != nil {
				golog.Errorf("could not connect to data node [%d]: %s", task.DataNode.DataNodeID, err.Error())
				return
			}

			frontend, err = pgproto.NewFrontend(conn, conn)
			if err != nil {
				golog.Errorf("could not setup frontend for data node [%d]: %s", task.DataNode.DataNodeID, err.Error())
				return
			}

			if err := frontend.Send(&pgproto.StartupMessage{
				ProtocolVersion: pgproto.ProtocolVersionNumber,
				Parameters: map[string]string{
					"user": "postgres",
				},
			}); err != nil {
				golog.Errorf("could not send startup message to data node [%d]: %s", task.DataNode.DataNodeID, err.Error())
				return
			}
			_, _ = frontend.Receive()

			if err := frontend.Send(&pgproto.Query{
				String: task.Query,
			}); err != nil {
				golog.Errorf("could not send query to data node [%d]: %s", task.DataNode.DataNodeID, err.Error())
				return
			}

		}(i, task)
	}

	for i := 0; i < len(plan.Tasks); i++ {
		frontend := <-responses
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
				goto Done
			case *pgproto.ReadyForQuery:
				if canExit {
					goto Done
				}
			default:
				// Do nothing
			}
		}
	Done:
		frontend.Send(&pgproto.Terminate{})
	}

	return nil
}
