package model

import (
	"database/sql"

	"github.com/dberstein/recanatid-go/typ"
)

func GetProfileUser(db *sql.DB, username string) (*typ.RegisterUser, error) {
	var user typ.RegisterUser
	row := db.QueryRow(`SELECT username, email, role FROM users WHERE username=?`, username)
	if err := row.Scan(&user.Username, &user.Email, &user.Role); err != nil {
		return nil, err
	}

	return &user, nil
}
