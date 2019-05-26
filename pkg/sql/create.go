package sql

import (
	"fmt"
	"github.com/readystock/golog"
	"strings"

	"github.com/elliotcourant/noahdb/pkg/ast"
	"github.com/elliotcourant/noahdb/pkg/core"
)

type createStmtPlanner struct {
	table core.Table
	tree  ast.CreateStmt
}

// NewCreateStatementPlan creates a new planner for create statements.
func newCreateStatementPlan(tree ast.CreateStmt) *createStmtPlanner {
	// ast.VariableShowStmt{}
	return &createStmtPlanner{
		table: core.Table{},
		tree:  tree,
	}
}

func (stmt *createStmtPlanner) getSimpleQueryPlan(s *session) (InitialPlan, bool, error) {
	// schemaName := *stmt.tree.Relation.Schemaname
	tableName := *stmt.tree.Relation.Relname
	tables, err := s.Colony().Tables().GetTables(tableName)
	if err != nil {
		return InitialPlan{}, false, fmt.Errorf("could not verify table doesn't exit: %v", err)
	}

	// If there was a table found with the same name and we are not being optimistic.
	if len(tables) > 0 && !stmt.tree.IfNotExists {
		return InitialPlan{}, false, fmt.Errorf("table with name [%s] already exists", tableName)
	}

	stmt.table.TableName = tableName

	if err := stmt.handleTableType(); err != nil {
		return InitialPlan{}, false, err
	}

	// We want to verify that if they are creating a tenant table that it is the only one.
	if stmt.table.TableType == core.TableType_Tenant {
		tenantTable, ok, err := s.Colony().Tables().GetTenantTable()
		if err != nil {
			return InitialPlan{}, false, fmt.Errorf("could not verify tenant table: %v", err)
		}
		if ok {
			return InitialPlan{}, false, fmt.Errorf("a tenant table already exists, named [%s]", tenantTable.TableName)
		}
	}

	if err := stmt.handleColumns(s); err != nil {
		return InitialPlan{}, false, err
	}

	compiledQuery, err := stmt.tree.Deparse(ast.Context_None)
	if err != nil {
		return InitialPlan{}, false, fmt.Errorf("could not recompile query: %v", err)
	}

	return InitialPlan{
		Target:  PlanTarget_STANDARD,
		ShardID: 0,
		Types: map[PlanType]InitialPlanTask{
			PlanType_WRITE: {
				Query: compiledQuery,
				Type:  stmt.tree.StatementType(),
			},
		},
	}, true, nil
}

func (stmt *createStmtPlanner) handleTableType() error {
	stmt.table.TableType = core.TableType_Global
	// If no tablespace is defined then we want to default to global.
	// This way tables declared normally will behave normally out of the box.
	if stmt.tree.Tablespacename == nil || *stmt.tree.Tablespacename == "" {
		return nil
	}

	tablespace := strings.ToLower(*stmt.tree.Tablespacename)

	// If the tablespace doesn't have the noah prefix then we want to just
	// use the specified tablespace.
	if !strings.HasPrefix(tablespace, "noah") {
		return nil
	}

	tablespace = strings.TrimPrefix(tablespace, "noah.")

	switch tablespace {
	case "tenants":
		stmt.table.TableType = core.TableType_Tenant
	case "global":
		stmt.table.TableType = core.TableType_Global
	case "sharded":
		stmt.table.TableType = core.TableType_Sharded
	default:
		return fmt.Errorf("tablespace [%s] is not a valid noahdb space", tablespace)
	}

	// If it was a valid noah tablespace then we want to clear the tablespace name
	// to make sure it's not included in the end query.
	stmt.tree.Tablespacename = nil

	return nil
}

func (stmt *createStmtPlanner) handleColumns(s *session) error {
	// verifyPrimaryKeyColumnType := func(column core.Column, typ core.Type) error {
	// 	switch typ {
	// 	case core.Type_int8, core.Type_int4, core.Type_int2:
	// 		// We only allow for these two types to be primary keys at this time.
	// 		return nil
	// 	default:
	// 		// At the moment noah only supports integer column sharding.
	// 		return fmt.Errorf("column [%s] cannot be a primary key, a primary key must be an integer column", column.ColumnName)
	// 	}
	// }

	if stmt.tree.TableElts.Items != nil && len(stmt.tree.TableElts.Items) > 0 {
		// columns := make([]core.Column, len(stmt.tree.TableElts.Items))

		for _, tableItem := range stmt.tree.TableElts.Items {
			switch col := tableItem.(type) {
			case ast.ColumnDef:
				column := core.Column{
					ColumnName: *col.Colname,
				}

				if col.TypeName != nil &&
					col.TypeName.Names.Items != nil &&
					len(col.TypeName.Names.Items) > 0 {

					names := make([]string, 0)
					for t := 0; t < len(col.TypeName.Names.Items); t++ {
						if typeName, ok := col.TypeName.Names.Items[t].(ast.String); ok && typeName.Str != "" {
							names = append(names, typeName.Str)
						}
					}

					typeName := strings.Join(names, ".")

					golog.Verbosef("table [%s] column [%s] type [%s]", stmt.table.TableName, column.ColumnName, typeName)

				}
			case ast.Constraint:

			default:
				panic(fmt.Sprintf("could not parse column item ast of type: %T", col))
			}
		}
	}

	return nil
}
