package cqrs

//EventBus defines the methods for manage the events publisher and consumer
type EventBus interface {
	Publish(event Event, bucket, subset string) error
}
