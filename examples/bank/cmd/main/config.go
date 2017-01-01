package main

import (
	"cqrs"
	"cqrs/commandbus"
	"cqrs/eventbus/nats"
	"cqrs/eventstore/mongo"
	"cqrs/examples/bank"
	"log"
)

func config() (cqrs.CommandBus, error) {
	//register events
	reg := cqrs.NewEventRegister()
	reg.Set(bank.AccountCreated{})

	//event store
	eventstore, err := mongo.NewClient("localhost", 27017, "bank")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	//eventbus
	//rabbit, err := rabbitmq.NewClient("guest", "guest", "localhost", 5672)
	nat, err := nats.NewClient("nats://ruser:T0pS3cr3t@localhost:4222", false)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	//repository
	repository := cqrs.NewRepository(eventstore, nat)

	//handlers
	commandHandler := cqrs.NewCommandHandler()
	accountHandler := bank.NewCommandHandler(repository)

	//add commands to commandhandler
	commandHandler.Add(bank.CreateAccount{}, accountHandler)

	//commandbus
	commandBus := async.NewBus(commandHandler, 30)
	return commandBus, nil
}
