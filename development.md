# Development Workflow for ana.world

[English](#english) | [Español](#español)

<a name="english"></a>
## Development Workflow and Best Practices

This document outlines the actual development workflow we've established for the ana.world project based on our experiences implementing features like PostgreSQL integration and Cerebras AI assistant.

### Development Tools

#### GitHub CLI (gh)

We use GitHub CLI extensively for repository operations to streamline our workflow:

```bash
# Clone the repository
gh repo clone juansgv/ana

# Create a pull request
gh pr create --title "Feature: PostgreSQL Integration" --body "Implements database connection and repositories."

# Check PR status
gh pr view

# Merge a pull request
gh pr merge 1 --merge --delete-branch

# Review a pull request
gh pr checkout 2
```

#### Feature Branch Workflow

Our established workflow follows these steps:

1. **Planning and Documentation**
   - Create detailed documentation for the feature in markdown format
   - Update next.md with specific implementation steps
   - Add references to external resources when necessary

2. **Feature Branch Creation**
   ```bash
   # Update main branch
   git checkout main
   git pull origin main
   
   # Create feature branch
   git checkout -b feature/feature-name
   ```

3. **Implementation Stages**
   - Create project structure first (directories, base files)
   - Implement core functionality with basic error handling
   - Add tests and improve error handling
   - Update documentation with implementation details

4. **Code Review and PR Process**
   ```bash
   # Commit changes in logical groups
   git add directory-or-files/
   git commit -m "feat(component): detailed description of changes
   
   - Add specific detail 1
   - Implement functionality 2
   - Fix issue 3"
   
   # Push changes
   git push origin feature/feature-name
   
   # Create pull request
   gh pr create --title "Feature: Feature Name" --body "Detailed description" --base main --head feature/feature-name
   ```

5. **Post-Merge Actions**
   - Sync changes to other branches if needed
   - Update documentation in flow.md with progress
   - Close related issues

### Implementation Approach

Our proven step-by-step implementation approach:

1. **Research and Planning (1-2 days)**
   - Research required libraries and best practices
   - Create detailed implementation plan in markdown
   - Identify potential challenges and solutions

2. **Core Implementation (2-3 days)**
   - Start with basic functionality without advanced error handling
   - Focus on the happy path first
   - Implement framework/structure for the entire feature

3. **Testing and Refinement (1-2 days)**
   - Add proper error handling and edge cases
   - Create test scripts as demonstrated in our test directories
   - Improve performance where needed

4. **Documentation and Integration (1 day)**
   - Update all relevant documentation
   - Ensure bilingual content (English/Spanish)
   - Prepare for PR and code review

### Real Examples from Our Development

#### PostgreSQL Integration Example

The PostgreSQL integration followed this process:

1. **Planning Phase**
   - Created database_setup.md with comprehensive setup instructions
   - Updated next.md with detailed implementation steps
   - Identified dependencies (GORM, PostgreSQL driver)

2. **Implementation Phase**
   ```bash
   # Create feature branch
   git checkout -b feature/postgresql-integration
   
   # Create directory structure
   mkdir -p internal/database internal/repositories db/migrations
   
   # Create initial migration file
   touch db/migrations/0001_init_tasks.sql
   ```

3. **Incremental Commits**
   ```bash
   # First commit: Database connection
   git add internal/database/db.go
   git commit -m "feat(db): implement PostgreSQL connection with GORM
   
   - Add database connection with environment configuration
   - Implement connection pooling
   - Add proper error handling and logging"
   
   # Second commit: Repository pattern
   git add internal/repositories/
   git commit -m "feat(repositories): implement task repository
   
   - Create task repository with CRUD operations
   - Add advanced query methods (FindTasksDueToday)
   - Implement error handling"
   ```

4. **PR and Review**
   ```bash
   # Create PR
   gh pr create --title "feat: PostgreSQL integration and repositories" --body "Implements database connection with GORM and repository pattern for task management."
   ```

#### Cerebras AI Integration Example

The Cerebras AI assistant implementation:

1. **Planning Phase**
   - Researched Cerebras API documentation
   - Created ai_integration.md with setup instructions
   - Updated next.md with implementation steps

2. **Implementation Phase**
   ```bash
   # Setup directory structure
   mkdir -p internal/ai
   touch internal/ai/cerebras_client.go
   touch internal/handlers/ai_cerebras_handler.go
   ```

3. **Progressive Implementation**
   - First created the client with API connection
   - Then implemented the handler with request/response handling
   - Finally added the frontend integration

4. **Troubleshooting**
   - Fixed imports between packages
   - Resolved token handling issues
   - Improved error messages and fallbacks

### Troubleshooting Common Issues

Based on our experience, here are solutions to common issues:

#### Import Path Issues

**Problem**: Package imports failing with "package not found" errors.

**Solution**:
```go
// Instead of this (which can cause issues)
import "github.com/juansgv/ana/internal/ai"

// Use this (matching the go.mod module name)
import "github.com/sebae/ana/internal/ai"
```

#### Database Connection Failures

**Problem**: Unable to connect to PostgreSQL database.

**Solution**:
1. Verify PostgreSQL is running: `systemctl status postgresql`
2. Check environment variables in .env file
3. Ensure database and user exist: `sudo -u postgres psql -c "\du"`
4. Verify connection string format

#### OAuth Token Issues

**Problem**: OAuth authentication fails or tokens expire.

**Solution**:
1. Implement token refresh mechanism
2. Store tokens securely (encrypted)
3. Add clear error messages for authentication failures
4. Implement proper scopes with least privilege

#### Frontend-Backend Integration

**Problem**: HTMX requests not working as expected.

**Solution**:
1. Check browser console for errors
2. Verify CORS settings in router.go
3. Ensure proper Content-Type headers
4. Test API endpoints directly with curl before frontend integration

### Code Conventions

Based on our implementations so far, we've established these conventions:

#### Go Conventions

1. **File Naming**: Lowercase with underscores for multiple words
   ```
   task_repository.go
   cerebras_client.go
   ```

2. **Package Structure**:
   ```
   internal/         # Internal packages
     models/         # Data models
     repositories/   # Data access layer
     handlers/       # HTTP handlers
     server/         # Server configuration
     database/       # Database connection
     ai/             # AI integration
   ```

3. **Function Naming**:
   - Public functions: PascalCase (e.g., `GetTasks`, `NewTaskRepository`)
   - Private functions: camelCase (e.g., `getEnv`, `parseToken`)

4. **Error Handling**:
   ```go
   if err != nil {
       log.Printf("Specific error message: %v", err)
       return nil, fmt.Errorf("readable error with context: %w", err)
   }
   ```

#### Commit Conventions

We follow a modified version of Conventional Commits:

```
feat(component): add new feature X

- Add specific functionality
- Fix related issues
- Update documentation
```

Categories:
- `feat`: New features
- `fix`: Bug fixes
- `docs`: Documentation changes
- `refactor`: Code refactoring
- `test`: Testing updates
- `chore`: Maintenance tasks

---

<a name="español"></a>
## Flujo de Desarrollo y Mejores Prácticas

Este documento describe el flujo de desarrollo real que hemos establecido para el proyecto ana.world basado en nuestras experiencias implementando características como la integración de PostgreSQL y el asistente de IA Cerebras.

### Herramientas de Desarrollo

#### GitHub CLI (gh)

Utilizamos GitHub CLI extensivamente para operaciones de repositorio para agilizar nuestro flujo de trabajo:

```bash
# Clonar el repositorio
gh repo clone juansgv/ana

# Crear un pull request
gh pr create --title "Característica: Integración de PostgreSQL" --body "Implementa conexión a base de datos y repositorios."

# Verificar estado del PR
gh pr view

# Fusionar un pull request
gh pr merge 1 --merge --delete-branch

# Revisar un pull request
gh pr checkout 2
```

#### Flujo de Trabajo de Ramas de Características

Nuestro flujo de trabajo establecido sigue estos pasos:

1. **Planificación y Documentación**
   - Crear documentación detallada para la característica en formato markdown
   - Actualizar next.md con pasos específicos de implementación
   - Agregar referencias a recursos externos cuando sea necesario

2. **Creación de Ramas de Características**
   ```bash
   # Actualizar rama principal
   git checkout main
   git pull origin main
   
   # Crear rama de característica
   git checkout -b feature/nombre-caracteristica
   ```

3. **Etapas de Implementación**
   - Crear primero la estructura del proyecto (directorios, archivos base)
   - Implementar funcionalidad central con manejo básico de errores
   - Agregar pruebas y mejorar el manejo de errores
   - Actualizar documentación con detalles de implementación

4. **Revisión de Código y Proceso de PR**
   ```bash
   # Confirmar cambios en grupos lógicos
   git add directorio-o-archivos/
   git commit -m "feat(componente): descripción detallada de cambios
   
   - Agregar detalle específico 1
   - Implementar funcionalidad 2
   - Corregir problema 3"
   
   # Enviar cambios
   git push origin feature/nombre-caracteristica
   
   # Crear pull request
   gh pr create --title "Característica: Nombre de Característica" --body "Descripción detallada" --base main --head feature/nombre-caracteristica
   ```

5. **Acciones Post-Fusión**
   - Sincronizar cambios a otras ramas si es necesario
   - Actualizar documentación en flow.md con el progreso
   - Cerrar issues relacionados

### Enfoque de Implementación

Nuestro enfoque de implementación paso a paso probado:

1. **Investigación y Planificación (1-2 días)**
   - Investigar bibliotecas requeridas y mejores prácticas
   - Crear un plan de implementación detallado en markdown
   - Identificar desafíos potenciales y soluciones

2. **Implementación Central (2-3 días)**
   - Comenzar con funcionalidad básica sin manejo avanzado de errores
   - Enfocarse primero en el camino feliz
   - Implementar marco/estructura para toda la característica

3. **Pruebas y Refinamiento (1-2 días)**
   - Agregar manejo adecuado de errores y casos límite
   - Crear scripts de prueba como se demuestra en nuestros directorios de prueba
   - Mejorar el rendimiento donde sea necesario

4. **Documentación e Integración (1 día)**
   - Actualizar toda la documentación relevante
   - Asegurar contenido bilingüe (Inglés/Español)
   - Preparar para PR y revisión de código

### Ejemplos Reales de Nuestro Desarrollo

#### Ejemplo de Integración de PostgreSQL

La integración de PostgreSQL siguió este proceso:

1. **Fase de Planificación**
   - Creación de database_setup.md con instrucciones completas de configuración
   - Actualización de next.md con pasos detallados de implementación
   - Identificación de dependencias (GORM, controlador de PostgreSQL)

2. **Fase de Implementación**
   ```bash
   # Crear rama de característica
   git checkout -b feature/postgresql-integration
   
   # Crear estructura de directorios
   mkdir -p internal/database internal/repositories db/migrations
   
   # Crear archivo de migración inicial
   touch db/migrations/0001_init_tasks.sql
   ```

3. **Commits Incrementales**
   ```bash
   # Primer commit: Conexión a base de datos
   git add internal/database/db.go
   git commit -m "feat(db): implementar conexión PostgreSQL con GORM
   
   - Agregar conexión a base de datos con configuración de entorno
   - Implementar agrupación de conexiones
   - Agregar manejo adecuado de errores y registro"
   
   # Segundo commit: Patrón de repositorio
   git add internal/repositories/
   git commit -m "feat(repositories): implementar repositorio de tareas
   
   - Crear repositorio de tareas con operaciones CRUD
   - Agregar métodos de consulta avanzados (FindTasksDueToday)
   - Implementar manejo de errores"
   ```

4. **PR y Revisión**
   ```bash
   # Crear PR
   gh pr create --title "feat: Integración de PostgreSQL y repositorios" --body "Implementa conexión a base de datos con GORM y patrón de repositorio para gestión de tareas."
   ```

#### Ejemplo de Integración de IA Cerebras

La implementación del asistente de IA Cerebras:

1. **Fase de Planificación**
   - Investigación de documentación de API de Cerebras
   - Creación de ai_integration.md con instrucciones de configuración
   - Actualización de next.md con pasos de implementación

2. **Fase de Implementación**
   ```bash
   # Configurar estructura de directorios
   mkdir -p internal/ai
   touch internal/ai/cerebras_client.go
   touch internal/handlers/ai_cerebras_handler.go
   ```

3. **Implementación Progresiva**
   - Primero se creó el cliente con conexión a API
   - Luego se implementó el manejador con procesamiento de solicitudes/respuestas
   - Finalmente se agregó la integración frontend

4. **Solución de Problemas**
   - Corrección de importaciones entre paquetes
   - Resolución de problemas de manejo de tokens
   - Mejora de mensajes de error y alternativas

### Solución de Problemas Comunes

Basado en nuestra experiencia, aquí hay soluciones a problemas comunes:

#### Problemas de Ruta de Importación

**Problema**: Fallos en importaciones de paquetes con errores "package not found".

**Solución**:
```go
// En lugar de esto (que puede causar problemas)
import "github.com/juansgv/ana/internal/ai"

// Usar esto (coincidiendo con el nombre del módulo en go.mod)
import "github.com/sebae/ana/internal/ai"
```

#### Fallos de Conexión a Base de Datos

**Problema**: No se puede conectar a la base de datos PostgreSQL.

**Solución**:
1. Verificar que PostgreSQL esté en ejecución: `systemctl status postgresql`
2. Comprobar variables de entorno en el archivo .env
3. Asegurar que la base de datos

# Development Guidelines for ana.world

This document outlines the standards, patterns, and practices to maintain consistency throughout the development of the ana.world project.

## Documentation Standards

### General Documentation

- **Bilingual Requirement**: All user-facing documentation must be in both English and Spanish.
- **Documentation Files**: 
  - `README.md` - Project overview and setup instructions
  - `plan.md` - Comprehensive project plan 
  - `next.md` - Current implementation steps
  - `flow.md` - Project evolution and development process
  - `development.md` - Development standards (this document)

### Code Documentation

- **Function Comments**: All exported functions must have comments in the following format:
  ```go
  // FunctionName does X when given Y.
  // It returns Z or an error if something goes wrong.
  func FunctionName() {}
  ```

- **Package Documentation**: Each package should have a package comment:
  ```go
  // Package name provides functionality for...
  package name
  ```

- **Complex Logic**: Any complex logic should be commented with explanations.

### API Documentation 

- **Endpoint Documentation**: All API endpoints should be documented with:
  - HTTP method 
  - Route
  - Request parameters/body
  - Response format
  - Example request/response
  - Possible errors

## Development Workflow Patterns

### Feature Development Process

1. **Planning**:
   - Review `plan.md` for feature requirements
   - Update `next.md` with detailed implementation steps if needed

2. **Branch Creation**:
   - Create a feature branch: `git checkout -b feature/feature-name`

3. **Implementation**:
   - Backend implementation first (API endpoints)
   - Frontend implementation second
   - Tests implementation third

4. **Documentation**:
   - Update relevant documentation
   - Add bilingual user documentation if needed

5. **Testing**:
   - Run unit tests
   - Manual testing of the feature
   - Verify both Spanish and English interfaces

6. **Pull Request**:
   - Submit a pull request with a comprehensive description
   - Ensure all CI checks pass

7. **Review & Merge**:
   - Address review comments
   - Merge to main branch once approved

### Issue Management

- Create detailed issue descriptions including:
  - Problem statement
  - Expected behavior
  - Steps to reproduce
  - Suggested implementation (if applicable)

## Commit Message Conventions

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

### Format

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

### Types

- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation changes
- `style`: Changes that don't affect code functionality (formatting, etc.)
- `refactor`: Code changes that neither fix bugs nor add features
- `perf`: Performance improvements
- `test`: Adding or correcting tests
- `chore`: Changes to build process, tools, etc.

### Examples

```
feat(tasks): add due date filtering functionality

fix(auth): resolve token refresh issue

docs: update PostgreSQL setup instructions in Spanish and English
```

## File Organization Rules

### Directory Structure

Maintain the established directory structure as outlined in the README:

```
ana/
├── cmd/                    # Application entry points
│   └── server/             # Main server binary
├── internal/               # Internal packages
│   ├── handlers/           # HTTP handlers
│   ├── models/             # Data models
│   ├── server/             # Server configuration
│   ├── database/           # Database connection
│   └── repositories/       # Data access layer
├── web/                    # Frontend static files
├── functions/              # Netlify serverless functions
├── db/                     # Database migrations and schemas
│   └── migrations/         # SQL migration files
└── config/                 # Configuration files
```

### Naming Conventions

- **Files**: Use snake_case for filenames (e.g., `task_handler.go`)
- **Go Packages**: Use lowercase, single-word names (e.g., `models`)
- **Go Structs/Interfaces**: Use PascalCase (e.g., `TaskRepository`)
- **Frontend Components**: Use kebab-case for file names (e.g., `task-list.html`)
- **Database Tables**: Use snake_case (e.g., `user_projects`)

## Quality Assurance Steps

### Code Review Checklist

- [ ] Code follows established patterns and naming conventions
- [ ] New functions/methods have appropriate documentation
- [ ] Error handling is implemented correctly
- [ ] Tests are included for new functionality
- [ ] Security considerations have been addressed
- [ ] Performance impact has been considered
- [ ] Bilingual support is maintained where applicable

### Testing Requirements

- **Unit Tests**: All backend packages should have unit tests with >70% coverage
- **Integration Tests**: API endpoints should have integration tests
- **Frontend Testing**: Manual testing of HTMX interactions
- **Database Testing**: Test migration scripts in a test database before applying to development

### Security Checks

- Input validation on all API endpoints
- Parameterized SQL queries
- Authentication/authorization checks
- Secure handling of sensitive data

## Project Consistency Guidelines

### Coding Style

#### Go

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `go fmt` before committing code
- Prefer standard library packages when possible
- Keep functions small and focused on a single responsibility

#### SQL

- Use uppercase for SQL keywords
- Use parameterized queries
- Document complex queries

#### Frontend

- Follow HTMX patterns for dynamic content loading
- Maintain consistent Tailwind CSS usage for styling
- Use Spanish language in user-facing elements with English comments

### Configuration Management

- Use environment variables for configuration
- Document all configuration options in both languages
- Provide sensible defaults for development

### Version Control

- Make small, focused commits
- Use feature branches for development
- Rebase feature branches before merging
- Squash commits when appropriate

## Implementation Phases

Follow the implementation phases as defined in the project plan:

1. **Phase 1**: Setup and Core Features
2. **Phase 2**: Project Management and Agenda
3. **Phase 3**: Google Calendar Integration
4. **Phase 4**: Supplier Search
5. **Phase 5**: Dashboards and Enhancements
6. **Phase 6**: Deployment and Maintenance

Document phase transitions in the `flow.md` file.

### Current Implementation Status

#### Completed Features

1. **Core Infrastructure (Phase 1)**
   - ✓ PostgreSQL integration with GORM
   - ✓ Repository pattern implementation
   - ✓ Basic CRUD operations
   - ✓ Error handling and logging
   - ✓ Environment-based configuration

2. **AI Integration (Phase 2)**
   - ✓ Cerebras AI integration
   - ✓ Dual model support (QWen-3B-32B and QWen-2.5-Vision)
   - ✓ Context-aware bilingual responses
   - ✓ Direct answer mode (/no_think)
   - ✓ Model selection dropdown in frontend

3. **Monitoring Stack**
   - ✓ Prometheus metrics integration
   - ✓ Health check endpoints
   - ✓ Error tracking system
   - ✓ Performance monitoring
   - ✓ Grafana dashboards

#### Current Focus

1. **Google Integration (Phase 3)**
   - Calendar API integration for task scheduling
   - Gmail integration for notifications
   - OAuth2 authentication flow
   - Token refresh mechanism
   - User: anais.villamarinj@gmail.com

2. **Testing Infrastructure**
   - Improving test coverage for database operations
   - Adding integration tests for AI endpoints
   - Implementing end-to-end tests
   - SQLite-based test isolation

3. **Documentation Refinement**
   - Updating API documentation
   - Expanding troubleshooting guides
   - Adding architecture diagrams
   - Maintaining bilingual format

#### Next Steps

1. **Short-term Goals**
   - Complete Google Calendar OAuth2 setup
   - Implement task-to-calendar sync
   - Add email notification system
   - Enhance test coverage to >80%

2. **Medium-term Goals**
   - Location-based supplier search
   - Financial tracking dashboard
   - Document management system
   - Performance optimization

3. **Long-term Goals**
   - Mobile responsiveness improvements
   - Offline capability
   - Advanced AI features
   - Analytics dashboard

## Deployment Practices

- Test in development environment before deploying to production
- Document deployment steps in both languages
- Maintain a deployment checklist
- Use environment-specific configuration

---

By following these guidelines, we ensure a consistent, high-quality, and well-documented development process for ana.world. These standards will evolve as the project grows, with updates documented in this file.


## Development Session
Reference: https://app.warp.dev/session/b660fd8a-f765-449c-a70c-f8c7b971e3c4?pwd=e9ccd7cb-d8be-494e-a2f2-35469f726896
Last Updated: Sat May 17 07:34:44 AM CEST 2025
