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
)

func TestBackendReceiveInterrupted(t *testing.T) {
	t.Parallel()

	server := &interruptReader{}
	server.push([]byte{'Q', 0, 0, 0, 6})

	backend, err := pgproto.NewBackend(server, nil)
	if err != nil {
		t.Fatal(err)
	}

	msg, err := backend.Receive()
	if err == nil {
		t.Fatal("expected err")
	}
	if msg != nil {
		t.Fatalf("did not expect msg, but %v", msg)
	}

	server.push([]byte{'I', 0})

	msg, err = backend.Receive()
	if err != nil {
		t.Fatal(err)
	}
	if msg, ok := msg.(*pgproto.Query); !ok || msg.String != "I" {
		t.Fatalf("unexpected msg: %v", msg)
	}
}
