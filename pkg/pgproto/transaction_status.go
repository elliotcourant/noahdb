package pgproto

type TransactionStatus byte

const (
	// TransactionStatus_Idle means the session is outside of a transaction.
	TransactionStatus_Idle TransactionStatus = 'I'
	// TransactionStatus_In means the session is inside a transaction.
	TransactionStatus_In TransactionStatus = 'T'
	// TransactionStatus_InFailed means the session is inside a transaction, but the
	// transaction is in the Aborted state.
	TransactionStatus_InFailed TransactionStatus = 'E'
)
