package main

//Account is a alias for string
//May be it will get replaced with a struct later
type Account string

//Tx is structure used to hold the transaction information
type Tx struct {
	From  Account `json:"from"`
	To    Account `json:"to"`
	Value uint    `json:"value"`
	Data  string  `json:"data"`
}

//IsReward is used to check whether a transaction is a reward
func (t Tx) IsReward() bool {
	return t.Data == "reward"
}

func main() {
	//TODO:generate and distributeclear blockchain token
	//TODO:develop CLI controlled DB
	//make the DB immutable using a secure cryptographic hash function
}
