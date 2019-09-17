package mock

import (
	"github.com/mishudark/eventhus"
)

type Client struct {
	// Maps aggregate ids to their event arrays
	events map[string][]eventhus.Event
}

/*
Creates a new mocked event store.

This does implementation is rather crude and does not perfectly mimic the behaviour of
the other implementations.

Therefore this mocked store should only be used to try out eventhus without
having to connect to an event store or an event bus.
*/
func NewClient() *Client {
	return &Client{
		events: make(map[string][]eventhus.Event),
	}
}

func (c *Client) Save(events []eventhus.Event, version int) error {
	for _, event := range events {
		aggregateId := event.AggregateID
		c.events[aggregateId] = append(c.events[aggregateId], event)
	}

	return nil
}

func (c *Client) SafeSave(events []eventhus.Event, version int) error {
	panic("SafeSave is not yet implemented on the mocked event store!")
}

func (c *Client) Load(aggregateID string) ([]eventhus.Event, error) {
	return c.events[aggregateID], nil
}
