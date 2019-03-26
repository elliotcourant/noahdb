package pgwire

import (
	"bytes"
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

func NewTestWire() (io.ReadWriter, *wireServer, error) {
	buf := bytes.NewBuffer(nil)
	wire, err := newWire(buf)
	return buf, wire, err
}
