package commands

import (
	"github.com/elliotcourant/noahdb/pkg/ast"
	"github.com/elliotcourant/noahdb/pkg/types"
	"github.com/elliotcourant/noahdb/pkg/util/queryutil"
)

type ExecuteStatement struct {
	Statement ast.Stmt
}

// Command Implements the command interface
func (ExecuteStatement) Command() {}

type PrepareStatement struct {
	Name string
	// Stmt can be nil, in which case executing it should produce an "empty query
	// response" message.
	TypeHints queryutil.PlaceholderTypes
	// RawTypeHints is the representation of type hints exactly as specified by
	// the client.
	RawTypeHints []types.OID

	Statement ast.Stmt
}

// Command Implements the command interface
func (PrepareStatement) Command() {}

type DescribeStatement struct {
}

// Command Implements the command interface
func (DescribeStatement) Command() {}

type BindStatement struct {
}

// Command Implements the command interface
func (BindStatement) Command() {}

type DeletePreparedStatement struct {
}

// Command Implements the command interface
func (DeletePreparedStatement) Command() {}
