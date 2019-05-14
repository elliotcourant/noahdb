package core

import (
	"github.com/readystock/golog"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	golog.SetLevel("verbose")
	res := m.Run()
	os.Exit(res)
}

func newTestColony(joinAddresses ...string) (Colony, func()) {
	tempdir, err := ioutil.TempDir("", "core-temp")
	if err != nil {
		panic(err)
	}

	joins := strings.Join(joinAddresses, ",")
	colony, _, err := NewColony(tempdir, joins, ":")
	if err != nil {
		panic(err)
	}

	return colony, func() {
		if err := os.RemoveAll(tempdir); err != nil {
			panic(err)
		}
	}
}
