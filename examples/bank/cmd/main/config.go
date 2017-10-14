package main

import (
	"github.com/mishudark/eventhus"
	"github.com/mishudark/eventhus/commandhandler/basic"
	"github.com/mishudark/eventhus/config"
	"github.com/mishudark/eventhus/examples/bank"
)

func getConfig() (eventhus.CommandBus, error) {
	//register events
	reg := eventhus.NewEventRegister()
	reg.Set(bank.AccountCreated{})
	reg.Set(bank.DepositPerformed{})
	reg.Set(bank.WithdrawalPerformed{})

	//eventbus
	//rabbit, err := rabbitmq.NewClient("guest", "guest", "localhost", 5672)

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
