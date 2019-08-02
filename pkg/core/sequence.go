package core

import (
	"github.com/elliotcourant/noahdb/pkg/frunk"
)

type sequenceContext struct {
	*base
}

type SequenceContext interface {
	GetSequenceChunk(name string) (*frunk.SequenceChunkResponse, error)
}

func (ctx *base) Sequences() SequenceContext {
	return &sequenceContext{
		ctx,
	}
}

func (s *sequenceContext) GetSequenceChunk(name string) (*frunk.SequenceChunkResponse, error) {
	return s.db.GetSequenceChunk(name)
}
