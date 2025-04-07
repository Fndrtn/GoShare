package storage

import (
	"database/sql"

	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() {
	var err error

	DB, err = sql.Open("sqlite", "users.db")

	if err != nil {
		log.Fatal("Cannot open SQL database: ", err)
	}

	createTable := ` 
		CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL
	);`

	_, err = DB.Exec(createTable)
	if err != nil {
		log.Fatal("Cannot create table: ", err)
	}

	fmt.Println("Connected to SQLite database successfully!")
}
