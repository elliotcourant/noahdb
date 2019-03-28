package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPluralize(t *testing.T) {
	assert.Equal(t, "", Pluralize(1))
	assert.Equal(t, "s", Pluralize(2))
}
