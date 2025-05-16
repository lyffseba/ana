# Cerebras AI Integration | Integración de IA Cerebras

[English](#english) | [Español](#español)

<a name="english"></a>
## Cerebras AI Assistant Integration

The ana.world platform integrates with Cerebras AI to provide intelligent assistance for architectural project management. This document explains how to set up, configure, and use the AI assistant features.

### Setup Instructions

1. **Get an API Key from Cerebras**
   - Visit [Cerebras AI](https://inference.cerebras.ai) and create an account
   - Generate an API key from your dashboard
   - Copy the API key for configuration

2. **Configure Environment Variables**
   - Add your Cerebras API key to the `.env` file:
     ```
     CEREBRAS_API_KEY=your_api_key_here
     CEREBRAS_API_URL=https://inference.cerebras.ai/v1/chat/completions
     ```

3. **Test the Integration**
   - Run the provided test script to verify API connectivity:
     ```bash
     bash scripts/test_cerebras_api.sh
     ```
   - The script will test both direct API access and the local server endpoint

### API Endpoint Documentation

The Cerebras AI assistant is accessible via a RESTful API endpoint:

- **Endpoint**: `/api/ai/cerebras`
- **Method**: POST
- **Content-Type**: application/json
- **Request Body**:
  ```json
  {
    "query": "Your question or request here"
  }
  ```
- **Response**:
  ```json
  {
    "response": "AI assistant's response text"
  }
  ```
- **Error Response**:
  ```json
  {
    "error": "Error message description"
  }
  ```

### Implementation Details

The integration consists of several components:

1. **Client Implementation**: `internal/ai/cerebras_client.go`
   - Handles API communication with Cerebras
   - Manages system context for architectural domain
   - Provides error handling and fallbacks

2. **Handler**: `internal/handlers/ai_cerebras_handler.go`
   - Processes HTTP requests to the AI endpoint
   - Validates inputs and formats responses
   - Logs errors and manages state

3. **Frontend Integration**: `web/index.html`
   - Provides a chat interface for interacting with the AI
   - Uses HTMX for real-time updates without page reload
   - Formats and displays AI responses

### Example Usage

#### Curl Example
```bash
curl -X POST http://localhost:8080/api/ai/cerebras \
  -H "Content-Type: application/json" \
  -d '{"query": "¿Cómo puedo organizar mejor mis proyectos de arquitectura?"}'
```

#### Frontend Example
The AI assistant is integrated into the main dashboard with a chat interface. Users can:
1. Type questions in the input box
2. Receive contextual answers about project management
3. Get assistance with scheduling, task organization, and architectural best practices

#### Example System Prompts
The AI is configured with architectural domain knowledge to assist with:
- Project scheduling and timeline management
- Building code compliance and regulations
- Material selection and sourcing
- Budget estimation and tracking
- Client communication strategies

### Troubleshooting

- **API Key Issues**: If you encounter authentication errors, verify your API key is correctly set in the .env file
- **Rate Limiting**: Cerebras API has rate limits; the system includes exponential backoff retry
- **Connectivity Problems**: Check your network and firewall settings if connections fail

---

<a name="español"></a>
## Integración del Asistente IA de Cerebras

La plataforma ana.world se integra con Cerebras AI para proporcionar asistencia inteligente para la gestión de proyectos arquitectónicos. Este documento explica cómo configurar y utilizar las funciones del asistente de IA.

### Instrucciones de Configuración

1. **Obtener una Clave API de Cerebras**
   - Visita [Cerebras AI](https://inference.cerebras.ai) y crea una cuenta
   - Genera una clave API desde tu panel de control
   - Copia la clave API para la configuración

2. **Configurar Variables de Entorno**
   - Añade tu clave API de Cerebras al archivo `.env`:
     ```
     CEREBRAS_API_KEY=tu_clave_api_aquí
     CEREBRAS_API_URL=https://inference.cerebras.ai/v1/chat/completions
     ```

3. **Probar la Integración**
   - Ejecuta el script de prueba proporcionado para verificar la conectividad de la API:
     ```bash
     bash scripts/test_cerebras_api.sh
     ```
   - El script probará tanto el acceso directo a la API como el endpoint del servidor local

### Documentación del Endpoint API

El asistente de IA Cerebras es accesible a través de un endpoint API RESTful:

- **Endpoint**: `/api/ai/cerebras`
- **Método**: POST
- **Content-Type**: application/json
- **Cuerpo de la Solicitud**:
  ```json
  {
    "query": "Tu pregunta o solicitud aquí"
  }
  ```
- **Respuesta**:
  ```json
  {
    "response": "Texto de respuesta del asistente de IA"
  }
  ```
- **Respuesta de Error**:
  ```json
  {
    "error": "Descripción del mensaje de error"
  }
  ```

### Detalles de Implementación

La integración consta de varios componentes:

1. **Implementación del Cliente**: `internal/ai/cerebras_client.go`
   - Maneja la comunicación API con Cerebras
   - Gestiona el contexto del sistema para el dominio arquitectónico
   - Proporciona manejo de errores y alternativas

2. **Manejador**: `internal/handlers/ai_cerebras_handler.go`
   - Procesa solicitudes HTTP al endpoint de IA
   - Valida entradas y formatea respuestas
   - Registra errores y gestiona el estado

3. **Integración Frontend**: `web/index.html`
   - Proporciona una interfaz de chat para interactuar con la IA
   - Utiliza HTMX para actualizaciones en tiempo real sin recarga de página
   - Formatea y muestra respuestas de IA

### Ejemplos de Uso

#### Ejemplo con Curl
```bash
curl -X POST http://localhost:8080/api/ai/cerebras \
  -H "Content-Type: application/json" \
  -d '{"query": "¿Cómo puedo organizar mejor mis proyectos de arquitectura?"}'
```

#### Ejemplo de Frontend
El asistente de IA está integrado en el dashboard principal con una interfaz de chat. Los usuarios pueden:
1. Escribir preguntas en el cuadro de entrada
2. Recibir respuestas contextuales sobre gestión de proyectos
3. Obtener asistencia con programación, organización de tareas y mejores prácticas arquitectónicas

#### Ejemplos de Prompts del Sistema
La IA está configurada con conocimiento del dominio arquitectónico para ayudar con:
- Programación de proyectos y gestión de cronogramas
- Cumplimiento de códigos de construcción y regulaciones
- Selección y abastecimiento de materiales
- Estimación y seguimiento de presupuestos
- Estrategias de comunicación con clientes

### Solución de Problemas

- **Problemas con la Clave API**: Si encuentras errores de autenticación, verifica que tu clave API esté correctamente configurada en el archivo .env
- **Límites de Tasa**: La API de Cerebras tiene límites de tasa; el sistema incluye reintentos con retroceso exponencial
- **Problemas de Conectividad**: Comprueba tu red y configuración de firewall si las conexiones fallan

