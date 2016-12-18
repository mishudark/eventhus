package cqrs

import (
	"fmt"
	"reflect"
	"sync"
)

//CommandHandle defines the contract to handle commands
type CommandHandle interface {
	Handle(command interface{}) error
}

//CommandHandlerRegister stores the handlers for commands
type CommandHandlerRegister interface {
	Add(command interface{}, handler CommandHandle)
	Get(command interface{}) (CommandHandle, error)
	Handlers() []string
}

//CommandHandler contains a registry of command-handler style
type CommandHandler struct {
	sync.RWMutex
	registry map[string]CommandHandle
	//repository Repository
}

//NewCommandHandler creates a new CommandHandler
func NewCommandHandler(repository Repository) *CommandHandler {
	return &CommandHandler{
		registry: make(map[string]CommandHandle),
		//repository: repository,
	}
}

//Add a new command with its handler
func (c *CommandHandler) Add(command interface{}, handler CommandHandle) {
	c.Lock()
	defer c.Unlock()

	rawType := reflect.TypeOf(command)
	name := rawType.String()
	c.registry[name] = handler
}

//Get the handler for a command
func (c *CommandHandler) Get(command interface{}) (CommandHandle, error) {
	rawType := reflect.TypeOf(command)
	name := rawType.String()

	handler, ok := c.registry[name]
	if !ok {
		return nil, fmt.Errorf("can't find %s in registry", name)
	}
	return handler, nil
}
