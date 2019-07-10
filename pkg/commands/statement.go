package commands

import (
	"github.com/elliotcourant/noahdb/pkg/ast"
	"github.com/elliotcourant/noahdb/pkg/pgwirebase"
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
	Name string
	Type pgwirebase.PrepareType
}

// Command Implements the command interface
func (DescribeStatement) Command() {}

// BindStatement is the Command for creating a portal from a prepared statement.
type BindStatement struct {
	PreparedStatementName string
	PortalName            string
	// OutFormats contains the requested formats for the output columns.
	// It either contains a bunch of format codes, in which case the number will
	// need to match the number of output columns of the portal, or contains a single
	// code, in which case that code will be applied to all columns.
	OutFormats []pgwirebase.FormatCode
	// Args are the arguments for the prepared statement.
	// They are passed in without decoding because decoding requires type
	// inference to have been performed.
	//
	// A nil element means a tree.DNull argument.
	Args [][]byte
	// ArgFormatCodes are the codes to be used to deserialize the Args.
	// It either contains a bunch of format codes, in which case the number will
	// need to match the number of arguments for the portal, or contains a single
	// code, in which case that code will be applied to all arguments.
	ArgFormatCodes []pgwirebase.FormatCode
}

// Command Implements the command interface
func (BindStatement) Command() {}

type DeletePreparedStatement struct {
}

// Command Implements the command interface
func (DeletePreparedStatement) Command() {}
