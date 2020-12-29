package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

//Account is a alias for string
//May be it will get replaced with a struct later
type Account string

//Tx is a structure used to hold the transaction information
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

//State is used to hold the current state in relation with transactions
type State struct {
	Balances  map[Account]uint
	txMempool []Tx
}

func (state *State) add(tx Tx) error {
	if err := state.apply(tx); err != nil {
		return err
	}

	state.txMempool = append(state.txMempool, tx)

	return nil
}

func (state *State) persist() error {
	length := len(state.txMempool)
	for i := 0; i < length; i++ {
		txJSON, err := json.Marshal(state.txMempool)
		if err != nil {
			return err
		}

		cwd, err := os.Getwd()
		filePath := filepath.Join(cwd, "database", "tx.db")
		dbFile, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			return err
		}
		if _, err := dbFile.Write(append(txJSON, '\n')); err != nil {
			return err
		}

		state.txMempool = append(state.txMempool[:i], state.txMempool[i+1:]...)
	}
	return nil
}

func (state *State) apply(tx Tx) error {
	if tx.IsReward() {
		state.Balances[tx.To] += tx.Value
	}

	if tx.Value > state.Balances[tx.From] {
		return errors.New("Insufficient Fund")
	}

	state.Balances[tx.From] -= tx.Value
	state.Balances[tx.To] += tx.Value

	return nil
}

func loadGenesis(genPath string) (map[Account]uint, error) {
	return nil, errors.New("error loading genesis")
}

//NewStateFromDisk loads the State from the local genesis file and updates the state based on transactions
func NewStateFromDisk() (*State, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	genpath := filepath.Join(cwd, "database", "genesis.json")
	gen, err := loadGenesis(genpath)
	if err != nil {
		return nil, err
	}

	balances := make(map[Account]uint)
	for account, balance := range gen {
		balances[account] = balance
	}

	cwd, err = os.Getwd()
	if err != nil {
		return nil, err
	}

	transactionPath := filepath.Join(cwd, "database", "tx.db")

	transactionFile, err := os.OpenFile(transactionPath, os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(transactionFile)
	state := &State{balances, make([]Tx, 0)}

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		var tx Tx
		json.Unmarshal(scanner.Bytes(), &tx)

		if err := state.apply(tx); err != nil {
			return nil, err
		}
	}
	return state, nil
}

func main() {
	//TODO:generate and distributeclear blockchain token
	//TODO:develop CLI controlled DB
	//make the DB immutable using a secure cryptographic hash function
}
