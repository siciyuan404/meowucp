package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	dsn := "host=localhost port=5432 user=postgres password=postgres sslmode=disable dbname=postgres"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to postgres: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping postgres: %v", err)
	}

	log.Println("Connected to postgres successfully")

	_, err = db.Exec("CREATE DATABASE meowucp")
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}

	log.Println("Database 'meowucp' created successfully")
}
