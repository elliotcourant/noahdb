package static

import (
	"database/sql"
	_ "github.com/elliotcourant/noahdb/pkg/drivers/sqlite"
	"github.com/stretchr/testify/assert"
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
	bytes, err := ioutil.ReadFile("files/internal_sql.sql")
	assert.NoError(t, err)
	_, err = db.Exec(string(bytes))
	assert.NoError(t, err)
	_, err = db.Exec("INSERT INTO tenants (tenant_id, shard_id) VALUES(1, 2);")
	assert.Error(t, err)
}

func TestCreateMemoryDataStoreWithSchema_Transaction(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	bytes, err := ioutil.ReadFile("files/internal_sql.sql")
	assert.NoError(t, err)
	_, err = db.Exec(string(bytes))
	assert.NoError(t, err)
	tx, err := db.Begin()
	assert.NoError(t, err)
	_, err = tx.Exec("INSERT INTO shards (shard_id) VALUES(121);")
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
