package database

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func InitDB(connectionString string) (*sql.DB, error) {
    log.Println("Attempting to connect to the database...")

    // Open database
    db, err := sql.Open("pgx", connectionString)
    if err != nil {
        log.Printf("Failed to open database connection: %v", err)
        return nil, err
    }

    // Test connection
    err = db.Ping()
    if err != nil {
        log.Printf("Failed to ping database: %v", err)
        return nil, err
    }

    // Set connection pool settings
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)

    log.Println("Database connected successfully")
    return db, nil
}