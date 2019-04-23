package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_SerializeType(t *testing.T) {
	bytes, err := serializeType(1321)
	assert.NoError(t, err)
	assert.NotEmpty(t, bytes)
}

func Test_DeserializeType(t *testing.T) {
	input := 1321
	bytes, err := serializeType(input)
	assert.NoError(t, err)
	assert.NotEmpty(t, bytes)
	result := 0
	err = deserializeType(bytes, &result)
	assert.NoError(t, err)
	assert.Equal(t, input, result)
}
