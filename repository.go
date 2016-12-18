package cqrs

//Repository is responsible to generate an Aggregate
//save events and publish it
type Repository struct {
	eventStore EventStore
	eventBus   interface{}
}

//NewRepository creates a repository wieh a eventstore and eventbus access
func NewRepository(store EventStore, bus interface{}) *Repository {
	return &Repository{
		store,
		bus,
	}
}

//Load restore the last state of an aggregate
func (r *Repository) Load(aggregate AggregateHandler, ID string) error {
	events, err := r.eventStore.Load(ID)

	if err != nil {
		return err
	}

	aggregate.LoadsFromHistory(events)
	return nil
}

//Save the events and publish it to eventbus
func (r *Repository) Save(aggregate AggregateHandler, version int) error {
	return r.eventStore.Save(aggregate.Uncommited(), version)
}

//SafeSave the events without check the version
func (r *Repository) SafeSave(aggregate AggregateHandler, version int) error {
	return r.eventStore.SafeSave(aggregate.Uncommited(), version)
}
