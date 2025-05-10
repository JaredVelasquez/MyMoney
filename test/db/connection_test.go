package db

import (
	"testing"

	"mi-app-backend/db/config"
)

func TestDatabaseConnection(t *testing.T) {
	// Intentar crear una conexión a la base de datos
	conn, err := config.NewConnection()
	if err != nil {
		t.Fatalf("Error connecting to database: %v", err)
	}
	defer conn.Close()

	// Verificar que la conexión funcione con una consulta simple
	db := conn.GetDB()
	var result int
	err = db.QueryRow("SELECT 1").Scan(&result)
	if err != nil {
		t.Fatalf("Error executing test query: %v", err)
	}

	if result != 1 {
		t.Errorf("Expected result to be 1, got %d", result)
	}

	t.Log("Database connection test passed")
}
