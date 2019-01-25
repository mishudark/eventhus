package basic

import (
	"errors"
	"reflect"

	"github.com/mishudark/eventhus/v2"
)

// ErrInvalidID missing initial event
var ErrInvalidID = errors.New("Invalid ID, initial event missign")

// Handler contains the info to manage commands
type Handler struct {
	repository     *eventhus.Repository
	aggregate      reflect.Type
	bucket, subset string
}

// NewCommandHandler return a handler
func NewCommandHandler(repository *eventhus.Repository, aggregate eventhus.AggregateHandler, bucket, subset string) eventhus.CommandHandler {
	return &Handler{
		repository: repository,
		aggregate:  reflect.TypeOf(aggregate).Elem(),
		bucket:     bucket,
		subset:     subset,
	}
}

// Handle a command
func (h *Handler) Handle(command eventhus.Command) error {
	var err error

	version := command.GetVersion()
	aggregate := reflect.New(h.aggregate).Interface().(eventhus.AggregateHandler)

	if version != 0 {
		if err = h.repository.Load(aggregate, command.GetAggregateID()); err != nil {
			return err
		}
	}

	// the aggregate can have errors trying to replay the previous events
	if aggregate.HasError() {
		return aggregate.GetError()
	}

	if err = aggregate.HandleCommand(command); err != nil {
		return err
	}

	// After to handle the command, the aggregate can have errors applying the new events
	if aggregate.HasError() {
		return aggregate.GetError()
	}

	// if not contain a valid ID,  the initial event (some like createAggreagate event) is missing
	if aggregate.GetID() == "" {
		return ErrInvalidID
	}

	// add the command id for traceability
	aggregate.AttachCommandID(command.GetID())

	// save the changes using the repository
	if err = h.repository.Save(aggregate, version); err != nil {
		return err
	}

	if err = h.repository.PublishEvents(aggregate, h.bucket, h.subset); err != nil {
		return err
	}

	return nil
}
