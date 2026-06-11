package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func connectDB(databaseURL string) *sql.DB {
	var db *sql.DB
	var err error

	for i := 0; i < 20; i++ {
		db, err = sql.Open("postgres", databaseURL)
		if err == nil {
			err = db.Ping()
		}
		if err == nil {
			log.Println("connected to database")
			return db
		}
		log.Printf("waiting for database... (%d/20): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	log.Fatal(fmt.Errorf("could not connect to database: %w", err))
	return nil
}
