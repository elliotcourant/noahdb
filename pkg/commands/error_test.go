package commands

import (
	"testing"
)

func TestSendError_Command(t *testing.T) {
	SendError{}.Command()
}
