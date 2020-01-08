package pgmock

import (
	"testing"
)

type (
	PostgresMock interface {
		Address() (address string, port int)
		Login() (username, password string)

		ExpectDatabaseConnection(database string, handler func(session Session))

		ExpectQuery(func(t *testing.T, query string))
	}

	Session interface {
		ExpectQuery(query string)

		ExpectQueryEx(query string, columns []string, rows ...[]interface{})
	}

	postgresBase struct {
		t *testing.T
	}
)

func NewPostgresMock(t *testing.T) PostgresMock {
	return nil
}
