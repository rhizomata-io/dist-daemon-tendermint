package types

import (
	"bytes"
	"container/list"
	"errors"
	"fmt"
	"sync"
	"time"
)

type EventScope string

const GlobalEventScope = EventScope("")

type EventPath string

func (path EventPath) HasPrefix(prefix EventPath) (ok bool) {
	return bytes.HasPrefix([]byte(path), []byte(prefix))
}

type Event interface {
	Path() EventPath
}

type EventHandler func(event Event)

type EventBus struct {
	sync.Mutex
	scope     EventScope
	listeners map[EventPath]map[string]EventHandler
	queue     *Queue
	processor *CommandProcessor
	started   bool
}

var (
	eventBusMap = make(map[EventScope]*EventBus)
)

func RegisterEventBus(scope EventScope) *EventBus {
	bus := &EventBus{
		scope:     scope,
		listeners: make(map[EventPath]map[string]EventHandler),
		queue:     newQueue(scope),
	}
	
	bus.processor = &CommandProcessor{queue: bus.queue}
	
	eventBusMap[scope] = bus
	
	bus.start()
	
	return bus
}

func (bus *EventBus) start() {
	if !bus.started {
		fmt.Printf(" - EventBus[%s] starting...\n", bus.scope)
		go bus.processor.start()
		bus.started = true
		fmt.Printf(" - EventBus[%s] started.\n", bus.scope)
	}
}

func (bus *EventBus) Subscribe(path EventPath, name string, handler EventHandler) error {
	bus.Lock()
	handlers, ok := bus.listeners[path]
	if !ok {
		handlers = make(map[string]EventHandler)
		bus.listeners[path] = handlers
	}
	
	if _, ok := handlers[name]; ok {
		err := errors.New(fmt.Sprintf("EventHandler[%s] at %s is already registered.", name, path))
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
	delete(handlers, name)
	bus.Unlock()
}

func (bus *EventBus) pushCommand(name string, handler EventHandler, event Event) {
	command := Command{
		name:    name,
		handler: handler,
		event:   event,
	}
	bus.queue.Push(&command)
}

func (bus *EventBus) Publish(event Event) {
	bus.Lock()
	
	eventPath := event.Path()
	
	// fmt.Println("- EventBus Publish ", bus.scope, eventPath, len(bus.listeners))
	for path, handlers := range bus.listeners {
		if eventPath.HasPrefix(path) {
			// fmt.Println("   # EventBus ", bus.scope , eventPath, "=", path )
			for name, handler := range handlers {
				// fmt.Println("     - handler =", name)
				bus.pushCommand(name, handler, event)
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
		return errors.New(fmt.Sprintf("Unknown Event Scope %s", scope))
	}
}

func Publish(scope EventScope, event Event) error {
	bus, ok := eventBusMap[scope]
	if ok {
		bus.Publish(event)
		return nil
	} else {
		return errors.New(fmt.Sprintf("Unknown Event Scope %s", scope))
	}
}

// Queue ...
type Queue struct {
	sync.Mutex
	scope     EventScope
	innerList *list.List
	lock      chan bool
	waiting   int64
}

type Command struct {
	name    string
	handler EventHandler
	event   Event
}

// NewQueue ...
func newQueue(scope EventScope) *Queue {
	queue := Queue{scope: scope, innerList: list.New(), lock: make(chan bool)}
	return &queue
}

// Size ..
func (queue *Queue) Size() int {
	return queue.innerList.Len()
}

// Push ..
func (queue *Queue) Push(value *Command) {
	queue.Lock()
	defer queue.Unlock()
	queue.innerList.PushBack(value)
	
	// fmt.Println("Queue Push A ", queue.waiting, queue.scope)
	// if queue.waiting > 0 {
	// 	fmt.Println("Queue Push lock false ", queue.waiting, queue.scope)
	// 	queue.lock <- false
	// }
	// fmt.Println("Queue Push B ", queue.waiting, queue.scope)
}

// Pop ..
func (queue *Queue) Pop() (value *Command) {
	value = queue._pop()
	// fmt.Println("Queue _pop", queue.scope)
	// for ; value == nil; value = queue._pop() {
	// 	queue.waiting++
	// 	fmt.Println("Queue _pop wait...", queue.waiting, queue.scope)
	// 	<-queue.lock
	// 	queue.waiting--
	// }
	
	return value
}

// Pop ..
func (queue *Queue) _pop() (value *Command) {
	queue.Lock()
	defer queue.Unlock()
	
	el := queue.innerList.Front()
	if el != nil {
		value = el.Value.(*Command)
		queue.innerList.Remove(el)
	}
	return value
}

type CommandProcessor struct {
	queue   *Queue
	running bool
}

func (proc *CommandProcessor) start() {
	go func() {
		proc.running = true
		proc.process()
	}()
}

func (proc *CommandProcessor) process() {
	for proc.running {
		command := proc.queue.Pop()
		if command != nil {
			command.handler(command.event)
		} else {
			time.Sleep(100 * time.Millisecond)
		}
		
	}
}
