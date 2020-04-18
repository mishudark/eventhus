package eventhus

// Command contains the methods to retrieve basic info about it
type Command interface {
	GetType() string
	GetAggregateID() string
	GetAggregateType() string
	IsValid() bool
	GetVersion() int
}

// BaseCommand contains the basic info
// that all commands should have
type BaseCommand struct {
	Type          string
	AggregateID   string
	AggregateType string
	Version       int
}

// GetAggregateID returns the command aggregate ID
func (b BaseCommand) GetAggregateID() string {
	return b.AggregateID
}

// GetType returns the command type
func (b BaseCommand) GetType() string {
	return b.Type
}

// GetAggregateType returns the command aggregate type
func (b BaseCommand) GetAggregateType() string {
	return b.AggregateType
}

// IsValid checks validates the command
func (b BaseCommand) IsValid() bool {
	return true
}

// GetVersion of the command
func (b BaseCommand) GetVersion() int {
	return b.Version
}
