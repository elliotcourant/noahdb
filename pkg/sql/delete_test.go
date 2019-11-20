package sql

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeleteStmtPlanner_GetNormalQueryPlan(t *testing.T) {
	planner := GetStatementPlanner(t, "DELETE FROM test WHERE 1=1")
	assert.IsType(t, &deleteStmtPlanner{}, planner)
	assert.NotNil(t, planner)
}
