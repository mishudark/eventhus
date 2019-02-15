package eventhus

import "testing"

type MockAggregate struct {
	BaseAggregate
}

func (m *MockAggregate) Reduce(event Event) error {
	return nil
}

func (m *MockAggregate) HandleCommand(command Command) error {
	return nil
}

func TestBaseAggregateUncommited(t *testing.T) {
	var mock MockAggregate
	event := Event{
		ID:            "bzvayj",
		AggregateID:   "kasdyui",
		AggregateType: "mock_aggregate",
		Type:          "event",
		Data: struct {
			Foo int
		}{
			Foo: 1,
		},
	}

	length := len(mock.Uncommited())
	if length != 0 {
		t.Error("expected 0, got", length)
	}

	Dispatch(&mock, event)

	length = len(mock.Uncommited())
	if length != 1 {
		t.Error("expected 1, got", length)
	}

	mock.ClearUncommited()
	length = len(mock.Uncommited())
	if length != 0 {
		t.Error("expected 0, got", length)
	}
}
