package bank

import "cqrs"

//CreateAccount assigned to an owner
type CreateAccount struct {
	cqrs.BaseCommand
	Owner string
}

//PerformDeposit to a given account
type PerformDeposit struct {
	cqrs.BaseCommand
	Ammount int
}

//ChangeOwner of an account
type ChangeOwner struct {
	cqrs.BaseCommand
	Owner string
}

//PerformWithdrawal to a given account
type PerformWithdrawal struct {
	cqrs.BaseCommand
	Ammount int
}
