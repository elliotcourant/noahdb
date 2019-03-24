/*
 * Copyright (c) 2019 Ready Stock
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing
 * permissions and limitations under the License.
 */

package numeric_test

import (
	"fmt"
	"math/big"
	"math/rand"
	"reflect"
	"testing"

	"github.com/elliotcourant/noahdb/pkg/sql/types"
	shopspring "github.com/elliotcourant/noahdb/pkg/sql/types/ext/shopspring-numeric"
	"github.com/elliotcourant/noahdb/pkg/sql/types/testutil"
	"github.com/shopspring/decimal"
)

func mustParseDecimal(t *testing.T, src string) decimal.Decimal {
	dec, err := decimal.NewFromString(src)
	if err != nil {
		t.Fatal(err)
	}
	return dec
}

func TestNumericNormalize(t *testing.T) {
	t.Skip("unsupported at this time")
	testutil.TestSuccessfulNormalizeEqFunc(t, []testutil.NormalizeTest{
		{
			SQL:   "select '0'::numeric",
			Value: &shopspring.Numeric{Decimal: mustParseDecimal(t, "0"), Status: types.Present},
		},
		{
			SQL:   "select '1'::numeric",
			Value: &shopspring.Numeric{Decimal: mustParseDecimal(t, "1"), Status: types.Present},
		},
		{
			SQL:   "select '10.00'::numeric",
			Value: &shopspring.Numeric{Decimal: mustParseDecimal(t, "10.00"), Status: types.Present},
		},
		{
			SQL:   "select '1e-3'::numeric",
			Value: &shopspring.Numeric{Decimal: mustParseDecimal(t, "0.001"), Status: types.Present},
		},
		{
			SQL:   "select '-1'::numeric",
			Value: &shopspring.Numeric{Decimal: mustParseDecimal(t, "-1"), Status: types.Present},
		},
		{
			SQL:   "select '10000'::numeric",
			Value: &shopspring.Numeric{Decimal: mustParseDecimal(t, "10000"), Status: types.Present},
		},
		{
			SQL:   "select '3.14'::numeric",
			Value: &shopspring.Numeric{Decimal: mustParseDecimal(t, "3.14"), Status: types.Present},
		},
		{
			SQL:   "select '1.1'::numeric",
			Value: &shopspring.Numeric{Decimal: mustParseDecimal(t, "1.1"), Status: types.Present},
		},
		{
			SQL:   "select '100010001'::numeric",
			Value: &shopspring.Numeric{Decimal: mustParseDecimal(t, "100010001"), Status: types.Present},
		},
		{
			SQL:   "select '100010001.0001'::numeric",
			Value: &shopspring.Numeric{Decimal: mustParseDecimal(t, "100010001.0001"), Status: types.Present},
		},
		{
			SQL: "select '4237234789234789289347892374324872138321894178943189043890124832108934.43219085471578891547854892438945012347981'::numeric",
			Value: &shopspring.Numeric{
				Decimal: mustParseDecimal(t, "4237234789234789289347892374324872138321894178943189043890124832108934.43219085471578891547854892438945012347981"),
				Status:  types.Present,
			},
		},
		{
			SQL: "select '0.8925092023480223478923478978978937897879595901237890234789243679037419057877231734823098432903527585734549035904590854890345905434578345789347890402348952348905890489054234237489234987723894789234'::numeric",
			Value: &shopspring.Numeric{
				Decimal: mustParseDecimal(t, "0.8925092023480223478923478978978937897879595901237890234789243679037419057877231734823098432903527585734549035904590854890345905434578345789347890402348952348905890489054234237489234987723894789234"),
				Status:  types.Present,
			},
		},
		{
			SQL: "select '0.000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000123'::numeric",
			Value: &shopspring.Numeric{
				Decimal: mustParseDecimal(t, "0.000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000123"),
				Status:  types.Present,
			},
		},
	}, func(aa, bb interface{}) bool {
		a := aa.(shopspring.Numeric)
		b := bb.(shopspring.Numeric)

		return a.Status == b.Status && a.Decimal.Equal(b.Decimal)
	})
}

func TestNumericTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscodeEqFunc(t, "numeric", []interface{}{
		&shopspring.Numeric{Decimal: mustParseDecimal(t, "0"), Status: types.Present},
		&shopspring.Numeric{Decimal: mustParseDecimal(t, "1"), Status: types.Present},
		&shopspring.Numeric{Decimal: mustParseDecimal(t, "-1"), Status: types.Present},
		&shopspring.Numeric{Decimal: mustParseDecimal(t, "100000"), Status: types.Present},

		&shopspring.Numeric{Decimal: mustParseDecimal(t, "0.1"), Status: types.Present},
		&shopspring.Numeric{Decimal: mustParseDecimal(t, "0.01"), Status: types.Present},
		&shopspring.Numeric{Decimal: mustParseDecimal(t, "0.001"), Status: types.Present},
		&shopspring.Numeric{Decimal: mustParseDecimal(t, "0.0001"), Status: types.Present},
		&shopspring.Numeric{Decimal: mustParseDecimal(t, "0.00001"), Status: types.Present},
		&shopspring.Numeric{Decimal: mustParseDecimal(t, "0.000001"), Status: types.Present},

		&shopspring.Numeric{Decimal: mustParseDecimal(t, "3.14"), Status: types.Present},
		&shopspring.Numeric{Decimal: mustParseDecimal(t, "0.00000123"), Status: types.Present},
		&shopspring.Numeric{Decimal: mustParseDecimal(t, "0.000000123"), Status: types.Present},
		&shopspring.Numeric{Decimal: mustParseDecimal(t, "0.0000000123"), Status: types.Present},
		&shopspring.Numeric{Decimal: mustParseDecimal(t, "0.00000000123"), Status: types.Present},
		&shopspring.Numeric{Decimal: mustParseDecimal(t, "0.00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001234567890123456789"), Status: types.Present},
		&shopspring.Numeric{Decimal: mustParseDecimal(t, "4309132809320932980457137401234890237489238912983572189348951289375283573984571892758234678903467889512893489128589347891272139.8489235871258912789347891235879148795891238915678189467128957812395781238579189025891238901583915890128973578957912385798125789012378905238905471598123758923478294374327894237892234"), Status: types.Present},
		&shopspring.Numeric{Status: types.Null},
	}, func(aa, bb interface{}) bool {
		a := aa.(shopspring.Numeric)
		b := bb.(shopspring.Numeric)

		return a.Status == b.Status && a.Decimal.Equal(b.Decimal)
	})

}

func TestNumericTranscodeFuzz(t *testing.T) {
	r := rand.New(rand.NewSource(0))
	max := &big.Int{}
	max.SetString("9999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999", 10)

	values := make([]interface{}, 0, 2000)
	for i := 0; i < 500; i++ {
		num := fmt.Sprintf("%s.%s", (&big.Int{}).Rand(r, max).String(), (&big.Int{}).Rand(r, max).String())
		negNum := "-" + num
		values = append(values, &shopspring.Numeric{Decimal: mustParseDecimal(t, num), Status: types.Present})
		values = append(values, &shopspring.Numeric{Decimal: mustParseDecimal(t, negNum), Status: types.Present})
	}

	testutil.TestSuccessfulTranscodeEqFunc(t, "numeric", values,
		func(aa, bb interface{}) bool {
			a := aa.(shopspring.Numeric)
			b := bb.(shopspring.Numeric)

			return a.Status == b.Status && a.Decimal.Equal(b.Decimal)
		})
}

