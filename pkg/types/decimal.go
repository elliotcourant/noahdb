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

package types

import (
	"fmt"
	"github.com/cockroachdb/apd"
	"github.com/kataras/go-errors"
	"github.com/readystock/noah/db/sql/pgwire/pgerror"
)

var (
	// DecimalCtx is the default context for decimal operations. Any change
	// in the exponent limits must still guarantee a safe conversion to the
	// postgres binary decimal format in the wire protocol, which uses an
	// int16. See pgwire/types.go.
	DecimalCtx = &apd.Context{
		Precision:   20,
		Rounding:    apd.RoundHalfUp,
		MaxExponent: 2000,
		MinExponent: -2000,
		Traps:       apd.DefaultTraps,
	}
	// ExactCtx is a decimal context with exact precision.
	ExactCtx = DecimalCtx.WithPrecision(0)
	// HighPrecisionCtx is a decimal context with high precision.
	HighPrecisionCtx = DecimalCtx.WithPrecision(2000)
	// IntermediateCtx is a decimal context with additional precision for
	// intermediate calculations to protect against order changes that can
	// happen in dist SQL. The additional 5 allows the stress test to pass.
	// See #13689 for more analysis and other algorithms.
	IntermediateCtx = DecimalCtx.WithPrecision(DecimalCtx.Precision + 5)
	// RoundCtx is a decimal context with high precision and RoundHalfEven
	// rounding.
	RoundCtx = func() *apd.Context {
		ctx := *HighPrecisionCtx
		ctx.Rounding = apd.RoundHalfEven
		return &ctx
	}()

	errScaleOutOfRange = pgerror.NewError(pgerror.CodeNumericValueOutOfRangeError, "scale out of range")
)

type Decimal Numeric

func (src *Decimal) GetApdDecimal() (*apd.Decimal, error) {
	floatyMcFloatyFace := float64(0.0)
	if err := src.AssignTo(&floatyMcFloatyFace); err != nil {
		return nil, err
	}
	roachDecimal := apd.Decimal{}
	if _, condition, err := HighPrecisionCtx.SetString(&roachDecimal, fmt.Sprintf("%v", floatyMcFloatyFace)); err != nil {
		return nil, err
	} else if condition != 0 {
		return nil, errors.New("YOU MESSED UP THE DECIMAL")
	}
	return &roachDecimal, nil
}

func (dst *Decimal) Set(src interface{}) error {
	return (*Numeric)(dst).Set(src)
}

func (dst *Decimal) Get() interface{} {
	return (*Numeric)(dst).Get()
}

func (src *Decimal) AssignTo(dst interface{}) error {
	return (*Numeric)(src).AssignTo(dst)
}

func (dst *Decimal) DecodeText(ci *ConnInfo, src []byte) error {
	return (*Numeric)(dst).DecodeText(ci, src)
}

func (dst *Decimal) DecodeBinary(ci *ConnInfo, src []byte) error {
	return (*Numeric)(dst).DecodeBinary(ci, src)
}

func (src *Decimal) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	return (*Numeric)(src).EncodeText(ci, buf)
}

func (src *Decimal) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	return (*Numeric)(src).EncodeBinary(ci, buf)
}
