package cqrs

import (
	"fmt"
	"reflect"
	"sync"
)

//Event stores the data for every event
type Event struct {
	ID            string
	AggregateID   string
	AggregateType string
	Version       int
	Type          string
	Data          interface{}
}

//EventTypeRegister defines the register for all the events that are Data field child of event struct
type EventTypeRegister interface {
	Register(source interface{})
	Get(name string) (interface{}, error)
	Count() int
	Events() []string
}

//EventType implements the EventyTypeRegister interface
type EventType struct {
	sync.RWMutex
	registry map[string]reflect.Type
}

//NewRegister gets a EventyTypeRegister interface
func NewRegister() EventTypeRegister {
	return &EventType{
		registry: make(map[string]reflect.Type),
	}
}

//Register a new type
func (e *EventType) Register(source interface{}) {
	e.Lock()
	defer e.Unlock()

	rawType := reflect.TypeOf(source)
	name := rawType.String()
	//we need to extract only the name without the package
	//name currently follows the format `package.StructName`
	//parts := strings.Split(name, ".")
	//Registry[parts[1]] = source

	e.registry[name] = rawType

}

//Get a type based on its name
func (e *EventType) Get(name string) (interface{}, error) {
	e.RLock()
	defer e.RUnlock()

	rawType, ok := e.registry[name]
	if !ok {
		return nil, fmt.Errorf("can't find %s in registry", name)
	}

	return reflect.New(rawType).Interface(), nil
}

//Count the quantity of events registered
func (e *EventType) Count() int {
	e.RLock()
	count := len(e.registry)
	e.RUnlock()

	return count
}

//Events registered
func (e *EventType) Events() []string {
	var i int
	values := make([]string, len(e.registry))

	e.RLock()
	defer e.RUnlock()

	for key := range e.registry {
		values[i] = key
		i++
	}

	return values
}
