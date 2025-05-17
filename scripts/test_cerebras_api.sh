#!/bin/bash
# Test script for the Cerebras AI API integration
# This script uses curl to test the AI assistant API endpoint

# Load environment variables from .env file
if [ -f ../.env ]; then
  export $(grep -v '^#' ../.env | xargs)
fi

CEREBRAS_API_KEY=${CEREBRAS_API_KEY:-"your_api_key_here"}
CEREBRAS_API_URL=${CEREBRAS_API_URL:-"https://inference.cerebras.ai/v1/chat/completions"}

# Set text colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Ana Cerebras AI API Test${NC}"
echo -e "${YELLOW}======================${NC}\n"

# Test direct Cerebras API access
echo -e "Testing direct Cerebras API connection..."

if [ -z "$CEREBRAS_API_KEY" ] || [ "$CEREBRAS_API_KEY" = "your_api_key_here" ]; then
  echo -e "${RED}Error: CEREBRAS_API_KEY not set correctly. Please update the .env file.${NC}"
  echo "Skipping direct API test."
else
  # Create a test payload
  payload=$(cat <<EOF
{
  "model": "cerebras/Cerebras-GPT-4o",
  "messages": [
    {
      "role": "system",
      "content": "You are an AI assistant for architects."
    },
    {
      "role": "user", 
      "content": "What are the key considerations for sustainable building design?"
    }
  ],
  "temperature": 0.7,
  "max_tokens": 500
}
EOF
)

  # Make the API request
  response=$(curl -s -w "\n%{http_code}" \
    -X POST "$CEREBRAS_API_URL" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $CEREBRAS_API_KEY" \
    -d "$payload")
  
  # Extract status code
  http_code=$(echo "$response" | tail -n1)
  response_body=$(echo "$response" | sed '$ d')
  
  if [ "$http_code" -eq 200 ]; then
    echo -e "${GREEN}Successfully connected to Cerebras API!${NC}"
    echo -e "Received response (truncated):"
    echo "${response_body}" | grep -o '"content":"[^"]*"' | head -n1
  else
    echo -e "${RED}Error connecting to Cerebras API.${NC}"
    echo -e "Status code: ${http_code}"
    echo -e "Response: ${response_body}"
  fi
fi

echo -e "\n${YELLOW}Testing Local API Server${NC}"
echo -e "${YELLOW}======================${NC}\n"

# Start the local server if not already running
echo "Checking if local server is running..."
if ! curl -s localhost:8080/api/health > /dev/null; then
  echo "Local server not responding. Please start the server in another terminal with:"
  echo "cd $(dirname "$0")/.."
  echo "go run cmd/server/main.go"
  exit 1
fi

# Test the AI assistant endpoint
echo "Testing AI assistant endpoint..."
response=$(curl -s -w "\n%{http_code}" \
  -X POST "http://localhost:8080/api/ai/assistant" \
  -H "Content-Type: application/json" \
  -d '{"query": "¿Cómo puedo organizar mejor mis proyectos de arquitectura?"}')

# Extract status code
http_code=$(echo "$response" | tail -n1)
response_body=$(echo "$response" | sed '$ d')

if [ "$http_code" -eq 200 ]; then
  echo -e "${GREEN}AI assistant endpoint test successful!${NC}"
  echo -e "Response: ${response_body}"
else
  echo -e "${RED}AI assistant endpoint test failed.${NC}"
  echo -e "Status code: ${http_code}"
  echo -e "Response: ${response_body}"
fi

echo -e "\n${GREEN}Testing completed!${NC}"

