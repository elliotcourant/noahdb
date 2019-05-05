package sql

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/ast"
	"github.com/readystock/golog"
	"time"
)

func (s *session) stageQueryToResult(statement ast.Stmt) error {
	planAndExpandTimestamp := time.Now()
	defer func() {
		golog.Debugf("planning and execution of statement took %s", time.Since(planAndExpandTimestamp))
	}()

	plan, err := func() (InitialPlan, error) {
		startTimestamp := time.Now()
		defer func() {
			golog.Debugf("initial planning of statement took %s", time.Since(startTimestamp))
		}()
		planner, err := getStatementHandler(statement)
		if err != nil {
			return InitialPlan{}, err
		}

		plan := InitialPlan{}

		// Check to see if the provided statement can target noah's internal query interface.
		if noahPlanner, ok := planner.(NoahQueryPlanner); ok {
			// Try to build a noah query plan, if the query that was provided does actually use noah
			// tables then this will skip the standard planner and jump to expand the initial query plan
			if plan, ok, err = noahPlanner.getNoahQueryPlan(s); err != nil {
				return InitialPlan{}, err
			} else if ok {
				return plan, nil
			}
		}

		// We want to try to address simple queries before we address standard ones, a simple query
		// can absolutely be handled in a standard query plan, but we want to try to return results
		// as fast as possible, so if a query is simple enough that it doesn't need to be sent to a
		// data node and can be addressed directly from noah then we want to prioritize that.
		if simplePlanner, ok := planner.(SimpleQueryPlanner); ok {
			if plan, ok, err = simplePlanner.getSimpleQueryPlan(s); err != nil {
				return InitialPlan{}, err
			} else if ok {
				return plan, nil
			}
		}

		// Check standard query plans.
		// Standard query plans are plans that target data nodes in the cluster.
		if standardPlanner, ok := planner.(StandardQueryPlanner); ok {
			// If a standard query planner is available then try to build a plan.
			if plan, ok, err = standardPlanner.getStandardQueryPlan(s); err != nil {
				return InitialPlan{}, err
			} else if ok {
				return plan, nil
			}
		}

		return InitialPlan{}, fmt.Errorf("could not generate plan for statement")
	}()

	expandedPlan, err := s.expandQueryPlan(plan)
	golog.Debugf("planning and expanding of statement took %s", time.Since(planAndExpandTimestamp))
	if err != nil {
		return err
	}

	return s.executeExpandedPlan(expandedPlan)
}
