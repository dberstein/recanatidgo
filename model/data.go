package model

import (
	"database/sql"

	"github.com/dberstein/recanatid-go/typ"
)

type Data struct {
	db       *sql.DB
	pageSize int
}

func NewData(db *sql.DB, pageSize int) *Data {
	return &Data{
		db:       db,
		pageSize: pageSize,
	}
}

func (d *Data) ListUsers(page int) ([]typ.RegisterUser, error) {
	rows, err := d.db.Query("SELECT username, email, role FROM users LIMIT ?,?", (page-1)*d.pageSize, d.pageSize)
	if err != nil {
		return nil, err
	}

	persons := []typ.RegisterUser{}
	for rows.Next() {
		var person typ.RegisterUser
		if err := rows.Scan(&person.Username, &person.Email, &person.Role); err != nil {
			return nil, err
		}
		persons = append(persons, person)
	}

	return persons, nil
}
