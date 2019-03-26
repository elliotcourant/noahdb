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
	"encoding/json"

	"github.com/readystock/pgx/pgio"
)

type ParameterStatus struct {
	Name  string
	Value string
}

func (*ParameterStatus) Backend() {}

func (dst *ParameterStatus) Decode(src []byte) error {
	buf := bytes.NewBuffer(src)

	b, err := buf.ReadBytes(0)
	if err != nil {
		return err
	}
	name := string(b[:len(b)-1])

	b, err = buf.ReadBytes(0)
	if err != nil {
		return err
	}
	value := string(b[:len(b)-1])

	*dst = ParameterStatus{Name: name, Value: value}
	return nil
}

func (src *ParameterStatus) Encode(dst []byte) []byte {
	dst = append(dst, 'S')
	sp := len(dst)
	dst = pgio.AppendInt32(dst, -1)

	dst = append(dst, src.Name...)
	dst = append(dst, 0)
	dst = append(dst, src.Value...)
	dst = append(dst, 0)

	pgio.SetInt32(dst[sp:], int32(len(dst[sp:])))

	return dst
}

func (ps *ParameterStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type  string
		Name  string
		Value string
	}{
		Type:  "ParameterStatus",
		Name:  ps.Name,
		Value: ps.Value,
	})
}
