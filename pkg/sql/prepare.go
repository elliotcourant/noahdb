package sql

import (
	"github.com/elliotcourant/noahdb/pkg/ast"
	"github.com/elliotcourant/noahdb/pkg/commands"
	"github.com/elliotcourant/noahdb/pkg/pgerror"
	"github.com/elliotcourant/noahdb/pkg/util/queryutil"
)

func (s *session) AddPreparedStatement(name string, stmt ast.Stmt, parseTypeHints queryutil.PlaceholderTypes) *PreparedStatement {
	prepared := &PreparedStatement{
		TypeHints: parseTypeHints,
		Statement: &stmt,
	}
	s.preparedStatements[name] = preparedStatementEntry{
		PreparedStatement: prepared,
		portals:           make(map[string]struct{}),
	}
	return prepared
}

func (s *session) HasPreparedStatement(name string) bool {
	_, ok := s.preparedStatements[name]
	return ok
}

func (s *session) DeletePreparedStatement(name string) {
	psEntry, ok := s.preparedStatements[name]
	if !ok {
		return
	}
	for portalName := range psEntry.portals {
		s.DeletePortal(portalName)
	}
	delete(s.preparedStatements, name)
}

func (s *session) DeletePortal(name string) {
	portalEntry, ok := s.portals[name]
	if !ok {
		return
	}
	delete(s.portals, name)
	delete(s.preparedStatements[portalEntry.psName].portals, name)
}

func (s *session) ExecutePrepare(prepare commands.PrepareStatement, result *commands.CommandResult) error {
	if prepare.Name != "" {
		if s.HasPreparedStatement(prepare.Name) {
			return pgerror.NewErrorf(pgerror.CodeDuplicatePreparedStatementError,
				"prepared statement %q already exists", prepare.Name)
		}
	} else {
		// Deallocate the unnamed statement, if it exists.
		s.DeletePreparedStatement("")
	}

	s.AddPreparedStatement(prepare.Name, prepare.Statement, prepare.TypeHints)

	return nil
}
