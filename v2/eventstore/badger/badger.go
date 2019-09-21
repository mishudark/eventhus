package badger

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/mishudark/eventhus/v2"
)

// AggregateDB defines version and id of an aggregate
type AggregateDB struct {
	ID      string
	Version int
}

// EventDB defines the structure of the events to be stored
type EventDB struct {
	ID            string
	Type          string
	AggregateID   string
	AggregateType string
	CommandID     string
	RawData       []byte
	Timestamp     time.Time
	Version       int
}

// Client for access to badger
type Client struct {
	session *badger.DB
}

var _ eventhus.EventStore = (*Client)(nil)

// NewClient generates a new client for access to badger using badgerhold
func NewClient(dbDir string) (*Client, error) {
	options := badger.DefaultOptions(dbDir)
	options.ValueDir = dbDir

	session, err := badger.Open(options)
	if err != nil {
		return nil, err
	}

	cli := &Client{
		session: session,
	}

	return cli, nil
}

// Close db connection
func (c *Client) Close() error {
	return c.session.Close()
}

func (c *Client) save(events []eventhus.Event, version int, safe bool) error {
	if len(events) == 0 {
		return nil
	}

	aggregateID := events[0].AggregateID

	txn := c.session.NewTransaction(true)
	defer txn.Discard()

	for _, event := range events {
		raw, err := encode(event.Data)
		if err != nil {
			return err
		}

		item := EventDB{
			ID:            event.ID,
			Type:          event.Type,
			AggregateID:   event.AggregateID,
			AggregateType: event.AggregateType,
			CommandID:     event.CommandID,
			RawData:       raw,
		}

		blob, err := encode(item)
		if err != nil {
			return err
		}

		// the id contains the aggregateID as prefix
		// aggregateID.eventID
		id := fmt.Sprintf("%s.%s", aggregateID, event.ID)
		err = txn.Set([]byte(id), blob)
		if err != nil {
			return err
		}
	}

	// Now that events are saved, aggregate version needs to be updated
	aggregate := AggregateDB{
		ID:      events[0].AggregateID,
		Version: version + len(events),
	}

	aggregateBlob, err := encode(aggregate)
	if err != nil {
		return err
	}

	item, err := txn.Get([]byte(aggregate.ID))
	if version == 0 {
		switch err {
		case nil:
			return fmt.Errorf("badger: %s, aggregate already exists", aggregate.ID)
		case badger.ErrKeyNotFound:
			err = txn.Set([]byte(aggregate.ID), aggregateBlob)
		default: // another error differente from key not found is not desirable
			return err
		}
	} else {
		var blob []byte
		_, err := item.ValueCopy(blob)
		if err != nil {
			return err
		}

		var payload AggregateDB
		err = decode(blob, &payload)
		if err != nil {
			return err
		}

		if payload.Version != version {
			return fmt.Errorf("badger: %s, aggregate version missmatch, wanted: %d, got: %d", aggregate.ID, version, payload.Version)
		}

		err = txn.Set([]byte(aggregate.ID), aggregateBlob)
	}

	if err != nil {
		return err
	}

	return txn.Commit()
}

// SafeSave store the events without check the current version
func (c *Client) SafeSave(events []eventhus.Event, version int) error {
	return c.save(events, version, true)
}

// Save the events ensuring the current version
func (c *Client) Save(events []eventhus.Event, version int) error {
	return c.save(events, version, false)
}

// Load the stored events for an AggregateID
func (c *Client) Load(aggregateID string) ([]eventhus.Event, error) {
	var (
		events   []eventhus.Event
		eventsDB []EventDB
	)

	c.session.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		// prexi has the format aggregateID.
		prefix := []byte(aggregateID + ".")
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			err := item.Value(func(v []byte) error {
				var event EventDB

				err := decode(v, &event)
				if err != nil {
					return err
				}

				eventsDB = append(eventsDB, event)
				return nil
			})

			if err != nil {
				return err
			}
		}

		return nil
	})

	events = make([]eventhus.Event, len(eventsDB))
	register := eventhus.NewEventRegister()

	for i, dbEvent := range eventsDB {
		dataType, err := register.Get(dbEvent.Type)
		if err != nil {
			return events, err
		}

		if err = decode(dbEvent.RawData, dataType); err != nil {
			return events, err
		}

		// Translate dbEvent to eventhus.Event
		events[i] = eventhus.Event{
			AggregateID:   aggregateID,
			AggregateType: dbEvent.AggregateType,
			CommandID:     dbEvent.CommandID,
			Version:       dbEvent.Version,
			Type:          dbEvent.Type,
			Data:          dataType,
		}
	}

	return events, nil
}

func encode(value interface{}) ([]byte, error) {
	var buff bytes.Buffer
	en := gob.NewEncoder(&buff)

	err := en.Encode(value)
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

func decode(data []byte, value interface{}) error {
	var buff bytes.Buffer
	de := gob.NewDecoder(&buff)

	_, err := buff.Write(data)
	if err != nil {
		return err
	}

	return de.Decode(value)
}
