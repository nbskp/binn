package binn

import (
	"errors"
	"time"
)

const (
	defaultStorageSize = 100
	defaultInterval    = 1 * time.Second
)

type queue struct {
	hs []handler
}

func newQueue() *queue {
	return &queue{
		hs: []handler{},
	}
}

func (q *queue) push(h handler) {
	q.hs = append(q.hs, h)
}

func (q *queue) pop() (handler, error) {
	if q.size() == 0 {
		return nil, errors.New("no channels")
	}
	h := q.hs[0]
	q.hs = q.hs[1:]
	return h, nil
}

func (q *queue) size() int {
	return len(q.hs)
}

type Binn struct {
	Storage  Keeper
	Interval time.Duration
	queue    *queue
	closed   chan struct{}
}

func New(storage Keeper, interval time.Duration) *Binn {
	bn := &Binn{
		Storage:  storage,
		Interval: interval,
		queue:    newQueue(),
		closed:   make(chan struct{}),
	}
	bn.Run()
	return bn
}

func Default() *Binn {
	return New(NewBottleStorage(defaultStorageSize), defaultInterval)
}

func (bn *Binn) Publish(b *Bottle) error {
	return bn.Storage.Add(b)
}

type handler func(*Bottle)

func (bn *Binn) Subscribe(fn handler) error {
	bn.queue.push(fn)
	return nil
}

func (bn *Binn) Run() {
	go bn.deliveryLoop()
}

func (bn *Binn) Close() {
	bn.closed <- struct{}{}
}

func (bn *Binn) deliveryLoop() {
Loop:
	for {
		select {
		case <-bn.closed:
			break Loop
		case <-time.After(bn.Interval):
			if bn.queue.size() == 0 {
				break
			}
			b, err := bn.Storage.Get()
			if err != nil {
				break
			}
			fn, err := bn.queue.pop()
			if err != nil {
				bn.Storage.Add(b)
				break
			}
			fn(b)
			bn.queue.push(fn)
		}
	}
}
