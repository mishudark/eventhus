package bank

import (
	"errors"
	"reflect"

	"github.com/mishudark/eventhus"
)

var ErrInvalidID = errors.New("Invalid ID, initial event missign")

type CommandHandler struct {
	repository     *eventhus.Repository
	aggregate      reflect.Type
	bucket, subset string
}

func NewCommandHandler(repository *eventhus.Repository, aggregate eventhus.AggregateHandler, bucket, subset string) *CommandHandler {
	return &CommandHandler{
		repository: repository,
		aggregate:  reflect.TypeOf(aggregate),
		bucket:     bucket,
		subset:     subset,
	}
}

func (c *CommandHandler) Handle(command eventhus.Command) error {
	var err error

	version := command.GetVersion()
	aggregate := reflect.New(c.aggregate).Interface().(eventhus.AggregateHandler)

	if version != 0 {
		if err = c.repository.Load(aggregate, command.GetAggregateID()); err != nil {
			return err
		}
	}

	if err = aggregate.HandleCommand(command); err != nil {
		return err
	}

	//if not contain a valid ID,  the initial event (some like createAggreagate event) is missing
	if aggregate.GetID() == "" {
		return ErrInvalidID
	}

	if err = c.repository.Save(aggregate, version); err != nil {
		return err
	}

	if err = c.repository.PublishEvents(aggregate, c.bucket, c.subset); err != nil {
		return err
	}

	return nil
}
