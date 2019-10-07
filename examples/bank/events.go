package bank

//AccountCreated event
type AccountCreated struct {
	Owner string `json:"owner"`
}

//DepositPerformed event
type DepositPerformed struct {
	Amount int `json:"amount"`
}

//OwnerChanged event
type OwnerChanged struct {
	Owner string `json:"owner"`
}

//WithdrawalPerformed event
type WithdrawalPerformed struct {
	Amount int `json:"amount"`
}
