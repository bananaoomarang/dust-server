package main

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/lib/pq"
)

var Db *sql.DB

const (
    host     = "localhost"
    port     = 5432
    user     = "dustapi"
    password = "dustapi"
    dbname   = "dust"
)

func ConnectDB() {
    connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname = %s sslmode=disable", host, port, user, password, dbname)
    db, err := sql.Open("postgres", connString)
    if err != nil {
        log.Printf("failed to connect to database: %v", err)
		panic(err)
    }
	Db = db
	fmt.Println("Successfully connected to database!")
}
