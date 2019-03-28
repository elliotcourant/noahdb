package commands

type SendError struct {
	Err error
}

// Command Implements the command interface
func (SendError) Command() {}
