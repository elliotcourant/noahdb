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

// GenericText is a placeholder for text format values that no other type exists
// to handle.
type GenericText Text

func (dst *GenericText) Set(src interface{}) error {
	return (*Text)(dst).Set(src)
}

func (dst *GenericText) Get() interface{} {
	return (*Text)(dst).Get()
}

func (src *GenericText) AssignTo(dst interface{}) error {
	return (*Text)(src).AssignTo(dst)
}

func (dst *GenericText) DecodeText(ci *ConnInfo, src []byte) error {
	return (*Text)(dst).DecodeText(ci, src)
}

func (src *GenericText) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	return (*Text)(src).EncodeText(ci, buf)
}

// Scan implements the database/sql Scanner interface.
func (dst *GenericText) Scan(src interface{}) error {
	return (*Text)(dst).Scan(src)
}

// Value implements the database/sql/driver Valuer interface.
func (src *GenericText) Value() (driver.Value, error) {
	return (*Text)(src).Value()
}
