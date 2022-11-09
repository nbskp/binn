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

func TestBinnWithReceiver(t *testing.T) {
	bn := New(NewBottleStorage(100), time.Duration(0))
	bn.Publish(&Bottle{
		Msg: "sample message",
	})
	r := NewReceiver()
	select {
	case <-time.After(100 * time.Millisecond):
		assert.Fail(t, "failed")
	case b := <-r.Receive(bn):
		assert.Equal(t, "sample message", b.Msg)
	}
}

func TestBinnFilledChan(t *testing.T) {
	bn := New(NewBottleStorage(100), time.Duration(0))
	bn.Publish(&Bottle{
		Msg: "sample message1",
	})
	bn.Publish(&Bottle{
		Msg: "sample message2",
	})
	ch := make(chan *Bottle, 1)
	err := bn.Subscribe(ch)
	time.Sleep(1 * time.Millisecond)
	require.NoError(t, err)
	select {
	case <-time.After(100 * time.Millisecond):
		assert.Fail(t, "failed")
	case b, ok := <-ch:
		assert.False(t, ok)
		assert.Nil(t, b)
	}
}

func TestBinnFilledChanWithReceiver(t *testing.T) {
	bn := New(NewBottleStorage(100), time.Duration(0))
	bn.Publish(&Bottle{
		Msg: "sample message1",
	})
	bn.Publish(&Bottle{
		Msg: "sample message2",
	})
	r := NewReceiver()
	select {
	case <-time.After(100 * time.Millisecond):
		assert.Fail(t, "failed")
	case b := <-r.Receive(bn):
		assert.Equal(t, "sample message1", b.Msg)
	}
}
