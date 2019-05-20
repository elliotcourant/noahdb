package pgtransport

import (
	"github.com/hashicorp/raft"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_SettingValueOfPointer(t *testing.T) {
	setVal := func(val interface{}) {
		if r, ok := val.(*raft.AppendEntriesResponse); ok {
			*r = raft.AppendEntriesResponse{
				LastLog: 10,
			}
		}
	}

	setValTop := func(v *raft.AppendEntriesResponse) {
		setVal(v)
	}

	input := raft.AppendEntriesResponse{
		LastLog: 0,
	}

	setValTop(&input)
	assert.Equal(t, uint64(10), input.LastLog)
}
