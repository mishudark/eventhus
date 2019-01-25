package eventhus

// BaseAggregate contains the basic info
// that all aggregates should have
type BaseAggregate struct {
	ID      string
	Type    string
	Version int
	Changes []Event
	Error   error
}

// AggregateHandler defines the methods to process commands
type AggregateHandler interface {
	// LoadsFromHistory(events []Event)
	Reduce(event Event) error
	//Dispatch(aggregate AggregateHandler, event Event)
	//ReduceHelper(aggregate AggregateHandler, event Event, commit bool)
	HandleCommand(Command) error
	AddEvent(Event)
	AttachCommandID(id string)
	Uncommited() []Event
	ClearUncommited()
	IncrementVersion()
	GetID() string
	GetVersion() int
	AddError(error)
	GetError() error
	HasError() bool
}

// Uncommited return the events to be saved
func (b *BaseAggregate) Uncommited() []Event {
	return b.Changes
}

// ClearUncommited the events
func (b *BaseAggregate) ClearUncommited() {
	b.Changes = []Event{}
}

// IncrementVersion ads 1 to the current version
func (b *BaseAggregate) IncrementVersion() {
	b.Version++
}

// Dispatch process the event and commit it
func Dispatch(aggregate AggregateHandler, event Event) {
	ReduceHelper(aggregate, event, true)
}

// ReduceHelper increments the version of an aggregate and apply the change itself
func ReduceHelper(aggregate AggregateHandler, event Event, commit bool) {
	// if aggregate already contains an error, nop the operation
	if aggregate.HasError() {
		return
	}

	// increments the version in event and aggregate
	aggregate.IncrementVersion()

	// apply the event itself
	if err := aggregate.Reduce(event); err != nil {
		// if there is  an error, add it to the errors
		aggregate.AddError(err)
	}

	if commit {
		event.Version = aggregate.GetVersion()
		_, event.Type = GetTypeName(event.Data)
		aggregate.AddEvent(event)
	}
}

// GetID of the current aggregate
func (b *BaseAggregate) GetID() string {
	return b.ID
}

// GetVersion of the current aggregate
func (b *BaseAggregate) GetVersion() int {
	return b.Version
}

// AddEvent to the aggregate
func (b *BaseAggregate) AddEvent(event Event) {
	b.Changes = append(b.Changes, event)
}

// AddError to the aggregate
func (b *BaseAggregate) AddError(err error) {
	b.Error = err
}

// GetError returns a list of errors
func (b *BaseAggregate) GetError() error {
	return b.Error
}

// HasError returns true if it contains at least one error
func (b *BaseAggregate) HasError() bool {
	return b.Error != nil
}

// AttachCommandID to every change for traceability
func (b *BaseAggregate) AttachCommandID(id string) {
	for i := range b.Changes {
		b.Changes[i].CommandID = id
	}
}
