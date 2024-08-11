package main

import (
	"database/sql"
	"sync"
)

type User struct {
	mu       sync.Mutex
	username string
	pwhash   string
	email    string
	db       *sql.DB
}

func (u *User) SetPassword(password string) error {
	hash, err := HashPassword(password)
	if err != nil {
		return err
	}
	u.pwhash = hash
	return nil
}
