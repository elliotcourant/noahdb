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
	"database/sql/driver"
)

// GenericBinary is a placeholder for binary format values that no other type exists
// to handle.
type GenericBinary Bytea

func (dst *GenericBinary) Set(src interface{}) error {
	return (*Bytea)(dst).Set(src)
}

func (dst *GenericBinary) Get() interface{} {
	return (*Bytea)(dst).Get()
}

func (src *GenericBinary) AssignTo(dst interface{}) error {
	return (*Bytea)(src).AssignTo(dst)
}

func (dst *GenericBinary) DecodeBinary(ci *ConnInfo, src []byte) error {
	return (*Bytea)(dst).DecodeBinary(ci, src)
}

func (src *GenericBinary) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	return (*Bytea)(src).EncodeBinary(ci, buf)
}

// Scan implements the database/sql Scanner interface.
func (dst *GenericBinary) Scan(src interface{}) error {
	return (*Bytea)(dst).Scan(src)
}

// Value implements the database/sql/driver Valuer interface.
func (src *GenericBinary) Value() (driver.Value, error) {
	return (*Bytea)(src).Value()
}
