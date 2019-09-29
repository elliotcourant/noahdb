package executor

type TransactionState int

const (
	TransactionStateNone TransactionState = iota + 1
	TransactionStateActive
)
