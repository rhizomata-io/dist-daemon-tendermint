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
	name      string
	listeners map[EventType]map[string]EventHandler
}

func NewEventBus(name string) *EventBus {
	return &EventBus{name: name, listeners: make(map[EventType]map[string]EventHandler)}
}

func (bus *EventBus) Subscribe(typ EventType, name string, handler EventHandler) {
	bus.Lock()
	handlers, ok := bus.listeners[typ]
	if !ok {
		handlers = make(map[string]EventHandler)
		bus.listeners[typ] = handlers
	}
	handlers[name] = handler
	bus.Unlock()
}

func (bus *EventBus) Unsubscribe(typ EventType, name string, handler EventHandler) {
	bus.Lock()
	handlers, ok := bus.listeners[typ]
	if !ok {
		handlers = make(map[string]EventHandler)
		bus.listeners[typ] = handlers
	}
	handlers[name] = handler
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
	globalEventBus = NewEventBus("global")
)

func Subscribe(typ EventType, name string, handler EventHandler) {
	globalEventBus.Subscribe(typ, name, handler)
}

func Publish(event Event) {
	globalEventBus.Publish(event)
}
