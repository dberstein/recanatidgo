package models

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type hasher struct{}

func NewHasher() *hasher {
	return &hasher{}
}

func (h *hasher) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (h *hasher) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (h *hasher) GetPwhashRole(db *sql.DB, username string) (string, string, error) {
	var pwhash string
	var role string

	row := db.QueryRow(`SELECT pwhash, role FROM users WHERE username=?`, username)
	err := row.Scan(&pwhash, &role)
	if err != nil {
		return pwhash, role, err
	}

	return pwhash, role, nil
}
