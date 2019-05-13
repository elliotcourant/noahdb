package rqliter

import (
	"database/sql/driver"
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/frunk"
	"github.com/kataras/go-errors"
	sdb "github.com/rqlite/rqlite/db"
	"reflect"
	"strconv"
	"time"
)

type Rows interface {
	Columns() []string
	ColumnTypes() []string
	Next() bool
	NextResultSet() bool
	HasNextResultSet() bool
	Scan(...interface{}) error
	Err() error
}

type rqlrows struct {
	resultIndex int
	rowIndex    int
	response    *frunk.QueryResponse
}

func NewRqlRows(response *frunk.QueryResponse) Rows {
	return &rqlrows{
		resultIndex: 0,
		rowIndex:    -1,
		response:    response,
	}
}

func (rows *rqlrows) Columns() []string {
	if !rows.anyResults() {
		return make([]string, 0)
	}
	return rows.results().Columns
}

func (rows *rqlrows) ColumnTypes() []string {
	if !rows.anyResults() {
		return make([]string, 0)
	}
	return rows.results().Types
}

func (rows *rqlrows) HasNextResultSet() bool {
	return len(rows.response.Rows) > rows.resultIndex+1
}

func (rows *rqlrows) NextResultSet() bool {
	if !rows.HasNextResultSet() {
		return false
	}
	rows.resultIndex++
	rows.rowIndex = -1
	return true
}

func (rows *rqlrows) Next() bool {
	if !rows.hasAnotherRow() {
		return false
	}
	rows.rowIndex++
	return true
}

func (rows *rqlrows) Err() (err error) {
	r := rows.results()
	if err = nil; r.Error != "" {
		err = errors.New(r.Error)
	}
	return err
}

func (rows *rqlrows) Scan(destinations ...interface{}) error {
	cols := len(rows.Columns())
	if cols != len(destinations) {
		return fmt.Errorf("expected %d destination arguments in Scan, not %d", cols, len(destinations))
	}
	values := rows.results().Values[rows.rowIndex]
	for i := 0; i < cols; i++ {
		if err := convertAssignRows(destinations[i], values[i]); err != nil {
			return err
		}
	}
	return nil
}

func (rows *rqlrows) hasAnotherRow() bool {
	res := rows.results()
	return len(res.Values) > rows.rowIndex+1
}

func (rows *rqlrows) anyResults() bool {
	return len(rows.response.Rows) >= rows.resultIndex+1
}

func (rows *rqlrows) results() *sdb.Rows {
	if !rows.anyResults() {
		return &sdb.Rows{}
	}
	return rows.response.Rows[rows.resultIndex]
}

