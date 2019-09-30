package executor

import (
	"github.com/elliotcourant/noahdb/pkg/pgproto"
)

type Result struct {
	rowDescription pgproto.RowDescription
	rows           []pgproto.DataRow
}
