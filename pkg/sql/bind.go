package sql

import (
	"github.com/elliotcourant/noahdb/pkg/commands"
	"github.com/elliotcourant/noahdb/pkg/pgerror"
)

func (s *session) ExecuteBind(bind commands.BindStatement, result *commands.CommandResult) error {

	if bind.PortalName != "" {
		if _, ok := s.portals[bind.PortalName]; ok {
			return pgerror.NewErrorf(
				pgerror.CodeDuplicateCursorError,
				"portal %q already exists", bind.PortalName)
		}
	}
	return nil
}
