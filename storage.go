package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)
 
type Storage interface {
  CreateAccount(account *Account) error
  DeleteAccount(id int) error
  UpdateAccount(account *Account) error
  GetAccounts() ([]*Account, error)
  GetAccountByID(id int) (*Account, error)
}

type PostgressStore struct {
  db *sql.DB
}

func (s *PostgressStore) Init() error {
  if err := s.createAccountTable(); err != nil {
    return err
  }

  return nil
}

func NewPostgressStore() (*PostgressStore, error) {
  connStr := "user=postgres dbname=postgres password=gobank sslmode=disable"
  db, err := sql.Open("postgres", connStr)

  if err != nil {
    return nil, err
  }

  if err := db.Ping(); err != nil {
    return nil, err
  }

  return &PostgressStore{
    db: db,
  }, nil
} 

func (s *PostgressStore) createAccountTable() error {
  query := `CREATE TABLE IF NOT EXISTS accounts (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    number SERIAL,
    balance INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
  )`

  _, err := s.db.Exec(query)

  return err
}

func (s *PostgressStore) CreateAccount(a *Account) error {
  query := `INSERT INTO accounts (
      first_name,
      last_name,
      balance,
      created_at
    ) VALUES ($1, $2, $3, $4)`

  resp, err := s.db.Query(
     query,
     a.FirstName,
     a.LastName,
     a.Balance,
     a.CreatedAt,
  )

  if err != nil {
    return err
  }

  fmt.Printf("resp: %v+\n", resp)
  return nil
}

func (s *PostgressStore) DeleteAccount(id int) error {
  // query := "DELETE FROM accounts WHERE id = $1"
  // _, err := s.db.Exec(query, id)
  return nil 
}

func (s *PostgressStore) UpdateAccount(account *Account) error {
  // query := "UPDATE accounts SET name = $1, balance = $2 WHERE id = $3"
  // _, err := s.db.Exec(query, account.Name, account.Balance, account.ID)
  return nil 
}

func (s *PostgressStore) GetAccountByID(id int) (*Account, error) {
  // query := "SELECT id, name, balance FROM accounts WHERE id = $1"
  // row := s.db.QueryRow(query, id)

  account := &Account{}
  // err := row.Scan(&account.ID, &account.Name, &account.Balance)
  // if err != nil {
  //   return nil, err
  // }

  return account, nil
}

func (s *PostgressStore) GetAccounts() ([]*Account, error) {
  rows, err := s.db.Query("SELECT * FROM accounts")
  if err != nil {
    return nil, err
  }

  accounts := []*Account{}

  for rows.Next() {
    account, err := scanIntoAccount(rows)

    if err != nil {
      return nil, err
    }

    accounts = append(accounts, account)
  }

  return accounts, nil
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
    account := new(Account)

    err := rows.Scan(
      &account.ID,
      &account.FirstName,
      &account.LastName,
      &account.Number,
      &account.Balance,
      &account.CreatedAt,
    )

    return account, err 
}
