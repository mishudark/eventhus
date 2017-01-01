package bank

//AccountCreated event
type AccountCreated struct {
	Owner string
}

//DepositPerformed event
type DepositPerformed struct {
	Ammount int
}

//OwnerChanged event
type OwnerChanged struct {
	Owner string
}

//WithdrawalPerformed event
type WithdrawalPerformed struct {
	Ammount int
}
