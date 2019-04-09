package pgwire

import (
	"bytes"
	"github.com/elliotcourant/noahdb/pkg/core"
	"github.com/readystock/golog"
	"io"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	golog.SetLevel("verbose")
	res := m.Run()
	os.Exit(res)
}

func NewTestWire(colony core.Colony) (io.ReadWriter, *wireServer, error) {
	buf := bytes.NewBuffer(nil)
	wire, err := newWire(colony, buf)
	return buf, wire, err
}
