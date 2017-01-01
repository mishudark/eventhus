package eventhus

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
	IncrementVersion()
}

//Uncommited return the events to be saved
func (b *BaseAggregate) Uncommited() []Event {
	return b.Changes
}

//ClearUncommited the events
func (b *BaseAggregate) ClearUncommited() {
	b.Changes = []Event{}
}

//IncrementVersion ads 1 to the current version
func (b *BaseAggregate) IncrementVersion() {
	b.Version++
}

//ApplyChange increments the version of an aggregate and apply the change itself
func (b *BaseAggregate) ApplyChange(aggregate AggregateHandler, event Event, commit bool) {
	//increments the version in event and aggregate
	b.IncrementVersion()

	//apply the event itself
	aggregate.ApplyChange(event, commit)
	if commit {
		event.Version = b.Version
		b.Changes = append(b.Changes, event)
	}
}
