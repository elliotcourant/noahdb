package pgwire

import (
	"github.com/elliotcourant/timber"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	timber.SetLevel(timber.Level_Trace)
	res := m.Run()
	os.Exit(res)
}
