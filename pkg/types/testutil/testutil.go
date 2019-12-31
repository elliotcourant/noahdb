package testutil

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/elliotcourant/timber"
	"github.com/jackc/pgx"
	"os"
	"reflect"
	"testing"

	"github.com/elliotcourant/noahdb/pkg/types"
	_ "github.com/lib/pq"
)

func MustConnectDatabaseSQL(t testing.TB, driverName string) *sql.DB {
	var sqlDriverName string
	switch driverName {
	case "github.com/lib/pq":
		sqlDriverName = "postgres"
	case "github.com/jackc/pgx/stdlib":
		sqlDriverName = "pgx"
	default:
		t.Fatalf("Unknown driver %v", driverName)
	}

	db, err := sql.Open(sqlDriverName, os.Getenv("PGX_TEST_DATABASE"))
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func MustConnectPgx(t testing.TB) *pgx.Conn {
	config, err := pgx.ParseConnectionString(os.Getenv("PGX_TEST_DATABASE"))
	if err != nil {
		t.Fatal(err)
	}
	config.LogLevel = 6

	conn, err := pgx.Connect(config)
	if err != nil {
		t.Fatal(err)
	}

	return conn
}

func MustClose(t testing.TB, conn interface {
	Close() error
}) {
	err := conn.Close()
	if err != nil {
		t.Fatal(err)
	}
}

type forceTextEncoder struct {
	e types.TextEncoder
}

func (f forceTextEncoder) EncodeText(ci *types.ConnInfo, buf []byte) ([]byte, error) {
	return f.e.EncodeText(ci, buf)
}

type forceBinaryEncoder struct {
	e types.BinaryEncoder
}

func (f forceBinaryEncoder) EncodeBinary(ci *types.ConnInfo, buf []byte) ([]byte, error) {
	return f.e.EncodeBinary(ci, buf)
}

func ForceEncoder(e interface{}, formatCode int16) interface{} {
	switch formatCode {
	case pgx.TextFormatCode:
		if e, ok := e.(types.TextEncoder); ok {
			return forceTextEncoder{e: e}
		}
	case pgx.BinaryFormatCode:
		if e, ok := e.(types.BinaryEncoder); ok {
			return forceBinaryEncoder{e: e.(types.BinaryEncoder)}
		}
	}
	return nil
}

func TestSuccessfulTranscode(t testing.TB, typesName string, values []interface{}) {
	TestSuccessfulTranscodeEqFunc(t, typesName, values, func(a, b interface{}) bool {
		return reflect.DeepEqual(a, b)
	})
}

func TestSuccessfulTranscodeEqFunc(t testing.TB, typesName string, values []interface{}, eqFunc func(a, b interface{}) bool) {
	TestPgxSuccessfulTranscodeEqFunc(t, typesName, values, eqFunc)
	TestPgxSimpleProtocolSuccessfulTranscodeEqFunc(t, typesName, values, eqFunc)
	for _, driverName := range []string{"github.com/lib/pq", "github.com/jackc/pgx/stdlib"} {
		TestDatabaseSQLSuccessfulTranscodeEqFunc(t, driverName, typesName, values, eqFunc)
	}
}

func TestPgxSuccessfulTranscodeEqFunc(t testing.TB, typesName string, values []interface{}, eqFunc func(a, b interface{}) bool) {
	t.Skip("noah does not support prepared statements at this time.")
	conn := MustConnectPgx(t)
	defer MustClose(t, conn)

	ps, err := conn.Prepare("test", fmt.Sprintf("select $1::%s", typesName))
	if err != nil {
		t.Fatal(err)
	}

	formats := []struct {
		name       string
		formatCode int16
	}{
		{name: "TextFormat", formatCode: pgx.TextFormatCode},
		{name: "BinaryFormat", formatCode: pgx.BinaryFormatCode},
	}

	for i, v := range values {
		for _, fc := range formats {
			ps.FieldDescriptions[0].FormatCode = fc.formatCode
			vEncoder := ForceEncoder(v, fc.formatCode)
			if vEncoder == nil {
				t.Logf("Skipping: %#v does not implement %v", v, fc.name)
				continue
			}
			// Derefence value if it is a pointer
			derefV := v
			refVal := reflect.ValueOf(v)
			if refVal.Kind() == reflect.Ptr {
				derefV = refVal.Elem().Interface()
			}

			result := reflect.New(reflect.TypeOf(derefV))
			err := conn.QueryRow("test", ForceEncoder(v, fc.formatCode)).Scan(result.Interface())
			if err != nil {
				t.Errorf("%v %d: %v", fc.name, i, err)
				timber.Errorf("%v %d: %v", fc.name, i, err)
			}

			if !eqFunc(result.Elem().Interface(), derefV) {
				t.Errorf("%v %d: expected %v, got %v", fc.name, i, derefV, result.Elem().Interface())
			}
		}
	}
}

func TestPgxSimpleProtocolSuccessfulTranscodeEqFunc(t testing.TB, typesName string, values []interface{}, eqFunc func(a, b interface{}) bool) {
	conn := MustConnectPgx(t)
	defer MustClose(t, conn)

	for i, v := range values {
		// Derefence value if it is a pointer
		derefV := v
		refVal := reflect.ValueOf(v)
		if refVal.Kind() == reflect.Ptr {
			derefV = refVal.Elem().Interface()
		}

		result := reflect.New(reflect.TypeOf(derefV))
		err := conn.QueryRowEx(
			context.Background(),
			fmt.Sprintf("select ($1)::%s", typesName),
			&pgx.QueryExOptions{SimpleProtocol: true},
			v,
		).Scan(result.Interface())
		if err != nil {
			t.Errorf("Simple protocol %d: %v", i, err)
		}

		if !eqFunc(result.Elem().Interface(), derefV) {
			t.Errorf("Simple protocol %d: expected %v, got %v", i, derefV, result.Elem().Interface())
		}
	}
}

func TestDatabaseSQLSuccessfulTranscodeEqFunc(t testing.TB, driverName, typesName string, values []interface{}, eqFunc func(a, b interface{}) bool) {
	conn := MustConnectDatabaseSQL(t, driverName)
	defer MustClose(t, conn)

	ps, err := conn.Prepare(fmt.Sprintf("select $1::%s", typesName))
	if err != nil {
		t.Fatal(err)
	}

	for i, v := range values {
		// Derefence value if it is a pointer
		derefV := v
		refVal := reflect.ValueOf(v)
		if refVal.Kind() == reflect.Ptr {
			derefV = refVal.Elem().Interface()
		}

		result := reflect.New(reflect.TypeOf(derefV))
		err := ps.QueryRow(v).Scan(result.Interface())
		if err != nil {
			t.Errorf("%v %d: %v", driverName, i, err)
		}

		if !eqFunc(result.Elem().Interface(), derefV) {
			t.Errorf("%v %d: expected %v, got %v", driverName, i, derefV, result.Elem().Interface())
		}
	}
}

type NormalizeTest struct {
	SQL   string
	Value interface{}
}

func TestSuccessfulNormalize(t testing.TB, tests []NormalizeTest) {
	TestSuccessfulNormalizeEqFunc(t, tests, func(a, b interface{}) bool {
		return reflect.DeepEqual(a, b)
	})
}

func TestSuccessfulNormalizeEqFunc(t testing.TB, tests []NormalizeTest, eqFunc func(a, b interface{}) bool) {
	TestPgxSuccessfulNormalizeEqFunc(t, tests, eqFunc)
	for _, driverName := range []string{"github.com/lib/pq", "github.com/jackc/pgx/stdlib"} {
		TestDatabaseSQLSuccessfulNormalizeEqFunc(t, driverName, tests, eqFunc)
	}
}

func TestPgxSuccessfulNormalizeEqFunc(t testing.TB, tests []NormalizeTest, eqFunc func(a, b interface{}) bool) {
	conn := MustConnectPgx(t)
	defer MustClose(t, conn)

	formats := []struct {
		name       string
		formatCode int16
	}{
		{name: "TextFormat", formatCode: pgx.TextFormatCode},
		{name: "BinaryFormat", formatCode: pgx.BinaryFormatCode},
	}

	for i, tt := range tests {
		for _, fc := range formats {
			psName := fmt.Sprintf("test%d", i)
			ps, err := conn.Prepare(psName, tt.SQL)
			if err != nil {
				t.Fatal(err)
			}

			ps.FieldDescriptions[0].FormatCode = fc.formatCode
			if ForceEncoder(tt.Value, fc.formatCode) == nil {
				t.Logf("Skipping: %#v does not implement %v", tt.Value, fc.name)
				continue
			}
			// Derefence value if it is a pointer
			derefV := tt.Value
			refVal := reflect.ValueOf(tt.Value)
			if refVal.Kind() == reflect.Ptr {
				derefV = refVal.Elem().Interface()
			}

			result := reflect.New(reflect.TypeOf(derefV))
			err = conn.QueryRow(psName).Scan(result.Interface())
			if err != nil {
				t.Errorf("%v %d: %v", fc.name, i, err)
			}

			if !eqFunc(result.Elem().Interface(), derefV) {
				t.Errorf("%v %d: expected %v, got %v", fc.name, i, derefV, result.Elem().Interface())
			}
		}
	}
}

func TestDatabaseSQLSuccessfulNormalizeEqFunc(t testing.TB, driverName string, tests []NormalizeTest, eqFunc func(a, b interface{}) bool) {
	conn := MustConnectDatabaseSQL(t, driverName)
	defer MustClose(t, conn)

	for i, tt := range tests {
		ps, err := conn.Prepare(tt.SQL)
		if err != nil {
			t.Errorf("%d. %v", i, err)
			continue
		}

		// Derefence value if it is a pointer
		derefV := tt.Value
		refVal := reflect.ValueOf(tt.Value)
		if refVal.Kind() == reflect.Ptr {
			derefV = refVal.Elem().Interface()
		}

		result := reflect.New(reflect.TypeOf(derefV))
		err = ps.QueryRow().Scan(result.Interface())
		if err != nil {
			t.Errorf("%v %d: %v", driverName, i, err)
		}

		if !eqFunc(result.Elem().Interface(), derefV) {
			t.Errorf("%v %d: expected %v, got %v", driverName, i, derefV, result.Elem().Interface())
		}
	}
}
