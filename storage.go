package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)
 
type Storage interface {
  CreateAccount(account *Account) (*Account, error)
  DeleteAccount(id int) error
  UpdateAccount(account *Account) error
  GetAccounts() ([]*Account, error)
  GetAccountByID(id int) (*Account, error)
  GetAccountByNumber(number int) (*Account, error)
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
  driverEnvVar := os.Getenv("DB_DRIVER")

  if driverEnvVar == "" {
    return nil, fmt.Errorf("DB_DRIVER not set")
  }

  hostEnvVar := os.Getenv("DB_HOST")

  if hostEnvVar == "" {
    return nil, fmt.Errorf("DB_HOST not set")
  }

  userEnvVar := os.Getenv("DB_USER")

  if userEnvVar == "" {
    return nil, fmt.Errorf("DB_USER not set")
  }

  passwordEnvVar := os.Getenv("DB_PASSWORD")

  if passwordEnvVar == "" {
    return nil, fmt.Errorf("DB_PASSWORD not set")
  }

  dbNameEnvVar := os.Getenv("DB_NAME")

  if dbNameEnvVar == "" {
    return nil, fmt.Errorf("DB_NAME not set")
  }

  sslModeEnvVar := os.Getenv("DB_SSLMODE")

  if sslModeEnvVar == "" {
    return nil, fmt.Errorf("DB_SSLMODE not set")
  }

  connStr := fmt.Sprintf(
    "host=%s user=%s password=%s dbname=%s sslmode=%s",
    hostEnvVar,
    userEnvVar,
    passwordEnvVar,
    dbNameEnvVar,
    sslModeEnvVar,
  )

  db, err := sql.Open(driverEnvVar, connStr)

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
    encrypted_password VARCHAR(100),
    number SERIAL,
    balance INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
  )`

  _, err := s.db.Exec(query)

  return err
}

func (s *PostgressStore) CreateAccount(a *Account) (*Account, error) {
  query := `INSERT INTO accounts (
      first_name,
      last_name,
      encrypted_password,
      balance,
      created_at
    ) VALUES ($1, $2, $3, $4, $5) RETURNING id, number`

  var id int64
  var number int64

  err := s.db.QueryRow(
     query,
     a.FirstName,
     a.LastName,
     a.EncryptedPassword,
     a.Balance,
     a.CreatedAt,
  ).Scan(&id, &number)

  if err != nil {
    return &Account{}, err
  }

  a.ID = int(id)
  a.Number = int(number)

  return a, nil
}

func (s *PostgressStore) DeleteAccount(id int) error {
  _, err := s.db.Query("DELETE FROM accounts WHERE id = $1", id)

  return err 
}

func (s *PostgressStore) UpdateAccount(account *Account) error {
  // query := "UPDATE accounts SET name = $1, balance = $2 WHERE id = $3"
  // _, err := s.db.Exec(query, account.Name, account.Balance, account.ID)
  return nil 
}

func (s *PostgressStore) GetAccountByNumber(number int) (*Account, error) {
  rows, err := s.db.Query("SELECT * FROM accounts WHERE number = $1", number)

  if err != nil {
    return nil, err
  }

  for rows.Next() {
    return scanIntoAccount(rows)
  }

  return nil, fmt.Errorf("account with number %d not found", number)
}

func (s *PostgressStore) GetAccountByID(id int) (*Account, error) {
  rows, err := s.db.Query("SELECT * FROM accounts WHERE id = $1", id)

  if err != nil {
    return nil, err
  }

  for rows.Next() {
    return scanIntoAccount(rows)
  }

  return nil, fmt.Errorf("account with id %d not found", id)
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
      &account.EncryptedPassword,
      &account.Number,
      &account.Balance,
      &account.CreatedAt,
    )

    return account, err 
}
