package config

import (
	"database/sql"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
)

var db *sql.DB

// ConnectDB connects to the PostgreSQL database
func ConnectDB() error {
	var err error
	connectionString := os.Getenv("CONNECTION_STRING")
	db, err = sql.Open("pgx", connectionString)
	if err != nil {
		return err
	}
	return nil
}

// DB returns the database connection
func DB() *sql.DB {
	return db
}

// CloseDB closes the database connection
func CloseDB() {
	db.Close()
}
