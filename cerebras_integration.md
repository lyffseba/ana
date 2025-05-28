# Cerebras API Integration

## Overview
This document outlines the integration of the Cerebras API into the ANA project. The integration allows the application to use Cerebras' AI models for generating responses to user queries.

## Key Components

### 1. CerebrasClient
The `CerebrasClient` in `internal/ai/cerebras_client.go` handles communication with the Cerebras API. It includes:
- API authentication
- Request/response handling
- Response caching
- Circuit breaker pattern for fault tolerance
- Concurrency limiting
- Metrics collection

### 2. API Endpoint
The API endpoint was updated from `https://inference.cerebras.ai/v1/chat/completions` to `https://api.cerebras.ai/v1/chat/completions` to match the current Cerebras API.

### 3. Response Processing
The `removeThinkingTags` function was added to strip out thinking tags (`<think>...</think>`) from the model's responses, ensuring that only the final response is shown to the user.

### 4. Environment Configuration
The integration uses the following environment variables:
- `CEREBRAS_API_KEY`: API key for authentication
- `CEREBRAS_API_URL`: Optional override for the API endpoint
- `CEREBRAS_CACHE_TTL`: Optional cache time-to-live setting
- `ENABLE_CEREBRAS_METRICS`: Optional flag to enable metrics collection

## Available Models
The integration supports the following Cerebras models:
- `llama-4-scout-17b-16e-instruct`
- `llama3.1-8b`
- `llama-3.3-70b`
- `qwen-3-32b`
- `deepseek-r1-distill-llama-70b` (private preview)

## Testing
The integration can be tested using the `test_cerebras.sh` script, which sends a sample query to the API and displays the response.

## Future Improvements
- Add support for streaming responses
- Implement more sophisticated caching strategies
- Add support for more advanced model parameters
- Improve error handling and retry logic
