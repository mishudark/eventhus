package main

import (
	"cqrs/examples/bank"
	"math/rand"
	"os"
	"time"

	"github.com/oklog/ulid"
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
			t := time.Now()
			entropy := rand.New(rand.NewSource(t.UnixNano()))
			uuid := ulid.MustNew(ulid.Timestamp(t), entropy).String()
			//uuid := "0000XSNJG0MQJHBF4QX1EFD6Y1"
			account.AggregateID = uuid
			commandBus.HandleCommand(account)
		}()
	}
	<-end
}
