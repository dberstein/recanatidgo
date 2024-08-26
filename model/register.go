package model

import (
	"database/sql"
	"errors"

	"github.com/dberstein/recanatid-go/typ"
)

type Register struct {
	db *sql.DB
}

func NewRegister(db *sql.DB) *Register {
	return &Register{
		db: db,
	}
}

func (r *Register) Validate(user *typ.RegisterUser) error {
	if user.Username == "" {
		return errors.New("missing: username")
	}
	if user.Password == "" {
		return errors.New("missing: password")
	}
	if user.Email == "" {
		return errors.New("missing: email")
	}
	// if user.Role == "" {
	// 	return errors.New("missing: role")
	// }

	return nil
}

func (r *Register) Insert(user *typ.RegisterUser, pwhash string) error {
	_, err := r.db.Exec(
		"INSERT INTO users (username, pwhash, email, role) VALUES (?, ?, ?, ?)",
		&user.Username, &pwhash, &user.Email, &user.Role,
	)

	return err
}
