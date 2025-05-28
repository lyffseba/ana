# Cerebras AI Models

This document provides information about the Cerebras AI models available for use in the ANA project.

## Available Models

### llama-4-scout-17b-16e-instruct
- **Size**: 17 billion parameters
- **Type**: Instruction-tuned model
- **Features**: General-purpose model with strong reasoning capabilities
- **Best for**: Complex reasoning tasks, multi-step problem solving, and detailed explanations

### llama3.1-8b
- **Size**: 8 billion parameters
- **Type**: General-purpose model
- **Features**: Efficient model with good performance for common tasks
- **Best for**: Everyday queries, content generation, and applications with resource constraints

### llama-3.3-70b
- **Size**: 70 billion parameters
- **Type**: Large general-purpose model
- **Features**: Advanced reasoning and knowledge capabilities
- **Best for**: Complex tasks requiring deep understanding and nuanced responses

### qwen-3-32b
- **Size**: 32 billion parameters
- **Type**: General-purpose model
- **Features**: Strong multilingual capabilities and technical knowledge
- **Best for**: Technical content, code generation, and multilingual applications

### deepseek-r1-distill-llama-70b (private preview)
- **Size**: 70 billion parameters
- **Type**: Specialized model
- **Features**: Advanced capabilities for specific domains
- **Best for**: Specialized applications requiring domain-specific knowledge
- **Note**: Available in private preview. Contact Cerebras for access.

## Usage Guidelines

### Model Selection
- For most general-purpose tasks, `qwen-3-32b` provides a good balance of performance and efficiency
- For resource-constrained environments, use `llama3.1-8b`
- For complex reasoning tasks, use `llama-4-scout-17b-16e-instruct` or `llama-3.3-70b`

### Parameter Optimization
- **Temperature**: Controls randomness (0.0-1.5)
  - Lower values (0.1-0.3): More deterministic, focused responses
  - Medium values (0.4-0.7): Balanced creativity and coherence
  - Higher values (0.8-1.5): More creative, varied responses
- **Max Tokens**: Limits response length
  - Set appropriately based on expected response length
  - Default: 500 tokens (approximately 375 words)

## API Reference
For more details, see the [Cerebras API documentation](https://inference-docs.cerebras.ai/api-reference/chat-completions).
