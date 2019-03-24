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

import "fmt"

// Message is the interface implemented by an object that can decode and encode
// a particular PostgreSQL message.
type Message interface {
	// Decode is allowed and expected to retain a reference to data after
	// returning (unlike encoding.BinaryUnmarshaler).
	Decode(data []byte) error

	// Encode appends itself to dst and returns the new buffer.
	Encode(dst []byte) []byte
}

type FrontendMessage interface {
	Message
	Frontend() // no-op method to distinguish frontend from backend methods
}

type BackendMessage interface {
	Message
	Backend() // no-op method to distinguish frontend from backend methods
}

type invalidMessageLenErr struct {
	messageType string
	expectedLen int
	actualLen   int
}

func (e *invalidMessageLenErr) Error() string {
	return fmt.Sprintf("%s body must have length of %d, but it is %d", e.messageType, e.expectedLen, e.actualLen)
}

type invalidMessageFormatErr struct {
	messageType string
}

func (e *invalidMessageFormatErr) Error() string {
	return fmt.Sprintf("%s body is invalid", e.messageType)
}
