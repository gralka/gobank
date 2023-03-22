package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
  seed := flag.Bool("seed", false, "seed the db");

  flag.Parse()

  if err := loadDotEnv(); err != nil {
    log.Fatal("Error loading .env file")
  }

  store, err := NewPostgressStore()

  if err != nil { log.Fatal(err) }
  if err := store.Init(); err != nil { log.Fatal(err) }

  if *seed {
    fmt.Println("Seeding the accounts table...")

    seedAccounts(store)
  }

  server := NewAPIServer(":3000", store)

  server.Run()
}

func loadDotEnv() error {
  if prod := os.Getenv("PROD"); prod != "true" {
    err := godotenv.Load();
    if err != nil {
      log.Println("Error loading .env file")
      return err;
    }
  }

  return nil
}

func seedAccounts(store Storage) {
  seedAccount(store, "John", "Doe", "password")
  seedAccount(store, "Jane", "Doe", "password")
}

func seedAccount(store Storage, fname string, lname string, pw string) *Account {
  acc, err := NewAccount(fname, lname, pw)

  if err != nil { log.Fatal(err) }

  newAccount, err := store.CreateAccount(acc)

  if err != nil { log.Fatal(err) }

  return newAccount 
}
