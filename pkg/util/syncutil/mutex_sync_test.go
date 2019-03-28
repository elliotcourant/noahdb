package syncutil

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMutex_AssertHeld(t *testing.T) {
	mut := new(Mutex)
	mut.Lock()
	mut.AssertHeld()
	mut.Unlock()
	assert.Panics(t, func() {
		mut.AssertHeld()
	})
}

func TestRWMutex_AssertHeld(t *testing.T) {
	mut := new(RWMutex)
	mut.Lock()
	mut.AssertHeld()
	mut.Unlock()
	assert.Panics(t, func() {
		mut.AssertHeld()
	})
}
