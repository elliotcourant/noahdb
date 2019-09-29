package plan

import (
	"github.com/elliotcourant/noahdb/pkg/ast"
	"github.com/elliotcourant/noahdb/pkg/pgwirebase"
)

type ExecutionPlan struct {
	Tasks           []Task
	OutFormats      []pgwirebase.FormatCode
	TransactionType TransactionType
}

type Task struct {
	DataNodeShardID uint64
	ReadOnly        bool
	Query           string
	Type            ast.StmtType
}
