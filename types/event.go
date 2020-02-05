package types

import (
	"sync"
)

type EventType string

type Event interface {
	Type() EventType
}

type EventHandler func(event Event)

type EventBus struct {
	sync.Mutex
	listeners map[EventType][]EventHandler
}

func NewEventBus() *EventBus {
	return &EventBus{listeners: make(map[EventType][]EventHandler)}
}

func (bus *EventBus) Subscribe(typ EventType, handler EventHandler) {
	bus.Lock()
	handlers, ok := bus.listeners[typ]
	if !ok {
		handlers = []EventHandler{}
	}
	handlers = append(handlers, handler)
	bus.Unlock()
}

func (bus *EventBus) Publish(event Event) {
	bus.Lock()
	handlers, ok := bus.listeners[event.Type()]
	if ok {
		for _, handler := range handlers {
			go func() { handler(event) }()
		}
	}
	bus.Unlock()
}


var (
	globalEventBus = NewEventBus()
)


func Subscribe(typ EventType, handler EventHandler) {
	globalEventBus.Subscribe(typ,handler )
}

func Publish(event Event) {
	globalEventBus.Publish(event)
}
