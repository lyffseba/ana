# ana.world

[English](#english) | [Español](#español)

<a name="english"></a>
## Project Management for Architects

A specialized application designed for architects managing multiple construction projects in Bogota, Colombia. The platform combines task management, project tracking, and supplier search with a user-friendly interface.

### Key Features

- **Daily Agenda**: View and manage tasks due today with clear prioritization
- **Project Management**: Track multiple architecture projects with comprehensive details
- **Google Calendar Integration**: Sync tasks and milestones with your Google Calendar
- **Supplier Search**: Find architecture material suppliers in Bogota with location-based search
- **Financial Tracking**: Monitor project budgets and expenses

---

<a name="español"></a>
## Gestión de Proyectos para Arquitectos

Una aplicación especializada diseñada para arquitectos que gestionan múltiples proyectos de construcción en Bogotá, Colombia. La plataforma combina gestión de tareas, seguimiento de proyectos y búsqueda de proveedores con una interfaz fácil de usar.

### Características Principales

- **Agenda Diaria**: Visualiza y gestiona tareas con vencimiento hoy con priorización clara
- **Gestión de Proyectos**: Seguimiento de múltiples proyectos arquitectónicos con detalles completos
- **Integración con Google Calendar**: Sincroniza tareas y eventos importantes con tu Google Calendar
- **Búsqueda de Proveedores**: Encuentra proveedores de materiales arquitectónicos en Bogotá con búsqueda basada en ubicación
- **Seguimiento Financiero**: Monitorea presupuestos y gastos de proyectos

---

## Technology Stack | Stack Tecnológico

- **Frontend**: HTMX + Tailwind CSS (minimalist JavaScript approach)
- **Backend**: Go with Gin framework
- **Database**: PostgreSQL
- **Hosting**: Netlify (frontend), with Go serverless functions or GCP (backend)
- **APIs**: Google Calendar, Google Places API

## Project Structure | Estructura del Proyecto

```
ana/
├── cmd/                    # Application entry points | Puntos de entrada de la aplicación
│   └── server/             # Main server binary | Binario principal del servidor
├── internal/               # Internal packages | Paquetes internos
│   ├── handlers/           # HTTP handlers | Manejadores HTTP
│   ├── models/             # Data models | Modelos de datos
│   ├── server/             # Server configuration | Configuración del servidor
│   ├── database/           # Database connection | Conexión a base de datos
│   └── repositories/       # Data access layer | Capa de acceso a datos
├── web/                    # Frontend static files | Archivos estáticos del frontend
├── functions/              # Netlify serverless functions | Funciones serverless para Netlify
├── db/                     # Database migrations and schemas | Migraciones y esquemas de base de datos
└── config/                 # Configuration files | Archivos de configuración
```

## Setup Instructions | Instrucciones de Configuración

### Prerequisites | Requisitos Previos

- Go 1.20+ 
- PostgreSQL 14+
- Node.js 18+ (optional, for Netlify CLI)
- git

### Development Setup | Configuración de Desarrollo

#### English:

1. Clone the repository:
   ```bash
   git clone https://github.com/juansgv/ana.git
   cd ana
   ```

2. Install Go dependencies:
   ```bash
   go mod download
   ```
3. Set up PostgreSQL (follow instructions in database_setup.md)

4. Run the server:
   ```bash
   go run cmd/server/main.go
   ```

5. Access the web interface at http://localhost:8080

#### Español:

1. Clona el repositorio:
   ```bash
   git clone https://github.com/juansgv/ana.git
   cd ana
   ```

2. Instala las dependencias de Go:
   ```bash
   go mod download
   ```
3. Configura PostgreSQL (sigue las instrucciones en database_setup.md)

4. Ejecuta el servidor:
   ```bash
   go run cmd/server/main.go
   ```

5. Accede a la interfaz web en http://localhost:8080

## Environment Variables | Variables de Entorno

### English:

The application uses the following environment variables:

- **Core Application**:
  - `ANA_ENVIRONMENT`: Deployment environment (`dev`, `test`, or `prod`)
  - `ANA_VERSION`: Application version
  - `DATABASE_URL`: PostgreSQL connection string

- **AI Integration**:
  - `CEREBRAS_API_KEY`: API key for Cerebras AI service
    - For development: Use a development API key
    - For testing: Use a test API key with limited quotas
    - For production: Use a production API key with full access
  - `CEREBRAS_API_URL`: (Optional) Custom endpoint for Cerebras API
  - `CEREBRAS_CACHE_TTL`: (Optional) Cache duration for AI responses

- **Google Services**:
  - `GOOGLE_API_KEY`: API key for Google services
  - `GOOGLE_OAUTH_CLIENT_ID`: OAuth client ID for Google Calendar integration

### Español:

La aplicación utiliza las siguientes variables de entorno:

- **Aplicación Principal**:
  - `ANA_ENVIRONMENT`: Entorno de despliegue (`dev`, `test`, o `prod`)
  - `ANA_VERSION`: Versión de la aplicación
  - `DATABASE_URL`: Cadena de conexión PostgreSQL

- **Integración de IA**:
  - `CEREBRAS_API_KEY`: Clave API para el servicio Cerebras AI
    - Para desarrollo: Usa una clave API de desarrollo
    - Para pruebas: Usa una clave API de prueba con cuotas limitadas
    - Para producción: Usa una clave API de producción con acceso completo
  - `CEREBRAS_API_URL`: (Opcional) Endpoint personalizado para la API de Cerebras
  - `CEREBRAS_CACHE_TTL`: (Opcional) Duración del caché para respuestas de IA

- **Servicios de Google**:
  - `GOOGLE_API_KEY`: Clave API para servicios de Google
  - `GOOGLE_OAUTH_CLIENT_ID`: ID de cliente OAuth para integración con Google Calendar

## Project Origin | Origen del Proyecto

This project was initially conceived in a Grok.com workspace named "ANA", where the comprehensive plan was developed. It aims to solve the specific challenges faced by architects managing multiple construction projects in Bogota, Colombia.
The development flow is documented in [flow.md](flow.md), and the detailed plan in [plan.md](plan.md). The next implementation steps are outlined in [next.md](next.md). Database setup instructions can be found in [database_setup.md](database_setup.md). For details about the AI assistant integration, see [ai_integration.md](ai_integration.md).

Este proyecto fue concebido inicialmente en un espacio de trabajo de Grok.com llamado "ANA", donde se desarrolló el plan integral. Su objetivo es resolver los desafíos específicos que enfrentan los arquitectos que gestionan múltiples proyectos de construcción en Bogotá, Colombia.

El flujo de desarrollo está documentado en [flow.md](flow.md), y el plan detallado en [plan.md](plan.md). Los próximos pasos de implementación se describen en [next.md](next.md). Las instrucciones para la configuración de la base de datos se encuentran en [database_setup.md](database_setup.md). Para obtener detalles sobre la integración del asistente de IA, consulta [ai_integration.md](ai_integration.md).

## Current Status | Estado Actual

The project has successfully progressed through Phase 1 and is now in Phase 2 of development. We have implemented: | El proyecto ha progresado exitosamente a través de la Fase 1 y está ahora en la Fase 2 de desarrollo. Hemos implementado:

### Core Infrastructure | Infraestructura Central
- Complete project structure with Go/Gin backend | Estructura completa del proyecto con backend Go/Gin
- HTMX and Tailwind CSS frontend integration | Integración frontend con HTMX y Tailwind CSS
- PostgreSQL database with GORM ORM | Base de datos PostgreSQL con GORM ORM
- Repository pattern for data access | Patrón repositorio para acceso a datos
- Environment-based configuration | Configuración basada en variables de entorno

### AI Integration | Integración de IA
- **Cerebras AI assistant** with dual-model support: | **Asistente de IA Cerebras** con soporte de modelos duales:
  - QWen-3B-32B for architectural expertise | QWen-3B-32B para experiencia arquitectónica
  - QWen-2.5-Vision for image analysis | QWen-2.5-Vision para análisis de imágenes
  - Context-aware responses in Spanish/English | Respuestas contextuales en Español/Inglés
  - Optimized caching and error handling | Caché optimizado y manejo de errores
  - See [ai_integration.md](ai_integration.md) and [cerebras_api.md](cerebras_api.md) for details | Ver [ai_integration.md](ai_integration.md) y [cerebras_api.md](cerebras_api.md) para detalles

### Testing Infrastructure | Infraestructura de Pruebas
- **Comprehensive test coverage**: | **Cobertura completa de pruebas**:
  - Unit tests for core functionality | Pruebas unitarias para funcionalidad central
  - Handler tests with mock implementations | Pruebas de manejadores con implementaciones simuladas
  - Repository tests with SQLite in-memory database | Pruebas de repositorio con base de datos SQLite en memoria
  - Integration tests for API endpoints | Pruebas de integración para endpoints API
  - See [tests.md](tests.md) for details | Ver [tests.md](tests.md) para detalles

### Monitoring & Operations | Monitoreo y Operaciones
- Prometheus metrics collection | Recolección de métricas con Prometheus
- Health check endpoints | Endpoints de verificación de salud
- Error tracking and logging | Seguimiento y registro de errores

### In Progress | En Progreso
- Google Calendar/Gmail integration | Integración con Google Calendar/Gmail
- Enhanced project management features | Características mejoradas de gestión de proyectos
- Location-based supplier search | Búsqueda de proveedores basada en ubicación

## Contact | Contacto

For questions or contributions, please contact:

Para preguntas o contribuciones, por favor contacta:

- **Project Owner**: Ana Architect
- **Email**: [ana@example.com](mailto:ana@example.com)
- **GitHub**: [@juansgv](https://github.com/juansgv)

## License | Licencia

This project is proprietary and intended for specific use. Not licensed for public redistribution.

Este proyecto es propietario y está destinado para uso específico. No está licenciado para redistribución pública.

