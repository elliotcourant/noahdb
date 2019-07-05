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
		golog.Verbosef("[%s] planning and execution of statement", time.Since(planAndExpandTimestamp))
	}()

	plan, err := func() (InitialPlan, error) {
		startTimestamp := time.Now()
		defer func() {
			golog.Verbosef("[%s] initial planning of statement", time.Since(startTimestamp))
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

		if simplePlanner, ok := planner.(SimpleQueryPlanner); ok {
			if plan, ok, err = simplePlanner.getSimpleQueryPlan(s); err != nil {
				return InitialPlan{}, err
			} else if ok {
				return plan, nil
			}
		}

		return InitialPlan{}, fmt.Errorf("could not generate plan for statement")
	}()

	if err != nil {
		return err
	}

	expandedPlan, err := s.expandQueryPlan(plan)
	golog.Verbosef("[%s] planning and expanding of statement", time.Since(planAndExpandTimestamp))
	if err != nil {
		return err
	}

	return s.executeExpandedPlan(expandedPlan)
}
