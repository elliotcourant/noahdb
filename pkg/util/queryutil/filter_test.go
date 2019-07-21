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

func Test_FindAccountIdsEx(t *testing.T) {
	t.Run("single account id", func(t *testing.T) {
		query := `SELECT p.id, p.account_id, p.sku FROM products p JOIN products p2 ON p2.id = p.id WHERE p.id = 1234 AND p2.account_id = 123421;`
		tree, _ := ast.Parse(query)
		ids, err := FindAccountIdsEx(tree, map[string]string{
			"products": "account_id",
		}, map[string][]string{
			"id":         {"product"},
			"account_id": {"product"},
			"sku":        {"sku"},
		})
		assert.NoError(t, err)
		assert.Equal(t, []uint64{123421}, ids, "the returned ids do not match expected values")
	})

	t.Run("multiple account ids", func(t *testing.T) {
		query := `SELECT p.id, p.account_id, p.sku FROM products p JOIN products p2 ON p2.id = p.id WHERE p.id = 1234 AND p2.account_id IN (12, 34, 21);`
		tree, _ := ast.Parse(query)
		ids, err := FindAccountIdsEx(tree, map[string]string{
			"products": "account_id",
		}, map[string][]string{
			"id":         {"product"},
			"account_id": {"product"},
			"sku":        {"sku"},
		})
		assert.NoError(t, err)
		assert.Equal(t, []uint64{12, 34, 21}, ids, "the returned ids do not match expected values")
	})

	t.Run("multiple account ids with bad IN", func(t *testing.T) {
		query := `SELECT p.id, p.account_id, p.sku FROM products p JOIN products p2 ON p2.id = p.id WHERE p.id = 1234 AND p2.account_id IN (12, 34, 21, '12a');`
		tree, _ := ast.Parse(query)
		ids, err := FindAccountIdsEx(tree, map[string]string{
			"products": "account_id",
		}, map[string][]string{
			"id":         {"product"},
			"account_id": {"product"},
			"sku":        {"sku"},
		})
		assert.EqualError(t, err, `strconv.ParseUint: parsing "12a": invalid syntax`)
		assert.Empty(t, ids)
	})
}
