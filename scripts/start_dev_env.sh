#!/bin/bash

# Script para iniciar el entorno de desarrollo con Docker Compose

# Colores para output
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Verificar si Docker está instalado
echo -e "${YELLOW}Verificando Docker...${NC}"
if ! command -v docker &> /dev/null; then
    echo -e "${RED}Docker no está instalado. Por favor, instálalo desde https://docs.docker.com/get-docker/${NC}"
    exit 1
fi
echo -e "${GREEN}Docker está instalado.${NC}"

# Verificar si Docker Compose está instalado
echo -e "${YELLOW}Verificando Docker Compose...${NC}"
if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}Docker Compose no está instalado. Por favor, instálalo desde https://docs.docker.com/compose/install/${NC}"
    exit 1
fi
echo -e "${GREEN}Docker Compose está instalado.${NC}"

# Ir al directorio de scripts
cd "$(dirname "$0")" || exit

# Iniciar los contenedores
echo -e "${YELLOW}Iniciando contenedores...${NC}"
docker-compose up -d

# Verificar si los contenedores están funcionando
echo -e "${YELLOW}Verificando estado de los contenedores...${NC}"
if [ "$(docker-compose ps -q | wc -l)" -gt 0 ]; then
    echo -e "${GREEN}Entorno de desarrollo iniciado correctamente.${NC}"
    echo -e "${YELLOW}Acceso a servicios:${NC}"
    echo -e "${GREEN}- PostgreSQL:${NC} localhost:5432 (Usuario: postgres, Contraseña: postgres, BD: mi_app_db)"
    echo -e "${GREEN}- pgAdmin:${NC} http://localhost:5050 (Email: admin@miapp.com, Contraseña: admin)"
else
    echo -e "${RED}Error al iniciar los contenedores. Revisa los logs con 'docker-compose logs'.${NC}"
    exit 1
fi

# Mensaje final
echo -e "${GREEN}¡Todo listo! Ya puedes ejecutar la aplicación con:${NC}"
echo -e "${YELLOW}make run${NC}" 