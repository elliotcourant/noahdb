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
	"encoding/hex"
	"encoding/json"

	"github.com/readystock/pgx/pgio"
)

type CopyData struct {
	Data []byte
}

func (*CopyData) Backend()  {}
func (*CopyData) Frontend() {}

func (dst *CopyData) Decode(src []byte) error {
	dst.Data = src
	return nil
}

func (src *CopyData) Encode(dst []byte) []byte {
	dst = append(dst, 'd')
	dst = pgio.AppendInt32(dst, int32(4+len(src.Data)))
	dst = append(dst, src.Data...)
	return dst
}

func (src *CopyData) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string
		Data string
	}{
		Type: "CopyData",
		Data: hex.EncodeToString(src.Data),
	})
}
