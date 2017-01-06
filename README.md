# Eventhus
CQRS/ES toolkit for Go.

**CQRS** stands for Command Query Responsibility Segregation. It's a pattern that I first heard described by Greg Young. At its heart is the notion that you can use a different model to update information than the model you use to read information.

The mainstream approach people use for interacting with an information system is to treat it as a CRUD datastore. By this I mean that we have mental model of some record structure where we can create new records, read records, update existing records, and delete records when we're done with them. In the simplest case, our interactions are all about storing and retrieving these records.

**Event Sourcing** ensure that every change to the state of an application is captured in an event object, and that these event objects are themselves stored in the sequence they were applied for the same lifetime as the application state itself.

# Usage
There are 3 basic units of work `event`, `command` and `aggregate` 

## Command
A command describe an **action** that should be performed, it's always named in the imperative tense such as  `PerformDeposit` `CreateAccount` 

Let’s start with some code:

```go

import "github.com/mishudark/eventhus"

//PerformDeposit to an account
type PerformDeposit struct {
	eventhus.BaseCommand
	Ammount int
}
```

At the beginning we create the `PerformDeposit` command,  it contains an anonymous struct field of type `eventhus.BaseCommand`. This means `PerformDeposit` automatically acquires all the methods of `eventhus.BaseCommand`.

Also you can define custom fields, in this case `Ammount` contains quantity to being depositen in an account.

##Event
An event is the notification that some happend in the past, you can view an event as the representation of reaction of **a command after being executed**. All events should be represented as verbs in the past tense such as `CustomerRelocated`, `CargoShipped` or `InventoryLossageRecorded`

```go
//DepositPerformed event
type DepositPerformed struct {
	Ammount int
}
```

We create the `DepositPerformed` event, it's a pure go struct, and it's the past equivalent to the previous command `PerformDeposit`

##Aggregate
The aggregate is a logical boundary for things that can change in a business transaction of a given context. In the eventhus context, it simple process the commands and produce events. 

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

We create the `Account` aggregate, it contains an anonymous struct field of type `eventhus.BaseAggregate`. This means `Account` automatically acquires all the methods of `eventhus.BaseAggregate`.

Additionally `Account` has the fields `Balance` and `Owner` that represents the basic info of this context

Now that we have our `aggregate`, we need to process the `PerformDeposit` command that we created earlier
 
```go
//HandleCommand create events and validate based on such command
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
			c.Ammount,
		}
	}

	a.BaseAggregate.ApplyChangeHelper(a, event, true)
	return nil
}

```
First, we create an `event` with the basic info `AggregateID` as an identifier and `AggregateType` with the same name as our `aggregate`. Next, we use a switch to determine the type of the `command` and produce and `event` as a consecuence.

Finally, the event should be applied to our aggregate, we use the helper `BaseAggregate.ApplyChangeHelper` with the params `aggregate`, `event` and the last argument set to `true` that means, it should be stored and published via `event store` and `event publisher` 
  
Note: `eventhus.BaseAggregate` has some helper methods to make our live easier, we use `HandleCommand` to process a `command` and produce the respective `event`

The last step in the aggregate journey is apply the `events` to our `aggregate`

```go
//ApplyChange to account
func (a *Account) ApplyChange(event eventhus.Event) {
	switch e := event.Data.(type) {
	case *AccountCreated:
		a.Owner = e.Owner
		a.ID = event.AggregateID
	case *DepositPerformed:
		a.Balance += e.Ammount
	}
}
```

Also we use a switch-case format to determine the type of the `event` (note that events are pointers), and apply the respective changes


Note: The aggregate is never save in it's current state, instead is stored as a series of `events` that can recreate the aggregate in it's last state.

Save events, publish it and recreate an aggregate from `event store` is made by **Eventhus** out of the box

# Config 
`Eventhus` needs to be configured to manage events, commands and to knows where to store and publish events, please refer to [config example](https://github.com/mishudark/eventhus/blob/master/examples/bank/cmd/main/config.go) for more info  

# Event Store
Currently it has support for `MongoDB`, `Rethinkdb` is in the scope to be added

# Event Publisher
`RabbitMQ` and `Nats.io` are supported

## Prior Art

- [looplab/eventhorizon](https://github.com/looplab/eventhorizon)
- [andrewwebber/cqrs](https://github.com/andrewwebber/cqrs)

