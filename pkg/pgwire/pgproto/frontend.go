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

type Frontend struct {
	cr *chunkreader.ChunkReader
	w  io.Writer

	// Backend message flyweights
	authentication       Authentication
	backendKeyData       BackendKeyData
	bindComplete         BindComplete
	closeComplete        CloseComplete
	commandComplete      CommandComplete
	copyBothResponse     CopyBothResponse
	copyData             CopyData
	copyInResponse       CopyInResponse
	copyOutResponse      CopyOutResponse
	dataRow              DataRow
	emptyQueryResponse   EmptyQueryResponse
	errorResponse        ErrorResponse
	functionCallResponse FunctionCallResponse
	noData               NoData
	noticeResponse       NoticeResponse
	notificationResponse NotificationResponse
	parameterDescription ParameterDescription
	parameterStatus      ParameterStatus
	parseComplete        ParseComplete
	readyForQuery        ReadyForQuery
	rowDescription       RowDescription

	bodyLen    int
	msgType    byte
	partialMsg bool
}

func NewFrontend(r io.Reader, w io.Writer) (*Frontend, error) {
	cr := chunkreader.NewChunkReader(r)
	return &Frontend{cr: cr, w: w}, nil
}

func (b *Frontend) Send(msg FrontendMessage) error {
	_, err := b.w.Write(msg.Encode(nil))
	return err
}

func (b *Frontend) Receive() (BackendMessage, error) {
	if !b.partialMsg {
		header, err := b.cr.Next(5)
		if err != nil {
			return nil, err
		}

		b.msgType = header[0]
		b.bodyLen = int(binary.BigEndian.Uint32(header[1:])) - 4
		b.partialMsg = true
	}

	var msg BackendMessage
	switch b.msgType {
	case '1':
		msg = &b.parseComplete
	case '2':
		msg = &b.bindComplete
	case '3':
		msg = &b.closeComplete
	case 'A':
		msg = &b.notificationResponse
	case 'C':
		msg = &b.commandComplete
	case 'd':
		msg = &b.copyData
	case 'D':
		msg = &b.dataRow
	case 'E':
		msg = &b.errorResponse
	case 'G':
		msg = &b.copyInResponse
	case 'H':
		msg = &b.copyOutResponse
	case 'I':
		msg = &b.emptyQueryResponse
	case 'K':
		msg = &b.backendKeyData
	case 'n':
		msg = &b.noData
	case 'N':
		msg = &b.noticeResponse
	case 'R':
		msg = &b.authentication
	case 'S':
		msg = &b.parameterStatus
	case 't':
		msg = &b.parameterDescription
	case 'T':
		msg = &b.rowDescription
	case 'V':
		msg = &b.functionCallResponse
	case 'W':
		msg = &b.copyBothResponse
	case 'Z':
		msg = &b.readyForQuery
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
