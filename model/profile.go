package model

import (
	"database/sql"

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
	row := p.db.QueryRow(`SELECT username, email, role FROM users WHERE username=?`, username)
	if err := row.Scan(&user.Username, &user.Email, &user.Role); err != nil {
		return nil, err
	}

	return &user, nil
}

func (p *profile) Update(db *sql.DB, user *typ.RegisterUser) error {
	if user.Email != "" {
		_, err := p.db.Exec(`UPDATE users SET email = ? WHERE username = ?`, &user.Email, &user.Username)
		if err != nil {
			return err
		}
	}

	if user.Password != "" {
		hash := NewHasher()
		pwhash, err := hash.HashPassword(user.Password)
		if err != nil {
			return err
		}

		_, err = db.Exec(`UPDATE users SET pwhash = ? WHERE username = ?`, &pwhash, &user.Username)
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