func TestNumericSet(t *testing.T) {
	type _int8 int8

	successfulTests := []struct {
		source interface{}
		result *shopspring.Numeric
	}{
		{source: float32(1), result: &shopspring.Numeric{Decimal: mustParseDecimal(t, "1"), Status: types.Present}},
		{source: float64(1), result: &shopspring.Numeric{Decimal: mustParseDecimal(t, "1"), Status: types.Present}},
		{source: int8(1), result: &shopspring.Numeric{Decimal: mustParseDecimal(t, "1"), Status: types.Present}},
		{source: int16(1), result: &shopspring.Numeric{Decimal: mustParseDecimal(t, "1"), Status: types.Present}},
		{source: int32(1), result: &shopspring.Numeric{Decimal: mustParseDecimal(t, "1"), Status: types.Present}},
		{source: int64(1), result: &shopspring.Numeric{Decimal: mustParseDecimal(t, "1"), Status: types.Present}},
		{source: int8(-1), result: &shopspring.Numeric{Decimal: mustParseDecimal(t, "-1"), Status: types.Present}},
		{source: int16(-1), result: &shopspring.Numeric{Decimal: mustParseDecimal(t, "-1"), Status: types.Present}},
		{source: int32(-1), result: &shopspring.Numeric{Decimal: mustParseDecimal(t, "-1"), Status: types.Present}},
		{source: int64(-1), result: &shopspring.Numeric{Decimal: mustParseDecimal(t, "-1"), Status: types.Present}},
		{source: uint8(1), result: &shopspring.Numeric{Decimal: mustParseDecimal(t, "1"), Status: types.Present}},
		{source: uint16(1), result: &shopspring.Numeric{Decimal: mustParseDecimal(t, "1"), Status: types.Present}},
		{source: uint32(1), result: &shopspring.Numeric{Decimal: mustParseDecimal(t, "1"), Status: types.Present}},
		{source: uint64(1), result: &shopspring.Numeric{Decimal: mustParseDecimal(t, "1"), Status: types.Present}},
		{source: "1", result: &shopspring.Numeric{Decimal: mustParseDecimal(t, "1"), Status: types.Present}},
		{source: _int8(1), result: &shopspring.Numeric{Decimal: mustParseDecimal(t, "1"), Status: types.Present}},
		{source: float64(1000), result: &shopspring.Numeric{Decimal: mustParseDecimal(t, "1000"), Status: types.Present}},
		{source: float64(1234), result: &shopspring.Numeric{Decimal: mustParseDecimal(t, "1234"), Status: types.Present}},
		{source: float64(12345678900), result: &shopspring.Numeric{Decimal: mustParseDecimal(t, "12345678900"), Status: types.Present}},
		{source: float64(1.25), result: &shopspring.Numeric{Decimal: mustParseDecimal(t, "1.25"), Status: types.Present}},
	}

	for i, tt := range successfulTests {
		r := &shopspring.Numeric{}
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !(r.Status == tt.result.Status && r.Decimal.Equal(tt.result.Decimal)) {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, r)
		}
	}
}

