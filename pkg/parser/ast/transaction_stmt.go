// Auto-generated from postgres/src/include/nodes/parsenodes.h - DO NOT EDIT

package pg_query

import "encoding/json"

/* ----------------------
 *		{Begin|Commit|Rollback} Transaction Statement
 * ----------------------
 */
type TransactionStmt struct {
	Kind    TransactionStmtKind `json:"kind"`    /* see above */
	Options List                `json:"options"` /* for BEGIN/START and savepoint commands */
	Gid     *string             `json:"gid"`     /* for two-phase-commit related commands */
}

var transactionCmds = map[TransactionStmtKind]string{
	TRANS_STMT_BEGIN:             "BEGIN",
	TRANS_STMT_START:             "BEGIN",
	TRANS_STMT_COMMIT:            "COMMIT",
	TRANS_STMT_ROLLBACK:          "ROLLBACK",
	TRANS_STMT_SAVEPOINT:         "SAVEPOINT",
	TRANS_STMT_RELEASE:           "RELEASE",
	TRANS_STMT_ROLLBACK_TO:       "ROLLBACK",
	TRANS_STMT_PREPARE:           "PREPARE TRANSACTION",
	TRANS_STMT_COMMIT_PREPARED:   "COMMIT PREPARED",
	TRANS_STMT_ROLLBACK_PREPARED: "ROLLBACK PREPARED",
}

func (node TransactionStmt) MarshalJSON() ([]byte, error) {
	type TransactionStmtMarshalAlias TransactionStmt
	return json.Marshal(map[string]interface{}{
		"TransactionStmt": (*TransactionStmtMarshalAlias)(&node),
	})
}

func (node *TransactionStmt) UnmarshalJSON(input []byte) (err error) {
	var fields map[string]json.RawMessage

	err = json.Unmarshal(input, &fields)
	if err != nil {
		return
	}

	if fields["kind"] != nil {
		err = json.Unmarshal(fields["kind"], &node.Kind)
		if err != nil {
			return
		}
	}

	if fields["options"] != nil {
		node.Options.Items, err = UnmarshalNodeArrayJSON(fields["options"])
		if err != nil {
			return
		}
	}

	if fields["gid"] != nil {
		err = json.Unmarshal(fields["gid"], &node.Gid)
		if err != nil {
			return
		}
	}

	return
}
