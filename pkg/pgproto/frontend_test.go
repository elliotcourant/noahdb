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

package pgproto_test

import (
	"github.com/readystock/noah/db/sql/pgwire/pgproto"
	"testing"

	"github.com/pkg/errors"
)

type interruptReader struct {
	chunks [][]byte
}

func (ir *interruptReader) Read(p []byte) (n int, err error) {
	if len(ir.chunks) == 0 {
		return 0, errors.New("no data")
	}

	n = copy(p, ir.chunks[0])
	if n != len(ir.chunks[0]) {
		panic("this test reader doesn't support partial reads of chunks")
	}

	ir.chunks = ir.chunks[1:]

	return n, nil
}

func (ir *interruptReader) push(p []byte) {
	ir.chunks = append(ir.chunks, p)
}

func TestFrontendReceiveInterrupted(t *testing.T) {
	t.Parallel()

	server := &interruptReader{}
	server.push([]byte{'Z', 0, 0, 0, 5})

	frontend, err := pgproto.NewFrontend(server, nil)
	if err != nil {
		t.Fatal(err)
	}

	msg, err := frontend.Receive()
	if err == nil {
		t.Fatal("expected err")
	}
	if msg != nil {
		t.Fatalf("did not expect msg, but %v", msg)
	}

	server.push([]byte{'I'})

	msg, err = frontend.Receive()
	if err != nil {
		t.Fatal(err)
	}
	if msg, ok := msg.(*pgproto.ReadyForQuery); !ok || msg.TxStatus != 'I' {
		t.Fatalf("unexpected msg: %v", msg)
	}
}
