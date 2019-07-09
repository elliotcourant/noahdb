package sql

import (
	"github.com/elliotcourant/noahdb/pkg/ast"
	"github.com/elliotcourant/noahdb/pkg/commands"
	"github.com/elliotcourant/noahdb/pkg/core"
	"github.com/elliotcourant/noahdb/pkg/pgerror"
	"github.com/elliotcourant/noahdb/pkg/pgproto"
	"github.com/elliotcourant/noahdb/pkg/util/queryutil"
	"github.com/readystock/golog"
)

const (
	DefaultColumnName = "?column?"
)

func (s *session) AddPreparedStatement(name string, stmt ast.Stmt, parseTypeHints queryutil.PlaceholderTypes) (*PreparedStatement, error) {
	prepared, err := s.prepare(stmt, parseTypeHints)
	if err != nil {
		return nil, err
	}
	s.preparedStatements[name] = preparedStatementEntry{
		PreparedStatement: prepared,
		portals:           make(map[string]struct{}),
	}
	return prepared, nil
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

func (s *session) prepare(
	stmt ast.Stmt, placeholderHints queryutil.PlaceholderTypes,
) (*PreparedStatement, error) {
	if placeholderHints == nil {
		placeholderHints = make(map[int]core.Type)
	}

	prepared := &PreparedStatement{
		TypeHints: placeholderHints,
		Statement: &stmt,
	}

	if stmt == nil {
		return prepared, nil
	}

	tableAliasMap := queryutil.GetExtendedTables(stmt)
	referenceColumns := queryutil.GetColumns(stmt)

	columns := make([]pgproto.FieldDescription, len(referenceColumns))

	for i, col := range referenceColumns {
		column := pgproto.FieldDescription{
			Name: DefaultColumnName,
		}
	WALK: // Used to handle nested references.
		switch colt := col.Val.(type) {
		case ast.TypeCast:
			typeName, _ := colt.TypeName.Deparse(ast.Context_None)
			typ, ok, err := s.Colony().Types().GetTypeByName(typeName)
			if err != nil {
				return nil, err
			}

			// If we do not recognize the type we want to be optimistic, and return text for now
			// If the type truly does not exist then postgres will throw an error.
			if !ok {
				typ = core.Type_text
			}

			column.DataTypeOID = typ.Uint32()
			col.Val = colt.Arg
			goto WALK
		case ast.ParamRef:
			if column.DataTypeOID > 0 {
				break
			}

			if hint, ok := placeholderHints[colt.Number]; ok {
				column.DataTypeOID = hint.Uint32()
			} else {
				column.DataTypeOID = core.Type_text.Uint32()
			}
		case ast.ColumnRef:
			colNames, err := colt.Fields.DeparseList(ast.Context_Operator)
			if err != nil {
				return nil, err
			}

			column.Name = colNames[len(colNames)-1]

			golog.Debugf("colname: %v", colNames)
		default:
			golog.Debugf("test %+v", colt)
		}

		if col.Name != nil {
			column.Name = *col.Name
		}

		columns[i] = column
	}

	golog.Debugf("test %+v %+v", tableAliasMap, columns)

	return prepared, nil
}
