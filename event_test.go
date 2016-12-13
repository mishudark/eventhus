package cqrs

import "testing"

type SubEvent struct {
	Name string
	SKU  string
}

func TestEventTypeRegister(t *testing.T) {
	var err error
	reg := NewRegister()
	count := reg.Count()
	if count != 0 {
		t.Error("expected: 0, got: ", count)
	}

	reg.Register(SubEvent{})
	count = reg.Count()
	if count != 1 {
		t.Error("expected: 1, got: ", count)
	}

	_, err = reg.Get("hola")
	if err == nil {
		t.Error("expected error, got nil")
	}

	for i := 0; i < 1000; i++ {
		go func() {
			_, err = reg.Get("cqrs.SubEvent")
			if err != nil {
				t.Error("expected nil, got", err)
			}
		}()
	}

}
