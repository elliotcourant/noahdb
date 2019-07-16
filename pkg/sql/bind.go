package sql

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/commands"
	"github.com/elliotcourant/noahdb/pkg/pgerror"
	"github.com/elliotcourant/noahdb/pkg/pgwirebase"
	"github.com/elliotcourant/noahdb/pkg/types"
	"github.com/elliotcourant/timber"
)

func (s *session) addPortal(
	portalName,
	preparedStatementName string,
	stmt *PreparedStatement,
	args []types.Value,
	columnFormatCodes []pgwirebase.FormatCode) error {
	if _, ok := s.portals[portalName]; ok {
		panic(fmt.Sprintf("portal already exists: %s", portalName))
	}

	portal := &PreparedPortal{
		Stmt:       stmt,
		Qargs:      args,
		OutFormats: columnFormatCodes,
	}
	s.portals[portalName] = portalEntry{
		PreparedPortal: portal,
	}
	return nil
}

func (s *session) deletePortal(name string) {
	_, ok := s.portals[name]
	if !ok {
		return
	}
	delete(s.portals, name)
}

func (s *session) executeBind(bind commands.BindStatement, result *commands.CommandResult) error {
	if bind.PortalName != "" {
		if _, ok := s.portals[bind.PortalName]; ok {
			return pgerror.NewErrorf(
				pgerror.CodeDuplicateCursorError,
				"portal %q already exists", bind.PortalName)
		}
	} else {
		s.deletePortal("")
	}

	ps, ok := s.preparedStatements[bind.PreparedStatementName]
	if !ok {
		return pgerror.NewErrorf(
			pgerror.CodeInvalidSQLStatementNameError,
			"unknown prepared statement %q", bind.PreparedStatementName)
	}

	numberOfArguments := len(ps.InferredTypes)

	args := make([]types.Value, numberOfArguments)
	argFormatCodes := bind.ArgFormatCodes

	if len(bind.Args) != numberOfArguments {
		return pgwirebase.NewProtocolViolationErrorf(
			"expected %d arguments, got %d", numberOfArguments, len(bind.Args))
	}

	if len(bind.ArgFormatCodes) != 1 && len(bind.ArgFormatCodes) != numberOfArguments {
		return pgwirebase.NewProtocolViolationErrorf(
			"wrong number of format codes specified: %d for %d arguments",
			len(bind.ArgFormatCodes), numberOfArguments)
	}

	if len(bind.ArgFormatCodes) == 1 && numberOfArguments > 1 {
		argFormatCodes = make([]pgwirebase.FormatCode, numberOfArguments)
		for i := range argFormatCodes {
			argFormatCodes[i] = bind.ArgFormatCodes[0]
		}
	}

	for i, arg := range bind.Args {
		t := ps.InferredTypes[i]
		if arg == nil {
			args[i] = nil
		} else if v, err := types.Decode(argFormatCodes[i], t, arg); err != nil {
			return err
		} else {
			args[i] = v
		}
	}

	numCols := len(ps.Columns)
	if (len(bind.OutFormats) > 1) && (len(bind.OutFormats) != numCols) {
		return pgwirebase.NewProtocolViolationErrorf(
			"expected 1 or %d for number of format codes, got %d",
			numCols, len(bind.OutFormats))
	}

	if err := s.addPortal(
		bind.PortalName,
		bind.PreparedStatementName,
		ps.PreparedStatement,
		args,
		bind.OutFormats); err != nil {
		return err
	}

	timber.Verbosef("created portal [%s] for prepared statement [%s]", bind.PortalName, bind.PreparedStatementName)
	return nil
}
