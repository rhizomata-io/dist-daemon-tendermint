package types

import (
	"bytes"
	"errors"
	"fmt"
	"sync"
)


type EventScope string

const GlobalEventScope  = EventScope("")

type EventPath string

func (path EventPath)HasPrefix (prefix EventPath) (ok bool){
	return bytes.HasPrefix([]byte(path), []byte(prefix))
}

type Event interface {
	Path() EventPath
}

type EventHandler func(event Event)

type EventBus struct {
	sync.Mutex
	scope      EventScope
	listeners map[EventPath]map[string]EventHandler
}

var (
	eventBusMap = make(map[EventScope]*EventBus)
)

func RegisterEventBus(scope EventScope) *EventBus {
	bus := &EventBus{scope: scope, listeners: make(map[EventPath]map[string]EventHandler)}
	eventBusMap[scope] = bus
	return bus
}

func (bus *EventBus) Subscribe(path EventPath, name string, handler EventHandler) error {
	bus.Lock()
	handlers, ok := bus.listeners[path]
	if !ok {
		handlers = make(map[string]EventHandler)
		bus.listeners[path] = handlers
	}
	
	if _,ok := handlers[name]; ok {
		err := errors.New(fmt.Sprintf("EventHandler[%s] at %s is already registered.", name,path))
		return err
	}
	handlers[name] = handler
	bus.Unlock()
	return nil
}

func (bus *EventBus) Unsubscribe(path EventPath, name string) {
	bus.Lock()
	handlers, ok := bus.listeners[path]
	if !ok {
		handlers = make(map[string]EventHandler)
		bus.listeners[path] = handlers
	}
	delete(handlers,name)
	bus.Unlock()
}

func (bus *EventBus) Publish(event Event) {
	bus.Lock()
	
	eventPath := event.Path()
	
	for path, handlers := range bus.listeners {
		// fmt.Println(" - EventBus Publish ", eventPath, "=", path, len(handlers))
		if eventPath.HasPrefix(path) {
			// fmt.Println("     - EventBus Publish * Match ", eventPath, "=", path )
			for _, handler := range handlers {
				go func() { handler(event) }()
			}
		}
	}
	
	bus.Unlock()
}

func Subscribe(scope EventScope, path EventPath, name string, handler EventHandler) error {
	bus, ok := eventBusMap[scope]
	if ok {
		return bus.Subscribe(path, name, handler)
	} else {
		return errors.New(fmt.Sprintf("Unknown Event Scope %s",scope))
	}
}

func Publish(scope EventScope, event Event) error {
	bus, ok := eventBusMap[scope]
	if ok {
		bus.Publish(event)
		return nil
	} else {
		return errors.New(fmt.Sprintf("Unknown Event Scope %s",scope))
	}
}
