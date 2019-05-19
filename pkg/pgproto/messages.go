package pgproto

const (
	/*
		General Authentication
	*/

	// PgAuthentication is just a general authentication byte prefix.
	PgAuthentication = 'R'

	/*
		PostgreSQL Backend Messages
	*/

	PgAuthenticationOk                = 'R'
	PgAuthenticationKerberosV5        = 'R'
	PgAuthenticationCleartextPassword = 'R'
	PgAuthenticationMD5Password       = 'R'
	PgAuthenticationSCMCredential     = 'R'
	PgAuthenticationGSS               = 'R'
	PgAuthenticationSSPI              = 'R'
	PgAuthenticationGSSContinue       = 'R'
	PgAuthenticationSASL              = 'R'
	PgAuthenticationSASLContinue      = 'R'
	PgAuthenticationSASLFinal         = 'R'
	PgBackendKeyData                  = 'K'
	PgBindComplete                    = '2'
	PgCloseComplete                   = '3'
	PgCommandComplete                 = 'C'
	PgCopyInResponse                  = 'G'
	PgCopyOutResponse                 = 'H'
	PgCopyBothResponse                = 'W'
	PgDataRow                         = 'D'
	PgEmptyQueryResponse              = 'I'
	PgErrorResponse                   = 'E'
	PgFunctionCallResponse            = 'B'
	PgNegotiateProtocolVersion        = 'v'
	PgNoData                          = 'n'
	PgNoticeResponse                  = 'N'
	PgNotificationResponse            = 'A'
	PgParameterDescription            = 't'
	PgParameterStatus                 = 'S'
	PgParseComplete                   = '1'
	PgPortialSuspended                = 's'
	PgReadyForQuery                   = 'Z'
	PgRowDescription                  = 'T'

	/*
		PostgreSQL Frontend Messages
	*/

	// PgBind is the byte prefix for bind messages from the client.
	PgBind                = 'B'
	PgClose               = 'C'
	PgCopyFail            = 'f'
	PgDescribe            = 'D'
	PgExecute             = 'E'
	PgFlush               = 'H'
	PgFunctionCall        = 'F'
	PgGSSResponse         = 'p'
	PgParse               = 'P'
	PgPasswordMessage     = 'p'
	PgQuery               = 'Q'
	PgSASLInitialResponse = 'p'
	PgSASLResponse        = 'p'
	PgSync                = 'S'
	PgTerminate           = 'X'

	/*
		Raft Frontend Messages
	*/

	// RaftAppendEntriesRequest is the byte prefix for append entries from the leader.
	RaftAppendEntriesRequest   = 'a'
	RaftRequestVoteRequest     = 'v'
	RaftInstallSnapshotRequest = 'i'

	/*
		Raft Backend Messages
	*/

	// RaftRPCResponse is a generic tag
	RaftRPCResponse           = 'r'
	RaftAppendEntriesResponse = 'Y'

	/*
		RPC Backend Messages
	*/

	// RpcSequenceResponse indicates sequence metadata.
	RpcSequenceResponse = 'S'
	RpcJoinRequest      = 'j'

	// Both
	PgCopyData = 'd'
	PgCopyDone = 'c'
)
