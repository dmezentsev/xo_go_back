package app

import (
	"sync"
	"time"
)

type EventType string

type IEvent interface {
	SetType(et EventType)
	GetPayload() interface{}
	GetType() EventType
}

type Event struct {
	Type EventType `json:"type"`
	Payload interface{} `json:"payload"`
}

func (e Event) SetType(et EventType) {
	e.Type = et
}

func (e Event) GetPayload() interface{} {
	return e.Payload
}

func (e Event) GetType() EventType {
	return e.Type
}

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
	Emitter   chan IEvent
	cancelling
	Latency time.Duration
}

type CallbackMetaType interface{}
type CallbackFnType func(CallbackArgs)

type Callback struct {
	Callback CallbackFnType
	Meta     CallbackMetaType
	cancelling
}

type CallbackArgs struct {
	Initiator interface{}
	Event IEvent
	Meta CallbackMetaType
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

func (b *Bus) NewEmitter(et EventType, initiator interface{}) Emitter {
	e := Emitter{
		eventType: et,
		Emitter:   make(chan IEvent),
	}
	b.RegisterEmitter(e, initiator)
	return e
}

func (b *Bus) RegisterEmitter(e Emitter, initiator interface{}) {
	b.emitterMux.Lock()
	defer b.emitterMux.Unlock()
	b.emitters = append(b.emitters, e)
	go e.Serve(initiator, b)
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

func (e *Emitter) Serve(initiator interface{}, b *Bus) {
	ticker := time.NewTicker(b.Latency)
	for {
		select {
		case event := <-e.Emitter:
			event.SetType(e.eventType)
			cbs, ok := b.subscribers[e.eventType]
			if !ok {
				continue
			}
			for _, cb := range cbs {
				if !cb.IsCancelled() {
					go cb.Callback(CallbackArgs{initiator, event, cb.Meta})
				}
			}
		case <-ticker.C:
			if e.Cancelled {
				return
			}
		}
	}
}
