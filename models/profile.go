package models

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/dberstein/recanatid-go/typ"
)

type profile struct {
	db *sql.DB
}

func NewProfile(db *sql.DB) *profile {
	return &profile{
		db: db,
	}
}

func (p *profile) Get(username string) (*typ.RegisterUser, error) {
	var user typ.RegisterUser
	row := p.db.QueryRow(`SELECT username, email, role, pwhash FROM users WHERE username=?`, username)
	if err := row.Scan(&user.Username, &user.Email, &user.Role, &user.Pwhash); err != nil {
		return nil, err
	}

	return &user, nil
}

func (p *profile) UpdatePassword(db *sql.DB, user *typ.RegisterUser) error {
	// change of password means updated `pwhash`
	hasher := NewHasher()
	pwhash, err := hasher.HashPassword(user.Password)
	if err != nil {
		return err
	}

	_, err = db.Exec(`UPDATE users SET pwhash = ? WHERE username = ?`, &pwhash, &user.Username)
	if err != nil {
		return err
	}

	return nil
}

func (p *profile) Update(db *sql.DB, user *typ.RegisterUser) error {
	if user.Email != "" {
		_, err := p.db.Exec(`UPDATE users SET email = ? WHERE username = ?`, &user.Email, &user.Username)
		if err != nil {
			return err
		}
	}

	if user.Password != "" {
		err := p.UpdatePassword(db, user)
		if err != nil {
			return err
		}
	}

	if user.Role != "" {
		_, err := db.Exec(`UPDATE users SET role = ? WHERE username = ?`, &user.Role, &user.Username)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *profile) Insert(user *typ.RegisterUser, pwhash string) error {
	_, err := p.db.Exec(
		"INSERT INTO users (username, pwhash, email, role) VALUES (?, ?, ?, ?)",
		&user.Username, &pwhash, &user.Email, &user.Role,
	)

	return err
}

func (p *profile) Validate(user *typ.RegisterUser) error {
	if user.Username == "" {
		return errors.New("missing: username")
	}
	if user.Password == "" {
		return errors.New("missing: password")
	}
	if user.Email == "" || strings.IndexRune(user.Email, '@') < 1 {
		return errors.New("missing: email")
	}
	// if user.Role == "" {
	// 	return errors.New("missing: role")
	// }

	return nil
}
