#!/bin/bash

# Test the Cerebras integration
echo "Testing Cerebras integration..."
curl -X POST http://localhost:8081/api/cerebras/assistant \
  -F "query=What are the key considerations for sustainable architecture?" \
  -F "model_type=qwen-3-32b"

echo -e "\n\nDone!"
