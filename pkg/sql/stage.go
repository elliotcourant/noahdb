package sql

import (
	"github.com/elliotcourant/noahdb/pkg/ast"
	"github.com/elliotcourant/noahdb/pkg/commands"
	"github.com/readystock/golog"
	"time"
)

func (s *session) stageQueryToResult(statement ast.Stmt, result *commands.CommandResult) error {
	startTimestamp := time.Now()
	defer func() {
		golog.Debugf("execution of statement took %s", time.Since(startTimestamp))
	}()

	planner, err := getStatementHandler(statement)
	if err != nil {
		return err
	}

	plan := InitialPlan{}

	// Check to see if the provided statement can target noah's internal query interface.
	if noahPlanner, ok := planner.(NoahQueryPlanner); ok {
		// Try to build a noah query plan, if the query that was provided does actually use noah
		// tables then this will skip the standard planner and jump to expand the initial query plan
		if plan, ok, err = noahPlanner.getNoahQueryPlan(s); err != nil {
			return err
		} else if ok {
			goto ExpandInitialPlan
		}
	}

	if standardPlanner, ok := planner.(StandardQueryPlanner); ok {

	}

ExpandInitialPlan:

	expandedPlan, err := s.expandQueryPlan(plan)
	if err != nil {
		return err
	}

	return s.executeExpandedPlan(expandedPlan)
}
