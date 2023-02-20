package db

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestBcrypt(t *testing.T) {
	passsord := "123456"
	hash, err := bcrypt.GenerateFromPassword([]byte(passsord), bcrypt.DefaultCost)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(hash))
}
