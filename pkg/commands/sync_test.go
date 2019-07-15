package commands

import (
	"testing"
)

func TestSync_Command(t *testing.T) {
	Sync{}.Command()
}

func TestSync_Sync(t *testing.T) {
	Sync{}.Sync()
}
