# PostgreSQL Database Setup for ana.world

[English](#english) | [Español](#español)

<a name="english"></a>
## Database Setup Instructions

This document provides step-by-step instructions for setting up the PostgreSQL database for the ana.world project.

### Prerequisites

- PostgreSQL 14+ installed on your system
- Basic knowledge of SQL commands
- Terminal/command prompt access

### 1. Database Creation

Connect to PostgreSQL as the postgres user:

```bash
# Login as postgres user
sudo -u postgres psql
```

Once connected, create the database:

```sql
CREATE DATABASE ana_world;
```

### 2. User Setup

While still connected as postgres, create a dedicated user for the application:

```sql
CREATE USER ana_user WITH ENCRYPTED PASSWORD 'your_secure_password';
GRANT ALL PRIVILEGES ON DATABASE ana_world TO ana_user;
```

Exit the PostgreSQL prompt:

```sql
\q
```

### 3. Environment Variables Configuration

Create a `.env` file in the project root (this file should be git-ignored):

```
DB_HOST=localhost
DB_PORT=5432
DB_NAME=ana_world
DB_USER=ana_user
DB_PASSWORD=your_secure_password
DB_SSLMODE=disable
```

For development, you can also set these environment variables directly:

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=ana_world
export DB_USER=ana_user
export DB_PASSWORD=your_secure_password
export DB_SSLMODE=disable
```

### 4. Database Migration Scripts

Create the initial migration file in the `db/migrations` directory:

```bash
mkdir -p db/migrations
```

Create the first migration file `db/migrations/0001_init_tasks.sql`:

```sql
-- Up migration
CREATE TABLE tasks (
  id SERIAL PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  description TEXT,
  due_date TIMESTAMP,
  priority VARCHAR(50),
  project_id INTEGER,
  status VARCHAR(50),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for frequently queried fields
CREATE INDEX idx_tasks_due_date ON tasks(due_date);
CREATE INDEX idx_tasks_project_id ON tasks(project_id);
CREATE INDEX idx_tasks_status ON tasks(status);

-- Down migration (rollback)
-- DROP TABLE tasks;
```

### 5. Applying Migrations

Connect to the database and apply the migration:

```bash
# Connect to the database
psql -U ana_user -d ana_world -h localhost

# Once connected, run the migration
\i db/migrations/0001_init_tasks.sql

# Exit when done
\q
```

Alternatively, you can run the migration directly:

```bash
psql -U ana_user -d ana_world -h localhost -f db/migrations/0001_init_tasks.sql
```

### 6. Connection Testing Queries

To verify your database is set up correctly, run these test queries:

```sql
-- Connect to the database
psql -U ana_user -d ana_world -h localhost

-- List all tables
\dt

-- Insert a test task
INSERT INTO tasks (title, description, due_date, priority, project_id, status)
VALUES ('Test Task', 'This is a test task', NOW() + INTERVAL '1 day', 'Medium', 1, 'To-Do');

-- Query the inserted task
SELECT * FROM tasks;

-- Exit
\q
```

### 7. Go Database Connection

Implement the Go database connection as outlined in `next.md`. Create the file `internal/database/db.go` with the following content:

```go
package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB initializes the database connection
func InitDB() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "ana_user"),
		getEnv("DB_PASSWORD", "your_secure_password"),
		getEnv("DB_NAME", "ana_world"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_SSLMODE", "disable"),
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection established")
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
```

---

<a name="español"></a>
## Instrucciones de Configuración de la Base de Datos

Este documento proporciona instrucciones paso a paso para configurar la base de datos PostgreSQL para el proyecto ana.world.

### Prerrequisitos

- PostgreSQL 14+ instalado en su sistema
- Conocimiento básico de comandos SQL
- Acceso a terminal/línea de comandos

### 1. Creación de la Base de Datos

Conéctese a PostgreSQL como usuario postgres:

```bash
# Iniciar sesión como usuario postgres
sudo -u postgres psql
```

Una vez conectado, cree la base de datos:

```sql
CREATE DATABASE ana_world;
```

### 2. Configuración de Usuario

Mientras sigue conectado como postgres, cree un usuario dedicado para la aplicación:

```sql
CREATE USER ana_user WITH ENCRYPTED PASSWORD 'su_contraseña_segura';
GRANT ALL PRIVILEGES ON DATABASE ana_world TO ana_user;
```

Salga del prompt de PostgreSQL:

```sql
\q
```

### 3. Configuración de Variables de Entorno

Cree un archivo `.env` en la raíz del proyecto (este archivo debe ser ignorado por git):

```
DB_HOST=localhost
DB_PORT=5432
DB_NAME=ana_world
DB_USER=ana_user
DB_PASSWORD=su_contraseña_segura
DB_SSLMODE=disable
```

Para desarrollo, también puede establecer estas variables de entorno directamente:

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=ana_world
export DB_USER=ana_user
export DB_PASSWORD=su_contraseña_segura
export DB_SSLMODE=disable
```

### 4. Scripts de Migración de Base de Datos

Cree el archivo de migración inicial en el directorio `db/migrations`:

```bash
mkdir -p db/migrations
```

Cree el primer archivo de migración `db/migrations/0001_init_tasks.sql`:

```sql
-- Migración ascendente
CREATE TABLE tasks (
  id SERIAL PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  description TEXT,
  due_date TIMESTAMP,
  priority VARCHAR(50),
  project_id INTEGER,
  status VARCHAR(50),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Crear índices para campos consultados frecuentemente
CREATE INDEX idx_tasks_due_date ON tasks(due_date);
CREATE INDEX idx_tasks_project_id ON tasks(project_id);
CREATE INDEX idx_tasks_status ON tasks(status);

-- Migración descendente (rollback)
-- DROP TABLE tasks;
```

### 5. Aplicando Migraciones

Conéctese a la base de datos y aplique la migración:

```bash
# Conectarse a la base de datos
psql -U ana_user -d ana_world -h localhost

# Una vez conectado, ejecute la migración
\i db/migrations/0001_init_tasks.sql

# Salir cuando termine
\q
```

Alternativamente, puede ejecutar la migración directamente:

```bash
psql -U ana_user -d ana_world -h localhost -f db/migrations/0001_init_tasks.sql
```

### 6. Consultas de Prueba de Conexión

Para verificar que su base de datos está configurada correctamente, ejecute estas consultas de prueba:

```sql
-- Conectarse a la base de datos
psql -U ana_user -d ana_world -h localhost

-- Listar todas las tablas
\dt

-- Insertar una tarea de prueba
INSERT INTO tasks (title, description, due_date, priority, project_id, status)
VALUES ('Tarea de Prueba', 'Esta es una tarea de prueba', NOW() + INTERVAL '1 day', 'Media', 1, 'Por-Hacer');

-- Consultar la tarea insertada
SELECT * FROM tasks;

-- Salir
\q
```

### 7. Conexión a la Base de Datos en Go

Implemente la conexión a la base de datos en Go como se describe en `next.md`. Cree el archivo `internal/database/db.go` con el siguiente contenido:

```go
package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB inicializa la conexión a la base de datos
func InitDB() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "ana_user"),
		getEnv("DB_PASSWORD", "su_contraseña_segura"),
		getEnv("DB_NAME", "ana_world"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_SSLMODE", "disable"),
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}

	log.Println("Conexión a la base de datos establecida")
}

// getEnv obtiene una variable de entorno o devuelve un valor predeterminado
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
```

### 8. Testing Configuration

For testing, we use SQLite instead of PostgreSQL. This ensures:
- Tests can run without a PostgreSQL installation
- Fast test execution
- No test data persistence
- Reliable CI/CD pipeline

To run tests:

1. No additional setup required - SQLite is used automatically
2. Tests will create an in-memory database
3. Each test runs in isolation
4. Environment variables are automatically managed

Example test run:
```bash
# Run all tests
go test ./... -v

# Run database tests specifically
go test ./internal/database -v
```

Note: While we use PostgreSQL in production, SQLite is sufficient for testing most database operations. If you need to test PostgreSQL-specific features, use the `postgres` build tag:

```bash
go test ./... -tags postgres -v
```


## Development Session
Reference: https://app.warp.dev/session/b660fd8a-f765-449c-a70c-f8c7b971e3c4?pwd=e9ccd7cb-d8be-494e-a2f2-35469f726896
Last Updated: Sat May 17 07:34:44 AM CEST 2025
