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
	"io"

	"github.com/pkg/errors"
	"github.com/readystock/pgx/chunkreader"
)

type Backend struct {
	cr *chunkreader.ChunkReader
	w  io.Writer

	// Frontend message flyweights
	bind            Bind
	_close          Close
	describe        Describe
	execute         Execute
	flush           Flush
	parse           Parse
	passwordMessage PasswordMessage
	query           Query
	startupMessage  StartupMessage
	sync            Sync
	terminate       Terminate

	bodyLen    int
	msgType    byte
	partialMsg bool
}

func NewBackend(r io.Reader, w io.Writer) (*Backend, error) {
	cr := chunkreader.NewChunkReader(r)
	return &Backend{cr: cr, w: w}, nil
}

func (b *Backend) Send(msg BackendMessage) error {
	_, err := b.w.Write(msg.Encode(nil))
	return err
}

func (b *Backend) ReceiveStartupMessage() (*StartupMessage, error) {
	buf, err := b.cr.Next(4)
	if err != nil {
		return nil, err
	}
	msgSize := int(binary.BigEndian.Uint32(buf) - 4)

	buf, err = b.cr.Next(msgSize)
	if err != nil {
		return nil, err
	}

	err = b.startupMessage.Decode(buf)
	if err != nil {
		return nil, err
	}

	return &b.startupMessage, nil
}

func (b *Backend) Receive() (FrontendMessage, error) {
	if !b.partialMsg {
		header, err := b.cr.Next(5)
		if err != nil {
			return nil, err
		}

		b.msgType = header[0]
		b.bodyLen = int(binary.BigEndian.Uint32(header[1:])) - 4
		b.partialMsg = true
	}

	var msg FrontendMessage
	switch b.msgType {
	case 'B':
		msg = &b.bind
	case 'C':
		msg = &b._close
	case 'D':
		msg = &b.describe
	case 'E':
		msg = &b.execute
	case 'H':
		msg = &b.flush
	case 'P':
		msg = &b.parse
	case 'p':
		msg = &b.passwordMessage
	case 'Q':
		msg = &b.query
	case 'S':
		msg = &b.sync
	case 'X':
		msg = &b.terminate
	default:
		return nil, errors.Errorf("unknown message type: %c", b.msgType)
	}

	msgBody, err := b.cr.Next(b.bodyLen)
	if err != nil {
		return nil, err
	}

	b.partialMsg = false

	err = msg.Decode(msgBody)
	return msg, err
}
