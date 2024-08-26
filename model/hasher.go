package model

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type Hasher struct{}

func NewHasher() *Hasher {
	return &Hasher{}
}

func (h *Hasher) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (h *Hasher) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (h *Hasher) GetPwhash(db *sql.DB, username string) (string, error) {
	var pwhash string

	row := db.QueryRow(`SELECT pwhash FROM users WHERE username=?`, username)
	err := row.Scan(&pwhash)
	if err != nil {
		return pwhash, err
	}

	return pwhash, nil
}
