package testutils

import (
	"bytes"
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/pgproto"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

type BufferHealth int

const (
	BufferHealthy   BufferHealth = 0
	BufferUnhealthy              = 1
)

func CreateTestBackend(t *testing.T) *pgproto.Backend {
	return CreateTextBackendEx(t, BufferHealthy)
}

func CreateTextBackendEx(t *testing.T, health BufferHealth) *pgproto.Backend {
	var buffer io.ReadWriter
	switch health {
	case BufferHealthy:
		buffer = bytes.NewBuffer(make([]byte, 0))
	case BufferUnhealthy:
		buffer = &badBuffer{}
	}
	b, err := pgproto.NewBackend(buffer, buffer)
	assert.NoError(t, err)
	return b
}

type badBuffer struct {
	// Praise Rickster
}

func (*badBuffer) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("this is a bad buffer")
}

func (*badBuffer) Write(p []byte) (n int, err error) {
	return 0, fmt.Errorf("this is a bad buffer")
}
