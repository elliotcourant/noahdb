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
	"encoding/json"
)

type ReadyForQuery struct {
	TxStatus byte
}

func (*ReadyForQuery) Backend() {}

func (dst *ReadyForQuery) Decode(src []byte) error {
	if len(src) != 1 {
		return &invalidMessageLenErr{messageType: "ReadyForQuery", expectedLen: 1, actualLen: len(src)}
	}

	dst.TxStatus = src[0]

	return nil
}

func (src *ReadyForQuery) Encode(dst []byte) []byte {
	return append(dst, 'Z', 0, 0, 0, 5, src.TxStatus)
}

func (src *ReadyForQuery) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type     string
		TxStatus string
	}{
		Type:     "ReadyForQuery",
		TxStatus: string(src.TxStatus),
	})
}
