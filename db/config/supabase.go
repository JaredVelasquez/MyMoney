package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// SupabaseClient represents a connection to Supabase PostgreSQL
type SupabaseClient struct {
	db *sql.DB
}

// NewSupabaseClient creates a new connection to Supabase
func NewSupabaseClient() (*SupabaseClient, error) {
	// Get Supabase connection parameters from environment variables
	supabaseHost := os.Getenv("SUPABASE_HOST")
	supabasePort := os.Getenv("SUPABASE_PORT")
	supabaseUser := os.Getenv("SUPABASE_USER")
	supabasePassword := os.Getenv("SUPABASE_PASSWORD")
	supabaseDbName := os.Getenv("SUPABASE_DBNAME")

	// Check required parameters
	if supabasePassword == "" {
		return nil, fmt.Errorf("SUPABASE_PASSWORD environment variable must be set")
	}

	// Set defaults if not provided
	if supabaseHost == "" {
		// Use the shared connection pooler for IPv4 compatibility
		supabaseHost = "aws-0-us-east-2.pooler.supabase.com"
	}
	if supabasePort == "" {
		supabasePort = "6543" // Transaction pooler port (IPv4 compatible)
	}
	if supabaseUser == "" {
		supabaseUser = "postgres.eddrisdxafllaekoaduy"
	}
	if supabaseDbName == "" {
		supabaseDbName = "postgres"
	}

	log.Printf("Connecting to Supabase pooler at: %s", supabaseHost)

	// Build connection string for transaction pooler
	// Use the format for Supabase transaction pooler
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=require",
		supabaseUser,
		supabasePassword,
		supabaseHost,
		supabasePort,
		supabaseDbName)

	// Log connection string (hide password)
	connectionLog := fmt.Sprintf("postgres://%s:****@%s:%s/%s?sslmode=require",
		supabaseUser, supabaseHost, supabasePort, supabaseDbName)
	log.Printf("Using connection string: %s", connectionLog)

	// Open connection to database with specific options for Supabase
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open Supabase connection: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	// Test connection with a simple query
	var result int
	err = db.QueryRow("SELECT 1").Scan(&result)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to execute test query: %w", err)
	}

	log.Println("Successfully connected to Supabase PostgreSQL database")
	return &SupabaseClient{db: db}, nil
}

// GetDB returns the database connection
func (c *SupabaseClient) GetDB() *sql.DB {
	return c.db
}

// Close closes the database connection
func (c *SupabaseClient) Close() error {
	return c.db.Close()
}

// ExecuteQuery executes a simple query to test the connection
func (c *SupabaseClient) ExecuteQuery() (int, error) {
	var result int
	err := c.db.QueryRow("SELECT 1").Scan(&result)
	if err != nil {
		return 0, fmt.Errorf("failed to execute test query: %w", err)
	}
	return result, nil
}
