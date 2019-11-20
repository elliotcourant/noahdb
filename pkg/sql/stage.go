package sql

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/ast"
	"github.com/elliotcourant/noahdb/pkg/pgwirebase"
	"github.com/elliotcourant/noahdb/pkg/util/queryutil"
	"time"
)

func (s *session) stageQueryToResult(
	statement ast.Stmt,
	placeholders queryutil.QueryArguments,
	outFormats []pgwirebase.FormatCode) error {
	// If there are placeholders present then we need to walk the syntax tree and add the
	// placeholders into the query manually, this is a bit weird and kind of expensive. But it's
	// the best solution I have at the moment for the query planner.
	if placeholders != nil {
		statement = queryutil.ReplaceArguments(statement, placeholders).(ast.Stmt)
	}

	planAndExpandTimestamp := time.Now()
	defer func() {
		s.log.Verbosef("[%s] planning and execution of statement", time.Since(planAndExpandTimestamp))
	}()

	plan, sendToNodes, err := func() (InitialPlan, bool, error) {
		startTimestamp := time.Now()
		defer func() {
			s.log.Verbosef("[%s] initial planning of statement", time.Since(startTimestamp))
		}()
		planner, err := getStatementHandler(statement)
		if err != nil {
			return InitialPlan{}, false, err
		}

		plan := InitialPlan{}

		if transactionPlanner, ok := planner.(TransactionQueryPlanner); ok {
			transactionPlan, sendToNodes, err := transactionPlanner.getTransactionQueryPlan(s)
			if err != nil {
				return transactionPlan, false, err
			}
			return transactionPlan, sendToNodes, nil
		}

		// Check to see if the provided statement can target noah's internal query interface.
		if noahPlanner, ok := planner.(NoahQueryPlanner); ok {
			// Try to build a noah query plan, if the query that was provided does actually use noah
			// tables then this will skip the standard planner and jump to expand the initial query plan
			if plan, ok, err = noahPlanner.getNoahQueryPlan(s); err != nil {
				return InitialPlan{}, false, err
			} else if ok {
				return plan, ok, nil
			}
		}

		if normalQueryPlanner, ok := planner.(QueryPlanner); ok {
			if plan, ok, err = normalQueryPlanner.GetQueryPlan(s); err != nil {
				return InitialPlan{}, false, err
			} else if ok {
				return plan, ok, nil
			}
		}

		return InitialPlan{}, false, fmt.Errorf("could not generate plan for statement")
	}()

	if err != nil {
		return err
	}

	if !sendToNodes {
		return nil
	}

	expandedPlan, err := s.expandQueryPlan(plan)
	s.log.Verbosef("[%s] planning and expanding of statement", time.Since(planAndExpandTimestamp))
	if err != nil {
		return err
	}

	expandedPlan.OutFormats = outFormats

	return s.executeExpandedPlan(expandedPlan)
}
