package binn

import (
	"context"
	"time"
)

const (
	defaultEngineStorageSize = 100
)

var (
	defaultDeliveryInterval = 5 * time.Second
)

type Engine struct {
	Cfg     *Config
	Storage Keeper
	inCh    chan *Bottle
	outCh   chan *Bottle
}

func NewEngine(cfg *Config, storage *BottleStorage) *Engine {
	return &Engine{
		Cfg:     cfg,
		Storage: storage,
		inCh:    make(chan *Bottle),
		outCh:   make(chan *Bottle),
	}
}

func DefaultEngine() *Engine {
	cfg := &Config{DeliveryInterval: defaultDeliveryInterval}
	storage := NewBottleStorage(defaultEngineStorageSize)
	return NewEngine(cfg, storage)
}

type Gateway struct {
	inCh  chan<- *Bottle
	outCh <-chan *Bottle
}

func (g *Gateway) Send(b *Bottle) {
	g.inCh <- b
}

func (g *Gateway) Receive() <-chan *Bottle {
	return g.outCh
}

func (e *Engine) NewGateway() *Gateway {
	return &Gateway{
		inCh:  e.inCh,
		outCh: e.outCh,
	}
}

func (e *Engine) Run(ctx context.Context) {
	go func() {
	Loop:
		for {
			select {
			case <-ctx.Done():
				break Loop
			case b := <-e.inCh:
				if err := e.Storage.Add(b); err != nil {
					break
				}
			default:
				break
			}
		}
	}()

	go func() {
		t := time.NewTicker(e.Cfg.DeliveryInterval)
		defer t.Stop()

	Loop:
		for {
			select {
			case <-ctx.Done():
				break Loop
			case <-t.C:
				b, err := e.Storage.Get()
				if err != nil {
					break
				}
				e.outCh <- b
			default:
				break
			}
		}
	}()
}
