package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
  err := godotenv.Load()

  if err != nil { log.Fatal("Error loading .env file") }

  store, err := NewPostgressStore()

  if err != nil { log.Fatal(err) }
  if err := store.Init(); err != nil { log.Fatal(err) }

  server := NewAPIServer(":3000", store)

  server.Run()
}

func loadDotEnv() error {
  if prod := os.Getenv("ENVIRONEMNT"); prod != "true" {
    err := godotenv.Load();
    if err != nil {
      log.Println("Error loading .env file")
      // return err;
    }
  }

  return nil
}

