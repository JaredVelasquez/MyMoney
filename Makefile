# Makefile para mi-app-backend

# Variables
APP_NAME=mi-app-backend
GO=$(shell which go)
GOFLAGS=-v
BINARY_NAME=main

# Colores para output
COLOR_RESET=\033[0m
COLOR_GREEN=\033[0;32m
COLOR_YELLOW=\033[0;33m
COLOR_BLUE=\033[0;34m

.PHONY: all build clean run test test-coverage lint migrate migrate-status migrate-rollback swagger install help

all: build

help:
	@echo "${COLOR_GREEN}Comandos disponibles:${COLOR_RESET}"
	@echo "  ${COLOR_BLUE}make build${COLOR_RESET}              - Compila la aplicación"
	@echo "  ${COLOR_BLUE}make run${COLOR_RESET}                - Ejecuta la aplicación en modo desarrollo"
	@echo "  ${COLOR_BLUE}make clean${COLOR_RESET}              - Elimina binarios y archivos temporales"
	@echo "  ${COLOR_BLUE}make test${COLOR_RESET}               - Ejecuta las pruebas unitarias"
	@echo "  ${COLOR_BLUE}make test-coverage${COLOR_RESET}      - Ejecuta las pruebas con reporte de cobertura"
	@echo "  ${COLOR_BLUE}make lint${COLOR_RESET}               - Ejecuta el linter para verificar problemas de estilo"
	@echo "  ${COLOR_BLUE}make migrate${COLOR_RESET}            - Ejecuta todas las migraciones pendientes"
	@echo "  ${COLOR_BLUE}make migrate-status${COLOR_RESET}     - Muestra el estado de las migraciones"
	@echo "  ${COLOR_BLUE}make migrate-rollback${COLOR_RESET}   - Revierte la última migración"
	@echo "  ${COLOR_BLUE}make swagger${COLOR_RESET}            - Genera la documentación de Swagger"
	@echo "  ${COLOR_BLUE}make install${COLOR_RESET}            - Instala las dependencias del proyecto"

install:
	@echo "${COLOR_YELLOW}Instalando dependencias...${COLOR_RESET}"
	@go mod download
	@echo "${COLOR_GREEN}Dependencias instaladas correctamente${COLOR_RESET}"

build:
	@echo "${COLOR_YELLOW}Compilando la aplicación...${COLOR_RESET}"
	@go build $(GOFLAGS) -o $(BINARY_NAME) ./cmd/
	@echo "${COLOR_GREEN}Compilación completada: $(BINARY_NAME)${COLOR_RESET}"

clean:
	@echo "${COLOR_YELLOW}Limpiando binarios y archivos temporales...${COLOR_RESET}"
	@rm -f $(BINARY_NAME)
	@go clean
	@echo "${COLOR_GREEN}Limpieza completada${COLOR_RESET}"

run:
	@echo "${COLOR_YELLOW}Ejecutando la aplicación en modo desarrollo...${COLOR_RESET}"
	@go run ./cmd/

test:
	@echo "${COLOR_YELLOW}Ejecutando pruebas unitarias...${COLOR_RESET}"
	@go test ./... -v
	@echo "${COLOR_GREEN}Pruebas completadas${COLOR_RESET}"

test-coverage:
	@echo "${COLOR_YELLOW}Ejecutando pruebas con reporte de cobertura...${COLOR_RESET}"
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out
	@echo "${COLOR_GREEN}Pruebas de cobertura completadas${COLOR_RESET}"

lint:
	@echo "${COLOR_YELLOW}Ejecutando linter...${COLOR_RESET}"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "${COLOR_YELLOW}golangci-lint no está instalado. Instálalo con:${COLOR_RESET}"; \
		echo "go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi
	@echo "${COLOR_GREEN}Lint completado${COLOR_RESET}"

migrate:
	@echo "${COLOR_YELLOW}Ejecutando migraciones...${COLOR_RESET}"
	@go run ./cmd/ migrate
	@echo "${COLOR_GREEN}Migraciones completadas${COLOR_RESET}"

migrate-status:
	@echo "${COLOR_YELLOW}Verificando estado de migraciones...${COLOR_RESET}"
	@go run ./cmd/ migrate status
	@echo "${COLOR_GREEN}Estado de migraciones completado${COLOR_RESET}"

migrate-rollback:
	@echo "${COLOR_YELLOW}Revirtiendo última migración...${COLOR_RESET}"
	@go run ./cmd/ migrate rollback
	@echo "${COLOR_GREEN}Rollback de migración completado${COLOR_RESET}"

swagger:
	@echo "${COLOR_YELLOW}Generando documentación Swagger...${COLOR_RESET}"
	@if command -v swag >/dev/null 2>&1; then \
		swag init -g cmd/main.go -o internal/infraestructure/inbound/httprest/docs; \
	else \
		echo "${COLOR_YELLOW}swag no está instalado. Instálalo con:${COLOR_RESET}"; \
		echo "go install github.com/swaggo/swag/cmd/swag@latest"; \
	fi
	@echo "${COLOR_GREEN}Documentación Swagger generada${COLOR_RESET}" 