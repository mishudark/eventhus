package main

import (
	"cqrs/examples/bank"
	"cqrs/utils"
	"os"
)

func main() {
	end := make(chan bool)
	commandBus, err := config()
	if err != nil {
		os.Exit(1)
	}

	//Create Account
	for i := 0; i < 3000; i++ {
		go func() {
			var account bank.CreateAccount
			uuid, err := utils.UUID()
			if err != nil {
				return
			}
			account.AggregateID = uuid
			commandBus.HandleCommand(account)
		}()
	}
	<-end
}
