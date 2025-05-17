# Cerebras API Integration for ana.world

[English](#english) | [Español](#español)

<a name="english"></a>
## Cerebras API Integration Documentation

This document provides detailed information about the integration of the Cerebras AI API in the ana.world project, specifically focusing on the QWen-3B-32B model used for architectural assistance.

### Model Specifications: QWen-3B-32B

The QWen-3B-32B model is a cutting-edge language model optimized for technical domains including architecture. Key specifications:

- **Model Size**: 32 billion parameters fine-tuned from a 3B base model
- **Token Speed**: ~2100 tokens/second
- **Context Window**: 4096 tokens
- **Response Quality**: Specialized for technical and domain-specific knowledge
- **Languages**: Multilingual with emphasis on English and Spanish

### Special Features

#### 1. Hybrid Reasoning Modes

QWen-3B-32B supports hybrid reasoning modes that can be activated in your prompts:

- **Standard Mode**: Default response style
- **Analytical Mode**: Triggered by asking for "step-by-step analysis" or "detailed breakdown"
- **Comparative Mode**: Activated by requesting "compare and contrast" or "pros and cons"

Example prompt for analytical mode:
```
Provide a step-by-step analysis of passive cooling techniques for tropical climates in Colombia.
```

#### 2. /no_think Command Support

The model supports the `/no_think` command prefix for queries that need direct, concise answers without explanatory text.

Example:
```
/no_think What is the current NSR-10 seismic coefficient for Bogotá?
```

This will return just the coefficient value without additional explanation.

#### 3. Technical Domain Expertise

The model has been fine-tuned with architectural knowledge, particularly:

- Colombian building codes and regulations
- Architectural design principles
- Project management frameworks
- Sustainable building techniques
- Cost estimation and budgeting

### API Usage

#### Request Format

```json
{
  "model": "cerebras/QWen-3B-32B",
  "messages": [
    {
      "role": "system",
      "content": "System instructions..."
    },
    {
      "role": "user",
      "content": "User query..."
    }
  ],
  "temperature": 0.7,
  "max_tokens": 500
}
```

#### Response Format

```json
{
  "id": "resp_01234567",
  "object": "chat.completion",
  "created": 1715987654,
  "model": "cerebras/QWen-3B-32B",
  "choices": [
    {
      "message": {
        "role": "assistant",
        "content": "Response text..."
      },
      "finish_reason": "stop"
    }
  ],
  "usage": {
    "prompt_tokens": 89,
    "completion_tokens": 156,
    "total_tokens": 245
  }
}
```

### Architecture-Specific Examples

#### Example 1: General Architectural Query

**User Query:**
```
¿Cuáles son las consideraciones principales para el diseño bioclimático en Bogotá?
```

**Model Response:**
```
Para el diseño bioclimático en Bogotá, debes considerar:

1. Clima: Bogotá tiene clima frío de montaña (promedio 14°C) con poca variación estacional pero grandes variaciones diurnas.

2. Orientación: La orientación este-oeste maximiza la captación solar en fachadas norte-sur debido a la posición ecuatorial de Colombia.

3. Radiación solar: Aprovecha la radiación directa para calentamiento pasivo mediante ventanales al oriente/occidente con protección adecuada.

4. Aislamiento térmico: La normativa NSR-10 y la Resolución 0549 de 2015 establecen parámetros para aislamiento térmico en cubiertas y muros exteriores.

5. Ventilación: Implementa ventilación cruzada controlada para evitar pérdidas de calor pero garantizar renovación de aire.

6. Inercia térmica: Utiliza materiales con alta inercia térmica como adobe mejorado o ladrillo macizo (tradicional en la construcción bogotana).

7. Protección contra lluvias: Diseña cubiertas con pendientes adecuadas y aleros suficientes para proteger de la alta pluviosidad.

Un ejemplo exitoso es el edificio de Postgrados de la Universidad Nacional, que utiliza muros Trombe, orientación optimizada y aprovechamiento de luz natural.
```

#### Example 2: Technical Code Query with /no_think

**User Query:**
```
/no_think ¿Cuál es la carga viva mínima para viviendas según NSR-10?
```

**Model Response:**
```
180 kg/m²
```

#### Example 3: Project Management Query

**User Query:**
```
¿Cómo aplicar la metodología BIM a un proyecto de vivienda multifamiliar en Bogotá?
```

**Model Response:**
```
Para aplicar BIM a un proyecto de vivienda multifamiliar en Bogotá, sigue estos pasos:

1. Establecimiento de requisitos BIM: Define el Plan de Ejecución BIM (BEP) según estándares colombianos como la Guía Nacional BIM (2020).

2. Modelado arquitectónico: Desarrolla el modelo arquitectónico en LOD 300 (nivel de detalle) incluyendo cumplimiento con normas POT y normativa de propiedad horizontal (Ley 675).

3. Coordinación estructural: Integra diseño estructural según NSR-10, considerando zona de amenaza sísmica intermedia de Bogotá.

4. Modelado MEP: Incorpora redes según RETIE y normativa de acueducto local (EAAB).

5. Análisis energético: Utiliza herramientas como Green Building Studio para verificar cumplimiento con Resolución 0549/2015 de construcción sostenible.

6. Revisión de interferencias: Realiza sesiones ICE (Integrated Concurrent Engineering) semanales con todo el equipo.

7. Programación 4D: Vincula el cronograma con hitos de licenciamiento constructivo de Curaduría Urbana.

8. Presupuestación 5D: Conecta con bases de datos de precios actualizados para Bogotá (referencia: CONSTRUDATA).

9. Entregables: Genera planos de licencia, documentación para propiedad horizontal y modelos para uso en obra.

Un caso exitil en Bogotá es el proyecto BD Bacatá que utilizó BIM para coordinar su compleja estructura y sistemas.
```

### Troubleshooting

#### Common Issues and Solutions

| Issue | Possible Cause | Solution |
|-------|---------------|----------|
| Empty or cut-off responses | Token limit reached | Increase max_tokens parameter or break query into smaller parts |
| Generic responses | Vague or broad questions | Be more specific, include technical details in your query |
| "Model not found" error | API version mismatch | Verify URL endpoint is correct and up-to-date |
| Rate limit errors | Too many requests | Implement exponential backoff retry logic |
| Authorization errors | Invalid API key | Check environment variables and API key validity |

#### Monitoring and Logging

The application logs all interactions with the Cerebras API in the standard Go log format. Important aspects to monitor:

- Response times
- Token usage
- Error rates
- Query patterns

---

<a name="español"></a>
## Documentación de Integración de API Cerebras

Este documento proporciona información detallada sobre la integración de la API de Cerebras AI en el proyecto ana.world, enfocándose específicamente en el modelo QWen-3B-32B utilizado para asistencia arquitectónica.

### Especificaciones del Modelo: QWen-3B-32B

El modelo QWen-3B-32B es un modelo de lenguaje de vanguardia optimizado para dominios técnicos, incluyendo arquitectura. Especificaciones clave:

- **Tamaño del Modelo**: 32 mil millones de parámetros ajustados a partir de un modelo base de 3B
- **Velocidad de Tokens**: ~2100 tokens/segundo
- **Ventana de Contexto**: 4096 tokens
- **Calidad de Respuesta**: Especializado en conocimiento técnico y específico del dominio
- **Idiomas**: Multilingüe con énfasis en inglés y español

### Características Especiales

#### 1. Modos de Razonamiento Híbrido

QWen-3B-32B admite modos de razonamiento híbrido que pueden activarse en tus prompts:

- **Modo Estándar**: Estilo de respuesta predeterminado
- **Modo Analítico**: Activado al solicitar "análisis paso a paso" o "desglose detallado"
- **Modo Comparativo**: Activado al solicitar "comparar y contrastar" o "pros y contras"

Ejemplo de prompt para modo analítico:
```
Proporciona un análisis paso a paso de técnicas de enfriamiento pasivo para climas tropicales en Colombia.
```

#### 2. Soporte para Comando /no_think

El modelo admite el prefijo de comando `/no_think` para consultas que necesitan respuestas directas y concisas sin texto explicativo.

Ejemplo:
```
/no_think ¿Cuál es el coeficiente sísmico actual del NSR-10 para Bogotá?
```

Esto devolverá solo el valor del coeficiente sin explicación adicional.

#### 3. Experiencia en Dominio Técnico

El modelo ha sido ajustado con conocimiento arquitectónico, particularmente:

- Códigos y regulaciones de construcción colombianos
- Principios de diseño arquitectónico
- Marcos de gestión de proyectos
- Técnicas de construcción sostenible
- Estimación de costos y presupuesto

### Uso de la API

#### Formato de Solicitud

```json
{
  "model": "cerebras/QWen-3B-32B",
  "messages": [
    {
      "role": "system",
      "content": "Instrucciones del sistema..."
    },
    {
      "role": "user",
      "content": "Consulta del usuario..."
    }
  ],
  "temperature": 0.7,
  "max_tokens": 500
}
```

#### Formato de Respuesta

```json
{
  "id": "resp_01234567",
  "object": "chat.completion",
  "created": 1715987654,
  "model": "cerebras/QWen-3B-32B",
  "choices": [
    {
      "message": {
        "role": "assistant",
        "content": "Texto de respuesta..."
      },
      "finish_reason": "stop"
    }
  ],
  "usage": {
    "prompt_tokens": 89,
    "completion_tokens": 156,
    "total_tokens": 245
  }
}
```

### Ejemplos Específicos de Arquitectura

#### Ejemplo 1: Consulta Arquitectónica General

**Consulta del Usuario:**
```
¿Cuáles son las consideraciones principales para el diseño bioclimático en Bogotá?
```

**Respuesta del Modelo:**
```
Para el diseño bioclimático en Bogotá, debes considerar:

1. Clima: Bogotá tiene clima frío de montaña (promedio 14°C) con poca variación estacional pero grandes variaciones diurnas.

2. Orientación: La orientación este-oeste maximiza la captación solar en fachadas norte-sur debido a la posición ecuatorial de Colombia.

3. Radiación solar: Aprovecha la radiación directa para calentamiento pasivo mediante ventanales al oriente/occidente con protección adecuada.

4. Aislamiento térmico: La normativa NSR-10 y la Resolución 0549 de 2015 establecen parámetros para aislamiento térmico en cubiertas y muros exteriores.

5. Ventilación: Implementa ventilación cruzada controlada para evitar pérdidas de calor pero garantizar renovación de aire.

6. Inercia térmica: Utiliza materiales con alta inercia térmica como adobe mejorado o ladrillo macizo (tradicional en la construcción bogotana).

7. Protección contra lluvias: Diseña cubiertas con pendientes adecuadas y aleros suficientes para proteger de la alta pluviosidad.

Un ejemplo exitoso es el edificio de Postgrados de la Universidad Nacional, que utiliza muros Trombe, orientación optimizada y aprovechamiento de luz natural.
```

#### Ejemplo 2: Consulta Técnica de Código con /no_think

**Consulta del Usuario:**
```
/no_think ¿Cuál es la carga viva mínima para viviendas según NSR-10?
```

**Respuesta del Modelo:**
```
180 kg/m²
```

#### Ejemplo 3: Consulta de Gestión de Proyectos

**Consulta del Usuario:**
```
¿Cómo aplicar la metodología BIM a un proyecto de vivienda multifamiliar en Bogotá?
```

**Respuesta del Modelo:**
```
Para aplicar BIM a un proyecto de vivienda multifamiliar en Bogotá, sigue estos pasos:

1. Establecimiento de requisitos BIM: Define el Plan de Ejecución BIM (BEP) según estándares colombianos como la Guía Nacional BIM (2020).

2. Modelado arquitectónico: Desarrolla el modelo arquitectónico en LOD 300 (nivel de detalle) incluyendo cumplimiento con normas POT y normativa de propiedad horizontal (Ley 675).

3. Coordinación estructural: Integra diseño estructural según NSR-10, considerando zona de amenaza sísmica intermedia de Bogotá.

4. Modelado MEP: Incorpora redes según RETIE y normativa de acueducto local (EAAB).

5. Análisis energético: Utiliza herramientas como Green Building Studio para verificar cumplimiento con Resolución 0549/2015 de construcción sostenible.

6. Revisión de interferencias: Realiza sesiones ICE (Integrated Concurrent Engineering) semanales con todo el equipo.

7. Programación 4D: Vincula el cronograma con hitos de licenciamiento constructivo de Curaduría Urbana.

8. Presupuestación 5D: Conecta con bases de datos de precios actualizados para Bogotá (referencia: CONSTRUDATA).

9. Entregables: Genera planos de licencia, documentación para propiedad horizontal y modelos para uso en obra.

Un caso exitil en Bogotá es el proyecto BD Bacatá que utilizó BIM para coordinar su compleja estructura y sistemas.
```

### Solución de Problemas

#### Problemas Comunes y Soluciones

|| Problema | Posible Causa | Solución |
||----------|--------------|----------|
|| Respuestas vacías o cortadas | Límite de tokens alcanzado | Aumentar el parámetro max_tokens o dividir la consulta en partes más pequeñas |
|| Respuestas genéricas | Preguntas vagas o demasiado amplias | Ser más específico, incluir detalles técnicos en tu consulta |
|| Error "Model not found" | Incompatibilidad de versión de API | Verificar que la URL del endpoint sea correcta y esté actualizada |
|| Errores de límite de tasa | Demasiadas solicitudes | Implementar lógica de reintento con retroceso exponencial |
|| Errores de autorización | Clave API inválida | Verificar variables de entorno y validez de la clave API |

#### Monitoreo y Registro

La aplicación registra todas las interacciones con la API de Cerebras en el formato de registro estándar de Go. Aspectos importantes a monitorear:

- Tiempos de respuesta
- Uso de tokens
- Tasas de error
- Patrones de consulta

## Recursos Adicionales

- [Documentación oficial de Cerebras API](https://inference-docs.cerebras.ai/introduction)
- [Guía para modelos QWen](https://inference-docs.cerebras.ai/models/qwen)
- [Ejemplos de código en Go](https://inference-docs.cerebras.ai/api-reference/chat-completions)


## Development Session
Reference: https://app.warp.dev/session/b660fd8a-f765-449c-a70c-f8c7b971e3c4?pwd=e9ccd7cb-d8be-494e-a2f2-35469f726896
Last Updated: Sat May 17 07:34:44 AM CEST 2025
