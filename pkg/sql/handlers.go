package sql

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/ast"
)

func getStatementHandler(tree ast.Stmt) (interface{}, error) {
	switch stmt := tree.(type) {
	// case ast.AlterCollationStmt:
	// case ast.AlterDatabaseSetStmt:
	// case ast.AlterDatabaseStmt:
	// case ast.AlterDefaultPrivilegesStmt:
	// case ast.AlterDomainStmt:
	// case ast.AlterEnumStmt:
	// case ast.AlterEventTrigStmt:
	// case ast.AlterExtensionContentsStmt:
	// case ast.AlterExtensionStmt:
	// case ast.AlterFdwStmt:
	// case ast.AlterForeignServerStmt:
	// case ast.AlterFunctionStmt:
	// case ast.AlterObjectDependsStmt:
	// case ast.AlterObjectSchemaStmt:
	// case ast.AlterOperatorStmt:
	// case ast.AlterOpFamilyStmt:
	// case ast.AlterOwnerStmt:
	// case ast.AlterPolicyStmt:
	// case ast.AlterPublicationStmt:
	// case ast.AlterRoleSetStmt:
	// case ast.AlterRoleStmt:
	// case ast.AlterSeqStmt:
	// case ast.AlterSubscriptionStmt:
	// case ast.AlterSystemStmt:
	// case ast.AlterTableMoveAllStmt:
	// case ast.AlterTableSpaceOptionsStmt:
	// case ast.AlterTableStmt:
	// case ast.AlterTSConfigurationStmt:
	// case ast.AlterTSDictionaryStmt:
	// case ast.AlterUserMappingStmt:
	// case ast.CheckPointStmt:
	// case ast.ClosePortalStmt:
	// case ast.ClusterStmt:
	// case ast.CommentStmt:
	//     // return nil, _comment.CreateCommentStatment(stmt, tree).HandleComment(ctx)
	// case ast.CompositeTypeStmt:
	// case ast.ConstraintsSetStmt:
	// case ast.CopyStmt:
	// case ast.CreateAmStmt:
	// case ast.CreateCastStmt:
	// case ast.CreateConversionStmt:
	// case ast.CreateDomainStmt:
	// case ast.CreateEnumStmt:
	// case ast.CreateEventTrigStmt:
	// case ast.CreateExtensionStmt:
	// case ast.CreateFdwStmt:
	// case ast.CreateForeignServerStmt:
	// case ast.CreateForeignTableStmt:
	// case ast.CreateFunctionStmt:
	// case ast.CreatePLangStmt:
	// case ast.CreatePolicyStmt:
	// case ast.CreatePublicationStmt:
	// case ast.CreateRangeStmt:
	// case ast.CreateRoleStmt:
	case ast.CreateSchemaStmt:
		return NewCreateSchemaStatementPlan(stmt), nil
	// case ast.CreateSeqStmt:
	// case ast.CreateStatsStmt:
	// case ast.CreateStmt:
	// 	return CreateCreateStatement(stmt), nil
	// case ast.CreateSubscriptionStmt:
	// case ast.CreateTableAsStmt:
	// case ast.CreateTableSpaceStmt:
	// case ast.CreateTransformStmt:
	// case ast.CreateTrigStmt:
	// case ast.CreateUserMappingStmt:
	// case ast.CreatedbStmt:
	// case ast.DeallocateStmt:
	// case ast.DeclareCursorStmt:
	// case ast.DefineStmt:
	// case ast.DeleteStmt:
	// 	return CreateDeleteStatement(stmt), nil
	// case nodes.DiscardStmt:
	// case nodes.DoStmt:
	// case nodes.DropOwnedStmt:
	// case nodes.DropRoleStmt:
	// case ast.DropStmt:
	// 	return CreateDropStatement(stmt), nil
	// case nodes.DropSubscriptionStmt:
	// case nodes.DropTableSpaceStmt:
	// case nodes.DropUserMappingStmt:
	// case nodes.DropdbStmt:
	// case nodes.ExecuteStmt:
	// case nodes.ExplainStmt:
	// case nodes.FetchStmt:
	// case nodes.GrantRoleStmt:
	// case nodes.ImportForeignSchemaStmt:
	// case nodes.IndexStmt:
	// case ast.InsertStmt:
	// 	return CreateInsertStatement(stmt), nil
	// case nodes.ListenStmt:
	// case nodes.LoadStmt:
	// case nodes.LockStmt:
	// case nodes.NotifyStmt:
	// case nodes.PrepareStmt:
	// case nodes.ReassignOwnedStmt:
	// case nodes.RefreshMatViewStmt:
	// case nodes.ReindexStmt:
	// case nodes.RenameStmt:
	// case nodes.ReplicaIdentityStmt:
	// case nodes.RuleStmt:
	// case nodes.SecLabelStmt:
	case ast.SelectStmt:
		return NewSelectStatementPlan(stmt), nil
	// case nodes.SetOperationStmt:
	// case ast.TransactionStmt:
	// 	return CreateTransactionStatement(stmt), nil
	// case nodes.TruncateStmt:
	// case nodes.UnlistenStmt:
	// case nodes.UpdateStmt:
	//     // return nil, _update.HandleUpdate(ctx, stmt)
	// case nodes.VacuumStmt:
	// case ast.VariableSetStmt:
	// 	return CreateVariableSetStatement(stmt), nil
	// case ast.VariableShowStmt:
	// 	return CreateVariableShowStatement(stmt), nil
	// case nodes.ViewStmt:
	default:
		return nil, fmt.Errorf("invalid or unsupported nodes type")
	}
}