func cloneBytes(b []byte) []byte {
	if b == nil {
		return nil
	}
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

func convertAssignRows(dest, src interface{}) error {
	switch s := src.(type) {
	case string:
		switch d := dest.(type) {
		case *string:
			if d == nil {
				return fmt.Errorf("")
			}
			*d = s
			return nil
		case *[]byte:
			if d == nil {
				return fmt.Errorf("")
			}
			*d = []byte(s)
			return nil
		}
	case []byte:
		switch d := dest.(type) {
		case *string:
			if d == nil {
				return fmt.Errorf("")
			}
			*d = string(s)
			return nil
		case *interface{}:
			if d == nil {
				return fmt.Errorf("")
			}
			*d = cloneBytes(s)
			return nil
		case *[]byte:
			if d == nil {
				return fmt.Errorf("")
			}
			*d = cloneBytes(s)
			return nil
		}
	case time.Time:
		switch d := dest.(type) {
		case *time.Time:
			*d = s
			return nil
		case *string:
			*d = s.Format(time.RFC3339Nano)
			return nil
		case *[]byte:
			if d == nil {
				return fmt.Errorf("")
			}
			*d = []byte(s.Format(time.RFC3339Nano))
			return nil
		}
	case nil:
		switch d := dest.(type) {
		case *interface{}:
			if d == nil {
				return fmt.Errorf("")
			}
			*d = nil
			return nil
		case *[]byte:
			if d == nil {
				return fmt.Errorf("")
			}
			*d = nil
			return nil
		}
	}
	var sv reflect.Value

	switch d := dest.(type) {
	case *string:
		sv = reflect.ValueOf(src)
		switch sv.Kind() {
		case reflect.Bool,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			*d = asString(src)
			return nil
		}
	case *[]byte:
		sv = reflect.ValueOf(src)
		if b, ok := asBytes(nil, sv); ok {
			*d = b
			return nil
		}
	case *bool:
		bv, err := driver.Bool.ConvertValue(src)
		if err == nil {
			*d = bv.(bool)
		}
		return err
	case *interface{}:
		*d = src
		return nil
	}

	dpv := reflect.ValueOf(dest)
	if dpv.Kind() != reflect.Ptr {
		return fmt.Errorf("destination not a pointer")
	}
	if dpv.IsNil() {
		return fmt.Errorf("")
	}

	if !sv.IsValid() {
		sv = reflect.ValueOf(src)
	}

	dv := reflect.Indirect(dpv)
	if sv.IsValid() && sv.Type().AssignableTo(dv.Type()) {
		switch b := src.(type) {
		case []byte:
			dv.Set(reflect.ValueOf(cloneBytes(b)))
		default:
			dv.Set(sv)
		}
		return nil
	}

	if dv.Kind() == sv.Kind() && sv.Type().ConvertibleTo(dv.Type()) {
		dv.Set(sv.Convert(dv.Type()))
		return nil
	}

	// The following conversions use a string value as an intermediate representation
	// to convert between various numeric types.
	//
	// This also allows scanning into user defined types such as "type Int int64".
	// For symmetry, also check for string destination types.
	switch dv.Kind() {
	case reflect.Ptr:
		if src == nil {
			dv.Set(reflect.Zero(dv.Type()))
			return nil
		}
		dv.Set(reflect.New(dv.Type().Elem()))
		return convertAssignRows(dv.Interface(), src)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		s := asString(src)
		i64, err := strconv.ParseInt(s, 10, dv.Type().Bits())
		if err != nil {
			return fmt.Errorf("converting driver.Value type %T (%q) to a %s: %v", src, s, dv.Kind(), err)
		}
		dv.SetInt(i64)
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		s := asString(src)
		u64, err := strconv.ParseUint(s, 10, dv.Type().Bits())
		if err != nil {
			return fmt.Errorf("converting driver.Value type %T (%q) to a %s: %v", src, s, dv.Kind(), err)
		}
		dv.SetUint(u64)
		return nil
	case reflect.Float32, reflect.Float64:
		s := asString(src)
		f64, err := strconv.ParseFloat(s, dv.Type().Bits())
		if err != nil {
			return fmt.Errorf("converting driver.Value type %T (%q) to a %s: %v", src, s, dv.Kind(), err)
		}
		dv.SetFloat(f64)
		return nil
	case reflect.String:
		switch v := src.(type) {
		case string:
			dv.SetString(v)
			return nil
		case []byte:
			dv.SetString(string(v))
			return nil
		}
	}
	return nil
}

func asString(src interface{}) string {
	switch v := src.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	}
	rv := reflect.ValueOf(src)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(rv.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(rv.Uint(), 10)
	case reflect.Float64:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 64)
	case reflect.Float32:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 32)
	case reflect.Bool:
		return strconv.FormatBool(rv.Bool())
	}
	return fmt.Sprintf("%v", src)
}

func asBytes(buf []byte, rv reflect.Value) (b []byte, ok bool) {
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.AppendInt(buf, rv.Int(), 10), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.AppendUint(buf, rv.Uint(), 10), true
	case reflect.Float32:
		return strconv.AppendFloat(buf, rv.Float(), 'g', -1, 32), true
	case reflect.Float64:
		return strconv.AppendFloat(buf, rv.Float(), 'g', -1, 64), true
	case reflect.Bool:
		return strconv.AppendBool(buf, rv.Bool()), true
	case reflect.String:
		s := rv.String()
		return append(buf, s...), true
	}
	return
}
