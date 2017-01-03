package main

import (
	"github.com/mishudark/eventhus/examples/bank"
	"github.com/mishudark/eventhus/utils"
	"os"
	"time"
)

func main() {
	end := make(chan bool)
	commandBus, err := config()
	if err != nil {
		os.Exit(1)
	}

	//Create Account
	for i := 0; i < 3; i++ {
		go func() {
			uuid, err := utils.UUID()
			if err != nil {
				return
			}

			//1) Create an account
			var account bank.CreateAccount
			account.AggregateID = uuid
			account.Owner = "mishudark"

			commandBus.HandleCommand(account)

			time.Sleep(time.Millisecond * 100)

			//2) Perform a deposit
			deposit := bank.PerformDeposit{
				Ammount: 300,
			}

			deposit.AggregateID = uuid
			deposit.Version = 1

			commandBus.HandleCommand(deposit)

			//3) Perform a withdrawl
			time.Sleep(time.Millisecond * 100)
			withdrawl := bank.PerformWithdrawal{
				Ammount: 249,
			}

			withdrawl.AggregateID = uuid
			withdrawl.Version = 2

			commandBus.HandleCommand(withdrawl)
		}()
	}
	<-end
}
