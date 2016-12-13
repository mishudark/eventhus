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
