package mongo

import (
	"eventhus"
	"math/rand"
	"testing"
	"time"

	"github.com/oklog/ulid"
)

type SubEvent2 struct {
	Name string
	SKU  string
}

func TestNewClient(t *testing.T) {
	_, err := NewClient("localhost", 27017, "grunt")
	if err != nil {
		t.Error("expected nil, got", err)
	}
}

func TestClientLoad(t *testing.T) {
	cli, err := NewClient("localhost", 27017, "grunt")
	if err != nil {
		t.Error("expected nil, got", err)
	}
	_, err = cli.Load("123")
	if err != nil {
		t.Error("expected nil, got", err)
	}
}

func TestClientSave(t *testing.T) {
	cli, err := NewClient("localhost", 27017, "grunt")
	if err != nil {
		t.Error("expected nil, got", err)
	}

	ta := time.Now()
	entropy := rand.New(rand.NewSource(ta.UnixNano()))
	aid := ulid.MustNew(ulid.Timestamp(ta), entropy)

	events := []eventhus.Event{
		eventhus.Event{
			AggregateID:   aid.String(),
			AggregateType: "order",
			Version:       1,
			Type:          "eventhus.SubEvent2",
			Data: SubEvent2{
				Name: "muñeca",
				SKU:  "123",
			},
		},
		eventhus.Event{
			AggregateID:   aid.String(),
			AggregateType: "order",
			Version:       1,
			Type:          "eventhus.SubEvent2",
			Data: SubEvent2{
				Name: "muñeca",
				SKU:  "123",
			},
		},
	}

	err = cli.Save(events, 0)
	if err != nil {
		t.Error("expected nil, got", err)
	}

	err = cli.Save(events, 0)
	if err == nil {
		t.Error("expected error, got nil")
	}
}
