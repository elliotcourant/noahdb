// Auto-generated - DO NOT EDIT

package pg_query

import (
	"fmt"
	"github.com/juju/errors"
	"strings"
)

func (node TransactionStmt) Deparse(ctx Context) (*string, error) {
	out := make([]string, 0)
	if kind, ok := transactionCmds[node.Kind]; !ok {
		return nil, errors.Errorf("couldn't deparse transaction kind: %d", node.Kind)
	} else {
		out = append(out, kind)
	}

	if node.Kind == TRANS_STMT_PREPARE ||
		node.Kind == TRANS_STMT_COMMIT_PREPARED ||
		node.Kind == TRANS_STMT_ROLLBACK_PREPARED {
		if node.Gid != nil {
			out = append(out, fmt.Sprintf("'%s'", *node.Gid))
		}
	} else {
		if node.Options.Items != nil && len(node.Options.Items) > 0 {

		}
	}

	result := strings.Join(out, " ")
	return &result, nil
}
