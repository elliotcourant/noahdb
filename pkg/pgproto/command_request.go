package pgproto

import (
	"bytes"
	"encoding/binary"
	"github.com/elliotcourant/noahdb/pkg/pgio"
)

type RpcCommandType int

const (
	RpcCommandType_Execute RpcCommandType = iota
	RpcCommandType_Query
	RpcCommandType_MetadataSet
	RpcCommandType_MetadataDelete
	RpcCommandType_Connect
	RpcCommandType_Disconnect
	RpcCommandType_KeyValueSet
	RpcCommandType_KeyValueDelete
)

type KeyValue struct {
	Key   []byte
	Value []byte
}

type CommandRequest struct {
	CommandType     RpcCommandType
	Queries         []string
	KeyValueSets    []KeyValue
	KeyValueDeletes [][]byte
}

func (CommandRequest) Frontend() {}

func (CommandRequest) RpcFrontend() {}

func (item *CommandRequest) Decode(src []byte) error {
	*item = CommandRequest{}
	buf := bytes.NewBuffer(src)

	item.CommandType = RpcCommandType(binary.BigEndian.Uint32(buf.Next(4)))

	numberOfQueries := int(binary.BigEndian.Uint16(buf.Next(2)))
	item.Queries = make([]string, numberOfQueries)
	for i := 0; i < numberOfQueries; i++ {
		next := buf.Next(4) // Get this particular queries length
		size := int(int32(binary.BigEndian.Uint32(next)))

		query := bytes.NewBuffer(buf.Next(size))
		item.Queries[i] = query.String()
	}

	numberOfKeyValueSets := int(binary.BigEndian.Uint16(buf.Next(2)))
	item.KeyValueSets = make([]KeyValue, numberOfKeyValueSets)
	for i := 0; i < numberOfKeyValueSets; i++ {
		next := buf.Next(4)
		size := int(int32(binary.BigEndian.Uint32(next)))

		keyValueSet := bytes.NewBuffer(buf.Next(size))

		kv := KeyValue{}

		keyLength := int(int32(binary.BigEndian.Uint32(keyValueSet.Next(4))))
		kv.Key = keyValueSet.Next(keyLength)

		valueLength := int(int32(binary.BigEndian.Uint32(keyValueSet.Next(4))))
		kv.Value = keyValueSet.Next(valueLength)

		item.KeyValueSets[i] = kv
	}

	numberOfDeletes := int(binary.BigEndian.Uint16(buf.Next(2)))
	item.KeyValueDeletes = make([][]byte, numberOfDeletes)
	for i := 0; i < numberOfDeletes; i++ {
		next := buf.Next(4)
		size := int(int32(binary.BigEndian.Uint32(next)))

		item.KeyValueDeletes[i] = buf.Next(size)
	}
	return nil
}

func (item *CommandRequest) Encode(dst []byte) []byte {
	dst = append(dst, RpcCommandRequest)
	sp := len(dst)
	dst = pgio.AppendInt32(dst, -1)

	dst = pgio.AppendInt32(dst, int32(item.CommandType))

	dst = pgio.AppendUint16(dst, uint16(len(item.Queries)))
	for _, query := range item.Queries {
		queryBytes := make([]byte, 0)

		queryBytes = append(queryBytes, query...)

		queryLength := int32(len(queryBytes))
		dst = pgio.AppendInt32(dst, queryLength)
		dst = append(dst, queryBytes...)
	}

	dst = pgio.AppendUint16(dst, uint16(len(item.KeyValueSets)))
	for _, keyValueSet := range item.KeyValueSets {
		setBytes := make([]byte, 0)

		setBytes = pgio.AppendInt32(setBytes, int32(len(keyValueSet.Key)))
		setBytes = append(setBytes, keyValueSet.Key...)

		setBytes = pgio.AppendInt32(setBytes, int32(len(keyValueSet.Value)))
		setBytes = append(setBytes, keyValueSet.Value...)

		setLength := int32(len(setBytes))
		dst = pgio.AppendInt32(dst, setLength)
		dst = append(dst, setBytes...)
	}

	dst = pgio.AppendUint16(dst, uint16(len(item.KeyValueDeletes)))
	for _, keyValueDelete := range item.KeyValueDeletes {
		deleteBytes := make([]byte, 0)

		deleteBytes = append(deleteBytes, keyValueDelete...)

		deleteLength := int32(len(deleteBytes))
		dst = pgio.AppendInt32(dst, deleteLength)
		dst = append(dst, deleteBytes...)
	}

	pgio.SetInt32(dst[sp:], int32(len(dst[sp:])))
	return dst
}
