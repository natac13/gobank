package main

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type LoginResponse struct {
	Token  string `json:"token"`
	Number int64  `json:"number"`
}

type LoginRequest struct {
	Number   int64  `json:"number"`
	Password string `json:"password"`
}

type TransferRequest struct {
	ToAccountID int64 `json:"toAccountID"`
	Amount      int64 `json:"amount"`
}

type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
}

type Account struct {
	ID                int       `json:"id"`
	FirstName         string    `json:"firstName"`
	LastName          string    `json:"lastName"`
	Number            int64     `json:"number"`
	EncryptedPassword string    `json:"-"`
	Balance           int64     `json:"balance"`
	CreatedAt         time.Time `json:"createdAt"`
}

func (a *Account) ValidatePassword(password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(a.EncryptedPassword), []byte(password))

	if err != nil {
		return false, err
	}

	return true, nil
}

func NewAccount(firstName, lastName string, password string) (*Account, error) {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}
	// create a new account
	return &Account{
		// ID:        rand.Intn(10_000),
		FirstName:         firstName,
		LastName:          lastName,
		Number:            int64(rand.Intn(1_000_0000)),
		Balance:           0, // go will do this implicitly
		CreatedAt:         time.Now().UTC(),
		EncryptedPassword: string(encryptedPassword),
	}, nil

}
