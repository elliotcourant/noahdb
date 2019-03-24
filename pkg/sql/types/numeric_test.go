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

package types_test

import (
	"math/big"
	"math/rand"
	"reflect"
	"testing"

	"github.com/elliotcourant/noahdb/pkg/sql/types"
	"github.com/elliotcourant/noahdb/pkg/sql/types/testutil"
)

// For test purposes only. Note that it does not normalize values. e.g. (Int: 1, Exp: 3) will not equal (Int: 1000, Exp: 0)
func numericEqual(left, right *types.Numeric) bool {
	return left.Status == right.Status &&
		left.Exp == right.Exp &&
		((left.Int == nil && right.Int == nil) || (left.Int != nil && right.Int != nil && left.Int.Cmp(right.Int) == 0))
}

// For test purposes only.
func numericNormalizedEqual(left, right *types.Numeric) bool {
	if left.Status != right.Status {
		return false
	}

	normLeft := &types.Numeric{Int: (&big.Int{}).Set(left.Int), Status: left.Status}
	normRight := &types.Numeric{Int: (&big.Int{}).Set(right.Int), Status: right.Status}

	if left.Exp < right.Exp {
		mul := (&big.Int{}).Exp(big.NewInt(10), big.NewInt(int64(right.Exp-left.Exp)), nil)
		normRight.Int.Mul(normRight.Int, mul)
	} else if left.Exp > right.Exp {
		mul := (&big.Int{}).Exp(big.NewInt(10), big.NewInt(int64(left.Exp-right.Exp)), nil)
		normLeft.Int.Mul(normLeft.Int, mul)
	}

	return normLeft.Int.Cmp(normRight.Int) == 0
}

func mustParseBigInt(t *testing.T, src string) *big.Int {
	i := &big.Int{}
	if _, ok := i.SetString(src, 10); !ok {
		t.Fatalf("could not parse big.Int: %s", src)
	}
	return i
}

func TestNumericNormalize(t *testing.T) {
	t.Skip("cannot perform these operations at this time")
	testutil.TestSuccessfulNormalize(t, []testutil.NormalizeTest{
		{
			SQL:   "select '0'::numeric",
			Value: &types.Numeric{Int: big.NewInt(0), Exp: 0, Status: types.Present},
		},
		{
			SQL:   "select '1'::numeric",
			Value: &types.Numeric{Int: big.NewInt(1), Exp: 0, Status: types.Present},
		},
		{
			SQL:   "select '10.00'::numeric",
			Value: &types.Numeric{Int: big.NewInt(1000), Exp: -2, Status: types.Present},
		},
		{
			SQL:   "select '1e-3'::numeric",
			Value: &types.Numeric{Int: big.NewInt(1), Exp: -3, Status: types.Present},
		},
		{
			SQL:   "select '-1'::numeric",
			Value: &types.Numeric{Int: big.NewInt(-1), Exp: 0, Status: types.Present},
		},
		{
			SQL:   "select '10000'::numeric",
			Value: &types.Numeric{Int: big.NewInt(1), Exp: 4, Status: types.Present},
		},
		{
			SQL:   "select '3.14'::numeric",
			Value: &types.Numeric{Int: big.NewInt(314), Exp: -2, Status: types.Present},
		},
		{
			SQL:   "select '1.1'::numeric",
			Value: &types.Numeric{Int: big.NewInt(11), Exp: -1, Status: types.Present},
		},
		{
			SQL:   "select '100010001'::numeric",
			Value: &types.Numeric{Int: big.NewInt(100010001), Exp: 0, Status: types.Present},
		},
		{
			SQL:   "select '100010001.0001'::numeric",
			Value: &types.Numeric{Int: big.NewInt(1000100010001), Exp: -4, Status: types.Present},
		},
		{
			SQL: "select '4237234789234789289347892374324872138321894178943189043890124832108934.43219085471578891547854892438945012347981'::numeric",
			Value: &types.Numeric{
				Int:    mustParseBigInt(t, "423723478923478928934789237432487213832189417894318904389012483210893443219085471578891547854892438945012347981"),
				Exp:    -41,
				Status: types.Present,
			},
		},
		{
			SQL: "select '0.8925092023480223478923478978978937897879595901237890234789243679037419057877231734823098432903527585734549035904590854890345905434578345789347890402348952348905890489054234237489234987723894789234'::numeric",
			Value: &types.Numeric{
				Int:    mustParseBigInt(t, "8925092023480223478923478978978937897879595901237890234789243679037419057877231734823098432903527585734549035904590854890345905434578345789347890402348952348905890489054234237489234987723894789234"),
				Exp:    -196,
				Status: types.Present,
			},
		},
		{
			SQL: "select '0.000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000123'::numeric",
			Value: &types.Numeric{
				Int:    mustParseBigInt(t, "123"),
				Exp:    -186,
				Status: types.Present,
			},
		},
	})
}

func TestNumericTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscodeEqFunc(t, "numeric", []interface{}{
		&types.Numeric{Int: big.NewInt(0), Exp: 0, Status: types.Present},
		&types.Numeric{Int: big.NewInt(1), Exp: 0, Status: types.Present},
		&types.Numeric{Int: big.NewInt(-1), Exp: 0, Status: types.Present},
		&types.Numeric{Int: big.NewInt(1), Exp: 6, Status: types.Present},

		// preserves significant zeroes
		&types.Numeric{Int: big.NewInt(10000000), Exp: -1, Status: types.Present},
		&types.Numeric{Int: big.NewInt(10000000), Exp: -2, Status: types.Present},
		&types.Numeric{Int: big.NewInt(10000000), Exp: -3, Status: types.Present},
		&types.Numeric{Int: big.NewInt(10000000), Exp: -4, Status: types.Present},
		&types.Numeric{Int: big.NewInt(10000000), Exp: -5, Status: types.Present},
		&types.Numeric{Int: big.NewInt(10000000), Exp: -6, Status: types.Present},

		&types.Numeric{Int: big.NewInt(314), Exp: -2, Status: types.Present},
		&types.Numeric{Int: big.NewInt(123), Exp: -7, Status: types.Present},
		&types.Numeric{Int: big.NewInt(123), Exp: -8, Status: types.Present},
		&types.Numeric{Int: big.NewInt(123), Exp: -9, Status: types.Present},
		&types.Numeric{Int: big.NewInt(123), Exp: -1500, Status: types.Present},
		&types.Numeric{Int: mustParseBigInt(t, "2437"), Exp: 23790, Status: types.Present},
		&types.Numeric{Int: mustParseBigInt(t, "243723409723490243842378942378901237502734019231380123"), Exp: 23790, Status: types.Present},
		&types.Numeric{Int: mustParseBigInt(t, "43723409723490243842378942378901237502734019231380123"), Exp: 80, Status: types.Present},
		&types.Numeric{Int: mustParseBigInt(t, "3723409723490243842378942378901237502734019231380123"), Exp: 81, Status: types.Present},
		&types.Numeric{Int: mustParseBigInt(t, "723409723490243842378942378901237502734019231380123"), Exp: 82, Status: types.Present},
		&types.Numeric{Int: mustParseBigInt(t, "23409723490243842378942378901237502734019231380123"), Exp: 83, Status: types.Present},
		&types.Numeric{Int: mustParseBigInt(t, "3409723490243842378942378901237502734019231380123"), Exp: 84, Status: types.Present},
		&types.Numeric{Int: mustParseBigInt(t, "913423409823409243892349028349023482934092340892390101"), Exp: -14021, Status: types.Present},
		&types.Numeric{Int: mustParseBigInt(t, "13423409823409243892349028349023482934092340892390101"), Exp: -90, Status: types.Present},
		&types.Numeric{Int: mustParseBigInt(t, "3423409823409243892349028349023482934092340892390101"), Exp: -91, Status: types.Present},
		&types.Numeric{Int: mustParseBigInt(t, "423409823409243892349028349023482934092340892390101"), Exp: -92, Status: types.Present},
		&types.Numeric{Int: mustParseBigInt(t, "23409823409243892349028349023482934092340892390101"), Exp: -93, Status: types.Present},
		&types.Numeric{Int: mustParseBigInt(t, "3409823409243892349028349023482934092340892390101"), Exp: -94, Status: types.Present},
		&types.Numeric{Status: types.Null},
	}, func(aa, bb interface{}) bool {
		a := aa.(types.Numeric)
		b := bb.(types.Numeric)

		return numericEqual(&a, &b)
	})

}

