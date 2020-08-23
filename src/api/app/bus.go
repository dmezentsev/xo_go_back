package app

import (
	"sync"
	"time"
)

type EventType string
type IEvent interface {}

type ICancelling interface {
	Cancel()
	IsCancelled() bool
}

type cancelling struct {
	Cancelled bool
}

func (c *cancelling) Cancel() {
	c.Cancelled = true
}

func (c *cancelling) IsCancelled() bool {
	return c.Cancelled
}

type Emitter struct {
	eventType EventType
	ch chan IEvent
	cancelling
	Latency time.Duration
}

type CallbackMetaType interface{}
type CallbackFnType func(e IEvent, m CallbackMetaType)

type Callback struct {
	Callback CallbackFnType
	Meta     CallbackMetaType
	cancelling
}

type Bus struct {
	emitters []Emitter
	subscribers map[EventType][]Callback
	emitterMux sync.RWMutex
	subscribersMux sync.RWMutex
	cancelling
	Latency time.Duration
}

const defaultLatency = 50 * time.Millisecond

func NewBus() *Bus {
	return &Bus{
		emitters: make([]Emitter, 0),
		subscribers: make(map[EventType][]Callback),
		Latency: defaultLatency,
	}
}

func (b *Bus) NewEmitter(et EventType) Emitter {
	e := Emitter{
		eventType: et,
		ch:        make(chan IEvent),
	}
	b.RegisterEmitter(e)
	return e
}

func (b *Bus) RegisterEmitter(e Emitter) {
	b.emitterMux.Lock()
	defer b.emitterMux.Unlock()
	b.emitters = append(b.emitters, e)
	go e.Serve(b)
}

func (b *Bus) NewCallback(et EventType, fn CallbackFnType, meta CallbackMetaType) Callback {
	cb := Callback{
		Callback: fn,
		Meta:     meta,
	}
	b.RegisterCallback(et, cb)
	return cb
}

func (b *Bus) RegisterCallback(et EventType, cb Callback) {
	b.subscribersMux.Lock()
	defer b.subscribersMux.Unlock()
	_, ok := b.subscribers[et]
	if !ok {
		b.subscribers[et] = make([]Callback, 0)
	}
	b.subscribers[et] = append(b.subscribers[et], cb)
}

func (b *Bus) Cancel() {
	b.Cancelled = true
	for _, e := range b.emitters {
		e.Cancel()
	}
	for _, cbs := range b.subscribers {
		for _, cb := range cbs {
			cb.Cancel()
		}
	}
}

func (e *Emitter) Serve(b *Bus) {
	ticker := time.NewTicker(b.Latency)
	for {
		select {
		case event := <-e.ch:
			cbs, ok := b.subscribers[e.eventType]
			if !ok {
				continue
			}
			for _, cb := range cbs {
				if !cb.IsCancelled() {
					cb.Callback(event, cb.Meta)
				}
			}
		case <-ticker.C:
			if e.Cancelled {
				return
			}
		}
	}
}
