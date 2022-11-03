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
	ch := make(chan *Bottle, 1)
	err := bn.Subscribe(ch)
	require.NoError(t, err)
	select {
	case <-time.After(100 * time.Millisecond):
		assert.Fail(t, "failed")
	case b, ok := <-ch:
		assert.True(t, ok)
		assert.Equal(t, "sample message", b.Msg)
	}
}

func TestBinnUnreceivedChannel(t *testing.T) {
	bn := New(NewBottleStorage(100), time.Duration(0))
	bn.Publish(&Bottle{
		Msg: "sample message1",
	})
	bn.Publish(&Bottle{
		Msg: "sample message2",
	})
	ch := make(chan *Bottle, 1)
	err := bn.Subscribe(ch)
	require.NoError(t, err)
	time.Sleep(1 * time.Millisecond)
	select {
	case b, ok := <-ch:
		assert.Nil(t, b)
		assert.False(t, ok)
	}
}
