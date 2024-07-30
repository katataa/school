package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Initialize() error {
	var err error
	DB, err = sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		return err
	}

	// Ensure the dbscripts directory and SQL files exist and are accessible
	if err = executeSQLFile(DB, "dbscripts/001_users_create_table.sql"); err != nil {
		return err
	}

	log.Println("Database initialized successfully.")
	return nil
}

func executeSQLFile(db *sql.DB, filepath string) error {
	script, err := os.ReadFile(filepath)
	if err != nil {
		log.Println("Error reading SQL file:", err)
		return err
	}

	_, err = db.Exec(string(script))
	if err != nil {
		log.Println("Error executing SQL script:", err)
		return err
	}

	return nil
}

func Close() {
	DB.Close()
	log.Println("Database connection closed.")
}
