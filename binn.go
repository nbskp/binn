package binn

import (
	"time"
)

const (
	defaultStorageSize = 100
	defaultInterval    = 1 * time.Second
)

type Binn struct {
	Storage  Keeper
	Interval time.Duration
	queue    []chan *Bottle
}

func New(storage Keeper, interval time.Duration) *Binn {
	bn := &Binn{
		Storage:  storage,
		Interval: interval,
		queue:    []chan *Bottle{},
	}
	bn.Run()
	return bn
}

func Default() *Binn {
	return New(NewBottleStorage(defaultStorageSize), defaultInterval)
}

func (bn *Binn) Add(b *Bottle) error {
	return bn.Storage.Add(b)
}

func (bn *Binn) Get() <-chan *Bottle {
	ch := make(chan *Bottle)
	bn.queue = append(bn.queue, ch)
	return ch
}

func (bn *Binn) Run() {
	go bn.publishLoop()
}

func (bn *Binn) publishLoop() {
	for {
		select {
		case <-time.After(bn.Interval):
			if len(bn.queue) == 0 {
				break
			}
			ch := bn.queue[0]
			bn.queue = bn.queue[1:]
			b, err := bn.Storage.Get()
			if err != nil {
				break
			}
			ch <- b
		}
	}
}
