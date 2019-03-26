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

type Parse struct {
	Name          string
	Query         string
	ParameterOIDs []uint32
}

func (*Parse) Frontend() {}

func (dst *Parse) Decode(src []byte) error {
	*dst = Parse{}

	buf := bytes.NewBuffer(src)

	b, err := buf.ReadBytes(0)
	if err != nil {
		return err
	}
	dst.Name = string(b[:len(b)-1])

	b, err = buf.ReadBytes(0)
	if err != nil {
		return err
	}
	dst.Query = string(b[:len(b)-1])

	if buf.Len() < 2 {
		return &invalidMessageFormatErr{messageType: "Parse"}
	}
	parameterOIDCount := int(binary.BigEndian.Uint16(buf.Next(2)))

	for i := 0; i < parameterOIDCount; i++ {
		if buf.Len() < 4 {
			return &invalidMessageFormatErr{messageType: "Parse"}
		}
		dst.ParameterOIDs = append(dst.ParameterOIDs, binary.BigEndian.Uint32(buf.Next(4)))
	}

	return nil
}

func (src *Parse) Encode(dst []byte) []byte {
	dst = append(dst, 'P')
	sp := len(dst)
	dst = pgio.AppendInt32(dst, -1)

	dst = append(dst, src.Name...)
	dst = append(dst, 0)
	dst = append(dst, src.Query...)
	dst = append(dst, 0)

	dst = pgio.AppendUint16(dst, uint16(len(src.ParameterOIDs)))
	for _, oid := range src.ParameterOIDs {
		dst = pgio.AppendUint32(dst, oid)
	}

	pgio.SetInt32(dst[sp:], int32(len(dst[sp:])))

	return dst
}

func (src *Parse) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type          string
		Name          string
		Query         string
		ParameterOIDs []uint32
	}{
		Type:          "Parse",
		Name:          src.Name,
		Query:         src.Query,
		ParameterOIDs: src.ParameterOIDs,
	})
}
