package binn

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorageAddBottle(t *testing.T) {
	s := NewBottleStorage(1)
	b := &Bottle{
		Msg: "This is a Test Message",
	}
	err := s.Add(b)
	assert.NoError(t, err)
}

func TestStorageGetBottle(t *testing.T) {
	s := NewBottleStorage(1)
	b := &Bottle{
		Msg: "This is a Test Message",
	}
	err := s.Add(b)
	require.NoError(t, err)

	b, err = s.Get()
	require.NoError(t, err)
	assert.Equal(t, "This is a Test Message", b.Msg)
}
