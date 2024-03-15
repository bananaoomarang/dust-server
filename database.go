package main

import (
    "database/sql"
    "fmt"
    "log"
	"os"

    _ "github.com/lib/pq"
)

var Db *sql.DB

func getDbUrl() string {
	if value, ok := os.LookupEnv("DATABASE_URL"); ok {
		return value
	}

	return "postgres://dustapi:dustapi@localhost:5432/dust_server?sslmode=disable"
}

func ConnectDB() {
    db, err := sql.Open("postgres", getDbUrl())
    if err != nil {
        log.Printf("failed to connect to database: %v", err)
		panic(err)
    }
	Db = db
	fmt.Println("Successfully connected to database!")
}
