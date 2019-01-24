package eventhus

import (
	"fmt"
	"sync"
)

// CommandHandler defines the contract to handle commands
type CommandHandler interface {
	Handle(command Command) error
}

// CommandHandlerRegister stores the handlers for commands
type CommandHandlerRegister interface {
	Add(command interface{}, handler CommandHandler)
	GetHandler(command interface{}) (CommandHandler, error)
}

// CommandRegister contains a registry of command-handler style
type CommandRegister struct {
	mu       sync.RWMutex
	registry map[string]CommandHandler
}

// NewCommandRegister creates a new CommandHandler
func NewCommandRegister() *CommandRegister {
	return &CommandRegister{
		registry: make(map[string]CommandHandler),
	}
}

// Add a new command with its handler
func (c *CommandRegister) Add(command interface{}, handler CommandHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, name := GetTypeName(command)
	c.registry[name] = handler
}

// GetHandler the handler for a command
func (c *CommandRegister) GetHandler(command interface{}) (CommandHandler, error) {
	_, name := GetTypeName(command)

	c.mu.RLock()
	handler, ok := c.registry[name]
	c.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("can't find %s in registry", name)
	}
	return handler, nil
}
