package sql

import (
	"fmt"
	"github.com/readystock/golinq"
	"github.com/readystock/golog"
	"strings"

	"github.com/elliotcourant/noahdb/pkg/ast"
	"github.com/elliotcourant/noahdb/pkg/core"
)

type createStmtPlanner struct {
	table   core.Table
	columns []core.Column
	tree    ast.CreateStmt
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
	verifyPrimaryKeyColumnType := func(column core.Column) error {
		switch column.Type {
		case core.Type_int8, core.Type_int4, core.Type_int2:
			// We only allow for these types to be primary keys at this time.
			return nil
		default:
			// At the moment noah only supports integer column sharding.
			return fmt.Errorf("column [%s] cannot be a primary key, a primary key must be an integer column", column.ColumnName)
		}
	}

	verifyForeignKeyColumn := func(column *core.Column, constraint ast.Constraint) error {
		if len(constraint.PkAttrs.Items) != 1 {
			return fmt.Errorf("currently noahdb only supports single column foreign keys")
		}

		referenceTableName := strings.ToLower(*constraint.Pktable.Relname)
		key := strings.ToLower(constraint.PkAttrs.Items[0].(ast.String).Str)

		referenceTable, ok, err := s.Colony().Tables().GetTable(referenceTableName)
		if err != nil {
			return fmt.Errorf("could not create constraint referencing table [%s]: %v", referenceTableName, err)
		}
		if !ok {
			return fmt.Errorf("could not create constraint referencing table [%s], it does not exist", referenceTableName)
		}

		referencePrimaryKey, ok, err := s.Colony().Tables().GetPrimaryKeyColumnByName(referenceTableName)
		if err != nil {
			return fmt.Errorf("could not verify primary key on reference table [%s]: %v", referenceTableName, err)
		}
		if !ok {
			return fmt.Errorf("could not create foreign key referencing table [%s], the table does not have a primary key", referenceTableName)
		}

		if key != referencePrimaryKey.ColumnName {
			return fmt.Errorf("referenced table [%s] has primary key [%s], cannot create a reference to column [%s]", referenceTableName, referencePrimaryKey.ColumnName, key)
		}

		column.ForeignColumnID = referencePrimaryKey.ColumnID

		// If this column is reference the tenant table, then this column is the shard key
		// for this new table.
		if referenceTable.TableType == core.TableType_Tenant {
			column.ShardKey = true
		}
		return nil
	}

	if stmt.tree.TableElts.Items != nil && len(stmt.tree.TableElts.Items) > 0 {
		stmt.columns = make([]core.Column, len(stmt.tree.TableElts.Items))

		hasPrimaryKey, hasShardKey := false, false

		for i, tableItem := range stmt.tree.TableElts.Items {
			switch col := tableItem.(type) {
			case ast.ColumnDef:
				column := core.Column{
					ColumnName: *col.Colname,
					Sort:       int32(i),
					Nullable:   !col.IsNotNull,
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

					// Handle serial types
					switch typeName {
					case "bigserial":
						typeName = "bigint"
						column.Serial = true
					case "serial":
						typeName = "int"
						column.Serial = true
					}

					col.TypeName.Names.Items = []ast.Node{ast.String{Str: typeName}}

					pgType, ok, err := s.Colony().Types().GetTypeByName(typeName)
					if err != nil {
						return err
					} else if !ok {
						return fmt.Errorf("could not resolve type [%s]", typeName)
					}
					column.Type = pgType
				}

				// Check to see if this column is the primary key, primary keys will be used for tables
				// like account tables. If someone tries to create a table without a primary key an
				// error will be returned at this time.
				if col.Constraints.Items != nil && len(col.Constraints.Items) > 0 {
					for _, c := range col.Constraints.Items {
						constraint := c.(ast.Constraint)
						switch constraint.Contype {
						case ast.CONSTR_PRIMARY:
							if hasPrimaryKey {
								return fmt.Errorf("cannot have more than 1 primary key per table")
							}

							if err := verifyPrimaryKeyColumnType(column); err != nil {
								return err
							}

							column.Nullable = false

							hasPrimaryKey, column.PrimaryKey = true, true
						case ast.CONSTR_FOREIGN:
							if err := verifyForeignKeyColumn(&column, constraint); err != nil {
								return err
							}

							if column.ShardKey {
								if hasShardKey {
									return fmt.Errorf("cannot have multiple foreign keys referencing the tenants table")
								}
								hasShardKey = true
							}
						}
					}
				}

				stmt.columns[i] = column
			case ast.Constraint:
				// Its possible for primary keys, foreign keys and identities to be defined
				// somewhere other than the column line itself, if this happens we still want to
				// handle it gracefully.
				switch col.Contype {
				case ast.CONSTR_PRIMARY:
					if hasPrimaryKey {
						return fmt.Errorf("cannot have more than 1 primary key per table")
					}

					if len(col.Keys.Items) != 1 {
						return fmt.Errorf("currently noah only supports single column primary keys")
					}

					// We want to search columns based on the column name.
					key := strings.ToLower(col.Keys.Items[0].(ast.String).Str)

					colIndex := linq.From(stmt.columns).IndexOf(func(i interface{}) bool {
						column, ok := i.(core.Column)
						return ok && column.ColumnName == key
					})

					if colIndex < 0 {
						return fmt.Errorf("could not use column [%s] as primary key, it is not defined in the create statement", key)
					}

					if err := verifyPrimaryKeyColumnType(stmt.columns[colIndex]); err != nil {
						return err
					}

					stmt.columns[colIndex].PrimaryKey = true
					stmt.columns[colIndex].Nullable = false
					hasPrimaryKey = true
				case ast.CONSTR_FOREIGN:
					if len(col.FkAttrs.Items) != 1 {
						return fmt.Errorf("only 1 column can be used in a foreign key constraint")
					}

					key := strings.ToLower(col.FkAttrs.Items[0].(ast.String).Str)

					colIndex := linq.From(stmt.columns).IndexOf(func(i interface{}) bool {
						column, ok := i.(core.Column)
						return ok && column.ColumnName == key
					})

					if err := verifyForeignKeyColumn(&stmt.columns[colIndex], col); err != nil {
						return err
					}

					if stmt.columns[colIndex].ShardKey {
						if hasShardKey {
							return fmt.Errorf("cannot have multiple foreign keys referencing the tenants table")
						}
						hasShardKey = true
					}
				default:
					return fmt.Errorf("could not handle contraint type [%s]", col.Contype)
				}
			default:
				panic(fmt.Sprintf("could not parse column item ast of type: %T", col))
			}
		}
	}

	return nil
}
