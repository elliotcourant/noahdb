package commands

type Sync struct {
}

// Command Implements the command interface
func (Sync) Command() {}

// Sync Implements the sync interface
func (Sync) Sync() {}
