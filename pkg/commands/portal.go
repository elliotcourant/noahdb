package commands

type ExecutePortal struct {
	Name  string
	Limit int
}

// Command Implements the command interface
func (ExecutePortal) Command() {}
