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

type PasswordMessage struct {
	Password string
}

func (*PasswordMessage) Frontend() {}

func (dst *PasswordMessage) Decode(src []byte) error {
	buf := bytes.NewBuffer(src)

	b, err := buf.ReadBytes(0)
	if err != nil {
		return err
	}
	dst.Password = string(b[:len(b)-1])

	return nil
}

func (src *PasswordMessage) Encode(dst []byte) []byte {
	dst = append(dst, 'p')
	dst = pgio.AppendInt32(dst, int32(4+len(src.Password)+1))

	dst = append(dst, src.Password...)
	dst = append(dst, 0)

	return dst
}

func (src *PasswordMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type     string
		Password string
	}{
		Type:     "PasswordMessage",
		Password: src.Password,
	})
}
