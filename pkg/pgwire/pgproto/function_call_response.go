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
	"encoding/hex"
	"encoding/json"

	"github.com/readystock/pgx/pgio"
)

type FunctionCallResponse struct {
	Result []byte
}

func (*FunctionCallResponse) Backend() {}

func (dst *FunctionCallResponse) Decode(src []byte) error {
	if len(src) < 4 {
		return &invalidMessageFormatErr{messageType: "FunctionCallResponse"}
	}
	rp := 0
	resultSize := int(binary.BigEndian.Uint32(src[rp:]))
	rp += 4

	if resultSize == -1 {
		dst.Result = nil
		return nil
	}

	if len(src[rp:]) != resultSize {
		return &invalidMessageFormatErr{messageType: "FunctionCallResponse"}
	}

	dst.Result = src[rp:]
	return nil
}

func (src *FunctionCallResponse) Encode(dst []byte) []byte {
	dst = append(dst, 'V')
	sp := len(dst)
	dst = pgio.AppendInt32(dst, -1)

	if src.Result == nil {
		dst = pgio.AppendInt32(dst, -1)
	} else {
		dst = pgio.AppendInt32(dst, int32(len(src.Result)))
		dst = append(dst, src.Result...)
	}

	pgio.SetInt32(dst[sp:], int32(len(dst[sp:])))

	return dst
}

func (src *FunctionCallResponse) MarshalJSON() ([]byte, error) {
	var formattedValue map[string]string
	var hasNonPrintable bool
	for _, b := range src.Result {
		if b < 32 {
			hasNonPrintable = true
			break
		}
	}

	if hasNonPrintable {
		formattedValue = map[string]string{"binary": hex.EncodeToString(src.Result)}
	} else {
		formattedValue = map[string]string{"text": string(src.Result)}
	}

	return json.Marshal(struct {
		Type   string
		Result map[string]string
	}{
		Type:   "FunctionCallResponse",
		Result: formattedValue,
	})
}
