#!/bin/bash

# Reemplazar todas las ocurrencias de mi-app-backend con MyMoneyBackend en archivos .go
find . -type f -name "*.go" -exec sed -i 's/mi-app-backend/MyMoneyBackend/g' {} \;

# Regenerar la documentación Swagger
swag init -g cmd/main.go --parseDependency -d . -o internal/infraestructure/inbound/httprest/docs

echo "Todas las importaciones han sido actualizadas a MyMoneyBackend" 