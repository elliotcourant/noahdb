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

// OIDValue (Object Identifier Type) is, according to
// https://www.postgresql.org/docs/current/static/datatype-OIDValue.html, used
// internally by PostgreSQL as a primary key for various system tables. It is
// currently implemented as an unsigned four-byte integer. Its definition can be
// found in src/include/postgres_ext.h in the PostgreSQL sources.
type OIDValue pguint32

// Set converts from src to dst. Note that as OIDValue is not a general
// number type Set does not do automatic type conversion as other number
// types do.
func (dst *OIDValue) Set(src interface{}) error {
	return (*pguint32)(dst).Set(src)
}

func (dst *OIDValue) Get() interface{} {
	return (*pguint32)(dst).Get()
}

// AssignTo assigns from src to dst. Note that as OIDValue is not a general number
// type AssignTo does not do automatic type conversion as other number types do.
func (src *OIDValue) AssignTo(dst interface{}) error {
	return (*pguint32)(src).AssignTo(dst)
}

func (dst *OIDValue) DecodeText(ci *ConnInfo, src []byte) error {
	return (*pguint32)(dst).DecodeText(ci, src)
}

func (dst *OIDValue) DecodeBinary(ci *ConnInfo, src []byte) error {
	return (*pguint32)(dst).DecodeBinary(ci, src)
}

func (src *OIDValue) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	return (*pguint32)(src).EncodeText(ci, buf)
}

func (src *OIDValue) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	return (*pguint32)(src).EncodeBinary(ci, buf)
}

// Scan implements the database/sql Scanner interface.
func (dst *OIDValue) Scan(src interface{}) error {
	return (*pguint32)(dst).Scan(src)
}

// Value implements the database/sql/driver Valuer interface.
func (src *OIDValue) Value() (driver.Value, error) {
	return (*pguint32)(src).Value()
}