func TestNumericAssignTo(t *testing.T) {
	type _int8 int8

	var i8 int8
	var i16 int16
	var i32 int32
	var i64 int64
	var i int
	var ui8 uint8
	var ui16 uint16
	var ui32 uint32
	var ui64 uint64
	var ui uint
	var pi8 *int8
	var _i8 _int8
	var _pi8 *_int8
	var f32 float32
	var f64 float64
	var pf32 *float32
	var pf64 *float64

	simpleTests := []struct {
		src      *shopspring.Numeric
		dst      interface{}
		expected interface{}
	}{
		{src: &shopspring.Numeric{Decimal: mustParseDecimal(t, "42"), Status: types.Present}, dst: &f32, expected: float32(42)},
		{src: &shopspring.Numeric{Decimal: mustParseDecimal(t, "42"), Status: types.Present}, dst: &f64, expected: float64(42)},
		{src: &shopspring.Numeric{Decimal: mustParseDecimal(t, "4.2"), Status: types.Present}, dst: &f32, expected: float32(4.2)},
		{src: &shopspring.Numeric{Decimal: mustParseDecimal(t, "4.2"), Status: types.Present}, dst: &f64, expected: float64(4.2)},
		{src: &shopspring.Numeric{Decimal: mustParseDecimal(t, "42"), Status: types.Present}, dst: &i16, expected: int16(42)},
		{src: &shopspring.Numeric{Decimal: mustParseDecimal(t, "42"), Status: types.Present}, dst: &i32, expected: int32(42)},
		{src: &shopspring.Numeric{Decimal: mustParseDecimal(t, "42"), Status: types.Present}, dst: &i64, expected: int64(42)},
		{src: &shopspring.Numeric{Decimal: mustParseDecimal(t, "42000"), Status: types.Present}, dst: &i64, expected: int64(42000)},
		{src: &shopspring.Numeric{Decimal: mustParseDecimal(t, "42"), Status: types.Present}, dst: &i, expected: int(42)},
		{src: &shopspring.Numeric{Decimal: mustParseDecimal(t, "42"), Status: types.Present}, dst: &ui8, expected: uint8(42)},
		{src: &shopspring.Numeric{Decimal: mustParseDecimal(t, "42"), Status: types.Present}, dst: &ui16, expected: uint16(42)},
		{src: &shopspring.Numeric{Decimal: mustParseDecimal(t, "42"), Status: types.Present}, dst: &ui32, expected: uint32(42)},
		{src: &shopspring.Numeric{Decimal: mustParseDecimal(t, "42"), Status: types.Present}, dst: &ui64, expected: uint64(42)},
		{src: &shopspring.Numeric{Decimal: mustParseDecimal(t, "42"), Status: types.Present}, dst: &ui, expected: uint(42)},
		{src: &shopspring.Numeric{Decimal: mustParseDecimal(t, "42"), Status: types.Present}, dst: &_i8, expected: _int8(42)},
		{src: &shopspring.Numeric{Status: types.Null}, dst: &pi8, expected: (*int8)(nil)},
		{src: &shopspring.Numeric{Status: types.Null}, dst: &_pi8, expected: (*_int8)(nil)},
	}

	for i, tt := range simpleTests {
		err := tt.src.AssignTo(tt.dst)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if dst := reflect.ValueOf(tt.dst).Elem().Interface(); dst != tt.expected {
			t.Errorf("%d: expected %v to assign %v, but result was %v", i, tt.src, tt.expected, dst)
		}
	}

	pointerAllocTests := []struct {
		src      *shopspring.Numeric
		dst      interface{}
		expected interface{}
	}{
		{src: &shopspring.Numeric{Decimal: mustParseDecimal(t, "42"), Status: types.Present}, dst: &pf32, expected: float32(42)},
		{src: &shopspring.Numeric{Decimal: mustParseDecimal(t, "42"), Status: types.Present}, dst: &pf64, expected: float64(42)},
	}

	for i, tt := range pointerAllocTests {
		err := tt.src.AssignTo(tt.dst)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if dst := reflect.ValueOf(tt.dst).Elem().Elem().Interface(); dst != tt.expected {
			t.Errorf("%d: expected %v to assign %v, but result was %v", i, tt.src, tt.expected, dst)
		}
	}

	errorTests := []struct {
		src *shopspring.Numeric
		dst interface{}
	}{
		{src: &shopspring.Numeric{Decimal: mustParseDecimal(t, "150"), Status: types.Present}, dst: &i8},
		{src: &shopspring.Numeric{Decimal: mustParseDecimal(t, "40000"), Status: types.Present}, dst: &i16},
		{src: &shopspring.Numeric{Decimal: mustParseDecimal(t, "-1"), Status: types.Present}, dst: &ui8},
		{src: &shopspring.Numeric{Decimal: mustParseDecimal(t, "-1"), Status: types.Present}, dst: &ui16},
		{src: &shopspring.Numeric{Decimal: mustParseDecimal(t, "-1"), Status: types.Present}, dst: &ui32},
		{src: &shopspring.Numeric{Decimal: mustParseDecimal(t, "-1"), Status: types.Present}, dst: &ui64},
		{src: &shopspring.Numeric{Decimal: mustParseDecimal(t, "-1"), Status: types.Present}, dst: &ui},
		{src: &shopspring.Numeric{Status: types.Null}, dst: &i32},
	}

	for i, tt := range errorTests {
		err := tt.src.AssignTo(tt.dst)
		if err == nil {
			t.Errorf("%d: expected error but none was returned (%v -> %v)", i, tt.src, tt.dst)
		}
	}
}
