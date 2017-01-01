package mongo

import (
	"eventhus"
	"fmt"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//AggregateDB defines the collection to store the aggregate with their events
type AggregateDB struct {
	ID      string    `bson:"_id"`
	Version int       `bson:"version"`
	Events  []EventDB `bson:"events"`
}

//EventDB defines the structure of the events to be stored
type EventDB struct {
	Type          string      `bson:"event_type"`
	AggregateID   string      `bson:"_id"`
	RawData       bson.Raw    `bson:"data,omitempty"`
	data          interface{} `bson:"-"`
	Timestamp     time.Time   `bson:"timestamp"`
	AggregateType string      `bson:"aggregate_type"`
	Version       int         `bson:"version"`
}

//Client for access to mongodb
type Client struct {
	db      string
	session *mgo.Session
}

//NewClient generates a new client to access to mongodb
func NewClient(host string, port int, db string) (eventhus.EventStore, error) {
	session, err := mgo.Dial(fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, err
	}

	session.SetMode(mgo.Monotonic, true)
	//session.SetSafe(&mgo.Safe{W: 1})

	cli := &Client{
		db,
		session,
	}

	return cli, nil
}

func (c *Client) save(events []eventhus.Event, version int, safe bool) error {
	if len(events) == 0 {
		return nil
	}

	sess := c.session.Copy()
	defer sess.Close()

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
			rawData, err := bson.Marshal(event.Data)
			if err != nil {
				return err
			}
			eventsDB[i].RawData = bson.Raw{Kind: 3, Data: rawData}
		}
	}

	// Either insert a new aggregate or append to an existing.
	if version == 0 {
		aggregate := AggregateDB{
			ID:      aggregateID,
			Version: len(eventsDB),
			Events:  eventsDB,
		}

		if err := sess.DB(c.db).C("events").Insert(aggregate); err != nil {
			return err
		}
	} else {
		// Increment aggregate version on insert of new event record, and
		// only insert if version of aggregate is matching (ie not changed
		// since loading the aggregate).
		query := bson.M{"_id": aggregateID}
		if !safe {
			query["version"] = version
		}

		if err := sess.DB(c.db).C("events").Update(
			query,
			bson.M{
				"$push": bson.M{"events": bson.M{"$each": eventsDB}},
				"$inc":  bson.M{"version": len(eventsDB)},
			},
		); err != nil {
			return err
		}
	}
	return nil
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

	sess := c.session.Copy()
	defer sess.Close()

	var aggregate AggregateDB
	err := sess.DB(c.db).C("events").FindId(aggregateID).One(&aggregate)
	if err == mgo.ErrNotFound {
		return events, nil
	} else if err != nil {
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

		// Manually decode the raw BSON event.
		if err := dbEvent.RawData.Unmarshal(dataType); err != nil {
			return events, err
		}

		// Set conrcete event and zero out the decoded event.
		dbEvent.data = dataType
		dbEvent.RawData = bson.Raw{}

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
