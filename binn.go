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

func (q *queue) pop() (chan *Bottle, error) {
	if q.size() == 0 {
		return nil, errors.New("no channels")
	}
	ch := q.chs[0]
	q.chs = q.chs[1:]
	return ch, nil
}

func (q *queue) size() int {
	return len(q.chs)
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

func (bn *Binn) Subscribe(ch chan *Bottle) error {
	if cap(ch) == 0 {
		return errors.New("channel capacity must be more than 0")
	}
	bn.queue.push(ch)
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
			ch, err := bn.queue.pop()
			if err != nil {
				bn.Storage.Add(b)
				break
			}
			// once, if a channel length is more than 0(parent doesn't receive a bottle ex. process reads channel is killed),
			// binn re-adds bottles in the channel and close the channel to notify parent of unsubscribing
			if len(ch) > 0 {
				bn.Storage.Add(b)
				bn.readd(ch)
				close(ch)
				break
			}
			ch <- b
			bn.queue.push(ch)
		}
	}
}

func (bn *Binn) readd(ch chan *Bottle) {
	for {
		select {
		case b := <-ch:
			bn.Storage.Add(b)
		default:
			if len(ch) == 0 {
				return
			}
		}
	}
}

type Receiver struct {
	ch   chan *Bottle
	once bool
}

func NewReceiver() *Receiver {
	return &Receiver{
		ch:   make(chan *Bottle, 1),
		once: false,
	}
}

func (r *Receiver) Receive(bn *Binn) chan *Bottle {
	if !r.once {
		bn.Subscribe(r.ch)
		r.once = true
		return r.ch
	}
	select {
	case b, ok := <-r.ch:
		if !ok {
			r.ch = make(chan *Bottle, 1)
			bn.Subscribe(r.ch)
			return r.ch
		}
		if b != nil {
			r.ch <- b
		}
		return r.ch
	default:
		return r.ch
	}
}
