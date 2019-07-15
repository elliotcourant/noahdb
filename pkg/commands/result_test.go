package commands

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/ast"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateBindStatementResult(t *testing.T) {
	CreateBindStatementResult(nil)
}

func TestCreateDescribeStatementResult(t *testing.T) {
	CreateDescribeStatementResult(nil)
}

func TestCreateErrorResult(t *testing.T) {
	CreateErrorResult(nil, fmt.Errorf("test"))
}

func TestCreateExecuteCommandResult(t *testing.T) {
	stmt, err := ast.Parse("SELECT 1")
	assert.NoError(t, err)
	CreateExecuteCommandResult(nil, stmt.Statements[0].(ast.RawStmt).Stmt.(ast.SelectStmt))
}

func TestCreatePreparedStatementResult(t *testing.T) {
	stmt, err := ast.Parse("SELECT 1")
	assert.NoError(t, err)
	CreatePreparedStatementResult(nil, stmt.Statements[0].(ast.RawStmt).Stmt.(ast.SelectStmt))
}

func TestCreateSyncCommandResult(t *testing.T) {
	CreateSyncCommandResult(nil)
}

func TestCommandResult_SetError(t *testing.T) {
	CreateSyncCommandResult(nil).SetError(fmt.Errorf("test"))
}

func TestCommandResult_Err(t *testing.T) {
	err := fmt.Errorf("test")
	command := CreateSyncCommandResult(nil)
	command.SetError(err)
	assert.Equal(t, err, command.Err())
}

func TestCommandResult_SetNoDataMessage(t *testing.T) {
	CreateSyncCommandResult(nil).SetNoDataMessage(true)
}
