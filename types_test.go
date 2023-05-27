package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAccount(t *testing.T) {
	acc, err := NewAccount("John", "Doe", "password123")

	assert.NoError(t, err)

	assert.Equal(t, "John", acc.FirstName)
}
