package binn

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBinn(t *testing.T) {
	bn := New(NewBottleStorage(100), time.Duration(0))
	bn.Add(&Bottle{
		Id:  "1234abc",
		Msg: "sample message",
	})
	select {
	case <-time.After(100 * time.Millisecond):
		assert.Fail(t, "failed")
	case b := <-bn.Get():
		assert.Equal(t, "1234abc", b.Id)
		assert.Equal(t, "sample message", b.Msg)
	}
}
