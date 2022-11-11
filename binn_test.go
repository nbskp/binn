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
	err := bn.Subscribe(func(b *Bottle) {
		ch <- struct{}{}
		assert.Equal(t, "sample message", b.Msg)
	})
	require.NoError(t, err)
	select {
	case <-time.After(100 * time.Millisecond):
		assert.Fail(t, "failed")
	case <-ch:
	}
}
