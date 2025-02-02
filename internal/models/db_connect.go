package models

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// Postgres represents a connection to a PostgreSQL database.
type Postgres struct {
	DB *sql.DB // The underlying database connection
}

// OpenDB opens a new connection to the PostgreSQL database using the provided connection string.
// It returns a Postgres object if the connection is successful or an error if the connection fails.
//
// Parameters:
//
//	dsn (string): The connection string for the PostgreSQL database.
//
// Returns:
//
//	*Postgres: A pointer to a Postgres object with the database connection.
//	error: An error if the connection or ping fails.
func OpenDB(dsn string) (*Postgres, error) {
	// Open a new PostgreSQL database connection
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Ping the database to ensure the connection is valid
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	// Return the Postgres object containing the DB connection
	return &Postgres{DB: db}, nil
}

// Close closes the PostgreSQL database connection gracefully.
// If the close operation fails, it logs a fatal error and terminates the program.
func (psql *Postgres) Close() {
	// Attempt to close the database connection
	err := psql.DB.Close()
	if err != nil {
		log.Fatalf("postgresql close failure: %v", err) // Fatal log if closing fails
	}
}
