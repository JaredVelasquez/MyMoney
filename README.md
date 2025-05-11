# MyMoneyBackend

API backend para la aplicación MyMoney, que proporciona gestión de usuarios, suscripciones, monedas y planes.

## Requisitos

- Go 1.18 o superior
- PostgreSQL 13+
- Docker (opcional)

## Instalación

Clone el repositorio:

```bash
git clone https://github.com/yourusername/MyMoneyBackend.git
cd MyMoneyBackend
```

Instale las dependencias:

```bash
make install
```

## Uso

El proyecto incluye un Makefile para facilitar las operaciones comunes:

```bash
# Ver comandos disponibles
make help

# Compilar la aplicación
make build

# Ejecutar en modo desarrollo
make run

# Ejecutar pruebas
make test

# Ejecutar pruebas con cobertura
make test-coverage

# Ejecutar el linter
make lint

# Gestión de migraciones
make migrate       # Ejecutar migraciones pendientes
make migrate-status    # Ver estado de migraciones
make migrate-rollback  # Revertir última migración

# Generar documentación Swagger
make swagger
```

## Estructura del Proyecto

```
MyMoneyBackend/
├── cmd/                # Punto de entrada de la aplicación
├── internal/           # Código principal de la aplicación
│   ├── application/    # Casos de uso y lógica de aplicación
│   ├── domain/         # Entidades y reglas de negocio
│   ├── infraestructure/# Implementaciones externas (BD, HTTP, etc.)
├── pkg/                # Bibliotecas que pueden ser utilizadas por aplicaciones externas
├── db/                 # Migraciones de base de datos
├── test/               # Pruebas adicionales (integración, E2E)
├── scripts/            # Scripts de utilidad
└── docs/               # Documentación
```

## API

La API está documentada usando Swagger. Después de iniciar la aplicación, puede acceder a la documentación Swagger en:

```
http://localhost:8080/swagger/index.html
```

### Principales Endpoints

- **Autenticación**: `/auth/register`, `/auth/login`
- **Usuarios**: `/users/me`, `/users/update`
- **Monedas**: `/currencies`
- **Planes**: `/plans`
- **Suscripciones**: `/subscriptions`

## Desarrollo

### Convenciones

- Seguimos la arquitectura hexagonal (puertos y adaptadores)
- Usamos Clean Code y SOLID principles
- Cada feature debe incluir pruebas unitarias

## Licencia

Este proyecto está licenciado bajo [MIT License](LICENSE). 