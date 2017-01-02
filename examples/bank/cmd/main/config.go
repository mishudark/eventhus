package main

import (
	"github.com/mishudark/eventhus"
	"github.com/mishudark/eventhus/commandbus"
	"github.com/mishudark/eventhus/eventbus/nats"
	"github.com/mishudark/eventhus/eventstore/mongo"
	"github.com/mishudark/eventhus/examples/bank"
	"log"
)

func config() (eventhus.CommandBus, error) {
	//register events
	reg := eventhus.NewEventRegister()
	reg.Set(bank.AccountCreated{})
	reg.Set(bank.DepositPerformed{})
	reg.Set(bank.WithdrawalPerformed{})

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
	repository := eventhus.NewRepository(eventstore, nat)

	//handlers
	commandHandler := eventhus.NewCommandHandler()
	accountHandler := bank.NewCommandHandler(repository)

	//add commands to commandhandler
	commandHandler.Add(bank.CreateAccount{}, accountHandler)
	commandHandler.Add(bank.PerformDeposit{}, accountHandler)
	commandHandler.Add(bank.PerformWithdrawal{}, accountHandler)

	//commandbus
	commandBus := async.NewBus(commandHandler, 30)
	return commandBus, nil
}
