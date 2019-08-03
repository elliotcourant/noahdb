package ast

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestA_ArrayExpr_MarshalJSON(t *testing.T) {
	j, err := A_ArrayExpr{}.MarshalJSON()
	assert.NoError(t, err)
	assert.NotEmpty(t, j)
}
