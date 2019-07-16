package sql

type NoahQueryPlanner interface {
	getNoahQueryPlan(*session) (InitialPlan, bool, error)
}

type NormalQueryPlanner interface {
	getNormalQueryPlan(*session) (InitialPlan, bool, error)
}

type StandardQueryPlanner interface {
	getStandardQueryPlan(*session) (InitialPlan, bool, error)
}
