package eventhus

// Repository is responsible to generate an Aggregate
// save events and publish it
type Repository struct {
	eventStore EventStore
	eventBus   EventBus
}

// NewRepository creates a repository wieh a eventstore and eventbus access
func NewRepository(store EventStore, bus EventBus) *Repository {
	return &Repository{
		eventStore: store,
		eventBus:   bus,
	}
}

// Load restore the last state of an aggregate
func (r *Repository) Load(aggregate AggregateHandler, ID string) error {
	events, err := r.eventStore.Load(ID)

	if err != nil {
		return err
	}

	for _, event := range events {
		ReduceHelper(aggregate, event, false)
	}
	return nil
}

// Save the events and publish it to eventbus
func (r *Repository) Save(aggregate AggregateHandler, version int) error {
	return r.eventStore.Save(aggregate.Uncommited(), version)
}

// PublishEvents to an eventBus
func (r *Repository) PublishEvents(aggregate AggregateHandler, bucket, subset string) error {
	var err error

	for _, event := range aggregate.Uncommited() {
		if err = r.eventBus.Publish(event, bucket, subset); err != nil {
			return err
		}
	}

	return nil
}

// PublishError to an eventBus
func (r *Repository) PublishError(err error, command Command, bucket, subset string) error {
	event := Event{
		ID:            GenerateUUID(),
		AggregateID:   command.GetAggregateID(),
		AggregateType: command.GetAggregateType(),
		CommandID:     command.GetID(),
		Version:       command.GetVersion(),
		Type:          "failure",
	}

	if failure, ok := err.(Failure); ok {
		event.Data = failure
	} else {
		event.Data = struct {
			Error error `json:"error"`
		}{
			Error: err,
		}
	}

	return r.eventBus.Publish(event, bucket, "errors")
}

// SafeSave the events without check the version
func (r *Repository) SafeSave(aggregate AggregateHandler, version int) error {
	return r.eventStore.SafeSave(aggregate.Uncommited(), version)
}
