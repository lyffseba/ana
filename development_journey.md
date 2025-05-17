# ana.world Development Journey

[English](#english) | [Español](#español)

<a name="english"></a>
## Development Journey Documentation

This document chronicles the complete development process of the ana.world project, from initial conception to current implementation, capturing both the technical and narrative aspects of our journey.

### Project Inception and Vision

The ana.world project was born in March 2025 out of a specific need identified by Ana, an architect in Bogotá, Colombia, who required a specialized tool to manage multiple architectural projects simultaneously. What started as a simple task management concept evolved into a comprehensive architectural project management platform.

Initial requirements gathering and concept development took place in a Grok.com workspace named "ANA", where we collaborated on defining the scope, features, and technology stack. This phase was crucial in establishing a shared vision for what would become ana.world.

### Development Timeline and Phases

#### Phase 1: Foundation (March-April 2025)

The first phase focused on establishing the core architecture and basic functionality:

1. **Project Structure Setup** (March 15-20)
   - Created initial GitHub repository architecture
   - Established Go + Gin framework backend structure
   - Set up static frontend with HTMX and Tailwind CSS
   - Implemented basic routing

2. **Core Task Management** (March 21-31)
   - Implemented initial task model and CRUD operations
   - Created basic API endpoints
   - Designed UI for task management

3. **Database Integration** (April 1-10)
   - Set up PostgreSQL database structure
   - Implemented GORM ORM integration
   - Created repository pattern for data access
   - Developed database migration strategy

#### Phase 2: AI Assistant Integration (April-May 2025)

This phase introduced intelligent assistance capabilities:

1. **Cerebras API Integration** (April 11-25)
   - Researched and selected Cerebras AI for domain expertise
   - Implemented initial AI client with QWen-3B-32B model
   - Created system context with architectural knowledge
   - Developed error handling and retry mechanisms

2. **Advanced AI Features** (April 26-May 10)
   - Added QWen-2.5-Vision model for image analysis
   - Implemented model switching in UI
   - Added support for direct responses with /no_think command
   - Enhanced architectural domain knowledge in system prompts

3. **Testing and Documentation** (May 11-17)
   - Developed comprehensive test suite with mocks
   - Created detailed documentation
   - Improved error handling
   - Enhanced user experience

#### Phase 3: Google Integration and Refinement (Planned May-June 2025)

The next planned phase includes:

1. **Google Calendar Integration**
   - OAuth2 authentication flow
   - Task synchronization with calendar events
   - Reminder system

2. **Gmail Integration**
   - Email notifications for tasks
   - Email-based task creation
   - Communication tracking

### Tools and Technologies Used

#### Development Tools

| Tool/Technology | Purpose | Implementation Details |
|----------------|---------|------------------------|
| **Go 1.22.2** | Backend language | Used with Gin framework for API development |
| **Gin Framework** | Web framework | Handles routing, middleware, and request processing |
| **GORM** | ORM | Object-relational mapping for PostgreSQL |
| **HTMX** | Frontend interactivity | Provides dynamic UI without heavy JavaScript |
| **Tailwind CSS** | Styling | Theme based on terra colors for architectural aesthetic |
| **PostgreSQL** | Database | Persistent storage for application data |
| **SQLite** | Testing | In-memory database for repository testing |
| **Cerebras AI** | AI assistance | QWen models for architectural assistance |

#### Collaboration and Version Control

| Tool | Purpose | How We Used It |
|------|---------|----------------|
| **GitHub** | Version control | Repository: https://github.com/juansgv/ana |
| **Grok.com** | Initial planning | Created comprehensive project plan and user stories |
| **Warp** | Terminal & development sessions | Recorded development process in shared sessions |

### Architectural Decisions and Implementation Details

#### 1. Backend Architecture

We chose a layered architecture for the Go backend:

```
ana/
├── cmd/                    # Application entry points
│   └── server/             # Main server binary
├── internal/               # Internal packages
│   ├── handlers/           # HTTP handlers
│   ├── models/             # Data models
│   ├── repositories/       # Data access layer
│   ├── ai/                 # AI integration 
│   ├── database/           # Database connection
│   └── server/             # Server configuration
```

**Key Decisions:**
- **Repository Pattern**: Isolates data access logic, making testing easier
- **Dependency Injection**: Used throughout for better testability
- **Middleware Approach**: For cross-cutting concerns like authentication and logging
- **Error Handling Strategy**: Consistent approach with user-friendly messages in Spanish

#### 2. Frontend Approach

We chose a minimalist JavaScript approach with HTMX:

**Key Decisions:**
- **HTMX over SPA**: For simplicity and reduced JavaScript complexity
- **Tailwind CSS**: For rapid UI development with consistent styling
- **Progressive Enhancement**: Basic functionality works without JavaScript
- **Form Design**: Multipart forms to support image uploads for AI vision model

#### 3. AI Integration Architecture

Our AI implementation evolved from a simple text model to a dual-model system:

**Key Decisions:**
- **Dual-Model Support**: Added both text and vision capabilities
- **Interface-Based Design**: Created interfaces for better testing
- **Error Resilience**: Comprehensive error handling with user-friendly messages
- **Domain Knowledge**: Specialized system prompts for architectural expertise

### Development Process and Methodology

#### Iterative Development Approach

We followed an iterative development process:

1. **Planning & Research Phase**
   - Define feature requirements
   - Research technical solutions
   - Create implementation plan

2. **Implementation Phase**
   - Develop backend functionality
   - Create frontend components
   - Integrate systems

3. **Testing Phase**
   - Unit testing with mocks
   - Integration testing
   - Manual testing

4. **Documentation & Refinement**
   - Update documentation
   - Refine implementation
   - Plan next iteration

#### Test-Driven Development

For critical components, we employed a test-driven development approach:

1. Write failing tests defining expected behavior
2. Implement minimal code to pass tests
3. Refactor while maintaining test coverage
4. Repeat for new features

### Complete Development Session Records

The full development process has been recorded in Warp terminal sessions, providing a comprehensive view of the entire development journey, including:

- Command history
- Error resolution process
- Decision making
- Implementation challenges

**Full session logs**: [https://app.warp.dev/session/4a2816eb-5f73-42e1-8254-cb8f93eaa418?pwd=7a4dd4cd-b725-4721-9a1a-6be54aacfc70](https://app.warp.dev/session/4a2816eb-5f73-42e1-8254-cb8f93eaa418?pwd=7a4dd4cd-b725-4721-9a1a-6be54aacfc70)

### Challenges and Solutions

#### Challenge 1: AI Model Selection and Integration

**Problem**: The initial AI implementation only supported text queries, limiting architectural assistance.

**Solution**: 
- Researched and identified QWen models that support both text and vision
- Implemented dual-model support with model switching
- Added image upload functionality
- Enhanced system prompts with architectural domain knowledge

#### Challenge 2: Testing Strategy

**Problem**: Testing HTTP handlers and AI integration without making real API calls.

**Solution**:
- Developed a comprehensive mocking strategy
- Created interfaces for dependency injection
- Implemented table-driven tests for multiple scenarios
- Used SQLite in-memory database for repository testing

#### Challenge 3: Bilingual Support

**Problem**: Supporting both English and Spanish throughout the application.

**Solution**:
- Implemented system prompts that understand English but respond in Spanish
- Created error messages in Spanish for better user experience
- Made documentation bilingual
- Designed UI to support both languages

### Future Development Plans

The next phases of development will focus on:

1. **Google Integration**
   - Calendar synchronization for project milestones
   - Email integration for notifications
   - Document storage for architectural plans

2. **Enhanced AI Capabilities**
   - Plan analysis from images
   - Cost estimation from architectural drawings
   - Construction timeline prediction

3. **Mobile Optimization**
   - Responsive design enhancements
   - Progressive Web App implementation
   - Offline capabilities

### Conclusions and Lessons Learned

The ana.world project has demonstrated several important lessons:

1. **Start with a Clear User Focus**: Understanding Ana's specific needs as an architect in Bogotá shaped the entire development process.

2. **Choose the Right Tools**: The combination of Go, HTMX, and Tailwind CSS allowed for rapid development with minimal complexity.

3. **Plan for Evolution**: The system architecture was designed to evolve, which proved valuable when adding new capabilities like vision AI.

4. **Test Comprehensively**: The investment in testing infrastructure paid dividends when implementing new features.

5. **Document Throughout**: Maintaining comprehensive documentation made onboarding and feature development smoother.

The development of ana.world continues to be an evolving journey, with each phase building upon the foundations laid in earlier work, always with the goal of creating the most effective architectural project management tool possible.

---

<a name="español"></a>
## Documentación del Proceso de Desarrollo

Este documento narra el proceso completo de desarrollo del proyecto ana.world, desde su concepción inicial hasta la implementación actual, capturando tanto los aspectos técnicos como narrativos de nuestro viaje.

### Inicio del Proyecto y Visión

El proyecto ana.world nació en marzo de 2025 a partir de una necesidad específica identificada por Ana, una arquitecta en Bogotá, Colombia, que requería una herramienta especializada para gestionar múltiples proyectos arquitectónicos simultáneamente. Lo que comenzó como un simple concepto de gestión de tareas evolucionó hasta convertirse en una plataforma integral de gestión de proyectos arquitectónicos.

La recopilación inicial de requisitos y el desarrollo del concepto se llevaron a cabo en un espacio de trabajo de Grok.com llamado "ANA", donde colaboramos en la definición del alcance, las características y la pila tecnológica. Esta fase fue crucial para establecer una visión compartida de lo que se convertiría en ana.world.

### Cronología y Fases de Desarrollo

#### Fase 1: Fundación (Marzo-Abril 2025)

La primera fase se centró en establecer la arquitectura central y la funcionalidad básica:

1. **Configuración de la Estructura del Proyecto** (15-20 de marzo)
   - Creación de la arquitectura inicial del repositorio GitHub
   - Establecimiento de la estructura backend Go + Gin framework
   - Configuración del frontend estático con HTMX y Tailwind CSS
   - Implementación de enrutamiento básico

2. **Gestión de Tareas Central** (21-31 de marzo)
   - Implementación del modelo de tareas inicial y operaciones CRUD
   - Creación de endpoints API básicos
   - Diseño de UI para gestión de tareas

3. **Integración de Base de Datos** (1-10 de abril)
   - Configuración de la estructura de base de datos PostgreSQL
   - Implementación de integración GORM ORM
   - Creación del patrón de repositorio para acceso a datos
   - Desarrollo de estrategia de migración de base de datos

#### Fase 2: Integración de Asistente de IA (Abril-Mayo 2025)

Esta fase introdujo capacidades de asistencia inteligente:

1. **Integración de API Cerebras** (11-25 de abril)
   - Investigación y selección de Cerebras AI por su experiencia en el dominio
   - Implementación de cliente AI inicial con modelo QWen-3B-32B
   - Creación de contexto de sistema con conocimiento arquitectónico
   - Desarrollo de manejo de errores y mecanismos de reintento

2. **Características Avanzadas de IA** (26 de abril-10 de mayo)
   - Adición del modelo QWen-2.5-Vision para análisis de imágenes
   - Implementación de cambio de modelo en la UI
   - Adición de soporte para respuestas directas con comando /no_think
   - Mejora del conocimiento del dominio arquitectónico en los prompts del sistema

3. **Pruebas y Documentación** (11-17 de mayo)
   - Desarrollo de suite de pruebas completa con mocks
   - Creación de documentación detallada
   - Mejora del manejo de errores
   - Mejora de la experiencia de usuario

#### Fase 3: Integración de Google y Refinamiento (Planificado Mayo-Junio 2025)

La siguiente fase planificada incluye:

1. **Integración de Google Calendar**
   - Flujo de autenticación OAuth2
   - Sincronización de tareas con eventos del calendario
   - Sistema de recordatorios

2. **Integración de Gmail**
   - Notificaciones por correo electrónico para tareas
   - Creación de tareas basada en correos electrónicos
   - Seguimiento de comunicaciones

### Herramientas y Tecnologías Utilizadas

#### Herramientas de Desarrollo

| Herramienta/Tecnología | Propósito | Detalles de Implementación |
|------------------------|-----------|----------------------------|
| **Go 1.22.2** | Lenguaje backend | Utilizado con framework Gin para desarrollo API |
| **Gin Framework** | Framework web | Maneja enrutamiento, middleware y procesamiento de solicitudes |
| **GORM** | ORM | Mapeo objeto-relacional para PostgreSQL |
| **HTMX** | Interactividad frontend | Proporciona UI dinámica sin JavaScript pesado |
| **Tailwind CSS** | Estilización | Tema basado en colores terra para estética arquitectónica |
| **PostgreSQL** | Base de datos | Almacenamiento persistente para datos de aplicación |
| **SQLite** | Pruebas | Base de datos en memoria para pruebas de repositorio |
| **Cerebras AI** | Asistencia de IA | Modelos QWen para asistencia arquitectónica |

#### Colaboración y Control de Versiones

| Herramienta | Propósito | Cómo La Usamos |
|-------------|-----------|----------------|
| **GitHub** | Control de versiones | Repositorio: https://github.com/juansgv/ana |
| **Grok.com** | Planificación inicial | Creación de plan de proyecto integral e historias de usuario |
| **Warp** | Terminal y sesiones de desarrollo | Registro del proceso de desarrollo en sesiones compartidas |

### Decisiones Arquitectónicas y Detalles de Implementación

#### 1. Arquitectura Backend

Elegimos una arquitectura en capas para el backend Go:

```
ana/
├── cmd/                    # Puntos de entrada de la aplicación
│   └── server/             # Binario principal del servidor
├── internal/               # Paquetes internos
│   ├── handlers/           # Manejadores HTTP
│   ├── models/             # Modelos de datos
│   ├── repositories/       # Capa de acceso a datos
│   ├── ai/                 # Integración de IA
│   ├── database/           # Conexión a base de datos
│   └── server/             # Configuración del servidor
```

**Decisiones Clave:**
- **Patrón Repositorio**: Aísla la lógica de acceso a datos, facilitando las pruebas
- **Inyección de Dependencias**: Utilizada en toda la aplicación para mejor capacidad de prueba
- **Enfoque de Middleware**: Para preocupaciones transversales como autenticación y registro
- **Estrategia de Manejo de Errores**: Enfoque

