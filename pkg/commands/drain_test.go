package commands

import (
	"testing"
)

func TestDrainRequest_Command(t *testing.T) {
	DrainRequest{}.Command()
}
