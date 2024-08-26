package db

import (
	"database/sql"
	"log"
)

func ensureSchema(db *sql.DB) error {
	if _, err := db.Exec(
		"CREATE TABLE IF NOT EXISTS users (username TEXT PRIMARY KEY, pwhash TEXT, email TEXT, role TEXT)",
	); err != nil {
		return err
	}
	return nil
}

func GetDb(dsn string) *sql.DB {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	err = ensureSchema(db)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
