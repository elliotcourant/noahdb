package sql

import (
	"fmt"
	"github.com/ahmetb/go-linq"
	"github.com/elliotcourant/noahdb/pkg/ast"
	"github.com/elliotcourant/noahdb/pkg/core"
)

type insertStmtPlanner struct {
	tree ast.InsertStmt
}

func newInsertStatementPlan(tree ast.InsertStmt) *insertStmtPlanner {
	return &insertStmtPlanner{
		tree: tree,
	}
}

func (stmt *insertStmtPlanner) getSimpleQueryPlan(s *session) (InitialPlan, bool, error) {
	tableName := *stmt.tree.Relation.Relname
	tables, err := s.Colony().Tables().GetTables(tableName)
	if err != nil {
		return InitialPlan{}, false, err
	}

	planType := PlanType_WRITE
	if len(stmt.tree.ReturningList.Items) > 0 {
		planType = PlanType_READWRITE
	}

	// Handle any number of returned tables.
	switch len(tables) {
	case 0:
		return InitialPlan{}, false, fmt.Errorf("could not resolve table [%s]", tableName)
	case 1: // The desired number of tables returned
	default:
		return InitialPlan{}, false, fmt.Errorf("found multiple tables with name [%s]", tableName)
	}

	table := tables[0] // We only want to work with one table.

	if table.HasSequence {
		// If the table has a sequence then we want to get the column that has the sequence and
		// handle it in the query.
		sequenceColumn, ok, err := s.Colony().Tables().GetSequenceColumnForTable(table.TableID)
		if err != nil {
			return InitialPlan{}, false, err
		}
		if !ok {
			return InitialPlan{}, false, fmt.Errorf("table metadata indicates a sequence column, but none could be found for table [%s]", table.TableName)
		}

		sequenceInsertIndex := linq.From(stmt.tree.Cols.Items).IndexOf(func(i interface{}) bool {
			resTarget, ok := i.(ast.ResTarget)
			return ok && *resTarget.Name == sequenceColumn.ColumnName
		})

		switch sequenceInsertIndex {
		case -1: // The column is not specified in the insert stmt, we need to add it.
			stmt.tree.Cols.Items = append(stmt.tree.Cols.Items, ast.ResTarget{
				Name: &sequenceColumn.ColumnName,
			})

			for i, row := range stmt.tree.SelectStmt.(ast.SelectStmt).ValuesLists {
				newId, err := s.Colony().Tables().NextSequenceID(table, sequenceColumn)
				if err != nil {
					return InitialPlan{}, false, err
				}
				stmt.tree.SelectStmt.(ast.SelectStmt).ValuesLists[i] = append(row, ast.A_Const{
					Val: ast.Integer{
						Ival: int64(newId),
					},
				})
			}
		default:
			for i, row := range stmt.tree.SelectStmt.(ast.SelectStmt).ValuesLists {
				sequenceCell := row[sequenceInsertIndex]
				if _, ok := sequenceCell.(ast.SetToDefault); !ok {
					return InitialPlan{}, false, fmt.Errorf("cannot manually set value of serialized column [%s]", sequenceColumn.ColumnName)
				}

				// Generate a new ID.
				newId, err := s.Colony().Tables().NextSequenceID(table, sequenceColumn)
				if err != nil {
					return InitialPlan{}, false, err
				}

				stmt.tree.SelectStmt.(ast.SelectStmt).ValuesLists[i][sequenceInsertIndex] = ast.A_Const{
					Val: ast.Integer{
						Ival: int64(newId),
					},
				}
			}
		}
	}

	switch table.TableType {
	case core.TableType_Noah:
		panic("not handling this yet")
	case core.TableType_Tenant:

		fallthrough
	case core.TableType_Global:
		recompiled, err := stmt.tree.Deparse(ast.Context_None)
		if err != nil {
			return InitialPlan{}, false, err
		}

		return InitialPlan{
			Target:  PlanTarget_STANDARD,
			ShardID: 0,
			Types: map[PlanType]InitialPlanTask{
				planType: {
					Query: recompiled,
					Type:  stmt.tree.StatementType(),
				},
			},
		}, true, nil
	case core.TableType_Sharded:

	}
	return InitialPlan{}, false, nil
}
