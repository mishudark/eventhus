package badger

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/mishudark/eventhus/v2"
)

//AggregateDB defines the collection to store the aggregate with their events
type AggregateDB struct {
	ID      string    `json:"_id"`
	Version int       `json:"version"`
	Events  []EventDB `json:"events"`
}

//EventDB defines the structure of the events to be stored
type EventDB struct {
	Type          string `json:"event_type"`
	AggregateID   string `json:"_id"`
	RawData       []byte `json:"data,omitempty"`
	data          interface{}
	Timestamp     time.Time `json:"timestamp"`
	AggregateType string    `json:"aggregate_type"`
	Version       int       `json:"version"`
}

//Client for access to boltdb
type Client struct {
	session *badger.DB
}

//NewClient generates a new client for access to BadgerDB
func NewClient(dbDir string) (eventhus.EventStore, error) {
	// Open the Badger database located in the /tmp/badger directory.
	// It will be created if it doesn't exist.
	opts := badger.DefaultOptions
	opts.Dir = dbDir
	opts.ValueDir = dbDir
	session, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	cli := &Client{
		session,
	}

	return cli, nil
}

// CloseClient closes the db connection
func (c *Client) CloseClient() error {
	return c.session.Close()
}

func (c *Client) save(events []eventhus.Event, version int, safe bool) error {
	if len(events) == 0 {
		return nil
	}

	// Build all event records, with incrementing versions starting from the
	// original aggregate version.
	eventsDB := make([]EventDB, len(events))
	aggregateID := events[0].AggregateID

	for i, event := range events {

		// Create the event record with timestamp.
		eventsDB[i] = EventDB{
			Type:          event.Type,
			AggregateID:   event.AggregateID,
			Timestamp:     time.Now(),
			AggregateType: event.AggregateType,
			Version:       1 + version + i,
		}

		// Marshal event data if there is any.
		if event.Data != nil {
			rawData, err := json.Marshal(event.Data)
			if err != nil {
				return err
			}
			eventsDB[i].RawData = rawData
		}
	}

	// Either insert a new aggregate or append to an existing.
	if version == 0 {
		aggregate := AggregateDB{
			ID:      aggregateID,
			Version: len(eventsDB),
			Events:  eventsDB,
		}

		err := c.session.Update(func(txn *badger.Txn) error {
			aggregateJSON, _ := json.Marshal(aggregate)
			err := txn.Set([]byte(aggregateID), aggregateJSON)
			return err
		})

		return err
	}
	// Increment aggregate version on insert of new event record, and
	// only insert if version of aggregate is matching (ie not changed
	// since loading the aggregate).

	err := c.session.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(aggregateID))
		if err != nil {
			return err
		}
		val, err := item.Value()
		if err != nil {
			return err
		}

		var aggregateDB AggregateDB
		for {
			err := json.Unmarshal(val, &aggregateDB)
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
		}

		if !safe && aggregateDB.Version != version {
			return fmt.Errorf("There was an concurrent update of  %s", aggregateID)
		}

		aggregateJSON, err := json.Marshal(aggregateDB)
		if err != nil {
			return err
		}

		aggregateDB.Version = len(eventsDB)
		return txn.Set([]byte(aggregateID), aggregateJSON)
	})

	return err
}

//SafeSave store the events without check the current version
func (c *Client) SafeSave(events []eventhus.Event, version int) error {
	return c.save(events, version, true)
}

//Save the events ensuring the current version
func (c *Client) Save(events []eventhus.Event, version int) error {
	return c.save(events, version, false)
}

//Load the stored events for an AggregateID
func (c *Client) Load(aggregateID string) ([]eventhus.Event, error) {
	var events []eventhus.Event

	var aggregate AggregateDB

	err := c.session.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(aggregateID))
		if err != nil {
			return err
		}

		aggregateJSON, err := item.Value()
		if err != nil {
			return err
		}

		return json.Unmarshal(aggregateJSON, &aggregate)
	})

	if err != nil {
		// events is empty in this case
		return events, err
	}

	events = make([]eventhus.Event, len(aggregate.Events))
	register := eventhus.NewEventRegister()

	for i, dbEvent := range aggregate.Events {
		// Create an event of the correct type.
		dataType, err := register.Get(dbEvent.Type)
		if err != nil {
			return events, err
		}

		if err := json.Unmarshal(dbEvent.RawData, dataType); err != nil {
			return events, err
		}

		// Set concrete event and zero out the decoded event.
		dbEvent.data = dataType
		dbEvent.RawData = []byte{}

		// Translate dbEvent to eventhus.Event
		events[i] = eventhus.Event{
			AggregateID:   aggregateID,
			AggregateType: dbEvent.AggregateType,
			Version:       dbEvent.Version,
			Type:          dbEvent.Type,
			Data:          dbEvent.data,
		}
	}

	return events, nil
}
