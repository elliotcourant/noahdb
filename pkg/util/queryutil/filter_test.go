package queryutil

import (
	"github.com/elliotcourant/noahdb/pkg/ast"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_FindAccountIds(t *testing.T) {
	query := `SELECT * FROM users INNER JOIN accounts ON users.account_id=accounts.account_id WHERE users.account_id = 5`
	tree, _ := ast.Parse(query)
	ids := FindAccountIds(tree, "account_id")
	assert.Equal(t, []uint64{5}, ids, "the returned ids do not match expected values")
}
