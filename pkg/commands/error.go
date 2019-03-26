package commands

type SendError struct {
	Err error
}

// Implements the command interface
func (SendError) Command() {}
