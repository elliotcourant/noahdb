package plan

type TransactionType int

const (
	TransactionTypeNone TransactionType = iota + 1
	TransactionTypeCommit
	TransactionTypeRollback
)
