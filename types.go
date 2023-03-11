package main

import "math/rand"

type Account struct {
  ID        int    `json:"id"`
  FirstName string `json:"first_name"`
  LastName  string `json:"last_name"`
  Number    int64  `json:"number"`
  Balance   int64  `json:"balance"`
}

func NewAccount (firstName string, lastName string) *Account {
  return &Account {
    ID: rand.Intn(100000),
    FirstName: firstName,
    LastName: lastName,
    Number: int64(rand.Intn(1000000)),
    Balance: 0, // We don't need to specify 0, but I think it's good practice.
  }
}

