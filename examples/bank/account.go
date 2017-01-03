package bank

import (
	"errors"
	"github.com/mishudark/eventhus"
)

//ErrBalanceOut when you don't have balance to perform the operation
var ErrBalanceOut = errors.New("balance out")

//Account of bank
type Account struct {
	eventhus.BaseAggregate
	Owner   string
	Balance int
}

//ApplyChange to account
func (a *Account) ApplyChange(event eventhus.Event) {
	switch e := event.Data.(type) {
	case *AccountCreated:
		a.Owner = e.Owner
		a.ID = event.AggregateID
	case *DepositPerformed:
		a.Balance += e.Ammount
	case *WithdrawalPerformed:
		a.Balance -= e.Ammount
	}
}

//Handle a command
func (a *Account) Handle(command eventhus.Command) error {
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

	case PerformWithdrawal:
		if a.Balance < c.Ammount {
			return ErrBalanceOut
		}

		event.Data = &WithdrawalPerformed{
			c.Ammount,
		}
	}

	a.BaseAggregate.ApplyChangeHelper(a, event, true)
	return nil
}
