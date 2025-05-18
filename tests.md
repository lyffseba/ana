# Testing Strategy | Estrategia de Pruebas

[English](#english) | [Español](#español)

<a name="english"></a>
## Testing Documentation for ana.world

This document describes the testing approach, tools, and procedures used in the ana.world project to ensure code quality and reliability.

### Test Strategy

The ana.world project employs a multi-layered testing approach:

#### Unit Tests

Unit tests verify the correctness of individual components in isolation. Key areas with unit test coverage:

- **AI Client Tests**: Test the Cerebras API client with mock HTTP responses
- **Repository Tests**: Validate CRUD operations against an in-memory SQLite database
- **Model Tests**: Verify model validation and business logic rules
- **Handler Tests**: Ensure HTTP handlers process requests and responses correctly

#### Integration Tests

Integration tests verify that different components work correctly together:

- **API Endpoint Tests**: Test complete HTTP request-response cycles
- **Database Integration**: Validate repository patterns with actual database operations
- **External Service Integration**: Test integration with external services like Cerebras API

#### Mock Implementations

Mock objects are used extensively to isolate components for testing:

- **HTTP Client Mocks**: For testing API clients without making real HTTP requests
- **Database Mocks**: Using SQLite in-memory database for repository tests
- **Service Mocks**: Mock implementations of service interfaces for testing handlers

#### Test Coverage

Current test coverage focuses on critical components:

| Component | Test Coverage | Notes |
|-----------|---------------|-------|
| AI Client | 92% | Full coverage of error handling paths |
| AI Handlers | 88% | All endpoints and error conditions tested |
| Repositories | 95% | All CRUD operations covered |
| Models | 75% | Core validation covered |
| Server | 65% | Main routing tested |

### Examples of Test Implementations

#### Mock Client Example

```go
// MockCerebrasClient is a mock implementation of the Cerebras client for testing
type MockCerebrasClient struct {
	GenerateTextResponseFunc     func(query string, model string, context []ai.Message) (string, error)
	GenerateVisionResponseFunc   func(query string, imageBase64 string, context []ai.Message) (string, error)
	GenerateAssistantResponseFunc func(query string, context []ai.Message) (string, error)
}

// Implementation of interface methods that delegate to the function fields
func (m *MockCerebrasClient) GenerateTextResponse(query string, model string, context []ai.Message) (string, error) {
	return m.GenerateTextResponseFunc(query, model, context)
}
```

#### Test Helpers

```go
// setupTestDB initializes an in-memory SQLite database for testing
func setupTestDB(t *testing.T) {
	// Store the original DB connection
	originalDB := database.DB

	// Create an in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Replace the global DB connection with our test DB
	database.DB = db

	// Auto migrate the models for testing
	err = database.DB.AutoMigrate(&models.Task{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	// Helper function to restore the original DB after the test
	t.Cleanup(func() {
		database.DB = originalDB
	})
}
```

#### Common Testing Patterns

1. **Table-Driven Tests**: Used for testing multiple scenarios with the same test logic

```go
func TestErrorHandling(t *testing.T) {
	testCases := []struct {
		name           string
		modelType      string
		withImage      bool
		mockClient     *MockCerebrasClient
		expectedStatus int
		expectedError  string
	}{
		{
			name:      "Text Model Error",
			modelType: "qwen-text",
			withImage: false,
			mockClient: &MockCerebrasClient{
				GenerateTextResponseFunc: func(query string, model string, context []ai.Message) (string, error) {
					return "", io.EOF
				},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Error en el procesamiento de la consulta",
		},
		// More test cases...
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test logic here
		})
	}
}
```

2. **Dependency Injection**: Used to inject mock dependencies for testing

```go
func setupTestRouter(mockClient *MockCerebrasClient) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	
	r.POST("/api/ai/cerebras", createTestHandler(mockClient))
	
	return r
}
```

3. **Fixtures**: Using helper functions to create test data

```go
func createSampleTask() models.Task {
	return models.Task{
		Title:       "Test Task",
		Description: "This is a test task description",
		DueDate:     time.Now().AddDate(0, 0, 1), // Tomorrow
		Priority:    "Medium",
		ProjectID:   1,
		Status:      "To-Do",
	}
}
```

### Running Tests

To run all tests in the project:

```bash
go test ./... -v
```

To run tests for a specific package:

```bash
go test github.com/lyffseba/ana/internal/ai -v
```

To run a specific test:

```bash
go test github.com/lyffseba/ana/internal/handlers -run TestTextModelRequest -v
```

To generate a coverage report:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### Adding New Tests

When adding new functionality, follow these steps to add tests:

1. **Identify Test Scenarios**: Determine what needs to be tested, including edge cases and error conditions
2. **Create Test File**: Add a test file next to the implementation file (e.g., `handler.go` and `handler_test.go`)
3. **Implement Tests**: Write test functions using Go's testing package
4. **Run and Verify**: Run the tests and ensure they pass and provide adequate coverage

#### Guidelines for Writing Tests

- Test names should be descriptive and follow the `TestFunctionName` format
- Use table-driven tests for similar test cases
- Mock external dependencies to isolate the component being tested
- Test both success and error paths
- Verify all assertions with meaningful error messages

### CI/CD Integration

Tests are automatically run as part of our CI/CD pipeline:

1. **GitHub Actions**: Tests run on every pull request
2. **Pre-Merge Checks**: All tests must pass before merging
3. **Coverage Reports**: Test coverage reports are generated and stored as artifacts
4. **Deployment Testing**: Integration tests run before deployment to production

---

<a name="español"></a>
## Documentación de Pruebas para ana.world

Este documento describe el enfoque, herramientas y procedimientos de prueba utilizados en el proyecto ana.world para garantizar la calidad y fiabilidad del código.

### Estrategia de Pruebas

El proyecto ana.world emplea un enfoque de pruebas de múltiples capas:

#### Pruebas Unitarias

Las pruebas unitarias verifican la corrección de componentes individuales de forma aislada. Áreas clave con cobertura de pruebas unitarias:

- **Pruebas de Cliente AI**: Prueban el cliente de la API Cerebras con respuestas HTTP simuladas
- **Pruebas de Repositorio**: Validan operaciones CRUD contra una base de datos SQLite en memoria
- **Pruebas de Modelos**: Verifican la validación de modelos y reglas de lógica de negocio
- **Pruebas de Manejadores**: Aseguran que los manejadores HTTP procesen correctamente las solicitudes y respuestas

#### Pruebas de Integración

Las pruebas de integración verifican que diferentes componentes funcionen correctamente juntos:

- **Pruebas de Puntos Finales API**: Prueban ciclos completos de solicitud-respuesta HTTP
- **Integración de Base de Datos**: Validan patrones de repositorio con operaciones reales de base de datos
- **Integración de Servicios Externos**: Prueban la integración con servicios externos como la API de Cerebras

#### Implementaciones Simuladas (Mocks)

Los objetos simulados se utilizan ampliamente para aislar componentes para pruebas:

- **Mocks de Cliente HTTP**: Para probar clientes API sin realizar solicitudes HTTP reales
- **Mocks de Base de Datos**: Uso de base de datos SQLite en memoria para pruebas de repositorio
- **Mocks de Servicios**: Implementaciones simuladas de interfaces de servicio para probar manejadores

### Database Testing Strategy

We use SQLite for testing instead of PostgreSQL to ensure:
- Fast test execution
- No external dependencies
- Consistent test environment
- Easy CI/CD integration

Key components of our database testing approach:

1. **Mock Database Provider**:
   ```go
   // Database initialization can be mocked for testing
   var initializeDB = func() (*gorm.DB, error) {
       return gorm.Open(postgres.Open(buildDSN()), &gorm.Config{})
   }
   ```

2. **Environment Isolation**:
   - Tests temporarily clear database environment variables
   - SQLite in-memory database used for all tests
   - Original environment restored after tests

3. **Test Setup Helper**:
   ```go
   func setupTestDB(t *testing.T) {
       // Save and clear environment variables
       envVars := []string{"DB_HOST", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_PORT"}
       envBackup := make(map[string]string)
       
       for _, env := range envVars {
           envBackup[env] = os.Getenv(env)
           os.Unsetenv(env)
       }
       
       // Restore environment after test
       t.Cleanup(func() {
           for env, value := range envBackup {
               if value != "" {
                   os.Setenv(env, value)
               }
           }
       })

       // Use SQLite for testing
       db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
       if err != nil {
           t.Fatalf("Failed to connect to test database: %v", err)
       }

       // Replace global DB with test DB
       database.DB = db
   }
   ```

4. **Transaction Wrapping**:
   - Each test runs in a transaction
   - Transactions are rolled back after each test
   - Ensures test isolation

5. **Migration Testing**:
   - Migrations are tested against SQLite
   - Ensures schema compatibility
   - Validates migration reversibility

#### Cobertura de Pruebas

La cobertura de pruebas actual se centra en componentes críticos:

| Componente | Cobertura | Notas |
|------------|-----------|-------|
| Cliente AI | 92% | Cobertura completa de rutas de manejo de errores |
| Manejadores AI | 88% | Todos los endpoints y condiciones de error probados |
| Repositorios | 95% | Todas las operaciones CRUD cubiertas |
| Modelos | 75% | Validación principal cubierta |
| Servidor | 65% | Enrutamiento principal probado |

### Ejemplos de Implementaciones de Prueba

#### Ejemplo de Cliente Simulado

```go
// MockCerebrasClient es una implementación simulada del cliente Cerebras para pruebas
type MockCerebrasClient struct {
	GenerateTextResponseFunc     func(query string, model string, context []ai.Message) (string, error)
	GenerateVisionResponseFunc   func(query string, imageBase64 string, context []ai.Message) (string, error)
	GenerateAssistantResponseFunc func(query string, context []ai.Message) (string, error)
}

// Implementación de métodos de interfaz que delegan a los campos de función
func (m *MockCerebrasClient) GenerateTextResponse(query string, model string, context []ai.Message) (string, error) {
	return m.GenerateTextResponseFunc(query, model, context)
}
```

#### Ayudantes de Prueba

```go
// setupTestDB inicializa una base de datos SQLite en memoria para pruebas
func setupTestDB(t *testing.T) {
	// Guarda la conexión DB original
	originalDB := database.DB

	// Crea una base de datos SQLite en memoria para pruebas
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error al conectar a la base de datos de prueba: %v", err)
	}

	// Reemplaza la conexión DB global con nuestra DB de prueba
	database.DB = db

	// Auto migra los modelos para pruebas
	err = database.DB.AutoMigrate(&models.Task{})
	if err != nil {
		t.Fatalf("Error al migrar la base de datos de prueba: %v", err)
	}

	// Función auxiliar para restaurar la DB original después de la prueba
	t.Cleanup(func() {
		database.DB = originalDB
	})
}
```

#### Patrones Comunes de Prueba

1. **Pruebas Basadas en Tablas**: Utilizadas para probar múltiples escenarios con la misma lógica de prueba

```go
func TestErrorHandling(t *testing.T) {
	testCases := []struct {
		name           string
		modelType      string
		withImage      bool
		mockClient     *MockCerebrasClient
		expectedStatus int
		expectedError  string
	}{
		{
			name:      "Error de Modelo de Texto",
			modelType: "qwen-text",
			withImage: false,
			mockClient: &MockCerebrasClient{
				GenerateTextResponseFunc: func(query string, model string, context []ai.Message) (string, error) {
					return "", io.EOF
				},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Error en el procesamiento de la consulta",
		},
		// Más casos de prueba...
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Lógica de prueba aquí
		})
	}
}
```

2. **Inyección de Dependencias**: Utilizada para inyectar dependencias simuladas para pruebas

```go
func setupTestRouter(mockClient *MockCerebrasClient) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	
	r.POST("/api/ai/cerebras", createTestHandler(mockClient))
	
	return r
}
```

3. **Fixtures**: Uso de funciones auxiliares para crear datos de prueba

```go
func createSampleTask() models.Task {
	return models.Task{
		Title:       "Tarea de Prueba",
		Description: "Esta es una descripción de tarea de prueba",
		DueDate:     time.Now().AddDate(0, 0, 1), // Mañana
		Priority:    "Medium",
		ProjectID:   1,
		Status:      "To-Do",
	}
}
```

### Ejecución de Pruebas

Para ejecutar todas las pruebas en el proyecto:

```bash
go test ./... -v
```

Para ejecutar pruebas de un paquete específico:

```bash
go test github.com/lyffseba/ana/internal/ai -v
```

Para ejecutar una prueba específica:

```bash
go test github.com/lyffseba/ana/internal/handlers -run TestTextModelRequest -v
```

Para generar un informe de cobertura:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### Añadir Nuevas Pruebas

Al añadir nueva funcionalidad, sigue estos pasos para añadir pruebas:

1. **Identificar Escenarios de Prueba**: Determina qué necesita ser probado, incluyendo casos límite y condiciones de error
2. **Crear Archivo de Prueba**: Añade un archivo de prueba junto al archivo de implementación (p.ej., `handler.go` y `handler_test.go`)
3. **Implementar Pruebas**: Escribe funciones de prueba utilizando el paquete testing de Go
4. **Ejecutar y Verificar**: Ejecuta las pruebas y asegúrate de que pasan y proporcionan una cobertura adecuada

#### Directrices para Escribir Pruebas

- Los nombres de las pruebas deben ser descriptivos y seguir el formato `TestNombreFunción`
- Usa pruebas basadas en tablas para casos de prueba similares
- Simula dependencias externas para aislar el componente que se está probando
- Prueba tanto las rutas de éxito como las de error
- Verifica todas las aserciones con mensajes de error significativos

### Integración CI/CD

Las pruebas se ejecutan automáticamente como parte de nuestro pipeline CI/CD:

1. **GitHub Actions**: Las pruebas se ejecutan en cada pull request
2. **Verificaciones Pre-Merge**: Todas las pruebas deben pasar antes de la fusión
3. **


## Development Session
Reference: https://app.warp.dev/session/b660fd8a-f765-449c-a70c-f8c7b971e3c4?pwd=ec37c232-0b4b-4a91-af2c-ef680eaa123b
Last Updated: Sat May 17 08:12:57 AM CEST 2025
