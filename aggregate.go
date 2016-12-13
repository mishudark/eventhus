package cqrs

//BaseAggregate contains the basic info
//that all aggregates should have
type BaseAggregate struct {
	ID      string
	Type    string
	Version int
	Changes []Event
}

//AggregateHandler defines the methods to process commands
type AggregateHandler interface {
	LoadsFromHistory(events []Event)
	ApplyChange(event Event, commit bool)
	Uncommited() []Event
	ClearUncommited()
}

//Uncommited return the events to be saved
func (b *BaseAggregate) Uncommited() []Event {
	return b.Changes
}

//ClearUncommited the events
func (b *BaseAggregate) ClearUncommited() {
	b.Changes = []Event{}
}
