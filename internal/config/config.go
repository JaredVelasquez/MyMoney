package config

import (
	"os"
	"strconv"

	"MyMoneyBackend/internal/domain"
)

// Load carga la configuración desde variables de entorno
func Load() (*domain.Config, error) {
	// Valores por defecto
	serverPort := 30000
	dbHost := "localhost"
	dbPort := 5432
	dbUser := "postgres"
	dbPassword := "postgres"
	dbName := "mi_app"
	dbSSLMode := "disable"

	// Sobreescribir con variables de entorno si existen
	if port := os.Getenv("SERVER_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			serverPort = p
		}
	}

	if host := os.Getenv("DB_HOST"); host != "" {
		dbHost = host
	}

	if port := os.Getenv("DB_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			dbPort = p
		}
	}

	if user := os.Getenv("DB_USER"); user != "" {
		dbUser = user
	}

	if password := os.Getenv("DB_PASSWORD"); password != "" {
		dbPassword = password
	}

	if name := os.Getenv("DB_NAME"); name != "" {
		dbName = name
	}

	if sslMode := os.Getenv("DB_SSL_MODE"); sslMode != "" {
		dbSSLMode = sslMode
	}

	return &domain.Config{
		Server: domain.ServerConfig{
			Port: serverPort,
		},
		Database: domain.DatabaseConfig{
			Host:     dbHost,
			Port:     dbPort,
			User:     dbUser,
			Password: dbPassword,
			Name:     dbName,
			SSLMode:  dbSSLMode,
		},
	}, nil
}
