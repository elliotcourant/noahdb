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

	// We want to try to address simple queries before we address standard ones, a simple query
	// can absolutely be handled in a standard query plan, but we want to try to return results
	// as fast as possible, so if a query is simple enough that it doesn't need to be sent to a
	// data node and can be addressed directly from noah then we want to prioritize that.
	if simplePlanner, ok := planner.(SimpleQueryPlanner); ok {
		if plan, ok, err = simplePlanner.getSimpleQueryPlan(s); err != nil {
			return err
		} else if ok {
			goto ExpandInitialPlan
		}
	}

	// Check standard query plans.
	// Standard query plans are plans that target data nodes in the cluster.
	if standardPlanner, ok := planner.(StandardQueryPlanner); ok {
		// If a standard query planner is available then try to build a plan.
		if plan, ok, err = standardPlanner.getStandardQueryPlan(s); err != nil {
			return err
		} else if ok {
			goto ExpandInitialPlan
		}
	}

ExpandInitialPlan:

	expandedPlan, err := s.expandQueryPlan(plan)
	if err != nil {
		return err
	}

	return s.executeExpandedPlan(expandedPlan)
}
