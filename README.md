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

## Project Origin | Origen del Proyecto

This project was initially conceived in a Grok.com workspace named "ANA", where the comprehensive plan was developed. It aims to solve the specific challenges faced by architects managing multiple construction projects in Bogota, Colombia.
The development flow is documented in [flow.md](flow.md), and the detailed plan in [plan.md](plan.md). The next implementation steps are outlined in [next.md](next.md). Database setup instructions can be found in [database_setup.md](database_setup.md).

Este proyecto fue concebido inicialmente en un espacio de trabajo de Grok.com llamado "ANA", donde se desarrolló el plan integral. Su objetivo es resolver los desafíos específicos que enfrentan los arquitectos que gestionan múltiples proyectos de construcción en Bogotá, Colombia.

El flujo de desarrollo está documentado en [flow.md](flow.md), y el plan detallado en [plan.md](plan.md). Los próximos pasos de implementación se describen en [next.md](next.md). Las instrucciones para la configuración de la base de datos se encuentran en [database_setup.md](database_setup.md).

## Current Status | Estado Actual

The project is currently in Phase 1 of development. We have implemented:

- Basic project structure
- Initial API endpoints for task management
- Frontend with HTMX and Tailwind CSS integration
- Simple routing for API and static files
- PostgreSQL database integration with GORM ORM
- Repository pattern for data access layer
- Environment-based configuration
- Documentation for database setup
- Planning for Cerebras AI assistant integration

El proyecto se encuentra actualmente en la Fase 1 de desarrollo. Hemos implementado:

- Estructura básica del proyecto
- Puntos finales iniciales de API para gestión de tareas
- Frontend con integración de HTMX y Tailwind CSS
- Enrutamiento simple para API y archivos estáticos
- Integración de base de datos PostgreSQL con GORM ORM
- Patrón de repositorio para capa de acceso a datos
- Configuración basada en variables de entorno
- Documentación para la configuración de la base de datos
- Planificación para la integración del asistente IA de Cerebras

## Contact | Contacto

For questions or contributions, please contact:

Para preguntas o contribuciones, por favor contacta:

- **Project Owner**: Ana Architect
- **Email**: [ana@example.com](mailto:ana@example.com)
- **GitHub**: [@juansgv](https://github.com/juansgv)

## License | Licencia

This project is proprietary and intended for specific use. Not licensed for public redistribution.

Este proyecto es propietario y está destinado para uso específico. No está licenciado para redistribución pública.

