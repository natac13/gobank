package main

import "math/rand"

type Account struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Number    int64  `json:"number"`
	Balance   int64  `json:"balance"`
}

func NewAccount(firstName, lastName string) *Account {
	// create a new account
	return &Account{
		ID:        rand.Intn(10_000),
		FirstName: firstName,
		LastName:  lastName,
		Number:    int64(rand.Intn(1_000_0000)),
		Balance:   0, // go will do this implicitly
	}

}
