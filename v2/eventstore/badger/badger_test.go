package badger

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/mishudark/eventhus/v2"
)

type TestEvent struct {
	Name string
	SKU  string
}

var (
	Aid = eventhus.GenerateUUID()
	cli *Client
)

func TestMain(m *testing.M) {
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		log.Fatalln(err)
	}

	cli, err = NewClient(tmpDir)
	log.Println(Aid)
	if err != nil {
		log.Fatalln(err)
	}

	result := m.Run()
	cli.Close()

	os.Exit(result)
}

func TestClientSave(t *testing.T) {
	events := []eventhus.Event{
		eventhus.Event{
			ID:            eventhus.GenerateUUID(),
			AggregateID:   Aid,
			AggregateType: "order",
			Version:       1,
			Type:          "test_event",
			Data: TestEvent{
				Name: "muñeca",
				SKU:  "123",
			},
		},
		eventhus.Event{
			ID:            eventhus.GenerateUUID(),
			AggregateID:   Aid,
			AggregateType: "order",
			Version:       1,
			Type:          "test_event",
			Data: TestEvent{
				Name: "muñeca",
				SKU:  "123",
			},
		},
	}

	err := cli.Save(events, 0)
	if err != nil {
		t.Error("expected nil, got", err)
	}

}

func TestClientLoad(t *testing.T) {
	reg := eventhus.NewEventRegister()
	reg.Set(&TestEvent{})

	events, err := cli.Load(Aid)
	if err != nil {
		t.Error("expected nil, got", err)
	}

	length := len(events)
	if length != 2 {
		t.Errorf("[events] expected: 2, got: %d", length)
	}
}
