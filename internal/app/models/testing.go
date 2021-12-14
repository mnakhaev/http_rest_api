package models

import "testing"

// TestUser helper will return already prepared user with valid data for tests
func TestUser(t *testing.T) *User {
	return &User{
		Email:    "user@example.org",
		Password: "password",
	}
}
