package commands_test

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/ast"
	"github.com/elliotcourant/noahdb/pkg/commands"
	"github.com/elliotcourant/noahdb/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateBindStatementResult(t *testing.T) {
	commands.CreateBindStatementResult(nil)
}

func TestCreateDescribeStatementResult(t *testing.T) {
	commands.CreateDescribeStatementResult(nil)
}

func TestCreateErrorResult(t *testing.T) {
	commands.CreateErrorResult(nil, fmt.Errorf("test"))
}

func TestCreateExecuteCommandResult(t *testing.T) {
	stmt, err := ast.Parse("SELECT 1")
	assert.NoError(t, err)
	commands.CreateExecuteCommandResult(nil, stmt.Statements[0].(ast.RawStmt).Stmt.(ast.SelectStmt))
}

func TestCreatePreparedStatementResult(t *testing.T) {
	stmt, err := ast.Parse("SELECT 1")
	assert.NoError(t, err)
	commands.CreatePreparedStatementResult(nil, stmt.Statements[0].(ast.RawStmt).Stmt.(ast.SelectStmt))
}

func TestCreateSyncCommandResult(t *testing.T) {
	commands.CreateSyncCommandResult(nil)
}

func TestCommandResult_SetError(t *testing.T) {
	commands.CreateSyncCommandResult(nil).SetError(fmt.Errorf("test"))
}

func TestCommandResult_Err(t *testing.T) {
	err := fmt.Errorf("test")
	command := commands.CreateSyncCommandResult(nil)
	command.SetError(err)
	assert.Equal(t, err, command.Err())
}

func TestCommandResult_SetNoDataMessage(t *testing.T) {
	commands.CreateSyncCommandResult(nil).SetNoDataMessage(true)
}

func TestNewCommandResult(t *testing.T) {
	commands.NewCommandResult(nil)
}

func TestCommandResult_Close(t *testing.T) {
	t.Run("panic with no close type", func(t *testing.T) {
		assert.Panics(t, func() {
			_ = commands.NewCommandResult(testutils.CreateTestBackend(t)).Close()
		})
	})

	t.Run("bind", func(t *testing.T) {
		err := commands.CreateBindStatementResult(testutils.CreateTestBackend(t)).Close()
		assert.NoError(t, err)
	})

	t.Run("describe", func(t *testing.T) {
		err := commands.CreateDescribeStatementResult(testutils.CreateTestBackend(t)).Close()
		assert.NoError(t, err)
	})

	t.Run("describe w/ no data message", func(t *testing.T) {
		cmd := commands.CreateDescribeStatementResult(testutils.CreateTestBackend(t))
		cmd.SetNoDataMessage(true)
		err := cmd.Close()
		assert.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		err := commands.CreateErrorResult(testutils.CreateTestBackend(t), fmt.Errorf("test")).
			Close()
		assert.NoError(t, err)
	})

	t.Run("execute", func(t *testing.T) {
		stmt, err := ast.Parse("SELECT 1")
		assert.NoError(t, err)
		err = commands.CreateExecuteCommandResult(
			testutils.CreateTestBackend(t),
			stmt.Statements[0].(ast.RawStmt).Stmt.(ast.SelectStmt)).Close()
		assert.NoError(t, err)
	})

	t.Run("prepared", func(t *testing.T) {
		stmt, err := ast.Parse("SELECT 1")
		assert.NoError(t, err)
		err = commands.CreatePreparedStatementResult(
			testutils.CreateTestBackend(t),
			stmt.Statements[0].(ast.RawStmt).Stmt.(ast.SelectStmt)).Close()
		assert.NoError(t, err)
	})

	t.Run("sync", func(t *testing.T) {
		err := commands.CreateSyncCommandResult(testutils.CreateTestBackend(t)).Close()
		assert.NoError(t, err)
	})

	t.Run("already closed", func(t *testing.T) {
		cmd := commands.CreateSyncCommandResult(testutils.CreateTestBackend(t))
		err := cmd.Close()
		assert.NoError(t, err)

		err = cmd.Close()
		assert.Error(t, err)
	})

	t.Run("already closed", func(t *testing.T) {
		cmd := commands.CreateSyncCommandResult(testutils.CreateTestBackend(t))
		err := cmd.Close()
		assert.NoError(t, err)

		err = cmd.CloseWithErr(fmt.Errorf("test"))
		assert.Error(t, err)
	})
}

func TestCommandResult_CloseWithErr(t *testing.T) {
	err := commands.NewCommandResult(testutils.CreateTestBackend(t)).
		CloseWithErr(fmt.Errorf("test"))
	assert.NoError(t, err)
}
