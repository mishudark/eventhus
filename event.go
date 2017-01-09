package eventhus

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

var (
	mu       sync.RWMutex
	registry = make(map[string]reflect.Type)
)

//Event stores the data for every event
type Event struct {
	ID            string      `json:"id"`
	AggregateID   string      `json:"aggregate_id"`
	AggregateType string      `json:"aggregate_type"`
	Version       int         `json:"version"`
	Type          string      `json:"type"`
	Data          interface{} `json:"data"`
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

	mu.Lock()
	registry[name] = rawType
	mu.Unlock()
}

//Get a type based on its name
func (e *EventType) Get(name string) (interface{}, error) {
	mu.RLock()
	rawType, ok := registry[name]
	mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("can't find %s in registry", name)
	}

	return reflect.New(rawType).Interface(), nil
}

//Count the quantity of events registered
func (e *EventType) Count() int {
	mu.RLock()
	count := len(registry)
	mu.RUnlock()

	return count
}

//Events registered
func (e *EventType) Events() []string {
	var i int
	values := make([]string, len(registry))

	mu.RLock()
	for key := range registry {
		values[i] = key
		i++
	}

	mu.RUnlock()

	return values
}

//GetTypeName of given struct
func GetTypeName(source interface{}) (reflect.Type, string) {
	rawType := reflect.TypeOf(source)

	//source is a pointer, convert to its value
	if rawType.Kind() == reflect.Ptr {
		rawType = rawType.Elem()
	}

	name := rawType.String()
	//we need to extract only the name without the package
	//name currently follows the format `package.StructName`
	parts := strings.Split(name, ".")
	return rawType, parts[1]
}
