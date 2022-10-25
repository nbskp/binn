package binn

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEngine(t *testing.T) {
	cfg := &Config{
		DeliveryInterval: 1 * time.Millisecond,
	}
	storage := NewBottleStorage(1)
	e := NewEngine(cfg, storage)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	e.Run(ctx)
	g := e.NewGateway()
	g.Send(&Bottle{
		Id:  "1c7a8201-cdf7-11ec-a9b3-0242ac110004",
		Msg: "This is a Test Message",
	})
	var b *Bottle
	select {
	case b = <-g.Receive():
	case <-time.After(1 * time.Second):
		require.Fail(t, "timeout")
	}
	assert.Equal(t, "1c7a8201-cdf7-11ec-a9b3-0242ac110004", b.Id)
	assert.Equal(t, "This is a Test Message", b.Msg)
}
