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

// func TestA_ArrayExpr_UnmarshalJSON(t *testing.T) {
// 	a := A_ArrayExpr{
// 		Location: 1,
// 	}
// 	j, err := a.MarshalJSON()
// 	assert.NoError(t, err)
// 	assert.NotEmpty(t, j)
//
// 	b := &A_ArrayExpr{}
// 	var fields map[string]json.RawMessage
// 	err = json.Unmarshal(j, &fields)
// 	assert.NoError(t, err)
// 	assert.NotEmpty(t, fields)
//
// 	err = b.UnmarshalJSON([]byte(fields["A_ArrayExpr"]))
// 	assert.NoError(t, err)
// 	assert.Equal(t, a, *b)
// }
