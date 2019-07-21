package sql_test

import (
	"database/sql"
	"fmt"
	"github.com/elliotcourant/noahdb/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSelect(t *testing.T) {
	colony, cleanup := testutils.NewPgTestColony(t)
	defer cleanup()
	func() {
		db, err := sql.Open("postgres", testutils.ConnectionString(colony.Addr()))
		if err != nil {
			panic(err)
		}
		defer db.Close()

		_, err = db.Exec(`CREATE TABLE accounts (id BIGSERIAL PRIMARY KEY, name TEXT) TABLESPACE "noah.tenants"`)
		if !assert.NoError(t, err) {
			panic(err)
		}

		_, err = db.Exec(`CREATE TABLE products (id BIGSERIAL PRIMARY KEY, account_id BIGINT NOT NULL REFERENCES accounts (id), sku TEXT) TABLESPACE "noah.sharded"`)
		if !assert.NoError(t, err) {
			panic(err)
		}

		// Create some accounts.
		rows, err := db.Query(`INSERT INTO accounts (name) VALUES ('account one'), ('account two') RETURNING id;`)
		if !assert.NoError(t, err) {
			panic(err)
		}
		accountIds := make([]uint64, 0)
		for rows.Next() {
			var id uint64
			if err := rows.Scan(&id); !assert.NoError(t, err) {
				panic(err)
			}
			accountIds = append(accountIds, id)
		}
		if err := rows.Err(); !assert.NoError(t, err) {
			panic(err)
		}

		for _, accountId := range accountIds {
			_, err = db.Exec(fmt.Sprintf(`INSERT INTO products (account_id, sku) VALUES (%d, 'SKU%d001'), (%d, 'SKU%d002');`, accountId, accountId, accountId, accountId))
			if !assert.NoError(t, err) {
				panic(err)
			}
		}

		type Product struct {
			ID        uint64
			AccountID uint64
			SKU       string
		}

		t.Run("single account", func(t *testing.T) {
			for _, accountId := range accountIds {
				rows, err := db.Query(fmt.Sprintf(`SELECT p.id, p.account_id, p.sku FROM products p JOIN products p2 ON p2.id = p.id WHERE p.account_id = %d;`, accountId))
				if !assert.NoError(t, err) {
					panic(err)
				}

				products := make([]Product, 0)
				// Read the product rows
				for rows.Next() {
					var product Product
					if err := rows.Scan(
						&product.ID,
						&product.AccountID,
						&product.SKU,
					); !assert.NoError(t, err) {
						panic(err)
					}
					products = append(products, product)
				}
				if err := rows.Err(); !assert.NoError(t, err) {
					panic(err)
				}
				assert.NotEmpty(t, products)
			}
		})

		t.Run("without account", func(t *testing.T) {
			_, err := db.Query(`SELECT p.id, p.account_id, p.sku FROM products p JOIN products p2 ON p2.id = p.id;`)
			if !assert.EqualError(t, err, "pq: cannot query sharded tables without specifying a tenant ID") {
				panic("should have encountered error")
			}
		})

		t.Run("with multiple accounts", func(t *testing.T) {
			_, err := db.Query(`SELECT p.id, p.account_id, p.sku FROM products p JOIN products p2 ON p2.id = p.id WHERE p.account_id IN (1, 2, 3, 4);`)
			if !assert.EqualError(t, err, "pq: cannot query sharded tables for multiple tenants") {
				panic("should have encountered error")
			}
		})
	}()
}
