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

type CIDR Inet

func (dst *CIDR) Set(src interface{}) error {
	return (*Inet)(dst).Set(src)
}

func (dst *CIDR) Get() interface{} {
	return (*Inet)(dst).Get()
}

func (src *CIDR) AssignTo(dst interface{}) error {
	return (*Inet)(src).AssignTo(dst)
}

func (dst *CIDR) DecodeText(ci *ConnInfo, src []byte) error {
	return (*Inet)(dst).DecodeText(ci, src)
}

func (dst *CIDR) DecodeBinary(ci *ConnInfo, src []byte) error {
	return (*Inet)(dst).DecodeBinary(ci, src)
}

func (src *CIDR) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	return (*Inet)(src).EncodeText(ci, buf)
}

func (src *CIDR) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	return (*Inet)(src).EncodeBinary(ci, buf)
}
