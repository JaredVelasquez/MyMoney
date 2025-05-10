#!/bin/bash

# Script para configurar el entorno de desarrollo

# Colores para output
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Verificar instalación de Go
echo -e "${YELLOW}Verificando instalación de Go...${NC}"
if ! command -v go &> /dev/null; then
    echo -e "${RED}Go no está instalado. Por favor, instálalo desde https://golang.org/dl/${NC}"
    exit 1
fi
echo -e "${GREEN}Go está instalado correctamente.${NC}"

# Verificar versión de Go
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo -e "${YELLOW}Versión de Go: ${GO_VERSION}${NC}"
# Comparar versiones (esta es una verificación simple)
REQUIRED_VERSION="1.18"
if [[ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_VERSION" ]]; then
    echo -e "${RED}Se requiere Go ${REQUIRED_VERSION} o superior.${NC}"
    exit 1
fi

# Instalar herramientas de desarrollo
echo -e "${YELLOW}Instalando herramientas de desarrollo...${NC}"

# golangci-lint para linting
echo -e "${YELLOW}Instalando golangci-lint...${NC}"
if ! command -v golangci-lint &> /dev/null; then
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    echo -e "${GREEN}golangci-lint instalado correctamente.${NC}"
else
    echo -e "${GREEN}golangci-lint ya está instalado.${NC}"
fi

# swag para Swagger
echo -e "${YELLOW}Instalando swag para Swagger...${NC}"
if ! command -v swag &> /dev/null; then
    go install github.com/swaggo/swag/cmd/swag@latest
    echo -e "${GREEN}swag instalado correctamente.${NC}"
else
    echo -e "${GREEN}swag ya está instalado.${NC}"
fi

# Verificar si existe el archivo .env
echo -e "${YELLOW}Verificando archivo .env...${NC}"
if [ ! -f .env ]; then
    echo -e "${YELLOW}Creando archivo .env a partir de .env.example...${NC}"
    if [ -f .env.example ]; then
        cp .env.example .env
        echo -e "${GREEN}Archivo .env creado. Por favor, edítalo con tus configuraciones.${NC}"
    else
        echo -e "${YELLOW}Creando archivo .env básico...${NC}"
        cat > .env << EOF
# Server configuration
PORT=8080
GIN_MODE=debug

# JWT configuration
JWT_SECRET=your_jwt_secret_here

# Database configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=mi_app_db
DB_SSL_MODE=disable
EOF
        echo -e "${GREEN}Archivo .env básico creado. Por favor, edítalo con tus configuraciones.${NC}"
    fi
else
    echo -e "${GREEN}El archivo .env ya existe.${NC}"
fi

# Descargar dependencias
echo -e "${YELLOW}Descargando dependencias...${NC}"
go mod download
echo -e "${GREEN}Dependencias descargadas correctamente.${NC}"

# Verificar si Docker está instalado
echo -e "${YELLOW}Verificando Docker...${NC}"
if command -v docker &> /dev/null; then
    echo -e "${GREEN}Docker está instalado. Puedes usar Docker para PostgreSQL:${NC}"
    echo -e "${YELLOW}docker run --name mi-app-postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=mi_app_db -p 5432:5432 -d postgres:13${NC}"
else
    echo -e "${YELLOW}Docker no está instalado. Si deseas usarlo, instálalo desde https://docs.docker.com/get-docker/${NC}"
fi

echo -e "${GREEN}¡Configuración completada! Puedes iniciar el desarrollo con:${NC}"
echo -e "${YELLOW}make run${NC}" 