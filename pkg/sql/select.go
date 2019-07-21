package sql

import (
	"fmt"
	"github.com/ahmetb/go-linq"
	"github.com/elliotcourant/noahdb/pkg/ast"
	"github.com/elliotcourant/noahdb/pkg/core"
	"github.com/elliotcourant/noahdb/pkg/util/queryutil"
	"github.com/elliotcourant/timber"
	"strings"
)

type selectStmtPlanner struct {
	tables []core.Table
	tree   ast.SelectStmt
}

func newSelectStatementPlan(tree ast.SelectStmt) *selectStmtPlanner {
	return &selectStmtPlanner{
		tree: tree,
	}
}

func (stmt *selectStmtPlanner) getNoahQueryPlan(s *session) (InitialPlan, bool, error) {
	tableNames := queryutil.GetTables(stmt.tree)
	if len(tableNames) == 0 {
		return InitialPlan{}, false, nil
	}

	linq.From(tableNames).Distinct().ToSlice(&tableNames)

	tables, err := s.Colony().Tables().GetTables(tableNames...)
	if err != nil {
		return InitialPlan{}, false, err
	}

	if len(tables) != len(tableNames) {
		// This means that there is a table missing.
		missingTables := make([]string, 0)
		linq.From(tableNames).
			ExceptBy(linq.From(tables), func(i interface{}) interface{} {
				if table, ok := i.(core.Table); ok {
					return table.TableName
				}
				return nil
			}).ToSlice(&missingTables)
		timber.Debugf("could not resolve tables: %s", strings.Join(missingTables, ", "))
		return InitialPlan{}, false, fmt.Errorf("could not resolve tables with names: %s", strings.Join(missingTables, ", "))
	}

	stmt.tables = tables

	numberOfNoahTables := linq.From(stmt.tables).CountWith(func(t interface{}) bool {
		if table, ok := t.(core.Table); ok {
			return table.TableType == core.TableType_Noah
		}
		return false
	})

	if numberOfNoahTables != len(tableNames) {
		// All the tables should be noah tables, or none of them should be.
		return InitialPlan{}, false, nil
		// fmt.Errorf("all tables in a query must be normal tables, or noah tables")
	}

	if numberOfNoahTables > 0 {
		// If there are noah tables in the query at this point then we want to handle that.
		return InitialPlan{}, true, nil
	}

	// So far we have processed the tables in the select query. Now we want to check and see if any
	// of the targets for the select query are actually function calls.
	functionCalls := stmt.getFunctionCalls()

	// If there are no function calls, then there is nothing for us to do here.
	if len(functionCalls) == 0 {
		return InitialPlan{}, false, nil
	}

	// TODO elliotcourant implement function call handling.

	return InitialPlan{}, false, nil
}

func (stmt *selectStmtPlanner) getNormalQueryPlan(s *session) (InitialPlan, bool, error) {
	// We don't need to retrieve tables here, since getNoahQuery is called first
	// the tables will have been setup there.
	query, err := stmt.tree.Deparse(ast.Context_None)
	if err != nil {
		return InitialPlan{}, true, err
	}

	// If there are no tables then we can simply recompile the query and send it through,
	// this will make queries like CURRENT_TIMESTAMP or 1 very fast
	if len(stmt.tables) == 0 {
		return InitialPlan{
			Target:  PlanTarget_STANDARD,
			ShardID: 0,
			Types: map[PlanType]InitialPlanTask{
				PlanType_READ: {
					Type:  stmt.tree.StatementType(),
					Query: query,
				},
			},
		}, true, nil
	}

	var shardedTablesInQuery []core.Table
	linq.From(stmt.tables).Where(func(i interface{}) bool {
		table, ok := i.(core.Table)
		return ok && table.TableType == core.TableType_Sharded
	}).ToSlice(&shardedTablesInQuery)

	// If any of the queried tables are shard tables then we
	// need to target a specific shard. If none of the tables
	// are sharded tables then we can target any node/shard.
	if len(shardedTablesInQuery) > 0 {
		shardColumnNames := map[string]string{}
		columnsAndTables := map[string][]string{}

		for _, table := range stmt.tables {
			columns, err := s.Colony().Tables().GetColumns(table.TableID)
			if err != nil {
				return InitialPlan{}, false, err
			}
			for _, column := range columns {
				if column.ShardKey {
					shardColumnNames[table.TableName] = column.ColumnName
				}

				if tables, ok := columnsAndTables[column.ColumnName]; ok {
					columnsAndTables[column.ColumnName] = append(tables, table.TableName)
				}
			}
		}

		tenantIds, err := queryutil.FindAccountIdsEx(stmt.tree, shardColumnNames, columnsAndTables)
		if err != nil {
			return InitialPlan{}, false, err
		}

		tenantId := uint64(0)

		switch len(tenantIds) {
		case 0: // No account IDs were found in the query
			return InitialPlan{}, false,
				fmt.Errorf("cannot query sharded tables without specifying a tenant ID")
		case 1: // We are only querying a single tenant
			tenantId = tenantIds[0]
			timber.Verbosef("query targets tenant ID [%d]", tenantId)
		default:
			return InitialPlan{}, false,
				fmt.Errorf("cannot query sharded tables for multiple tenants")
		}

		tenant, err := s.Colony().Tenants().GetTenant(tenantId)
		if err != nil {
			return InitialPlan{}, false, fmt.Errorf("could not generate query plan: %s", err.Error())
		}

		return InitialPlan{
			Target:  PlanTarget_STANDARD,
			ShardID: tenant.ShardID,
			Types: map[PlanType]InitialPlanTask{
				PlanType_READ: {
					Type:  stmt.tree.StatementType(),
					Query: query,
				},
			},
		}, true, nil
	}

	// If nothing else then we can just issue a standard query
	// plan that is read only. This query can target any node.
	return InitialPlan{
		Target:  PlanTarget_STANDARD,
		ShardID: 0,
		Types: map[PlanType]InitialPlanTask{
			PlanType_READ: {
				Type:  stmt.tree.StatementType(),
				Query: query,
			},
		},
	}, true, nil
}

// getFunctionCalls loops over the targets of the query and returns any functions calls.
func (stmt *selectStmtPlanner) getFunctionCalls() []ast.FuncCall {
	functionCalls := make([]ast.FuncCall, 0)
	linq.From(stmt.tree.TargetList.Items).Where(func(i interface{}) bool {
		if resTarget, ok := i.(ast.ResTarget); ok {
			_, ok := resTarget.Val.(ast.FuncCall)
			return ok
		} else {
			return false
		}
	}).Select(func(i interface{}) interface{} {
		functionCall, _ := i.(ast.ResTarget).Val.(ast.FuncCall)
		return functionCall
	}).ToSlice(&functionCalls)
	return functionCalls
}
