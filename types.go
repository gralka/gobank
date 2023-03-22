package main

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Account struct {
  ID                int       `json:"id"`
  FirstName         string    `json:"first_name"`
  LastName          string    `json:"last_name"`
  Number            int       `json:"number"`
  EncryptedPassword string    `json:"encrypted_password"`
  Balance           int64     `json:"balance"`
  CreatedAt         time.Time `json:"created_at"`
}

type LoginResponse struct {
  Token string `json:"token"`
}

func NewAccount (
  firstName string,
  lastName string,
  password string,
) (*Account, error) {
  encryptedPassword, err := bcrypt.GenerateFromPassword(
    []byte(password),
    bcrypt.DefaultCost,
  )

  if err != nil {
    return nil, err
  }

  return &Account {
    FirstName: firstName,
    LastName:  lastName,
    EncryptedPassword: string(encryptedPassword),
    CreatedAt: time.Now().UTC(),
  }, nil
}

type CreateAccountRequest struct {
  FirstName string `json:"first_name"`
  LastName  string `json:"last_name"`
  Password  string `json:"password"`
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
