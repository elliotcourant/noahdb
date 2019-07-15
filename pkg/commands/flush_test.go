package commands

import (
	"testing"
)

func TestFlush_Command(t *testing.T) {
	Flush{}.Command()
}
