#!/bin/bash
# ANA Development Plan Script
# Reference: https://app.warp.dev/session/b660fd8a-f765-449c-a70c-f8c7b971e3c4?pwd=e9ccd7cb-d8be-494e-a2f2-35469f726896

set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}Executing ANA development plan...${NC}"

# Setup development environment
echo -e "${BLUE}Setting up development environment...${NC}"
go mod tidy
go mod verify

# Run tests
echo -e "${BLUE}Running tests...${NC}"
go test ./... -v

# Build project
echo -e "${BLUE}Building project...${NC}"
go build -o bin/ana ./cmd/ana

# Run linting
echo -e "${BLUE}Running linters...${NC}"
golangci-lint run

# Generate documentation
echo -e "${BLUE}Generating documentation...${NC}"
go doc -all > docs/API.md

echo -e "${GREEN}Development plan executed successfully!${NC}"
