package bank

import (
	"cqrs"
	"errors"
)

//ErrBalanceOut when you don't have balance to perform the operation
var ErrBalanceOut = errors.New("balance out")

//Account of bank
type Account struct {
	cqrs.BaseAggregate
	Owner   string
	Balance int
}

//LoadsFromHistory restore the account to last status
func (a *Account) LoadsFromHistory(events []cqrs.Event) {
	for _, event := range events {
		a.BaseAggregate.ApplyChange(a, event, false)
	}
}

//ApplyChange to account
func (a *Account) ApplyChange(event cqrs.Event, commit bool) {
	switch e := event.Data.(type) {
	case *AccountCreated:
		a.Owner = e.Owner
	case *DepositPerformed:
		a.Balance += e.Ammount
	case *WithdrawalPerformed:
		a.Balance -= e.Ammount
	}
}

//Handle a command
func (a *Account) Handle(command cqrs.Command) error {
	event := cqrs.Event{
		AggregateID:   a.ID,
		AggregateType: "Account",
	}

	switch c := command.(type) {
	case CreateAccount:
		event.AggregateID = c.AggregateID
		event.Type = "AccountCreated"
		event.Data = &AccountCreated{"mishudark"}

	case PerformDeposit:
		event.Type = "DepositPerformed"
		event.Data = &DepositPerformed{
			c.Ammount,
		}

	case PerformWithdrawal:
		if a.Balance < c.Ammount {
			return ErrBalanceOut
		}

		event.Type = "WithdrawalPerformed"
		event.Data = &WithdrawalPerformed{
			c.Ammount,
		}
	}

	a.BaseAggregate.ApplyChange(a, event, true)
	return nil
}
