package store

import (
	"encoding/hex"
	"fmt"
	"github.com/golang/protobuf/proto"
	"testing"
)

func TestSetRequest_String(t *testing.T) {
	s := &Command{
		Operation: Operation_SET,
		Key:       []byte("Key"),
		Value:     []byte("Value"),
	}
	g := &Command{
		Operation: Operation_GET,
		Key:       []byte("Key"),
		Value:     []byte("Value"),
	}
	d := &Command{
		Operation: Operation_DELETE,
		Key:       []byte("Key"),
		Value:     []byte("Value"),
	}
	sb, err := proto.Marshal(s)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	gb, err := proto.Marshal(g)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	db, err := proto.Marshal(d)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	fmt.Printf("GET: %s", hex.Dump(sb))
	fmt.Printf("SET: %s", hex.Dump(gb))
	fmt.Printf("DEL: %s", hex.Dump(db))
}
