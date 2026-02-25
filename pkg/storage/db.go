package storage

import "database/sql"

var globalDB *sql.DB

func SetDB(db *sql.DB) {
	globalDB = db
}

func GetDB() *sql.DB {
	return globalDB
}
