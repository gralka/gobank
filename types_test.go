package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAccount(t *testing.T) {
  acc, err := NewAccount("John", "Doe", "password")

  assert.Nil(t, err)

  assert.Equal(t, "John", acc.FirstName)
  assert.Equal(t, "Doe", acc.LastName)
  assert.NotEqual(t, "password", acc.EncryptedPassword)
}
