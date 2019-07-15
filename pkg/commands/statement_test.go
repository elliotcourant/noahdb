package commands

import (
	"testing"
)

func TestBindStatement_Command(t *testing.T) {
	BindStatement{}.Command()
}

func TestDescribeStatement_Command(t *testing.T) {
	DescribeStatement{}.Command()
}

func TestExecuteStatement_Command(t *testing.T) {
	ExecuteStatement{}.Command()
}

func TestPrepareStatement_Command(t *testing.T) {
	PrepareStatement{}.Command()
}

func TestDeletePreparedStatement_Command(t *testing.T) {
	DeletePreparedStatement{}.Command()
}
