package db

import (
	"database/sql"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Init() error {
	var err error
	DB, err = sql.Open("sqlite", "file:data.db?_foreign_keys=1")
	if err != nil {
		return err
	}
	return createTables()
}

func createTables() error {
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS gcp_connections (
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			name        TEXT NOT NULL,
			bucket      TEXT NOT NULL,
			credentials TEXT NOT NULL,
			created_at  DATETIME NOT NULL
		)`)
	if err != nil {
		return err
	}
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS aws_connections (
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			name        TEXT NOT NULL,
			bucket      TEXT NOT NULL,
			credentials TEXT NOT NULL,
			created_at  DATETIME NOT NULL
		)`)
	if err != nil {
		return err
	}
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS huawei_connections (
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			name        TEXT NOT NULL,
			bucket      TEXT NOT NULL,
			credentials TEXT NOT NULL,
			created_at  DATETIME NOT NULL
		)`)
	return err
}
