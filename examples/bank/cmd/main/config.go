package main

import (
	"log"

	"github.com/mishudark/eventhus"
	"github.com/mishudark/eventhus/commandbus"
	"github.com/mishudark/eventhus/commandhandler/basic"
	"github.com/mishudark/eventhus/eventbus/nats"
	"github.com/mishudark/eventhus/eventstore/mongo"
	"github.com/mishudark/eventhus/examples/bank"
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
	commandRegister := eventhus.NewCommandRegister()
	commandHandler := basic.NewCommandHandler(repository, &bank.Account{}, "bank", "account")

	//add commands to commandhandler
	commandRegister.Add(bank.CreateAccount{}, commandHandler)
	commandRegister.Add(bank.PerformDeposit{}, commandHandler)
	commandRegister.Add(bank.PerformWithdrawal{}, commandHandler)

	//commandbus
	commandBus := async.NewBus(commandRegister, 30)
	return commandBus, nil
}
