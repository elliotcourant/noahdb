package static

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"testing"
)

func TestCreateMemoryDataStore(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	row, err := db.Query("SELECT 142132")
	assert.NoError(t, err)
	row.Next()
	intVal := 0
	row.Scan(&intVal)
	assert.Equal(t, 142132, intVal)
}

func TestCreateMemoryDataStoreWithSchema_FKTest(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	bytes, err := ioutil.ReadFile("files/00_internal_sql.sql")
	assert.NoError(t, err)
	_, err = db.Exec(string(bytes))
	assert.NoError(t, err)
	_, err = db.Exec("INSERT INTO tenants (tenant_id, shard_id) VALUES(1, 2);")
	assert.Error(t, err)
}

func TestCreateMemoryDataStoreWithSchema_Transaction(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	bytes, err := ioutil.ReadFile("files/00_internal_sql.sql")
	assert.NoError(t, err)
	_, err = db.Exec(string(bytes))
	assert.NoError(t, err)
	tx, err := db.Begin()
	assert.NoError(t, err)
	_, err = tx.Exec("INSERT INTO shards (shard_id, state) VALUES(121, 0);")
	assert.NoError(t, err)
	assert.NoError(t, tx.Commit())
	_, err = db.Exec("INSERT INTO tenants (tenant_id, shard_id) VALUES(1132, 8);")
	assert.Error(t, err)
	row, err := db.Query("SELECT shard_id FROM shards")
	assert.NoError(t, err)
	row.Next()
	intVal := 0
	row.Scan(&intVal)
	assert.Equal(t, 121, intVal)
}

func TestCreateJsonFile(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)

	bytes, err := ioutil.ReadFile("files/00_internal_sql.sql")
	assert.NoError(t, err)

	_, err = db.Exec(string(bytes))
	assert.NoError(t, err)

	type Thing struct {
		Oid     uint16   `json:"oid"`
		Name    string   `json:"name"`
		Aliases []string `json:"aliases"`
	}

	typeRow, err := db.Query("SELECT type_id, type_name FROM types")
	assert.NoError(t, err)

	things := make([]Thing, 0)
	for typeRow.Next() {
		x := Thing{Aliases: make([]string, 0)}
		typeRow.Scan(&x.Oid, &x.Name)
		if x.Oid == 0 {
			continue
		}
		things = append(things, x)
	}

	for i, x := range things {
		aliasRow, err := db.Query(fmt.Sprintf("SELECT alias_name FROM type_aliases WHERE type_id = %d", x.Oid))
		assert.NoError(t, err)

		for aliasRow.Next() {
			alias := ""
			aliasRow.Scan(&alias)
			things[i].Aliases = append(things[i].Aliases, alias)
		}
	}
	j, _ := yaml.Marshal(things)

	// j, _ := json.MarshalIndent(things, "  ", "  ")
	fmt.Println(string(j))
}
