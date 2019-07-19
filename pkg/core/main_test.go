package core_test

import (
	"github.com/elliotcourant/timber"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	timber.SetLevel(timber.Level_Error)
	res := m.Run()
	os.Exit(res)
}
