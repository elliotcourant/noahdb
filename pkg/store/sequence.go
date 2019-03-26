package store

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/readystock/golog"
	"github.com/readystock/raft"
	"sync"
)

const (
	SequenceRangeSize   = 1000
	SequencePartitions  = 5
	SequencePreretrieve = 50
)

func (store *Store) getSequenceChunk(sequenceName string) (*SequenceChunkResponse, error) {
	if store.raft.State() != raft.Leader { // Only the leader can manage sequences
		return store.clusterClient.getNextChunkInSequence(sequenceName)
	}
	store.sequenceCacheSync.Lock()
	defer store.sequenceCacheSync.Unlock()
	sequenceCache, ok := store.sequenceCache[sequenceName]
	path := []byte(fmt.Sprintf("%s%s", sequencePath, sequenceName))
	if !ok {
		seq, err := store.Get(path)
		if err != nil {
			return nil, err
		}
		if len(seq) == 0 {
			sequenceCache = &Sequence{
				CurrentValue:       0,
				LastPartitionIndex: 0,
				MaxPartitionIndex:  SequencePartitions - 1,
				Partitions:         SequencePartitions,
			}
		} else {
			sequenceCache = &Sequence{}
			err = proto.Unmarshal(seq, sequenceCache)
			if err != nil {
				return nil, err
			}
		}
		store.sequenceCache[sequenceName] = sequenceCache
	}
	if sequenceCache.LastPartitionIndex >= sequenceCache.MaxPartitionIndex {
		sequenceCache.CurrentValue += SequenceRangeSize
		sequenceCache.LastPartitionIndex = 0
		sequenceCache.MaxPartitionIndex = SequencePartitions - 1
	}
	index := sequenceCache.LastPartitionIndex
	sequenceCache.LastPartitionIndex++
	b, err := proto.Marshal(sequenceCache)
	if err != nil {
		return nil, err
	}
	err = store.Set(path, b)
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

func (store *Store) NextSequenceValueById(sequenceName string) (*uint64, error) {
	store.chunkMapMutex.Lock()
	defer store.chunkMapMutex.Unlock()
	chunk, ok := store.sequenceChunks[sequenceName]
	if !ok {
		chunk = &SequenceChunk{
			current:      nil,
			next:         nil,
			index:        1,
			sequenceName: sequenceName,
			sync:         new(sync.Mutex),
			Store:        *store,
		}
		store.sequenceChunks[sequenceName] = chunk
	}
	return chunk.Next()
}

func (store *Store) SequenceIndexById(sequenceName string) (uint64, error) {
	store.chunkMapMutex.Lock()
	defer store.chunkMapMutex.Unlock()
	chunk, ok := store.sequenceChunks[sequenceName]
	if !ok {
		chunk = &SequenceChunk{
			current:      nil,
			next:         nil,
			index:        1,
			sequenceName: sequenceName,
			sync:         new(sync.Mutex),
			Store:        *store,
		}
		store.sequenceChunks[sequenceName] = chunk
	}
	return chunk.GetSequenceIndex(), nil
}

func (sequence *SequenceChunk) GetSequenceIndex() uint64 {
	return sequence.index
}

func (sequence *SequenceChunk) Next() (*uint64, error) {
	sequence.sync.Lock()
	defer sequence.sync.Unlock()
	if sequence.current == nil {
		chunk, err := sequence.Store.getSequenceChunk(sequence.sequenceName)
		if err != nil {
			return nil, err
		}
		sequence.current = chunk
		sequence.next = nil
		sequence.index = 1
	}
NewId:
	nextId := sequence.current.Start + sequence.current.Offset + (sequence.current.Count * sequence.index) - (sequence.current.Count - 1)
	if nextId > sequence.current.End {
		golog.Debugf("moving next chunk into current sequence [%s]", sequence.sequenceName)
		if sequence.next != nil {
			sequence.current = sequence.next
		} else {
			chunk, err := sequence.Store.getSequenceChunk(sequence.sequenceName)
			if err != nil {
				return nil, err
			}
			sequence.current = chunk
		}
		sequence.next = nil
		sequence.index = 1
		goto NewId
	}
	if sequence.next == nil && float64(sequence.index*sequence.current.Count)/float64(sequence.current.End-sequence.current.Start) > (float64(SequencePreretrieve)/100) {
		golog.Debugf("requesting next chunk in sequence [%s] preemptive", sequence.sequenceName)
		chunk, err := sequence.Store.getSequenceChunk(sequence.sequenceName)
		if err != nil {
			return nil, err
		}
		sequence.next = chunk
	}
	sequence.index++
	return &nextId, nil
}
