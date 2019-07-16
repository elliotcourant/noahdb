package sql

import (
	"fmt"
	"github.com/ahmetb/go-linq"
	"github.com/elliotcourant/noahdb/pkg/ast"
	"github.com/elliotcourant/noahdb/pkg/commands"
	"github.com/elliotcourant/noahdb/pkg/pgerror"
	"github.com/elliotcourant/noahdb/pkg/pgproto"
	"github.com/elliotcourant/noahdb/pkg/types"
	"github.com/elliotcourant/noahdb/pkg/util/queryutil"
	"github.com/elliotcourant/timber"
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
	}
	return prepared, nil
}

func (s *session) HasPreparedStatement(name string) bool {
	_, ok := s.preparedStatements[name]
	return ok
}

func (s *session) DeletePreparedStatement(name string) {
	_, ok := s.preparedStatements[name]
	if !ok {
		return
	}
	delete(s.preparedStatements, name)
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
		placeholderHints = make(map[int]types.Type)
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

	tableNames := make([]string, 0)
	linq.From(tableAliasMap).Select(func(i interface{}) interface{} {
		return i.(linq.KeyValue).Value.(string)
	}).Distinct().ToSlice(&tableNames)

	columns := make([]pgproto.FieldDescription, len(referenceColumns))

	inferredTypes := make([]types.Type, 0)

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
				typ = types.Type_text
			}

			column.DataTypeOID = typ.Uint32()
			col.Val = colt.Arg
			goto WALK
		case ast.ParamRef:
			timber.Verbosef("found parameter reference [$%d]", colt.Number)
			if column.DataTypeOID > 0 {
				inferredTypes = append(inferredTypes, types.Type(column.DataTypeOID))
				break
			}

			if hint, ok := placeholderHints[colt.Number]; ok {
				column.DataTypeOID = hint.Uint32()
			} else {
				column.DataTypeOID = types.Type_text.Uint32()
			}

			inferredTypes = append(inferredTypes, types.Type(column.DataTypeOID))
		case ast.ColumnRef:
			colNames, err := colt.Fields.DeparseList(ast.Context_Operator)
			if err != nil {
				return nil, err
			}

			column.Name = colNames[len(colNames)-1]

			// If we already have the data type OID then we can just skip
			if column.DataTypeOID > 0 {
				break
			}

			var c string   // Column name to resolve
			var t []string // Potential tables the column belongs to.
			if len(colNames) == 1 {
				c, t = colNames[0], tableNames
			} else if tbl, ok := tableAliasMap[colNames[0]]; ok {
				c, t = column.Name, []string{tbl}
			} else {
				return nil, fmt.Errorf("could not resolve table [%s]", colNames[0])
			}

			cl, ok, err := s.Colony().Tables().GetColumnFromTables(c, t)
			if err != nil {
				return nil, err
			} else if !ok {
				return nil, fmt.Errorf("could not resolve column [%s]", colNames[0])
			}

			column.DataTypeOID = cl.Type.Uint32()
			column.TableOID = uint32(cl.TableID)
		default:
			timber.Debugf("test %+v", colt)
		}

		if col.Name != nil {
			column.Name = *col.Name
		}

		columns[i] = column
	}

	prepared.Columns = columns
	prepared.InferredTypes = inferredTypes
	return prepared, nil
}
