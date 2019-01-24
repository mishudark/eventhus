package eventbus

import (
	"log"

	"github.com/mishudark/eventhus/v2"
)

// Logger logs messages sent to the event bus.
type Logger struct {
	log *log.Logger
}

// NewLogger returns new logger struct.
func NewLogger(l *log.Logger) *Logger {
	return &Logger{
		log: l,
	}
}

// Publish logs event details out.
func (l *Logger) Publish(e eventhus.Event, b, s string) error {
	log.Printf("bucket: %s subset: %s event: %+v", b, s, e)
	return nil
}
