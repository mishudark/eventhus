package eventhus

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"sync"
)

var registry = make(map[string]reflect.Type)

// Event stores the data for every event
type Event struct {
	ID            string      `json:"id"`
	AggregateID   string      `json:"aggregate_id"`
	AggregateType string      `json:"aggregate_type"`
	Version       int         `json:"version"`
	Type          string      `json:"type"`
	Data          interface{} `json:"data"`
}

// GetTypeName of given struct
func GetTypeName(source interface{}) (reflect.Type, string) {
	rawType := reflect.TypeOf(source)

	// source is a pointer, convert to its value
	if rawType.Kind() == reflect.Ptr {
		rawType = rawType.Elem()
	}

	name := rawType.String()
	// we need to extract only the name without the package
	// name currently follows the format `package.StructName`
	parts := strings.Split(name, ".")
	return rawType, snakeCase(parts[1])
}

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func snakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// Register defines generic methods to create a registry
type Register interface {
	Set(source interface{})
	Get(name string) (interface{}, error)
	Count() int
}

// EventTypeRegister defines the register for all the events that are Data field child of event struct
type EventTypeRegister interface {
	Register
	Events() []string
}

// EventType implements the EventyTypeRegister interface
type EventType struct {
	mu sync.RWMutex
}

// NewEventRegister gets a EventyTypeRegister interface
func NewEventRegister() EventTypeRegister {
	return &EventType{}
}

// Set a new type
func (e *EventType) Set(source interface{}) {
	rawType, name := GetTypeName(source)

	e.mu.Lock()
	registry[name] = rawType
	e.mu.Unlock()
}

// Get a type based on its name
func (e *EventType) Get(name string) (interface{}, error) {
	e.mu.RLock()
	rawType, ok := registry[name]
	e.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("can't find %s in registry", name)
	}

	return reflect.New(rawType).Interface(), nil
}

// Count the quantity of events registered
func (e *EventType) Count() int {
	e.mu.RLock()
	count := len(registry)
	e.mu.RUnlock()

	return count
}

// Events registered
func (e *EventType) Events() []string {
	var i int
	values := make([]string, len(registry))

	e.mu.RLock()
	for key := range registry {
		values[i] = key
		i++
	}

	e.mu.RUnlock()
	return values
}
