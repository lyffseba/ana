# Google Calendar & Gmail Integration for ana.world

[English](#english) | [Español](#español)

<a name="english"></a>

## Introduction

This document outlines the integration plan for Google Calendar and Gmail services with ana.world, specifically configured for the email account `anais.villamarinj@gmail.com`. This integration will enable architects to automatically sync project tasks with their Google Calendar and receive email notifications for important deadlines and updates.

## Prerequisites

### Google Cloud Project Setup

1. **Create a Google Cloud Project**
   - Go to [Google Cloud Console](https://console.cloud.google.com/)
   - Create a new project named "Ana World Architecture"
   - Note your Project ID for configuration files

2. **Required API Enablement**
   - Enable Google Calendar API
   - Enable Gmail API
   - Enable Google Drive API (optional, for future document integration)

3. **OAuth 2.0 Configuration**
   - Set up OAuth consent screen:
     - Application name: Ana World
     - User support email: anais.villamarinj@gmail.com
     - Developer contact information: anais.villamarinj@gmail.com
   - Create OAuth 2.0 Client ID credentials for web application
   - Add authorized redirect URIs:
     - `http://localhost:8080/api/auth/google/callback` (development)
     - `https://ana.world/api/auth/google/callback` (production)
   - Download credentials file as `credentials.json`

## Email Configuration for anais.villamarinj@gmail.com

### Gmail Setup

1. **Account Configuration**
   - Sign in to Gmail with anais.villamarinj@gmail.com
   - In Gmail settings, ensure POP/IMAP access is enabled
   - Configure "Allow less secure apps" if required, or preferably use OAuth 2.0

2. **Application Permissions**
   - The first time the application connects, anais.villamarinj@gmail.com will need to approve the following permissions:
     - View and manage Google Calendar
     - Send emails through Gmail
     - View and create emails (for creating task-related emails)

3. **Email Templates Setup**
   - Create email templates for:
     - Task assignments
     - Task reminders
     - Project milestone notifications
     - Weekly summary reports

## Integration Features

### 1. Calendar Event Syncing

- **Automatic Task to Event Conversion**
  - Each task with a due date will be converted to a Google Calendar event
  - Task priorities will be reflected in event colors:
    - High priority: Red
    - Medium priority: Yellow
    - Low priority: Blue
  - Task locations will be added to event locations

- **Bidirectional Synchronization**
  - Changes made in ana.world will update Google Calendar events
  - Changes to event times in Google Calendar will update task due dates in ana.world
  - Deleted tasks will remove calendar events and vice versa

- **Recurring Task Support**
  - Support for recurring tasks will create recurring events in Google Calendar
  - Frequency options: daily, weekly, monthly, custom

### 2. Email Notifications

- **Automated Email Notifications**
  - Task assignment notifications
  - Due date reminder emails (1 day, 1 hour before)
  - Overdue task notifications
  - Task completion confirmations

- **Customizable Email Preferences**
  - Users can select which notifications to receive
  - Adjustable timing for reminders
  - Option for daily or weekly digest emails

- **Email to Task Conversion**
  - Forward emails to a specific address to create tasks
  - Use email subject for task title
  - Use email body for task description
  - Parse dates in email to set due dates

### 3. Two-Way Task Synchronization

- **Google Tasks Integration**
  - Sync ana.world tasks with Google Tasks
  - Allow managing tasks from Google environment
  - Maintain consistent task status across platforms

- **Batch Synchronization**
  - Daily scheduled sync of all tasks
  - On-demand sync through UI button
  - Real-time sync for critical task changes

- **Conflict Resolution**
  - Smart merging of conflicting changes
  - Last-update-wins strategy with history tracking
  - Notification of sync conflicts requiring manual resolution

## Implementation Timeline and Milestones

### Phase 1: Setup and Authentication (Week 1)

- [  ] Set up Google Cloud Project
- [  ] Enable required APIs
- [  ] Configure OAuth 2.0
- [  ] Implement authentication flow in backend
- [  ] Create token storage with encryption

### Phase 2: Calendar Integration (Week 2)

- [  ] Implement basic Calendar API communication
- [  ] Create task-to-event conversion
- [  ] Implement one-way sync (ana.world to Calendar)
- [  ] Add event color mapping for task priorities
- [  ] Test with anais.villamarinj@gmail.com account

### Phase 3: Gmail Integration (Week 3)

- [  ] Set up Gmail API communication
- [  ] Implement email notification templates
- [  ] Create notification triggering system
- [  ] Add email preference settings
- [  ] Test notification system

### Phase 4: Two-Way Sync and Advanced Features (Week 4)

- [  ] Implement bidirectional synchronization
- [  ] Add conflict resolution logic
- [  ] Create email-to-task conversion
- [  ] Implement recurring task support
- [  ] Add Google Tasks integration

### Phase 5: Testing and Optimization (Week 5)

- [  ] Comprehensive end-to-end testing
- [  ] Performance optimization
- [  ] Security audit and improvements
- [  ] Documentation updates
- [  ] User acceptance testing

## Security Considerations

- Token storage will use secure encryption
- Scope permissions will be limited to only what's needed
- Regular token rotation will be implemented
- All API communications will use HTTPS
- User data will be stored in compliance with privacy regulations

## Technical Implementation Details

### Go Packages Required

```
go get -u google.golang.org/api/calendar/v3
go get -u google.golang.org/api/gmail/v1
go get -u golang.org/x/oauth2/google
```

### Environment Variables

```
GOOGLE_CLIENT_ID=your_client_id
GOOGLE_CLIENT_SECRET=your_client_secret
GOOGLE_REDIRECT_URI=http://localhost:8080/api/auth/google/callback
GOOGLE_PRIMARY_EMAIL=anais.villamarinj@gmail.com
```

### Testing Scripts

A testing script will be provided to verify the integration:

```bash
#!/bin/bash
# Test Google Calendar and Gmail integration

# Test OAuth authentication flow
go run cmd/test/google_auth.go

# Test Calendar API connection and event creation
go run cmd/test/calendar_test.go

# Test Gmail API connection and email sending
go run cmd/test/gmail_test.go

# Test bidirectional sync
go run cmd/test/sync_test.go
```

---

<a name="español"></a>

# Integración de Google Calendar y Gmail para ana.world

## Introducción

Este documento describe el plan de integración de los servicios de Google Calendar y Gmail con ana.world, configurado específicamente para la cuenta de correo electrónico `anais.villamarinj@gmail.com`. Esta integración permitirá a los arquitectos sincronizar automáticamente las tareas del proyecto con su Google Calendar y recibir notificaciones por correo electrónico para fechas límite y actualizaciones importantes.

## Requisitos Previos

### Configuración del Proyecto de Google Cloud

1. **Crear un Proyecto de Google Cloud**
   - Ir a [Google Cloud Console](https://console.cloud.google.com/)
   - Crear un nuevo proyecto llamado "Ana World Architecture"
   - Anotar el ID del Proyecto para los archivos de configuración

2. **Habilitación de APIs Requeridas**
   - Habilitar Google Calendar API
   - Habilitar Gmail API
   - Habilitar Google Drive API (opcional, para futura integración de documentos)

3. **Configuración de OAuth 2.0**
   - Configurar la pantalla de consentimiento de OAuth:
     - Nombre de la aplicación: Ana World
     - Correo electrónico de soporte: anais.villamarinj@gmail.com
     - Información de contacto del desarrollador: anais.villamarinj@gmail.com
   - Crear credenciales de ID de cliente OAuth 2.0 para aplicación web
   - Agregar URIs de redirección autorizados:
     - `http://localhost:8080/api/auth/google/callback` (desarrollo)
     - `https://ana.world/api/auth/google/callback` (producción)
   - Descargar el archivo de credenciales como `credentials.json`

## Configuración de Correo Electrónico para anais.villamarinj@gmail.com

### Configuración de Gmail

1. **Configuración de la Cuenta**
   - Iniciar sesión en Gmail con anais.villamarinj@gmail.com
   - En la configuración de Gmail, asegurarse de que el acceso POP/IMAP esté habilitado
   - Configurar "Permitir aplicaciones menos seguras" si es necesario, o preferiblemente usar OAuth 2.0

2. **Permisos de Aplicación**
   - La primera vez que la aplicación se conecte, anais.villamarinj@gmail.com deberá aprobar los siguientes permisos:
     - Ver y administrar Google Calendar
     - Enviar correos electrónicos a través de Gmail
     - Ver y crear correos electrónicos (para crear correos relacionados con tareas)

3. **Configuración de Plantillas de Correo**
   - Crear plantillas de correo electrónico para:
     - Asignaciones de tareas
     - Recordatorios de tareas
     - Notificaciones de hitos del proyecto
     - Informes resumidos semanales

## Características de Integración

### 1. Sincronización de Eventos de Calendario

- **Conversión Automática de Tarea a Evento**
  - Cada tarea con fecha de vencimiento se convertirá en un evento de Google Calendar
  - Las prioridades de las tareas se reflejarán en los colores de los eventos:
    - Prioridad alta: Rojo
    - Prioridad media: Amarillo
    - Prioridad baja: Azul
  - Las ubicaciones de las tareas se agregarán a las ubicaciones de los eventos

- **Sincronización Bidireccional**
  - Los cambios realizados en ana.world actualizarán los eventos de Google Calendar
  - Los cambios en los tiempos de eventos en Google Calendar actualizarán las fechas de vencimiento de las tareas en ana.world
  - Las tareas eliminadas eliminarán los eventos del calendario y viceversa

- **Soporte para Tareas Recurrentes**
  - El soporte para tareas recurrentes creará eventos recurrentes en Google Calendar
  - Opciones de frecuencia: diaria, semanal, mensual, personalizada

### 2. Notificaciones por Correo Electrónico

- **Notificaciones Automatizadas por Correo**
  - Notificaciones de asignación de tareas
  - Correos de recordatorio de fechas de vencimiento (1 día, 1 hora antes)
  - Notificaciones de tareas vencidas
  - Confirmaciones de finalización de tareas

- **Preferencias de Correo Personalizables**
  - Los usuarios pueden seleccionar qué notificaciones recibir
  - Tiempo ajustable para los recordatorios
  - Opción para correos electrónicos de resumen diario o semanal

- **Conversión de Correo a Tarea**
  - Reenviar correos electrónicos a una dirección específica para crear tareas
  - Usar el asunto del correo para el título de la tarea
  - Usar el cuerpo del correo para la descripción de la tarea
  - Analizar fechas en el correo para establecer fechas de vencimiento

### 3. Sincronización Bidireccional de Tareas

- **Integración con Google Tasks**
  - Sincronizar tareas de ana.world con Google Tasks
  - Permitir gestionar tareas desde el entorno de Google
  - Mantener el estado de las tareas consistente en todas las plataformas

- **Sincronización por Lotes**
  - Sincronización diaria programada de todas las tareas
  - Sincronización bajo demanda a través de botón en la interfaz
  - Sincronización en tiempo real para cambios críticos de tareas

- **Resolución de Conflictos**
  - Combinación inteligente de cambios conflictivos
  - Estrategia de última actualización prevalece con seguimiento de historial
  - Notificación de conflictos de sincronización que requieren resolución manual

## Cronograma de Implementación e Hitos

### Fase 1: Configuración y Autenticación (Semana 1)

- [  ] Configurar Proyecto de Google Cloud
- [  ] Habilitar APIs requeridas
- [  ] Configurar OAuth 2.0
- [  ] Implementar flujo de autenticación en backend
- [  ] Crear almacenamiento de tokens con encriptación

### Fase 2: Integración de Calendario (Semana 2)

- [  ] Implementar comunicación básica con Calendar API
- [  ] Crear conversión de tarea a evento
- [  ] Implementar sincronización unidireccional (ana.world a Calendar)
- [  ] Agregar mapeo de colores de eventos para prioridades de tareas
- [  ] Probar con la cuenta anais.villamarinj@gmail.com

### Fase 3: Integración de Gmail (Semana 3)

- [  ] Configurar comunicación con Gmail API
- [  ] Implementar plantillas de notificación por correo
- [  ] Crear sistema de activación de notificaciones
- [  ] Agregar configuración de preferencias de correo
- [  ] Probar sistema de notificaciones

### Fase 4: Sincronización Bidireccional y Características Avanzadas (Semana 4)

- [  ] Implementar sincronización bidireccional
- [  ] Agregar lógica de resolución de conflictos
- [  ] Crear conversión de correo electrónico a tarea
- [  ] Implementar soporte para tareas recurrentes
- [  ] Agregar integración con Google Tasks

### Fase 5: Pruebas y Optimización (Semana 5)

- [  ] Pruebas exhaustivas de extremo a extremo
- [  ] Optimización de rendimiento
- [  ] Auditoría de seguridad y mejoras
- [  ] Actualizaciones de documentación
- [  ] Pruebas de aceptación de usuario

## Consideraciones de Seguridad

- El almacenamiento de tokens utilizará encriptación segura
- Los permisos de alcance se limitarán solo a lo necesario
- Se implementará rotación regular de tokens
- Todas las comunicaciones API utilizarán HTTPS
- Los datos de los usuarios se almacenarán en cumplimiento con las regulaciones de privacidad

## Detalles Técnicos de Implementación

### Paquetes Go Requeridos

```
go get -u google.golang.org/api/calendar/v3
go get -u google.golang.org/api/gmail/v1
go get -u golang.org/x/oauth2/google
```

### Variables de Entorno

```
GOOGLE_CLIENT_ID=your_client_id
GOOGLE_CLIENT_SECRET=your_client_secret
GOOGLE_REDIRECT_URI=http://localhost:8080/api/auth/google/callback
GOOGLE_PRIMARY_EMAIL=anais.villamarinj@


## Development Session
Reference: https://app.warp.dev/session/b660fd8a-f765-449c-a70c-f8c7b971e3c4?pwd=e9ccd7cb-d8be-494e-a2f2-35469f726896
Last Updated: Sat May 17 07:34:44 AM CEST 2025
