package main

import (
	"math/rand"
	"time"
)

type TransferRequest struct {
	ToAccountID int64 `json:"toAccountID"`
	Amount      int64 `json:"amount"`
}

type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type Account struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Number    int64     `json:"number"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewAccount(firstName, lastName string) *Account {
	// create a new account
	return &Account{
		// ID:        rand.Intn(10_000),
		FirstName: firstName,
		LastName:  lastName,
		Number:    int64(rand.Intn(1_000_0000)),
		Balance:   0, // go will do this implicitly
		CreatedAt: time.Now().UTC(),
	}

}
