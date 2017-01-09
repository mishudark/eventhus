package bank

//AccountCreated event
type AccountCreated struct {
	Owner string `json:"owner"`
}

//DepositPerformed event
type DepositPerformed struct {
	Ammount int `json:"ammount"`
}

//OwnerChanged event
type OwnerChanged struct {
	Owner string `json:"owner"`
}

//WithdrawalPerformed event
type WithdrawalPerformed struct {
	Ammount int `json:"ammount"`
}
