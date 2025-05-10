package domain

// Config contiene toda la configuración de la aplicación
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

// ServerConfig contiene la configuración del servidor HTTP
type ServerConfig struct {
	Port int
}

// DatabaseConfig contiene la configuración de la base de datos
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}
