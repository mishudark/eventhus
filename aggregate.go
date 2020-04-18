package eventhus

// BaseAggregate contains the basic info
// that all aggregates should have
type BaseAggregate struct {
	ID      string
	Type    string
	Version int
	Changes []Event
}

// AggregateHandler defines the methods to process commands
type AggregateHandler interface {
	// LoadsFromHistory(events []Event)
	ApplyChange(event Event)
	ApplyChangeHelper(aggregate AggregateHandler, event Event, commit bool)
	HandleCommand(Command) error
	Uncommitted() []Event
	ClearUncommitted()
	IncrementVersion()
	GetID() string
}

// Uncommitted return the events to be saved
func (b *BaseAggregate) Uncommitted() []Event {
	return b.Changes
}

// ClearUncommitted the events
func (b *BaseAggregate) ClearUncommitted() {
	b.Changes = []Event{}
}

// IncrementVersion ads 1 to the current version
func (b *BaseAggregate) IncrementVersion() {
	b.Version++
}

// ApplyChangeHelper increments the version of an aggregate and apply the change itself
func (b *BaseAggregate) ApplyChangeHelper(aggregate AggregateHandler, event Event, commit bool) {
	// increments the version in event and aggregate
	b.IncrementVersion()

	// apply the event itself
	aggregate.ApplyChange(event)
	if commit {
		event.Version = b.Version
		_, event.Type = GetTypeName(event.Data)
		b.Changes = append(b.Changes, event)
	}
}

// GetID of the current aggregate
func (b *BaseAggregate) GetID() string {
	return b.ID
}
