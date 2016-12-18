package cqrs

import "testing"

type MockAggregate struct {
	BaseAggregate
}

func TestBaseAggregateUncommited(t *testing.T) {
	var mock MockAggregate
	var event Event

	length := len(mock.Uncommited())
	if length != 0 {
		t.Error("expected 0, got", length)
	}

	mock.Changes = append(mock.Changes, event)

	length = len(mock.Uncommited())
	if length == 0 {
		t.Error("expected 0, got", length)
	}

	mock.ClearUncommited()
	length = len(mock.Uncommited())
	if length != 0 {
		t.Error("expected 0, got", length)
	}
}
