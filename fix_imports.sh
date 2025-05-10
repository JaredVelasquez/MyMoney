#!/bin/bash

# Script para corregir las importaciones en todo el proyecto
find . -type f -name "*.go" -exec sed -i 's|github.com/user/mi-app-backend|mi-app-backend|g' {} \;

echo "Importaciones corregidas" 