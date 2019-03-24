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
	"encoding/json"

	"github.com/readystock/pgx/pgio"
)

type BackendKeyData struct {
	ProcessID uint32
	SecretKey uint32
}

func (*BackendKeyData) Backend() {}

func (dst *BackendKeyData) Decode(src []byte) error {
	if len(src) != 8 {
		return &invalidMessageLenErr{messageType: "BackendKeyData", expectedLen: 8, actualLen: len(src)}
	}

	dst.ProcessID = binary.BigEndian.Uint32(src[:4])
	dst.SecretKey = binary.BigEndian.Uint32(src[4:])

	return nil
}

func (src *BackendKeyData) Encode(dst []byte) []byte {
	dst = append(dst, 'K')
	dst = pgio.AppendUint32(dst, 12)
	dst = pgio.AppendUint32(dst, src.ProcessID)
	dst = pgio.AppendUint32(dst, src.SecretKey)
	return dst
}

func (src *BackendKeyData) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type      string
		ProcessID uint32
		SecretKey uint32
	}{
		Type:      "BackendKeyData",
		ProcessID: src.ProcessID,
		SecretKey: src.SecretKey,
	})
}
