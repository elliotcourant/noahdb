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
	defaultColumnName = "?column?"
)

func (s *session) addPreparedStatement(name string, stmt ast.Stmt, parseTypeHints queryutil.PlaceholderTypes) (*PreparedStatement, error) {
	prepared, err := s.prepare(stmt, parseTypeHints)
	if err != nil {
		return nil, err
	}
	s.preparedStatements[name] = preparedStatementEntry{
		PreparedStatement: prepared,
	}
	return prepared, nil
}

func (s *session) hasPreparedStatement(name string) bool {
	_, ok := s.preparedStatements[name]
	return ok
}

func (s *session) deletePreparedStatement(name string) {
	_, ok := s.preparedStatements[name]
	if !ok {
		return
	}
	delete(s.preparedStatements, name)
}

func (s *session) executePrepare(prepare commands.PrepareStatement, result *commands.CommandResult) error {
	if prepare.Name != "" {
		if s.hasPreparedStatement(prepare.Name) {
			return pgerror.NewErrorf(pgerror.CodeDuplicatePreparedStatementError,
				"prepared statement %q already exists", prepare.Name)
		}
	} else {
		// Deallocate the unnamed statement, if it exists.
		s.deletePreparedStatement("")
	}

	_, err := s.addPreparedStatement(prepare.Name, prepare.Statement, prepare.TypeHints)
	return err
}

func (s *session) prepare(
	stmt ast.Stmt, placeholderHints queryutil.PlaceholderTypes,
) (*PreparedStatement, error) {
	if placeholderHints == nil {
		placeholderHints = make(map[int]types.Type)
	}

	prepared := &PreparedStatement{
		TypeHints: placeholderHints,
		Statement: stmt,
	}

	if stmt == nil {
		return prepared, nil
	}

	tableAliasMap := queryutil.GetExtendedTables(stmt)
	referenceColumns := queryutil.GetColumns(stmt)
	tableNames := make([]string, 0)

	// Take all of the actual table names from the table alias map and create a distinct
	// array. We do this because if they are joining on the same table multiple times then
	// we only need to reference that one table once.
	linq.From(tableAliasMap).Select(func(i interface{}) interface{} {
		return i.(linq.KeyValue).Value.(string)
	}).Distinct().ToSlice(&tableNames)

	// Infer the type info for each of the columns that will be returned.
	columns, err := s.getPreparedStatementColumns(
		referenceColumns,
		tableNames,
		tableAliasMap,
		placeholderHints)
	if err != nil {
		return nil, err
	}
	prepared.Columns = columns

	// Infer the type for each param provided.
	inferredTypes, err := s.getInferredPreparedStatementParamTypes(
		prepared,
		placeholderHints)
	if err != nil {
		return nil, err
	}
	prepared.InferredTypes = inferredTypes

	return prepared, nil
}

func (s *session) getPreparedStatementColumns(
	resTargets []ast.ResTarget,
	tableNames []string,
	tableAliases map[string]string,
	placeholderHints queryutil.PlaceholderTypes) ([]pgproto.FieldDescription, error) {
	columns := make([]pgproto.FieldDescription, len(resTargets))

	for i, col := range resTargets {
		column := pgproto.FieldDescription{
			Name: defaultColumnName,
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
				break
			}

			if hint, ok := placeholderHints[colt.Number]; ok {
				column.DataTypeOID = hint.Uint32()
			} else {
				column.DataTypeOID = types.Type_text.Uint32()
			}
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
			} else if tbl, ok := tableAliases[colNames[0]]; ok {
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

	return columns, nil
}

func (s *session) getInferredPreparedStatementParamTypes(
	statement *PreparedStatement,
	placeholderHints queryutil.PlaceholderTypes) ([]types.Type, error) {
	params := queryutil.GetArguments(statement.Statement)
	inferredTypes := make([]types.Type, len(params))
	for i, n := range params {
		t, ok := placeholderHints[n]
		if ok {
			inferredTypes[i] = t
		} else {
			inferredTypes[i] = types.Type_text
		}
	}
	return inferredTypes, nil
}
