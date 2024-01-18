package tests

import (
	"backend/util"
	"testing"
)

func TestPassword(t *testing.T) {
	password := "password123"

	hashed, err := util.HashPassword(password)
	if err != nil {
		t.Fatal(err)
	}

	if hashed == "" {
		t.Fatal("Hashed password is unexpectedly empty")
	}

	err2 := util.CheckPassword(password, hashed)
	if err2 != nil {
		t.Fatal(err)
	}
}
