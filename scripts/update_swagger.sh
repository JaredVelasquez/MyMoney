#!/bin/bash

# Directorio base del proyecto
BASE_DIR=$(pwd)
DOCS_DIR="$BASE_DIR/internal/infraestructure/inbound/httprest/docs"

# Verificar que el directorio de documentación existe
if [ ! -d "$DOCS_DIR" ]; then
    mkdir -p "$DOCS_DIR"
fi

# Copiar el archivo principal de Swagger
cp "$BASE_DIR/internal/infraestructure/inbound/httprest/docs/swagger_main.yaml" "$DOCS_DIR/swagger.yaml"

echo "Documentación Swagger actualizada correctamente."
echo "Para generar el archivo swagger.json, ejecuta: swag init -g internal/infraestructure/inbound/httprest/server.go -o internal/infraestructure/inbound/httprest/docs"

# Instrucciones para el usuario
echo ""
echo "Para instalar swag si no lo tienes:"
echo "go install github.com/swaggo/swag/cmd/swag@latest"
echo ""
echo "Asegúrate de tener anotaciones Swagger en tu código para que la generación funcione correctamente." 