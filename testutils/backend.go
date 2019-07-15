package testutils

import (
	"bytes"
	"github.com/elliotcourant/noahdb/pkg/pgproto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func CreateTestBackend(t *testing.T) *pgproto.Backend {
	buffer := bytes.NewBuffer(make([]byte, 0))
	b, err := pgproto.NewBackend(buffer, buffer)
	assert.NoError(t, err)
	return b
}
