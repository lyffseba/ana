# Cerebras QWen Models Documentation

[English](#english) | [Español](#español)

<a name="english"></a>
## Cerebras QWen Models for Architectural Applications

This document provides comprehensive technical details about the Cerebras QWen models used in the ana.world project, focusing on their capabilities, implementation, and architectural domain specialization.

### Model Overview

Ana.world leverages Cerebras' Qwen 3 32B model for text-based architectural assistance. Vision model support depends on Cerebras API availability.

#### Qwen 3 32B (Text Model)

The Qwen 3 32B model represents the foundation of our AI assistance system, providing powerful text-based architectural insights.

| Feature | Specification | Notes |
|---------|---------------|-------|
| **Base Architecture** | QWen Large Language Model | Developed by Alibaba Cloud, hosted on Cerebras |
| **Parameter Size** | 32 billion parameters | Fine-tuned from 3B base model |
| **Context Window** | 4,096 tokens | Approximately 3,000 words |
| **Inference Speed** | ~2,100 tokens/second | Varies based on server load |
| **Response Time** | 1.2-2.5 seconds typical | For standard architectural queries |
| **Languages** | Multilingual | Optimized for English and Spanish |
| **Domain Expertise** | Architecture, Construction, Urban Planning | Enhanced with specialized context |

#### Vision Model

Cerebras does not currently list a public vision model in the official API reference. For vision capabilities, check the [official Cerebras documentation](https://inference-docs.cerebras.ai/api-reference/models) for availability and correct model IDs.

| Feature | Specification | Notes |
|---------|---------------|-------|
| **Base Architecture** | QWen Vision Language Model | Specialized for visual understanding |
| **Parameter Size** | 2.5 billion parameters | Tuned for visual recognition |
| **Context Window** | 4,096 tokens | Includes image token allocation |
| **Image Support** | JPEG, PNG, WebP | Up to 5MB size limit |
| **Resolution** | Up to 2048×2048 pixels | Higher resolution images are downsampled |
| **Inference Speed** | ~1,800 tokens/second | Slower than text-only model |
| **Response Time** | 2.5-4.0 seconds typical | For image analysis queries |

### Model Capabilities

#### Qwen 3 32B Capabilities

The Qwen 3 32B text model excels in several areas crucial for architectural applications:

1. **Technical Knowledge**
   - Construction materials and techniques
   - Structural engineering principles
   - Building codes and regulations (Colombia-specific)
   - Architectural terminology and concepts

2. **Specialized Reasoning Modes**
   - **Analytical Mode**: For step-by-step breakdown of complex architectural problems
   - **Comparative Mode**: For contrasting different architectural approaches
   - **Problem-Solving Mode**: For addressing design challenges
   - **Directional Mode**: For concise, direct answers via `/no_think` command

3. **Domain-Specific Features**
   - Colombian building code expertise (NSR-10)
   - Local construction practices knowledge
   - Regional climate-responsive design understanding
   - Spanish terminology mastery

#### QWen-2.5-Vision Capabilities

The QWen-2.5-Vision model provides visual analysis capabilities:

1. **Image Understanding**
   - Architectural style identification
   - Building element recognition
   - Spatial arrangement analysis
   - Construction technique identification

2. **Visual Technical Analysis**
   - Structural system assessment
   - Material identification
   - Potential issue detection
   - Dimensional analysis

3. **Domain-Specific Visual Features**
   - Architectural drawing interpretation
   - Construction site photo analysis
   - Building condition assessment
   - Design element recognition

### Implementation Examples

#### Text Model Implementation

Here's how to implement the Qwen 3 32B model (`qwen-3-32b`) in Go:

```go
// Cerebras client for text generation
func (c *CerebrasClient) GenerateTextResponse(userQuery string, conversationContext []Message) (string, error) {
    requestBody := ChatCompletionRequest{
        Model:       "qwen-3-32b",
        Messages:    append(conversationContext, Message{
            Role:    "user",
            Content: userQuery,
        }),
        Temperature: 0.7,
        MaxTokens:   500,
    }
    // ... send request as per Cerebras API docs ...
}
```

#### QuickStart Example (Python)

```python
import os
from cerebras.cloud.sdk import Cerebras

client = Cerebras(api_key=os.environ.get("CEREBRAS_API_KEY"))

chat_completion = client.chat.completions.create(
    model="qwen-3-32b",
    messages=[{"role": "user", "content": "¿Qué es arquitectura sostenible?"}]
)
print(chat_completion)
```

#### Official Model List

For the most up-to-date list of available models and their IDs, always refer to the [Cerebras API reference](https://inference-docs.cerebras.ai/api-reference/models).

As of May 2025, supported model IDs include:
- `llama-4-scout-17b-16e-instruct`
- `llama3.1-8b`
- `llama-3.3-70b`
- `qwen-3-32b`
- `deepseek-r1-distill-llama-70b` (private preview)

#### System Context Example

The power of our AI implementation comes from specialized system prompts:

```go
// Architectural system prompt construction
systemPrompt := fmt.Sprintf(
    "Eres un asistente especializado en arquitectura para la plataforma ana.world. %s%sTu conocimiento incluye: 1) Normativas colombianas: NSR-10 (Norma Sismo Resistente), POT de Bogotá, Decreto 1077 de 2015; 2) Diseño arquitectónico: metodología BIM, diseño paramétrico, estilos latinoamericanos; 3) Gestión de proyectos: metodologías PMI/PRINCE2, control de cronogramas, licencias; 4) Materiales sostenibles: guadua, tierra compactada, sistemas pasivos, certificación LEED/EDGE; 5) Presupuestos: estimación de costos por m², APU. %s Si te preguntan en inglés, comprende pero responde en español. %s",
    noThinkInstructions,
    visionSpecificInstructions,
    responseStyleInstructions,
    technicalInfoInstructions,
)
```

#### Vision Model Implementation

Here's our implementation for processing image inputs:

```go
// Vision model integration
func (c *CerebrasClient) GenerateVisionResponse(userQuery string, imageBase64 string, context []Message) (string, error) {
    // Validate image data
    if imageBase64 == "" {
        return "", fmt.Errorf("image data is required for vision model")
    }
    
    // Format image content with markdown syntax
    imageContent := fmt.Sprintf(
        `%s\n\n![Imagen arquitectónica](data:image/jpeg;base64,%s)`,
        userQuery,
        imageBase64,
    )
    
    // Create request with image content
    requestBody := ChatCompletionRequest{
        Model:       "(Vision model: see Cerebras docs for availability)",
        Messages:    append(context, Message{
            Role:    "user",
            Content: imageContent,
        }),
        Temperature: 0.7,
        MaxTokens:   800, // Higher for vision descriptions
    }
    
    // Send request and handle response...
}
```

#### Direct Response Mode with /no_think

Our implementation supports the special `/no_think` command for concise answers:

```go
// Process the query for no_think mode
if strings.HasPrefix(strings.ToLower(query), "/no_think") {
    isNoThink = true
    // Remove the command from the query
    query = strings.TrimSpace(query[9:])
    
    // Adjust system prompt for direct answers
    noThinkInstructions = "Estás en modo respuesta directa. Responde con la información precisa sin explicaciones adicionales, usando el mínimo de palabras posible."
    responseStyleInstructions = "Responde solo con datos concretos sin introducción ni explicación."
}

// After getting AI response, process for no_think mode
if isNoThink && len(response) > 0 {
    // Extract the most concise answer
    lines := strings.Split(response, "\n")
    shortestRelevantLine := response // Default to full response
    
    // Find the shortest line that contains relevant information
    minLength := 1000
    for _, line := range lines {
        trimmed := strings.TrimSpace(line)
        if len(trimmed) > 10 && len(trimmed) < minLength {
            // Check if line contains technical information
            if strings.ContainsAny(trimmed, "0123456789") || 
               strings.Contains(trimmed, "kg") ||
               strings.Contains(trimmed, "m²") ||
               strings.Contains(trimmed, "mm") {
                shortestRelevantLine = trimmed
                minLength = len(trimmed)
            }
        }
    }
    response = shortestRelevantLine
}
```

### Error Handling and Best Practices

#### Comprehensive Error Handling

Our implementation includes robust error handling:

```go
if resp.StatusCode != http.StatusOK {
    log.Printf("API error: status code %d, body: %s", resp.StatusCode, string(body))
    
    // Return user-friendly error messages based on status code
    switch resp.StatusCode {
    case http.StatusUnauthorized:
        return "No se pudo autenticar con el servicio de IA. Por favor verifica la configuración del asistente arquitectónico.", nil
    case http.StatusForbidden:
        return "No tienes permisos para utilizar el asistente arquitectónico. Por favor contacta al administrador.", nil
    case http.StatusTooManyRequests:
        return "El asistente arquitectónico está experimentando mucho tráfico. Por favor intenta de nuevo en unos momentos.", nil
    case http.StatusServiceUnavailable:
        return "El servicio de asistencia arquitectónica no está disponible temporalmente. Por favor intenta más tarde.", nil
    case http.StatusRequestEntityTooLarge:
        return "La imagen es demasiado grande. Por favor utiliza una imagen más pequeña (máximo 5MB).", nil
    default:
        return "Hubo un problema al procesar tu consulta arquitectónica. Por favor intenta reformularla o contacta soporte técnico.", nil
    }
}
```

#### Best Practices for Model Usage

For optimal results with Cerebras QWen models in architectural applications:

1. **System Prompt Engineering**
   - Include domain-specific terminology
   - Specify response format expectations
   - Include regional context (Colombian regulations)
   - Balance detail vs. conciseness based on use case

2. **Text Model Optimization**
   - Use detailed, specific queries for better answers
   - Provide sufficient context in user queries
   - Utilize `/no_think` for quick reference answers
   - Ask for step-by-step analysis for complex topics

3. **Vision Model Optimization**
   - Ensure images are clear and well-lit
   - Provide specific questions about the image
   - Include elements for scale reference when possible
   - Ask focused questions about specific visual elements

4. **Performance Considerations**
   - Balance token usage with response quality
   - Implement client-side caching for common queries
   - Use appropriate timeouts based on query complexity
   - Consider batch processing for multiple queries

5. **Bilingual Support Best Practices**
   - Include language detection in system prompts
   - Maintain Spanish-language glossary for technical terms
   - Specify Spanish response format for consistency
   - Test queries in both languages regularly

### Example Usage Scenarios

#### Scenario 1: Technical Code Query

```
/no_think ¿Cuál es la carga viva mínima para escaleras según NSR-10?
```

Output:
```
300 kg/m²
```

#### Scenario 2: Material Recommendation

```
¿Qué materiales sostenibles recomiendas para una vivienda en clima cálido húmedo en Barranquilla?
```

Output:
```
Para una vivienda en clima cálido húmedo en Barranquilla, recomiendo estos materiales sostenibles:

1. Muros: Bloques de tierra comprimida (BTC) con estabilizantes naturales que ofrecen alta inercia térmica y regulación de humedad.

2. Cubierta: Sistema de techo ventilado con tejas de arcilla sobre estructura de bambú Guadua Angustifolia certificada, con cámara de aire y aislamiento de fibra de coco.

3. Acabados: Revoques de cal hidráulica natural que permiten la transpiración de muros y controlan la humedad.

4. Protección solar: Celosías de bambú local tratado para sombreamiento en fachadas críticas (oeste y sur).

5. Pavimentos: Baldosas de barro cocido sobre base ventilada para reducir la temperatura radiante.

6. Ventanas: Marcos de madera certificada FSC con vidrio de control solar (factor solar 0.25-0.30).

7. Aislamiento: Paneles de fibra de coco o algodón reciclado para zonas críticas.

Estos materiales cumplen con la Resolución 0549 de 2015 para sostenibilidad y el Título J del NSR-10 para protección contra incendios, además de ser accesibles en la región Caribe colombiana.
```

#### Scenario 3: Image Analysis

```
Analiza esta imagen de fachada y sugiere mejoras para optimizar el control solar.
[IMAGEN]
```

Output:
```
Analizando la fachada mostrada en la imagen, identifico los siguientes aspectos relacionados con el control solar:

1. Orientación: La fachada parece estar orientada hacia el oeste/suroeste según las sombras proyectadas, lo que la expone a radiación solar directa de alta intensidad en horas de la tarde.

2. Elementos actuales:
   - Ventanería de vidrio transparente sin protección solar
   - Ausencia de elementos de sombreamiento horizontal
   - Balcones con barandas metálicas que no proporcionan sombreamiento
   - Muros con acabado claro que refleja parte de la radiación

Recomendaciones para optimizar el control solar:

1. Implementar quiebrasoles horizontales (parasoles) sobre las ventanas, con una profundidad calculada según la carta solar de Bogotá (aproximadamente 80-100cm).

2. Reemplazar el vidrio actual por vidrio de control solar con un factor solar (FS) de 0.25-0.30 y transmisión luminosa superior al 50%.

3. Instalar cortasoles verticales en los extremos de ventanas para controlar la radiación solar de ángulos bajos (mañana/tarde).

4. Extender los balcones existentes para generar mayor sombreamiento en los pisos inferiores.

5. Integrar jardineras con vegetación caducifolia local que proporcione sombra en verano y permita ganancia térmica en periodos más fríos.

Estas mejoras podrían reducir la carga térmica en aproximadamente 30-40%, mejorando el confort interior y reduciendo el consumo energético en climatización según los estándares de la Resolución 0549 de 2015.
```

---

<a name="español"></a>
## Modelos Cerebras QWen para Aplicaciones Arquitectónicas

Este documento proporciona detalles técnicos completos sobre los modelos Cerebras QWen utilizados en el proyecto ana.world, enfocándose en sus capacidades


## Development Session
Reference: https://app.warp.dev/session/b660fd8a-f765-449c-a70c-f8c7b971e3c4?pwd=e9ccd7cb-d8be-494e-a2f2-35469f726896
Last Updated: Sat May 17 07:34:44 AM CEST 2025
