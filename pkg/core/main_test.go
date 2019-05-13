package core

import (
	"github.com/readystock/golog"
	"io/ioutil"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	golog.SetLevel("verbose")
	res := m.Run()
	os.Exit(res)
}

func newTestColony() (Colony, func()) {
	tempdir, err := ioutil.TempDir("", "core-temp")
	if err != nil {
		panic(err)
	}

	colony, _, err := NewColony(tempdir, "", ":")
	if err != nil {
		panic(err)
	}

	return colony, func() {
		if err := os.RemoveAll(tempdir); err != nil {
			panic(err)
		}
	}
}