func TestNumericTranscodeFuzz(t *testing.T) {
	r := rand.New(rand.NewSource(0))
	max := &big.Int{}
	max.SetString("9999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999", 10)

	values := make([]interface{}, 0, 2000)
	for i := 0; i < 10; i++ {
		for j := -50; j < 50; j++ {
			num := (&big.Int{}).Rand(r, max)
			negNum := &big.Int{}
			negNum.Neg(num)
			values = append(values, &types.Numeric{Int: num, Exp: int32(j), Status: types.Present})
			values = append(values, &types.Numeric{Int: negNum, Exp: int32(j), Status: types.Present})
		}
	}

	testutil.TestSuccessfulTranscodeEqFunc(t, "numeric", values,
		func(aa, bb interface{}) bool {
			a := aa.(types.Numeric)
			b := bb.(types.Numeric)

			return numericNormalizedEqual(&a, &b)
		})
}

func TestNumericSet(t *testing.T) {
	successfulTests := []struct {
		source interface{}
		result *types.Numeric
	}{
		{source: float32(1), result: &types.Numeric{Int: big.NewInt(1), Status: types.Present}},
		{source: float64(1), result: &types.Numeric{Int: big.NewInt(1), Status: types.Present}},
		{source: int8(1), result: &types.Numeric{Int: big.NewInt(1), Status: types.Present}},
		{source: int16(1), result: &types.Numeric{Int: big.NewInt(1), Status: types.Present}},
		{source: int32(1), result: &types.Numeric{Int: big.NewInt(1), Status: types.Present}},
		{source: int64(1), result: &types.Numeric{Int: big.NewInt(1), Status: types.Present}},
		{source: int8(-1), result: &types.Numeric{Int: big.NewInt(-1), Status: types.Present}},
		{source: int16(-1), result: &types.Numeric{Int: big.NewInt(-1), Status: types.Present}},
		{source: int32(-1), result: &types.Numeric{Int: big.NewInt(-1), Status: types.Present}},
		{source: int64(-1), result: &types.Numeric{Int: big.NewInt(-1), Status: types.Present}},
		{source: uint8(1), result: &types.Numeric{Int: big.NewInt(1), Status: types.Present}},
		{source: uint16(1), result: &types.Numeric{Int: big.NewInt(1), Status: types.Present}},
		{source: uint32(1), result: &types.Numeric{Int: big.NewInt(1), Status: types.Present}},
		{source: uint64(1), result: &types.Numeric{Int: big.NewInt(1), Status: types.Present}},
		{source: "1", result: &types.Numeric{Int: big.NewInt(1), Status: types.Present}},
		{source: _int8(1), result: &types.Numeric{Int: big.NewInt(1), Status: types.Present}},
		{source: float64(1000), result: &types.Numeric{Int: big.NewInt(1), Exp: 3, Status: types.Present}},
		{source: float64(1234), result: &types.Numeric{Int: big.NewInt(1234), Exp: 0, Status: types.Present}},
		{source: float64(12345678900), result: &types.Numeric{Int: big.NewInt(123456789), Exp: 2, Status: types.Present}},
		{source: float64(12345.678901), result: &types.Numeric{Int: big.NewInt(12345678901), Exp: -6, Status: types.Present}},
	}

	for i, tt := range successfulTests {
		r := &types.Numeric{}
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !numericEqual(r, tt.result) {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, r)
		}
	}
}

