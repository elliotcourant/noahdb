package sql

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/ast"
)

type transactionStmtPlanner struct {
	tree ast.TransactionStmt
}

// NewCreateStatementPlan creates a new planner for create statements.
func newTransactionStatementPlan(tree ast.TransactionStmt) *transactionStmtPlanner {
	return &transactionStmtPlanner{
		tree: tree,
	}
}

// getTransactionQueryPlan will return false if nothing needs to be sent to any data nodes at this
// time. It will return true and include a plan with the distributed plan type if queries need to be
// distributed to the data nodes.
func (stmt *transactionStmtPlanner) getTransactionQueryPlan(s *session) (InitialPlan, bool, error) {
	switch stmt.tree.Kind {
	case ast.TRANS_STMT_BEGIN, ast.TRANS_STMT_START:
		switch s.GetTransactionState() {
		case TransactionState_None:
			s.SetTransactionState(TransactionState_Active)
		default:
			// Already in a transaction
			return InitialPlan{}, false, fmt.Errorf("transaction already active")
		}
		return InitialPlan{}, false, nil
	case ast.TRANS_STMT_COMMIT:
		switch s.GetTransactionState() {
		case TransactionState_Active:
			return InitialPlan{
				Types: map[PlanType]InitialPlanTask{
					PlanType_WRITE: {
						Query: "COMMIT",
						Type:  stmt.tree.StatementType(),
					},
				},
				ShardID:      0,
				DistPlanType: DistributedPlanType_COMMIT,
			}, true, nil
		default:
			// No transaction
			return InitialPlan{}, false, fmt.Errorf("no active transaction")
		}
	case ast.TRANS_STMT_ROLLBACK:
		switch s.GetTransactionState() {
		case TransactionState_Active:
			return InitialPlan{
				Types: map[PlanType]InitialPlanTask{
					PlanType_WRITE: {
						Query: "ROLLBACK",
						Type:  stmt.tree.StatementType(),
					},
				},
				ShardID:      0,
				DistPlanType: DistributedPlanType_ROLLBACK,
			}, true, nil
		default:
			// No transaction
			return InitialPlan{}, false, fmt.Errorf("no active transaction")
		}
	default:
		return InitialPlan{}, false, fmt.Errorf("could not handle transaction type [%s]", stmt.tree.Kind)
	}
}
