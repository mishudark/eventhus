package badger

import (
	"math/rand"
	"testing"
	"time"

	"github.com/mishudark/eventhus"
	"github.com/oklog/ulid"
)

type SomeEvent struct {
	Name string
	SKU  string
}

var Aid ulid.ULID

func TestNewClient(t *testing.T) {
	eventStore, err := NewClient()
	cli := eventStore.(*Client)
	defer cli.CloseClient()
	if err != nil {
		t.Error("expected nil, got", err)
	}
}

func TestClientSave(t *testing.T) {
	eventStore, err := NewClient()
	cli := eventStore.(*Client)
	defer cli.CloseClient()
	if err != nil {
		t.Error("expected nil, got", err)
	}

	ta := time.Now()
	entropy := rand.New(rand.NewSource(ta.UnixNano()))
	Aid = ulid.MustNew(ulid.Timestamp(ta), entropy)

	events := []eventhus.Event{
		eventhus.Event{
			AggregateID:   Aid.String(),
			AggregateType: "order",
			Version:       1,
			Type:          "SomeEvent",
			Data: SomeEvent{
				Name: "muñeca",
				SKU:  "123",
			},
		},
		eventhus.Event{
			AggregateID:   Aid.String(),
			AggregateType: "order",
			Version:       1,
			Type:          "SomeEvent",
			Data: SomeEvent{
				Name: "muñeca",
				SKU:  "123",
			},
		},
	}

	err = eventStore.Save(events, 0)
	if err != nil {
		t.Error("expected nil, got", err)
	}

}

func TestClientLoad(t *testing.T) {
	reg := eventhus.NewEventRegister()
	reg.Set(SomeEvent{})

	eventStore, err := NewClient()
	cli := eventStore.(*Client)
	defer cli.CloseClient()

	if err != nil {
		t.Error("expected nil, got", err)
	}
	_, err = eventStore.Load(Aid.String())
	if err != nil {
		t.Error("expected nil, got", err)
	}
}
