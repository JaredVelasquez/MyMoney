package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// Connection represents a PostgreSQL connection
type Connection struct {
	db *sql.DB
}

// NewConnection creates a new PostgreSQL connection using Supabase credentials
func NewConnection() (*Connection, error) {
	// Try to use SupabaseClient first
	supabaseClient, err := NewSupabaseClient()
	if err == nil {
		log.Println("Using Supabase client connection")
		return &Connection{db: supabaseClient.GetDB()}, nil
	}

	// If Supabase client fails, fall back to direct connection
	log.Printf("Supabase client connection failed: %v", err)
	log.Println("Attempting direct PostgreSQL connection...")

	// Get connection parameters from environment variables
	host := os.Getenv("SUPABASE_HOST")
	port := os.Getenv("SUPABASE_PORT")
	user := os.Getenv("SUPABASE_USER")
	password := os.Getenv("SUPABASE_PASSWORD")
	dbname := os.Getenv("SUPABASE_DBNAME")
	sslmode := os.Getenv("SUPABASE_SSLMODE")
	apiKey := os.Getenv("SUPABASE_API_KEY")

	// Check required parameters
	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
		return nil, fmt.Errorf("missing required database connection parameters")
	}

	if sslmode == "" {
		sslmode = "require" // Default to require SSL for Supabase
	}

	var connStr string

	// If API Key is provided, use it for connection
	if apiKey != "" {
		// For Supabase with API key
		log.Println("Using Supabase API Key for connection")
		connStr = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&options=apikey=%s",
			user, password, host, port, dbname, sslmode, apiKey)
	} else {
		// Standard PostgreSQL connection
		log.Println("Using standard PostgreSQL connection")
		connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			host, port, user, password, dbname, sslmode)
	}

	// Log connection info (without password/api key for security)
	log.Printf("Connecting to PostgreSQL: host=%s port=%s dbname=%s", host, port, dbname)

	// Open connection to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	// Test connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Connected to PostgreSQL database")
	return &Connection{db: db}, nil
}

// GetDB returns the database connection
func (c *Connection) GetDB() *sql.DB {
	return c.db
}

// Close closes the database connection
func (c *Connection) Close() error {
	return c.db.Close()
}
