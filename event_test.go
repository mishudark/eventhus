package eventhus

import "testing"

type SubEvent struct {
	Name string
	SKU  string
}

func TestEventTypeRegister(t *testing.T) {
	var err error
	reg := NewEventRegister()
	count := reg.Count()
	if count != 0 {
		t.Error("expected: 0, got: ", count)
	}

	reg.Set(SubEvent{})
	count = reg.Count()
	if count != 1 {
		t.Error("expected: 1, got: ", count)
	}

	_, err = reg.Get("hola")
	if err == nil {
		t.Error("expected error, got nil")
	}

	_, err = reg.Get("SubEvent")
	if err != nil {
		t.Error("expected nil, got", err)
	}

}
