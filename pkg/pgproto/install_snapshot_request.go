package pgproto

import (
	"bytes"
	"encoding/binary"
	"github.com/elliotcourant/noahdb/pkg/pgio"
	"github.com/hashicorp/raft"
	"io"
)

type InstallSnapshotRequest struct {
	raft.InstallSnapshotRequest

	SnapshotData []byte
}

func (installSnapshot *InstallSnapshotRequest) Reader() io.Reader {
	return bytes.NewBuffer(installSnapshot.SnapshotData)
}

func (InstallSnapshotRequest) Frontend() {}

func (InstallSnapshotRequest) Raft() {}

func (installSnapshot *InstallSnapshotRequest) Decode(src []byte) error {
	*installSnapshot = InstallSnapshotRequest{}
	buf := bytes.NewBuffer(src)
	installSnapshot.ProtocolVersion = raft.ProtocolVersion(binary.BigEndian.Uint32(buf.Next(4)))
	installSnapshot.SnapshotVersion = raft.SnapshotVersion(binary.BigEndian.Uint32(buf.Next(4)))
	installSnapshot.Term = binary.BigEndian.Uint64(buf.Next(8))

	leaderLen := int(binary.BigEndian.Uint32(buf.Next(4)))
	installSnapshot.Leader = buf.Next(leaderLen)

	installSnapshot.LastLogIndex = binary.BigEndian.Uint64(buf.Next(8))
	installSnapshot.LastLogTerm = binary.BigEndian.Uint64(buf.Next(8))

	peersLen := int(binary.BigEndian.Uint32(buf.Next(4)))
	installSnapshot.Peers = buf.Next(peersLen)

	configLen := int(binary.BigEndian.Uint32(buf.Next(4)))
	installSnapshot.Configuration = buf.Next(configLen)

	installSnapshot.ConfigurationIndex = binary.BigEndian.Uint64(buf.Next(8))
	installSnapshot.Size = int64(binary.BigEndian.Uint64(buf.Next(8)))

	installSnapshot.SnapshotData = buf.Next(int(installSnapshot.Size))
	return nil
}

func (installSnapshot *InstallSnapshotRequest) Encode(dst []byte) []byte {
	dst = append(dst, RaftInstallSnapshotRequest)
	sp := len(dst)
	dst = pgio.AppendInt32(dst, -1)

	dst = pgio.AppendInt32(dst, int32(installSnapshot.ProtocolVersion))
	dst = pgio.AppendInt32(dst, int32(installSnapshot.SnapshotVersion))

	dst = pgio.AppendUint64(dst, installSnapshot.Term)

	dst = pgio.AppendInt32(dst, int32(len(installSnapshot.Leader)))
	dst = append(dst, installSnapshot.Leader...)

	dst = pgio.AppendUint64(dst, installSnapshot.LastLogIndex)
	dst = pgio.AppendUint64(dst, installSnapshot.LastLogTerm)

	dst = pgio.AppendInt32(dst, int32(len(installSnapshot.Peers)))
	dst = append(dst, installSnapshot.Peers...)

	dst = pgio.AppendInt32(dst, int32(len(installSnapshot.Configuration)))
	dst = append(dst, installSnapshot.Configuration...)

	dst = pgio.AppendUint64(dst, installSnapshot.ConfigurationIndex)
	dst = pgio.AppendInt64(dst, installSnapshot.Size)

	dst = append(dst, installSnapshot.SnapshotData...)

	pgio.SetInt32(dst[sp:], int32(len(dst[sp:])))
	return dst
}
