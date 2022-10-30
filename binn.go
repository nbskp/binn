package binn

import (
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
	chs     []chan *Bottle
}

type Binn struct {
	Storage  Keeper
	Interval time.Duration
	chs      []chan *Bottle
	cnt      uint64
}

func New(storage Keeper, interval time.Duration) *Binn {
	bn := &Binn{
		Storage:  storage,
		Interval: interval,
		chs:      []chan *Bottle{},
		cnt:      0,
	}
	bn.Run()
	return bn
}

func Default() *Binn {
	return New(NewBottleStorage(100), 1*time.Second)
}

func (bn *Binn) Add(b *Bottle) error {
	return bn.Storage.Add(b)
}

func (bn *Binn) Subscribe(ch chan *Bottle) {
	bn.chs = append(bn.chs, ch)
}

func (bn *Binn) Run() {
	go bn.publishLoop()
}

func (bn *Binn) publishLoop() {
	for {
		select {
		case <-time.After(bn.Interval):
			if len(bn.chs) == 0 {
				break
			}
			idx := bn.cnt % uint64(len(bn.chs))
			ch := bn.chs[idx]
			b, err := bn.Storage.Get()
			if err != nil {
				break
			}
			ch <- b
			bn.cnt++
		}
	}
}
