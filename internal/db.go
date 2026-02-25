package internal

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"

	"awesomeProject/pkg/storage"
)

var DB *sql.DB

func InitDB() error {
	dbDir := filepath.Join(os.ExpandEnv("$HOME"), ".mcpclaw")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return fmt.Errorf("failed to create db directory: %w", err)
	}

	dbPath := filepath.Join(dbDir, "data.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	DB = db
	storage.SetDB(db)

	if err := createTables(); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	return nil
}

func createTables() error {
	schema := `
	CREATE TABLE IF NOT EXISTS memories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		category TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := DB.Exec(schema)
	return err
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
