# Eventhus

CQRS/ES toolkit for Go.

**CQRS** stands for Command Query Responsibility Segregation. It's a pattern that I first heard described by Greg Young. At its heart is the notion that you can use a different model to update information than the model you use to read information.

The mainstream approach people use for interacting with an information system is to treat it as a CRUD datastore. By this I mean that we have mental model of some record structure where we can create new records, read records, update existing records, and delete records when we're done with them. In the simplest case, our interactions are all about storing and retrieving these records.

**Event Sourcing** ensures that every change to the state of an application is captured in an event object, and that these event objects are themselves stored in the sequence they were applied for the same lifetime as the application state itself.

# Examples

[bank account](https://github.com/mishudark/eventhus/blob/master/examples/bank) shows a full example with `deposits` and `withdrawls`.

# Usage

There are 3 basic units of work: `event`, `command` and `aggregate`.

## Command

A command describes an **action** that should be performed; it's always named in the imperative tense such as `PerformDeposit` or `CreateAccount`.

Let’s start with some code:

```go
import "github.com/mishudark/eventhus"

// PerformDeposit to an account
type PerformDeposit struct {
	eventhus.BaseCommand
	Amount int
}
```

At the beginning, we create the `PerformDeposit` command. It contains an anonymous struct field of type `eventhus.BaseCommand`. This means `PerformDeposit` automatically acquires all the methods of `eventhus.BaseCommand`.

You can also define custom fields, in this case `Amount` contains a quantity to be deposited into an account.

## Event

An event is the notification that something happened in the past. You can view an event as the representation of the reaction to **a command after being executed**. All events should be represented as verbs in the past tense such as `CustomerRelocated`, `CargoShipped` or `InventoryLossageRecorded`.

We create the `DepositPerformed` event; it's a pure go struct, and it's the past equivalent to the previous command `PerformDeposit`:

```go
// DepositPerformed event
type DepositPerformed struct {
	Amount int
}
```

## Aggregate

The aggregate is a logical boundary for things that can change in a business transaction of a given context. In the **Eventhus** context, it simplifies the process the commands and produce events.

Show me the code!

```go
import "github.com/mishudark/eventhus"

//Account of bank
type Account struct {
	eventhus.BaseAggregate
	Owner   string
	Balance int
}
```

We create the `Account` aggregate. It contains an anonymous struct field of type `eventhus.BaseAggregate`. This means `Account` automatically acquires all the methods of `eventhus.BaseAggregate`.

Additionally `Account` has the fields `Balance` and `Owner` that represent the basic info of this context.

Now that we have our `aggregate`, we need to process the `PerformDeposit` command that we created earlier:

```go
// HandleCommand create events and validate based on such command
func (a *Account) HandleCommand(command eventhus.Command) error {
	event := eventhus.Event{
		AggregateID:   a.ID,
		AggregateType: "Account",
	}

	switch c := command.(type) {
	case CreateAccount:
		event.AggregateID = c.AggregateID
		event.Data = &AccountCreated{c.Owner}

	case PerformDeposit:
		event.Data = &DepositPerformed{
			c.Amount,
		}
	}

	a.BaseAggregate.ApplyChangeHelper(a, event, true)
	return nil
}
```

First, we create an `event` with the basic info `AggregateID` as an identifier and `AggregateType` with the same name as our `aggregate`. Next, we use a switch to determine the type of the `command` and produce an `event` as a result.

Finally, the event should be applied to our aggregate; we use the helper `BaseAggregate.ApplyChangeHelper` with the params `aggregate`, `event` and the last argument set to `true`, meaning it should be stored and published via `event store` and `event publisher`.

Note: `eventhus.BaseAggregate` has some helper methods to make our life easier, we use `HandleCommand` to process a `command` and produce the respective `event`.

The last step in the aggregate journey is to apply the `events` to our `aggregate`:

```go
// ApplyChange to account
func (a *Account) ApplyChange(event eventhus.Event) {
	switch e := event.Data.(type) {
	case *AccountCreated:
		a.Owner = e.Owner
		a.ID = event.AggregateID
	case *DepositPerformed:
		a.Balance += e.Amount
	}
}
```

Also, we use a switch-case format to determine the type of the `event` (note that events are pointers), and apply the respective changes.

Note: The aggregate is never saved in its current state. Instead, it is stored as a series of `events` that can recreate the aggregate in its last state.

Saving the events, publishing them, and recreating an aggregate from `event store` is made by **Eventhus** out of the box.

# Config

`Eventhus` needs to be configured to manage events and commands, and to know where to store and publish events.

## Event Store

Currently, it has support for `MongoDB`. `Rethinkdb` is in the scope to be added. A [mock implementation for development](#mocked-event-store) is available.

We create an `event store` with `config.Mongo`; it accepts `host`, `port` and `table` as arguments:

```go
import "github.com/mishudark/eventhus/config"
...

config.Mongo("localhost", 27017, "bank") // event store
```

## Event Publisher

`RabbitMQ` and `Nats.io` are supported. A [mock implementation for development](#mocked-event-bus) is available.

We create an `eventbus` with `config.Nats`, it accepts `url data config` and `useSSL` as arguments:

```go
import 	"github.com/mishudark/eventhus/config"
...

config.Nats("nats://ruser:T0pS3cr3t@localhost:4222", false) // event bus
```

## Wire it all together

Now that we have all the pieces, we can register our `events`, `commands` and `aggregates`:

```go
import (
	"github.com/mishudark/eventhus"
	"github.com/mishudark/eventhus/commandhandler/basic"
	"github.com/mishudark/eventhus/config"
	"github.com/mishudark/eventhus/examples/bank"
)

func getConfig() (eventhus.CommandBus, error) {
	// register events
	reg := eventhus.NewEventRegister()
	reg.Set(bank.AccountCreated{})
	reg.Set(bank.DepositPerformed{})
	reg.Set(bank.WithdrawalPerformed{})

    // wire all parts together
	return config.NewClient(
		config.Mongo("localhost", 27017, "bank"),                    // event store
		config.Nats("nats://ruser:T0pS3cr3t@localhost:4222", false), // event bus
		config.AsyncCommandBus(30),                                  // command bus
		config.WireCommands(
			&bank.Account{},          // aggregate
			basic.NewCommandHandler,  // command handler
			"bank",                   // event store bucket
			"account",                // event store subset
			bank.CreateAccount{},     // command
			bank.PerformDeposit{},    // command
			bank.PerformWithdrawal{}, // command
		),
	)
}
```

Now you are ready to process commands:

```go
uuid, _ := utils.UUID()

// 1) Create an account
var account bank.CreateAccount
account.AggregateID = uuid
account.Owner = "mishudark"

commandBus.HandleCommand(account)
```

First, we generate a new `UUID`. This is because is a new account and we need a unique identifier. After we created the basic structure of our `CreateAccount` command, we only need to send it using the `commandbus` created in our config.

## Event consumer

You should listen to your `eventbus`, the format of the event is always the same, only the `data` key changes in the function of your event struct.

```json
{
  "id": "0000XSNJG0SB2WDBTATBYEC51P",
  "aggregate_id": "0000XSNJG0N0ZVS3YXM4D7ZZ9Z",
  "aggregate_type": "Account",
  "version": 1,
  "type": "AccountCreated",
  "data": {
    "owner": "mishudark"
  }
}
```

## Mock Classes

There are mock implementations for development without an actual event store or event bus.

For Details, have a look inside the `mock` packages (example: `eventhus/eventbus/mock`)

### Mocked Event Store

```go
import "github.com/mishudark/eventhus/config"
...

config.MockEventStore()
```

This requires no configuration and stores the events internally in a map, mapping the aggregate ids to an array of events.

### Mocked Event Bus

```go
import "github.com/mishudark/eventhus/config"
...

config.MockEventBus()
```

This implementation simply does nothing with the events.

## Prior Art

- [looplab/eventhorizon](https://github.com/looplab/eventhorizon)
- [andrewwebber/cqrs](https://github.com/andrewwebber/cqrs)
