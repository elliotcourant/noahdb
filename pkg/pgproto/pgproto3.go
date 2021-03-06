package pgproto

import "fmt"

// Message is the interface implemented by an object that can decode and encode
// a particular PostgreSQL message.
type Message interface {
	// Decode is allowed and expected to retain a reference to data after
	// returning (unlike encoding.BinaryUnmarshaler).
	Decode(data []byte) error

	// Encode appends itself to dst and returns the new buffer.
	Encode(dst []byte) []byte
}

type InitialMessage interface {
	Message
	Initial()
}

type FrontendMessage interface {
	Message
	Frontend() // no-op method to distinguish frontend from backend methods
}

type RaftMessage interface {
	Message
	Raft() // no-op method to distinguish raft frontend from backend methods
}

type RpcFrontendMessage interface {
	Message
	Frontend()    // no-op method to distinguish frontend from backend methods
	RpcFrontend() // no-op method to distinguish raft frontend from backend methods
}

type BackendMessage interface {
	Message
	Backend() // no-op method to distinguish frontend from backend methods
}

type BackendMessages []BackendMessage

func (msgs BackendMessages) Decode(data []byte) error {
	return nil
}

func (msgs BackendMessages) Encode(dst []byte) []byte {
	res := make([]byte, 0)
	for _, msg := range msgs {
		res = append(res, msg.Encode(dst)...)
	}
	return res
}

// Implement backend.
func (BackendMessages) Backend() {}

type invalidMessageLenErr struct {
	messageType string
	expectedLen int
	actualLen   int
}

func (e *invalidMessageLenErr) Error() string {
	return fmt.Sprintf("%s body must have length of %d, but it is %d", e.messageType, e.expectedLen, e.actualLen)
}

type invalidMessageFormatErr struct {
	messageType string
}

func (e *invalidMessageFormatErr) Error() string {
	return fmt.Sprintf("%s body is invalid", e.messageType)
}
