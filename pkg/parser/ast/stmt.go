package pg_query

type Stmt interface {
	StatementType() StmtType
	StatementTag() string
	Deparse(ctx Context) (*string, error)
}

func (node CreateStmt) StatementType() StmtType { return DDL }

func (node CreateStmt) StatementTag() string { return "CREATE TABLE" }

func (node DeleteStmt) StatementType() StmtType {
	if node.ReturningList.Items != nil && len(node.ReturningList.Items) > 0 {
		return Rows
	} else {
		return RowsAffected
	}
}

func (node DeleteStmt) StatementTag() string { return "DELETE" }

func (node DropStmt) StatementType() StmtType { return DDL }

func (node DropStmt) StatementTag() string { return "DROP TABLE" }

func (node InsertStmt) StatementType() StmtType {
	if node.ReturningList.Items != nil && len(node.ReturningList.Items) > 0 {
		return Rows
	} else {
		return RowsAffected
	}
}

func (node InsertStmt) StatementTag() string { return "INSERT" }

func (node SelectStmt) StatementType() StmtType { return Rows }

func (node SelectStmt) StatementTag() string { return "SELECT" }

func (node UpdateStmt) StatementType() StmtType {
	if node.ReturningList.Items != nil && len(node.ReturningList.Items) > 0 {
		return Rows
	} else {
		return RowsAffected
	}
}

func (node UpdateStmt) StatementTag() string { return "UPDATE" }

func (node TransactionStmt) StatementType() StmtType { return Ack }

func (node TransactionStmt) StatementTag() string { return transactionCmds[node.Kind] }

func (node VariableSetStmt) StatementType() StmtType { return Ack }

func (node VariableSetStmt) StatementTag() string { return "SET" }

func (node VariableShowStmt) StatementType() StmtType { return Ack }

func (node VariableShowStmt) StatementTag() string { return "SHOW" }
