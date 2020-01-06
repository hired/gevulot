package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerEventHasFired(t *testing.T) {
	e := NewEvent()

	assert.False(t, e.HasFired())
	assert.True(t, e.Fire())
	assert.True(t, e.HasFired())
	assert.False(t, e.Fire())
}

func TestServerEventDoneChannel(t *testing.T) {
	e := NewEvent()

	select {
	case <-e.Done():
		assert.FailNow(t, "Done channel is closed")
	default:
	}

	e.Fire()

	select {
	case <-e.Done():
	default:
		assert.FailNow(t, "Done channel is open")
	}
}
