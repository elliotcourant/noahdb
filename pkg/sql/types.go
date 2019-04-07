package sql

type NoahQueryPlanner interface {
	getNoahQueryPlan(*session) (InitialPlan, bool, error)
}

type SimpleQueryPlanner interface {
	getSimpleQueryPlan(*session) (InitialPlan, bool, error)
}

type StandardQueryPlanner interface {
	getStandardQueryPlan(*session) (InitialPlan, bool, error)
}
