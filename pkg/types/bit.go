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

type Bit Varbit

func (dst *Bit) Set(src interface{}) error {
	return (*Varbit)(dst).Set(src)
}

func (dst *Bit) Get() interface{} {
	return (*Varbit)(dst).Get()
}

func (src *Bit) AssignTo(dst interface{}) error {
	return (*Varbit)(src).AssignTo(dst)
}

func (dst *Bit) DecodeBinary(ci *ConnInfo, src []byte) error {
	return (*Varbit)(dst).DecodeBinary(ci, src)
}

func (src *Bit) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	return (*Varbit)(src).EncodeBinary(ci, buf)
}

// Scan implements the database/sql Scanner interface.
func (dst *Bit) Scan(src interface{}) error {
	return (*Varbit)(dst).Scan(src)
}

// Value implements the database/sql/driver Valuer interface.
func (src *Bit) Value() (driver.Value, error) {
	return (*Varbit)(src).Value()
}
