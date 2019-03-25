package pgwire

import (
	"github.com/readystock/golog"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	golog.SetLevel("verbose")
	res := m.Run()
	os.Exit(res)
}
