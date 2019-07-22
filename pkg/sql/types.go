package sql

type NoahQueryPlanner interface {
	getNoahQueryPlan(s *session) (InitialPlan, bool, error)
}

type NormalQueryPlanner interface {
	getNormalQueryPlan(s *session) (InitialPlan, bool, error)
}

type StandardQueryPlanner interface {
	getStandardQueryPlan(s *session) (InitialPlan, bool, error)
}
