package main

import (
	"time"
)

type Account struct {
  ID        int       `json:"id"`
  FirstName string    `json:"first_name"`
  LastName  string    `json:"last_name"`
  Number    int       `json:"number"`
  Balance   int64     `json:"balance"`
  CreatedAt time.Time `json:"created_at"`
}

func NewAccount (firstName string, lastName string) *Account {
  return &Account {
    FirstName: firstName,
    LastName:  lastName,
    CreatedAt: time.Now().UTC(),
  }
}

type CreateAccountRequest struct {
  FirstName string `json:"first_name"`
  LastName  string `json:"last_name"`
}

type TransferRequest struct {
  FromAccountID int `json:"from_account_id"`
  ToAccountID   int `json:"to_account_id"`
  Amount        int `json:"amount"`
}

type LoginRequest struct {
  Number   int64  `json:"number"`
  Password string `json:"password"`
}
