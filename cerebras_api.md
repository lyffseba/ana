# Cerebras API Reference

This document provides a reference for using the Cerebras API in the ANA project.

## Client Initialization

```go
// Create a new Cerebras client
client := ai.NewCerebrasClient()
```

The client reads configuration from environment variables:
- `CEREBRAS_API_KEY`: Required for authentication
- `CEREBRAS_API_URL`: Optional override for API endpoint
- `CEREBRAS_CACHE_TTL`: Optional cache time-to-live setting
- `ENABLE_CEREBRAS_METRICS`: Optional flag to enable metrics collection

## Generating Text Responses

```go
// Define system context
systemContext := []ai.Message{
    {
        Role:    "system",
        Content: "You are a helpful assistant specialized in architecture.",
    },
}

// Generate a response
response, err := client.GenerateTextResponse(
    "What are the key considerations for sustainable architecture?",
    "qwen-3-32b",
    systemContext,
)
if err != nil {
    log.Printf("Error generating response: %v", err)
    return
}

fmt.Println(response)
```

## Caching

The client automatically caches responses to improve performance and reduce API calls:

```go
// Check if a response is cached
cachedResponse, isCached := client.GetCachedResponse(modelName, systemContext)
if isCached {
    fmt.Println("Using cached response:", cachedResponse)
    return
}

// Generate a new response
response, err := client.GenerateTextResponse(query, modelName, systemContext)
if err != nil {
    log.Printf("Error: %v", err)
    return
}

// Response is automatically cached
```

## Circuit Breaker

The client implements a circuit breaker pattern to prevent cascading failures:

```go
// Check circuit breaker before making a request
if err := client.CheckCircuitBreaker(); err != nil {
    log.Printf("Circuit breaker is open: %v", err)
    return
}

// Make the request
response, err := client.GenerateTextResponse(query, modelName, systemContext)
if err != nil {
    // Record failure to potentially open the circuit breaker
    client.RecordFailure()
    log.Printf("Error: %v", err)
    return
}

// Record success to reset failure count
client.RecordSuccess()
```

## Health and Metrics

```go
// Check API status
apiStatus := client.GetAPIStatus()
if apiStatus != "ok" {
    log.Printf("API status: %s", apiStatus)
}

// Get cache size
cacheSize := client.GetCacheSize()
log.Printf("Cache size: %d", cacheSize)

// Get circuit breaker state
circuitState := client.GetCircuitState()
log.Printf("Circuit state: %s", circuitState)
```

## Available Models

- `llama-4-scout-17b-16e-instruct`: 17B parameter instruction-tuned model
- `llama3.1-8b`: 8B parameter general-purpose model
- `llama-3.3-70b`: 70B parameter large general-purpose model
- `qwen-3-32b`: 32B parameter model with strong multilingual capabilities
- `deepseek-r1-distill-llama-70b`: 70B parameter specialized model (private preview)

For more details on models, see [cerebras_models.md](cerebras_models.md).

## Error Handling

The client provides user-friendly error messages based on HTTP status codes:

- 401 Unauthorized: Authentication issue
- 403 Forbidden: Permission issue
- 429 Too Many Requests: Rate limiting
- 503 Service Unavailable: Temporary service outage

Example:

```go
response, err := client.GenerateTextResponse(query, modelName, systemContext)
if err != nil {
    log.Printf("Error: %v", err)
    // Handle error
    return
}
```

## Best Practices

1. **Use appropriate models** for different tasks
2. **Implement caching** to reduce API calls and improve performance
3. **Handle errors gracefully** to provide a good user experience
4. **Monitor usage** to optimize costs and performance
5. **Use circuit breaker** to prevent cascading failures
