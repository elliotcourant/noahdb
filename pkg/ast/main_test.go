package ast

import (
	"github.com/readystock/golog"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	golog.SetLevel(os.Getenv("GOLOG_LEVEL"))
	retCode := m.Run()
	os.Exit(retCode)
}
