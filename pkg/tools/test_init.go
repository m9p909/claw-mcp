package tools

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"

	"awesomeProject/pkg/storage"
)

func initTestDB() error {
	// Create in-memory SQLite database for testing
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return err
	}

	// Create tables
	schema := `
	CREATE TABLE IF NOT EXISTS memories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		category TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	if _, err := db.Exec(schema); err != nil {
		return err
	}

	storage.SetDB(db)
	return nil
}

func init() {
	initTestDB()
}
