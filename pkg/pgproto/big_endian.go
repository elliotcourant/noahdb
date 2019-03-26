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

package pgproto

import (
	"encoding/binary"
)

type BigEndianBuf [8]byte

func (b BigEndianBuf) Int16(n int16) []byte {
	buf := b[0:2]
	binary.BigEndian.PutUint16(buf, uint16(n))
	return buf
}

func (b BigEndianBuf) Uint16(n uint16) []byte {
	buf := b[0:2]
	binary.BigEndian.PutUint16(buf, n)
	return buf
}

func (b BigEndianBuf) Int32(n int32) []byte {
	buf := b[0:4]
	binary.BigEndian.PutUint32(buf, uint32(n))
	return buf
}

func (b BigEndianBuf) Uint32(n uint32) []byte {
	buf := b[0:4]
	binary.BigEndian.PutUint32(buf, n)
	return buf
}

func (b BigEndianBuf) Int64(n int64) []byte {
	buf := b[0:8]
	binary.BigEndian.PutUint64(buf, uint64(n))
	return buf
}
