package mock

import (
	"github.com/mishudark/eventhus"
)

type Client struct{}

/*
Creates a new mocked event bus.

This does implementation is rather crude and does not perfectly mimic the behaviour of
the other implementations.

Warning: All events published by this client are simply thrown away and not used in any way!

Therefore this should only be used to try out eventhus without
having to connect to an event store or an event bus.
*/
func NewClient() *Client {
	return &Client{}
}

func (c *Client) Publish(event eventhus.Event, bucket, subset string) error {
	return nil
}
