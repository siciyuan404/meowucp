package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	dsn := "host=localhost port=5432 user=postgres password=postgres sslmode=disable dbname=meowucp"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Database connected successfully")

	_, err = db.Exec("CREATE EXTENSION IF NOT EXISTS pg_trgm")
	if err != nil {
		log.Fatalf("Failed to create extension: %v", err)
	}

	log.Println("Extension 'pg_trgm' created successfully")
}
