package eventhus

import "testing"

type MockAggregate struct {
	BaseAggregate
}

func TestBaseAggregateUncommitted(t *testing.T) {
	var mock MockAggregate
	var event Event

	length := len(mock.Uncommitted())
	if length != 0 {
		t.Error("expected 0, got", length)
	}

	mock.Changes = append(mock.Changes, event)

	length = len(mock.Uncommitted())
	if length == 0 {
		t.Error("expected 0, got", length)
	}

	mock.ClearUncommitted()
	length = len(mock.Uncommitted())
	if length != 0 {
		t.Error("expected 0, got", length)
	}
}
