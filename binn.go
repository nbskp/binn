package binn

import (
	"time"
)

const (
	defaultStorageSize = 100
	defaultInterval    = 1 * time.Second
)

type queue struct {
	chs []chan *Bottle
}

func newQueue() *queue {
	return &queue{
		chs: []chan *Bottle{},
	}
}

func (q *queue) push(ch chan *Bottle) {
	q.chs = append(q.chs, ch)
}

func (q *queue) pop() chan *Bottle {
	ch := q.chs[0]
	q.chs = q.chs[1:]
	return ch
}

func (q *queue) size() int {
	return len(q.chs)
}

type Binn struct {
	Storage  Keeper
	Interval time.Duration
	queue    *queue
}

func New(storage Keeper, interval time.Duration) *Binn {
	bn := &Binn{
		Storage:  storage,
		Interval: interval,
		queue:    newQueue(),
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
	bn.queue.push(ch)
	return ch
}

func (bn *Binn) Run() {
	go bn.publishLoop()
}

func (bn *Binn) publishLoop() {
	for {
		select {
		case <-time.After(bn.Interval):
			if bn.queue.size() == 0 {
				break
			}
			b, err := bn.Storage.Get()
			if err != nil {
				break
			}
			ch := bn.queue.pop()
			ch <- b
		}
	}
}
