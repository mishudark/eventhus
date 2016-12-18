package cqrs

//BaseCommand contains the basic info
//that all commands should have
type BaseCommand struct {
	Type          string
	AggregateID   string
	AggregateType string
}
