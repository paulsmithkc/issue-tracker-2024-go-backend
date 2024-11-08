package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var _db *sql.DB;

func OpenDB() (*sql.DB, error) {
  if (_db == nil) {
    log.Println("Connecting to database")

    connectionString := os.Getenv("PG_CONNECTION_STRING")
    db, err := sql.Open("postgres", connectionString)
    if err != nil {
      return nil, fmt.Errorf("Database.Open: %w", err)
    }

    err = db.Ping()
    if err != nil {
      return nil, fmt.Errorf("Database.Ping: %w", err)
    }

    _db = db;
  }
  return _db, nil;
}
