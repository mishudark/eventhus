# Eventhus
CQRS/ES toolkit for Go.

**CQRS** stands for Command Query Responsibility Segregation. It's a pattern that I first heard described by Greg Young. At its heart is the notion that you can use a different model to update information than the model you use to read information.

The mainstream approach people use for interacting with an information system is to treat it as a CRUD datastore. By this I mean that we have mental model of some record structure where we can create new records, read records, update existing records, and delete records when we're done with them. In the simplest case, our interactions are all about storing and retrieving these records.

**Event Sourcing** ensure that every change to the state of an application is captured in an event object, and that these event objects are themselves stored in the sequence they were applied for the same lifetime as the application state itself.

# Usage
There are 3 basic units of work `event`, `command` and `aggregate` 

## Command
A command describe an **action** that should be performed, it's always named in the imperative tense,  `PerformDeposit` and `CreateAccount` are examples of how to name a command

Let’s start with some code:

```go

import "github.com/mishudark/eventhus"

//CreateAccount assigned to an owner
type CreateAccount struct {
	eventhus.BaseCommand
	Owner string
}
```

At the beginning we create the `CreateAccount` command,  it contains an anonymous struct field of type `eventhus.BaseCommand`. This means `CreateAccount` automatically acquires all the methods of `eventhus.BaseCommand`.


# Event Store
Currently it has support for `MongoDB`, `Rethinkdb` is in the scope to be added

# Event Publisher
`RabbitMQ` and `Nats.io` are supported

## Prior Art

- [looplab/eventhorizon](https://github.com/looplab/eventhorizon)
- [andrewwebber/cqrs](https://github.com/andrewwebber/cqrs)

