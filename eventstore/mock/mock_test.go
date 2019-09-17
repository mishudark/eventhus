package mock

import (
	"github.com/mishudark/eventhus"
	"github.com/oklog/ulid"
	"math/rand"
	"testing"
	"time"
)

type SubEvent2 struct {
	Name string
	SKU  string
}

func TestClient_LoadEmptyShouldReturn0Events(t *testing.T) {
	store := NewClient()
	events, err := store.Load("new-aggregate")

	if err != nil {
		t.Error("Unexpected error", err)
	}

	if events != nil {
		t.Error("Function returned nil!")
	}

	if l := len(events); l != 0 {
		t.Error("Expected 0 events, returned were", l)
	}
}

func TestClient_SaveAndLoad(t *testing.T) {
	store := NewClient()

	ta := time.Now()
	entropy := rand.New(rand.NewSource(ta.UnixNano()))
	aid := ulid.MustNew(ulid.Timestamp(ta), entropy)

	aggregateId := aid.String()
	events := []eventhus.Event{
		{
			AggregateID:   aggregateId,
			AggregateType: "order",
			Version:       1,
			Type:          "SubEvent2",
			Data: SubEvent2{
				Name: "muñeca",
				SKU:  "123",
			},
		},
		{
			AggregateID:   aggregateId,
			AggregateType: "order",
			Version:       1,
			Type:          "SubEvent2",
			Data: SubEvent2{
				Name: "muñeca",
				SKU:  "123",
			},
		},
	}

	err := store.Save(events, 0)

	if err != nil {
		t.Error("Unexpected error", err)
	}

	events, err = store.Load(aggregateId)

	if err != nil {
		t.Error("Unexpected error", err)
	}

	if l := len(events); l != 2 {
		t.Error("Expected 2 events, returned were", l)
	}
}
