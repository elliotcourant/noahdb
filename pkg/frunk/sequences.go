package frunk

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/pgproto"
	"github.com/readystock/golog"
	"net"
	"sync"
)

type SequenceChunkResponse struct {
	SequenceName string
	Start        uint64
	End          uint64
	Offset       uint64
	Count        uint64
}

const (
	SequenceRangeSize   = 1000
	SequencePartitions  = 5
	SequencePreretrieve = 50
)

func (s *Store) getNextChunkInSequence(sequenceName string) error {
	leaderAddress := s.LeaderAddr()
	addr, _ := net.ResolveTCPAddr("tcp", leaderAddress)
	conn, err := net.DialTCP("tpc", nil, addr)
	if err != nil {
		return err
	}
	frontend, err := pgproto.NewFrontend(conn, conn)
	if err != nil {
		return err
	}
	if err := frontend.Send(&pgproto.RpcStartupMessage{}); err != nil {
		return err
	}

	return nil
}

func (s *Store) getSequenceChunk(sequenceName string) (*SequenceChunkResponse, error) {
	if !s.IsLeader() { // Only the leader can manage sequences
		// return s.clusterClient.getNextChunkInSequence(sequenceName)
	}
	s.sequenceCacheSync.Lock()
	defer s.sequenceCacheSync.Unlock()
	sequenceCache, ok := s.sequenceCache[sequenceName]
	path := []byte(fmt.Sprintf("%s%s", sequencePath, sequenceName))
	if !ok {
		seq, err := s.Get(path)
		if err != nil {
			return nil, err
		}
		if len(seq) == 0 {
			sequenceCache = &pgproto.SequenceResponse{
				CurrentValue:       0,
				LastPartitionIndex: 0,
				MaxPartitionIndex:  SequencePartitions - 1,
				Partitions:         SequencePartitions,
			}
		} else {
			sequenceCache = &pgproto.SequenceResponse{}
			err := sequenceCache.Decode(seq)
			if err != nil {
				return nil, err
			}
		}
		s.sequenceCache[sequenceName] = sequenceCache
	}
	if sequenceCache.LastPartitionIndex >= sequenceCache.MaxPartitionIndex {
		sequenceCache.CurrentValue += SequenceRangeSize
		sequenceCache.LastPartitionIndex = 0
		sequenceCache.MaxPartitionIndex = SequencePartitions - 1
	}
	index := sequenceCache.LastPartitionIndex
	sequenceCache.LastPartitionIndex++
	b := sequenceCache.Encode(nil)
	err := s.Set(path, b)
	if err != nil {
		return nil, err
	}
	return &SequenceChunkResponse{
		Start:  sequenceCache.CurrentValue,
		End:    sequenceCache.CurrentValue + SequenceRangeSize,
		Offset: index,
		Count:  sequenceCache.MaxPartitionIndex,
	}, nil
}

type SequenceChunk struct {
	Store
	current      *SequenceChunkResponse
	next         *SequenceChunkResponse
	index        uint64
	sync         *sync.Mutex
	sequenceName string
}

func (s *Store) NextSequenceValueById(sequenceName string) (uint64, error) {
	s.chunkMapMutex.Lock()
	defer s.chunkMapMutex.Unlock()
	chunk, ok := s.sequenceChunks[sequenceName]
	if !ok {
		chunk = &SequenceChunk{
			current:      nil,
			next:         nil,
			index:        1,
			sequenceName: sequenceName,
			sync:         new(sync.Mutex),
			Store:        *s,
		}
		s.sequenceChunks[sequenceName] = chunk
	}
	return chunk.Next()
}

func (s *Store) SequenceIndexById(sequenceName string) (uint64, error) {
	s.chunkMapMutex.Lock()
	defer s.chunkMapMutex.Unlock()
	chunk, ok := s.sequenceChunks[sequenceName]
	if !ok {
		chunk = &SequenceChunk{
			current:      nil,
			next:         nil,
			index:        1,
			sequenceName: sequenceName,
			sync:         new(sync.Mutex),
			Store:        *s,
		}
		s.sequenceChunks[sequenceName] = chunk
	}
	return chunk.GetSequenceIndex(), nil
}

func (sequence *SequenceChunk) GetSequenceIndex() uint64 {
	return sequence.index
}

func (sequence *SequenceChunk) Next() (uint64, error) {
	sequence.sync.Lock()
	defer sequence.sync.Unlock()
	if sequence.current == nil {
		chunk, err := sequence.Store.getSequenceChunk(sequence.sequenceName)
		if err != nil {
			return 0, err
		}
		sequence.current = chunk
		sequence.next = nil
		sequence.index = 1
	}
NewId:
	nextId := sequence.current.Start + sequence.current.Offset + (sequence.current.Count * sequence.index) - (sequence.current.Count - 1)
	if nextId > sequence.current.End {
		golog.Verbosef("moving next chunk into current sequence [%s]", sequence.sequenceName)
		if sequence.next != nil {
			sequence.current = sequence.next
		} else {
			chunk, err := sequence.Store.getSequenceChunk(sequence.sequenceName)
			if err != nil {
				return 0, err
			}
			sequence.current = chunk
		}
		sequence.next = nil
		sequence.index = 1
		goto NewId
	}
	if sequence.next == nil && float64(sequence.index*sequence.current.Count)/float64(sequence.current.End-sequence.current.Start) > (float64(SequencePreretrieve)/100) {
		golog.Verbosef("requesting next chunk in sequence [%s] preemptive", sequence.sequenceName)
		chunk, err := sequence.Store.getSequenceChunk(sequence.sequenceName)
		if err != nil {
			return 0, err
		}
		sequence.next = chunk
	}
	sequence.index++
	return nextId, nil
}
