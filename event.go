package eventhus

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

var registry = make(map[string]reflect.Type)

//Event stores the data for every event
type Event struct {
	ID            string
	AggregateID   string
	AggregateType string
	Version       int
	Type          string
	Data          interface{}
}

//Register defines generic methods to create a registry
type Register interface {
	Set(source interface{})
	Get(name string) (interface{}, error)
	Count() int
}

//EventTypeRegister defines the register for all the events that are Data field child of event struct
type EventTypeRegister interface {
	Register
	Events() []string
}

//EventType implements the EventyTypeRegister interface
type EventType struct {
	sync.RWMutex
}

//NewEventRegister gets a EventyTypeRegister interface
func NewEventRegister() EventTypeRegister {
	return &EventType{}
}

//Set a new type
func (e *EventType) Set(source interface{}) {
	rawType, name := GetTypeName(source)

	e.Lock()
	registry[name] = rawType
	e.Unlock()
}

//Get a type based on its name
func (e *EventType) Get(name string) (interface{}, error) {
	e.RLock()
	defer e.RUnlock()

	rawType, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("can't find %s in registry", name)
	}

	return reflect.New(rawType).Interface(), nil
}

//Count the quantity of events registered
func (e *EventType) Count() int {
	e.RLock()
	count := len(registry)
	e.RUnlock()

	return count
}

//Events registered
func (e *EventType) Events() []string {
	var i int
	values := make([]string, len(registry))

	e.RLock()
	defer e.RUnlock()

	for key := range registry {
		values[i] = key
		i++
	}

	return values
}

//GetTypeName of given struct
func GetTypeName(source interface{}) (reflect.Type, string) {
	rawType := reflect.TypeOf(source)
	name := rawType.String()
	//we need to extract only the name without the package
	//name currently follows the format `package.StructName`
	parts := strings.Split(name, ".")
	return rawType, parts[1]
}
