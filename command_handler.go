package eventhus

import (
	"fmt"
	"reflect"
	"sync"
)

// CommandHandle defines the contract to handle commands
type CommandHandle interface {
	Handle(command Command) error
}

// CommandHandlerRegister stores the handlers for commands
type CommandHandlerRegister interface {
	Add(command interface{}, handler CommandHandle)
	Get(command interface{}) (CommandHandle, error)
	// Handlers() []string
}

// CommandRegister contains a registry of command-handler style
type CommandRegister struct {
	sync.RWMutex
	registry map[string]CommandHandle
	// repository *Repository
}

// NewCommandRegister creates a new CommandHandler
func NewCommandRegister() *CommandRegister {
	return &CommandRegister{
		registry: make(map[string]CommandHandle),
		// repository: repository,
	}
}

// Add a new command with its handler
func (c *CommandRegister) Add(command interface{}, handler CommandHandle) {
	c.Lock()
	defer c.Unlock()

	rawType := reflect.TypeOf(command)
	name := rawType.String()
	c.registry[name] = handler
}

// Get the handler for a command
func (c *CommandRegister) Get(command interface{}) (CommandHandle, error) {
	rawType := reflect.TypeOf(command)
	name := rawType.String()

	handler, ok := c.registry[name]
	if !ok {
		return nil, fmt.Errorf("can't find %s in registry", name)
	}
	return handler, nil
}
