package api

import (
	"testing"
	"time"
)

func TestCsrfTokenValid(t *testing.T) {
	s := State{
		csrfJwtSecret: []byte("secret"),
		csrfJwtAlgo:   "HS256",
		csrfJwtTTL:    10,
	}
	token, err := s.genCsrfJwtToken([]string{"test"}, []string{"127.0.0.1"})
	if err != nil {
		t.Error(err)
	}
	valid, err := s.validateCsrfJwtToken(token, "test", "127.0.0.1")
	if err != nil {
		t.Error(err)
	}
	if !valid {
		t.Error("token should be valid")
	}
}

func TestCsrfTokenWrongIP(t *testing.T) {
	s := State{
		csrfJwtSecret: []byte("secret"),
		csrfJwtAlgo:   "HS256",
		csrfJwtTTL:    10,
	}
	token, err := s.genCsrfJwtToken([]string{"test"}, []string{"127.0.0.1"})
	if err != nil {
		t.Error(err)
	}
	valid, err := s.validateCsrfJwtToken(token, "test", "127.0.0.2")
	if err == nil {
		t.Error("token should be invalid")
	}
	if valid {
		t.Error("token should be invalid")
	}
}

func TestCsrfExpire(t *testing.T) {
	s := State{
		csrfJwtSecret: []byte("secret"),
		csrfJwtAlgo:   "HS256",
		csrfJwtTTL:    1,
	}
	token, err := s.genCsrfJwtToken([]string{"test"}, []string{"127.0.0.1"})
	if err != nil {
		t.Error(err)
	}
	time.Sleep(2 * time.Second)
	_, err = s.validateCsrfJwtToken(token, "test", "127.0.0.1")
	if err == nil {
		t.Error("token should be expired")
	}
}

func TestCsrfWrongAud(t *testing.T) {
	s := State{
		csrfJwtSecret: []byte("secret"),
		csrfJwtAlgo:   "HS256",
		csrfJwtTTL:    10,
	}
	token, err := s.genCsrfJwtToken([]string{"test"}, []string{""})
	if err != nil {
		t.Error(err)
	}
	valid, err := s.validateCsrfJwtToken(token, "wrong", "")
	if err == nil {
		t.Error(err)
	}
	if valid {
		t.Error("token should be invalid")
	}
}
