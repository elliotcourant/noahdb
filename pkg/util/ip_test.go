package util

import (
	"github.com/readystock/golog"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ExternalIP(t *testing.T) {
	ExternalIP() // Idk how to really do anything here, its hard to simulate all the scenarios.
}

func Test_ResolvedLocalAddress(t *testing.T) {
	input := ":5433"
	addr, err := ResolveAddress(input)
	assert.NoError(t, err)
	golog.Infof("resulting address: %s", addr)
}
