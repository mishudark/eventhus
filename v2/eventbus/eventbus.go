package eventbus

import (
	"fmt"

	"github.com/mishudark/eventhus/v2"
)

// MultiPublisherError is returned from publish when
// there is an error from a publisher.
type MultiPublisherError struct {
	Errors []error
}

// Error will produce an error message out of all errors.
func (e MultiPublisherError) Error() string {
	o := "A few errors occured:"

	for i, err := range e.Errors {
		o = fmt.Sprintf("%s\n\t%d) %s", o, i+1, err)
	}

	return o
}

// Add will push a new error on the slice.
func (e *MultiPublisherError) Add(err error) {
	if err != nil {
		e.Errors = append(e.Errors, err)
	}
}

// Len will return total amount of errors in error slice.
func (e *MultiPublisherError) Len() int {
	return len(e.Errors)
}

// MultiPublisher ...
type MultiPublisher struct {
	publishers []eventhus.EventBus
}

// NewMultiPublisher ...
func NewMultiPublisher(all ...eventhus.EventBus) *MultiPublisher {
	return &MultiPublisher{
		publishers: all,
	}
}

// Publish an event through all registered publishers.
func (c MultiPublisher) Publish(event eventhus.Event, bucket, subset string) error {
	errs := MultiPublisherError{}

	for _, p := range c.publishers {
		errs.Add(p.Publish(event, bucket, subset))
	}

	if errs.Len() > 0 {
		return errs
	}

	return nil
}
