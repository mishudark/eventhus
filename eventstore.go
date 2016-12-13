package cqrs

//EventStore saves the events from an aggregate
type EventStore interface {
	Save(events []Event, version int) error
	SafeSave(events []Event, version int) error
	Load(aggregateID string) ([]Event, error)
}
