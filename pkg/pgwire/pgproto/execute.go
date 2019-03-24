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
	"bytes"
	"encoding/binary"
	"encoding/json"

	"github.com/readystock/pgx/pgio"
)

type Execute struct {
	Portal  string
	MaxRows uint32
}

func (*Execute) Frontend() {}

func (dst *Execute) Decode(src []byte) error {
	buf := bytes.NewBuffer(src)

	b, err := buf.ReadBytes(0)
	if err != nil {
		return err
	}
	dst.Portal = string(b[:len(b)-1])

	if buf.Len() < 4 {
		return &invalidMessageFormatErr{messageType: "Execute"}
	}
	dst.MaxRows = binary.BigEndian.Uint32(buf.Next(4))

	return nil
}

func (src *Execute) Encode(dst []byte) []byte {
	dst = append(dst, 'E')
	sp := len(dst)
	dst = pgio.AppendInt32(dst, -1)

	dst = append(dst, src.Portal...)
	dst = append(dst, 0)

	dst = pgio.AppendUint32(dst, src.MaxRows)

	pgio.SetInt32(dst[sp:], int32(len(dst[sp:])))

	return dst
}

func (src *Execute) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type    string
		Portal  string
		MaxRows uint32
	}{
		Type:    "Execute",
		Portal:  src.Portal,
		MaxRows: src.MaxRows,
	})
}
