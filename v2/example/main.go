package main

import (
	"fmt"
	ev "github.com/mishudark/eventhus/v2"
)

type Reducer interface {
	Reduce(command ev.Command)
}

type Category string

const (
	Tacos  Category = "tacos"
	Tortas Category = "tortas"
)

type Location struct {
	Latitude  float64
	Longitude float64
}

type addVenuePayload struct {
	Name        string
	Description string
	Price       float64
	User        string
	Location    Location
	Category    Category
	URL         string
	Rating      int
	Monday      string
	Tuesday     string
	Wenesday    string
	Thursday    string
	Friday      string
	Saturday    string
	Sunday      string
}

type AddVenue struct {
	ev.Command
	addVenuePayload
}

type VenueAddedUnplubish struct {
	ev.Event
	addVenuePayload
}

type Venue struct {
	ev.BaseAggregate
	AddVenue
	URL      string
	ShortURL string
}

func (v *Venue) Reduce(event ev.Event) error {
	return nil
}

func (v *Venue) HandleCommand(command ev.Command) error {
	event := ev.Event{
		AggregateID:   v.ID,
		AggregateType: "venue",
	}

	switch command.GetType() {
	case "add_venue":
		if c, ok := command.(*AddVenue); ok {
			event.Data = c.addVenuePayload
		} else {
			return fmt.Errorf("%s: can't cast to the given Command type", command.GetType())
		}
	case "rate_venue":
		if c, ok := command.(*RateVenue); ok {
			event.Data = struct {
				Rate int
			}{
				c.Rate,
			}
		} else {
			return fmt.Errorf("%s: can't cast to the given Command type", command.GetType())
		}
	}

	ev.Dispatch(v, event)
	return nil
}

type RateVenue struct {
	ev.BaseCommand
	Rate int
}

func main() {
	// eventbus  := nats.NewClient()
	// eventStore := postgres.NewClient()
	// redux.Register(Venue{},
	// "add_venue",
	// "rate_venue",
	// "rate_burn",
	// )
}
