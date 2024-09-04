package store

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Storage interface {
	GetDB() *sql.DB
}

type store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *store {
	return &store{db: db}
}

func (s *store) GetDB() *sql.DB {
	return s.db
}
