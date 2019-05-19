package pgproto

import (
	"bytes"
	"encoding/binary"
	"github.com/elliotcourant/noahdb/pkg/pgio"
	"github.com/hashicorp/raft"
)

type AppendEntriesRequest struct {
	raft.AppendEntriesRequest
}

func (AppendEntriesRequest) Frontend() {}

func (AppendEntriesRequest) Raft() {}

func (appendEntries *AppendEntriesRequest) Decode(src []byte) error {
	*appendEntries = AppendEntriesRequest{}
	buf := bytes.NewBuffer(src)
	appendEntries.ProtocolVersion = raft.ProtocolVersion(binary.BigEndian.Uint32(buf.Next(4)))
	appendEntries.Term = binary.BigEndian.Uint64(buf.Next(8))

	leaderLen := int(binary.BigEndian.Uint32(buf.Next(4)))
	appendEntries.Leader = buf.Next(leaderLen)

	appendEntries.PrevLogEntry = binary.BigEndian.Uint64(buf.Next(8))
	appendEntries.PrevLogTerm = binary.BigEndian.Uint64(buf.Next(8))

	numberOfEntries := int(binary.BigEndian.Uint16(buf.Next(2)))
	appendEntries.Entries = make([]*raft.Log, numberOfEntries)
	for i := 0; i < numberOfEntries; i++ {
		next := buf.Next(4)
		entrySize := int(int32(binary.BigEndian.Uint32(next)))

		// null
		if entrySize == -1 {
			appendEntries.Entries[i] = nil
		} else {
			var entry raft.Log
			entryBuf := bytes.NewBuffer(buf.Next(entrySize))
			entry.Index = binary.BigEndian.Uint64(entryBuf.Next(8))
			entry.Term = binary.BigEndian.Uint64(entryBuf.Next(8))
			entry.Type = raft.LogType(entryBuf.Next(1)[0])

			dataLen := int(binary.BigEndian.Uint32(entryBuf.Next(4)))
			entry.Data = entryBuf.Next(dataLen)
			appendEntries.Entries[i] = &entry
		}
	}

	appendEntries.LeaderCommitIndex = binary.BigEndian.Uint64(buf.Next(8))

	return nil
}

func (appendEntries *AppendEntriesRequest) Encode(dst []byte) []byte {
	dst = append(dst, RaftAppendEntriesRequest)
	sp := len(dst)
	dst = pgio.AppendInt32(dst, -1)

	dst = pgio.AppendInt32(dst, int32(appendEntries.ProtocolVersion))
	dst = pgio.AppendUint64(dst, appendEntries.Term)

	dst = pgio.AppendInt32(dst, int32(len(appendEntries.Leader)))
	dst = append(dst, appendEntries.Leader...)

	dst = pgio.AppendUint64(dst, appendEntries.PrevLogEntry)
	dst = pgio.AppendUint64(dst, appendEntries.PrevLogTerm)

	dst = pgio.AppendUint16(dst, uint16(len(appendEntries.Entries)))
	for _, entry := range appendEntries.Entries {
		if entry == nil {
			dst = pgio.AppendInt32(dst, -1)
			continue
		}

		logBytes := make([]byte, 0)
		logBytes = pgio.AppendUint64(logBytes, entry.Index)
		logBytes = pgio.AppendUint64(logBytes, entry.Term)
		logBytes = append(logBytes, byte(entry.Type))

		logBytes = pgio.AppendInt32(logBytes, int32(len(entry.Data)))
		logBytes = append(logBytes, entry.Data...)

		lenLog := int32(len(logBytes))
		dst = pgio.AppendInt32(dst, lenLog)
		dst = append(dst, logBytes...)
	}

	dst = pgio.AppendUint64(dst, appendEntries.LeaderCommitIndex)

	pgio.SetInt32(dst[sp:], int32(len(dst[sp:])))

	return dst
}
