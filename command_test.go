package cqrs

import (
	"encoding/json"
	"testing"
)

func TestBaseCommand(t *testing.T) {
	type newProductCommand struct {
		BaseCommand
		SKU string
	}

	var n newProductCommand
	blob := []byte(`{"type":"product", "aggregate_id":"asdas", "aggregate_type": "order", "sku":"AF12R"}`)
	if err := json.Unmarshal(blob, &n); err != nil {
		t.Error("expected: nil, got: ", err)
	}

	if n.Type != "product" {
		t.Error("expected: product, got: ", n.Type)
	}
}
