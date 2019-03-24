package sessiondata

import (
	"strings"
)

// PgDatabaseName is the name of the default postgres system database.
const PgDatabaseName = "postgres"

// DefaultDatabaseName is the name ofthe default CockroachDB database used
// for connections without a current db set.
const DefaultDatabaseName = "defaultdb"

// PgCatalogName is the name of the pg_catalog system schema.
const PgCatalogName = "pg_catalog"

// SearchPath represents a list of namespaces to search builtins in.
// The names must be normalized (as per Name.Normalize) already.
type SearchPath struct {
	paths             []string
	containsPgCatalog bool
}

// MakeSearchPath returns a new SearchPath struct.
func MakeSearchPath(paths []string) SearchPath {
	containsPgCatalog := false
	for _, e := range paths {
		if e == PgCatalogName {
			containsPgCatalog = true
			break
		}
	}
	return SearchPath{
		paths:             paths,
		containsPgCatalog: containsPgCatalog,
	}
}

// Iter returns an iterator through the search path. We must include the
// implicit pg_catalog at the beginning of the search path, unless it has been
// explicitly set later by the user.
// "The system catalog schema, pg_catalog, is always searched, whether it is
// mentioned in the path or not. If it is mentioned in the path then it will be
// searched in the specified order. If pg_catalog is not in the path then it
// will be searched before searching any of the path items."
// - https://www.postgresql.org/docs/9.1/static/runtime-config-client.html
func (s SearchPath) Iter() func() (next string, ok bool) {
	i := -1
	if s.containsPgCatalog {
		i = 0
	}
	return func() (next string, ok bool) {
		if i == -1 {
			i++
			return PgCatalogName, true
		}
		if i < len(s.paths) {
			i++
			return s.paths[i-1], true
		}
		return "", false
	}
}

// IterWithoutImplicitPGCatalog is the same as Iter, but does not include the implicit pg_catalog.
func (s SearchPath) IterWithoutImplicitPGCatalog() func() (next string, ok bool) {
	i := 0
	return func() (next string, ok bool) {
		if i < len(s.paths) {
			i++
			return s.paths[i-1], true
		}
		return "", false
	}
}

func (s SearchPath) String() string {
	return strings.Join(s.paths, ", ")
}
