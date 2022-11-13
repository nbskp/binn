package binn

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBinn(t *testing.T) {
	bn := New(NewBottleStorage(100), time.Duration(0))
	bn.Publish(&Bottle{
		Msg: "sample message",
	})
	ch := make(chan struct{})
	err := bn.Subscribe(func(b *Bottle) bool {
		ch <- struct{}{}
		assert.Equal(t, "sample message", b.Msg)
		return true
	})
	require.NoError(t, err)
	select {
	case <-time.After(100 * time.Millisecond):
		assert.Fail(t, "failed")
	case <-ch:
	}
}

func TestBinnHandlerIncomplete(t *testing.T) {
	s := NewBottleStorage(100)
	bn := New(s, time.Duration(0))
	bn.Publish(&Bottle{
		Msg: "sample message",
	})
	ch := make(chan struct{})
	err := bn.Subscribe(func(b *Bottle) bool {
		ch <- struct{}{}
		return false
	})
	require.NoError(t, err)
	select {
	case <-time.After(100 * time.Millisecond):
		assert.Fail(t, "failed")
	case <-ch:
	}
	b, err := s.Get()
	require.NoError(t, err)
	assert.Equal(t, "sample message", b.Msg)
}
