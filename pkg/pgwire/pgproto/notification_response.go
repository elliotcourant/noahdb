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

type NotificationResponse struct {
	PID     uint32
	Channel string
	Payload string
}

func (*NotificationResponse) Backend() {}

func (dst *NotificationResponse) Decode(src []byte) error {
	buf := bytes.NewBuffer(src)

	pid := binary.BigEndian.Uint32(buf.Next(4))

	b, err := buf.ReadBytes(0)
	if err != nil {
		return err
	}
	channel := string(b[:len(b)-1])

	b, err = buf.ReadBytes(0)
	if err != nil {
		return err
	}
	payload := string(b[:len(b)-1])

	*dst = NotificationResponse{PID: pid, Channel: channel, Payload: payload}
	return nil
}

func (src *NotificationResponse) Encode(dst []byte) []byte {
	dst = append(dst, 'A')
	sp := len(dst)
	dst = pgio.AppendInt32(dst, -1)

	dst = append(dst, src.Channel...)
	dst = append(dst, 0)
	dst = append(dst, src.Payload...)
	dst = append(dst, 0)

	pgio.SetInt32(dst[sp:], int32(len(dst[sp:])))

	return dst
}

func (src *NotificationResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type    string
		PID     uint32
		Channel string
		Payload string
	}{
		Type:    "NotificationResponse",
		PID:     src.PID,
		Channel: src.Channel,
		Payload: src.Payload,
	})
}