func TestNumericAssignTo(t *testing.T) {
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
		src      *types.Numeric
		dst      interface{}
		expected interface{}
	}{
		{src: &types.Numeric{Int: big.NewInt(42), Status: types.Present}, dst: &f32, expected: float32(42)},
		{src: &types.Numeric{Int: big.NewInt(42), Status: types.Present}, dst: &f64, expected: float64(42)},
		{src: &types.Numeric{Int: big.NewInt(42), Exp: -1, Status: types.Present}, dst: &f32, expected: float32(4.2)},
		{src: &types.Numeric{Int: big.NewInt(42), Exp: -1, Status: types.Present}, dst: &f64, expected: float64(4.2)},
		{src: &types.Numeric{Int: big.NewInt(42), Status: types.Present}, dst: &i16, expected: int16(42)},
		{src: &types.Numeric{Int: big.NewInt(42), Status: types.Present}, dst: &i32, expected: int32(42)},
		{src: &types.Numeric{Int: big.NewInt(42), Status: types.Present}, dst: &i64, expected: int64(42)},
		{src: &types.Numeric{Int: big.NewInt(42), Exp: 3, Status: types.Present}, dst: &i64, expected: int64(42000)},
		{src: &types.Numeric{Int: big.NewInt(42), Status: types.Present}, dst: &i, expected: int(42)},
		{src: &types.Numeric{Int: big.NewInt(42), Status: types.Present}, dst: &ui8, expected: uint8(42)},
		{src: &types.Numeric{Int: big.NewInt(42), Status: types.Present}, dst: &ui16, expected: uint16(42)},
		{src: &types.Numeric{Int: big.NewInt(42), Status: types.Present}, dst: &ui32, expected: uint32(42)},
		{src: &types.Numeric{Int: big.NewInt(42), Status: types.Present}, dst: &ui64, expected: uint64(42)},
		{src: &types.Numeric{Int: big.NewInt(42), Status: types.Present}, dst: &ui, expected: uint(42)},
		{src: &types.Numeric{Int: big.NewInt(42), Status: types.Present}, dst: &_i8, expected: _int8(42)},
		{src: &types.Numeric{Int: big.NewInt(0), Status: types.Null}, dst: &pi8, expected: (*int8)(nil)},
		{src: &types.Numeric{Int: big.NewInt(0), Status: types.Null}, dst: &_pi8, expected: (*_int8)(nil)},
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
		src      *types.Numeric
		dst      interface{}
		expected interface{}
	}{
		{src: &types.Numeric{Int: big.NewInt(42), Status: types.Present}, dst: &pf32, expected: float32(42)},
		{src: &types.Numeric{Int: big.NewInt(42), Status: types.Present}, dst: &pf64, expected: float64(42)},
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
		src *types.Numeric
		dst interface{}
	}{
		{src: &types.Numeric{Int: big.NewInt(150), Status: types.Present}, dst: &i8},
		{src: &types.Numeric{Int: big.NewInt(40000), Status: types.Present}, dst: &i16},
		{src: &types.Numeric{Int: big.NewInt(-1), Status: types.Present}, dst: &ui8},
		{src: &types.Numeric{Int: big.NewInt(-1), Status: types.Present}, dst: &ui16},
		{src: &types.Numeric{Int: big.NewInt(-1), Status: types.Present}, dst: &ui32},
		{src: &types.Numeric{Int: big.NewInt(-1), Status: types.Present}, dst: &ui64},
		{src: &types.Numeric{Int: big.NewInt(-1), Status: types.Present}, dst: &ui},
		{src: &types.Numeric{Int: big.NewInt(0), Status: types.Null}, dst: &i32},
	}

	for i, tt := range errorTests {
		err := tt.src.AssignTo(tt.dst)
		if err == nil {
			t.Errorf("%d: expected error but none was returned (%v -> %v)", i, tt.src, tt.dst)
		}
	}
}

func TestNumericEncodeDecodeBinary(t *testing.T) {
	ci := types.NewConnInfo()
	tests := []interface{}{
		123,
		0.000012345,
		1.00002345,
	}

	for i, tt := range tests {
		toString := func(n *types.Numeric) string {
			ci := types.NewConnInfo()
			text, err := n.EncodeText(ci, nil)
			if err != nil {
				t.Errorf("%d: %v", i, err)
			}
			return string(text)
		}
		numeric := &types.Numeric{}
		numeric.Set(tt)

		encoded, err := numeric.EncodeBinary(ci, nil)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}
		decoded := &types.Numeric{}
		decoded.DecodeBinary(ci, encoded)

		text0 := toString(numeric)
		text1 := toString(decoded)

		if text0 != text1 {
			t.Errorf("%d: expected %v to equal to %v, but doesn't", i, text0, text1)
		}
	}
}
