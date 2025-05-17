#!/bin/bash
set -e

# Colors for terminal output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[0;33m'
RED='\033[0;31m'
NC='\033[0m'

# Environment
ENVIRONMENT=${1:-production}
DEPLOY_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
DEPLOY_DIR="$(pwd)"
REPO_NAME=$(basename "$DEPLOY_DIR")

echo -e "${BLUE}Starting deployment of $REPO_NAME to $ENVIRONMENT${NC}"

# Detect project type and perform appropriate deployment steps
if [ -f "requirements.txt" ]; then
    echo -e "${YELLOW}Python project detected${NC}"
    
    # Activate virtual environment or create if it doesn't exist
    if [ ! -d "venv" ]; then
        python3 -m venv venv
    fi
    source venv/bin/activate
    
    # Install dependencies
    pip install -r requirements.txt
    
    # Run tests if they exist
    if [ -d "tests" ]; then
        echo -e "${BLUE}Running tests...${NC}"
        pytest
    fi
    
    # Create deployment marker
    echo "$ENVIRONMENT" > .deployed
    echo "$DEPLOY_TIME" >> .deployed
    
elif [ -f "package.json" ]; then
    echo -e "${YELLOW}Node.js project detected${NC}"
    
    # Install dependencies
    npm install
    
    # Run build if it exists
    if grep -q "\"build\"" package.json; then
        npm run build
    fi
    
    # Run tests if they exist
    if grep -q "\"test\"" package.json; then
        echo -e "${BLUE}Running tests...${NC}"
        npm test
    fi
    
    # Create deployment marker
    echo "$ENVIRONMENT" > .deployed
    echo "$DEPLOY_TIME" >> .deployed
    
elif [ -f "go.mod" ]; then
    echo -e "${YELLOW}Go project detected${NC}"
    
    # Get dependencies
    go mod tidy
    go mod verify
    
    # Run tests
    echo -e "${BLUE}Running tests...${NC}"
    go test ./...
    
    # Build the project
    go build -o bin/app
    
    # Create deployment marker
    echo "$ENVIRONMENT" > .deployed
    echo "$DEPLOY_TIME" >> .deployed
else
    echo -e "${YELLOW}Generic project detected${NC}"
    
    # Create deployment marker
    echo "$ENVIRONMENT" > .deployed
    echo "$DEPLOY_TIME" >> .deployed
fi

# Update deployment configuration
cat > .deployment.yaml << YAML
environment: $ENVIRONMENT
last_deployed: $DEPLOY_TIME
status: deployed
deployer: lyff_admin
repository: $REPO_NAME
YAML

# Notify Discord if webhook is configured
if [ -n "$DISCORD_WEBHOOK_URL" ]; then
    curl -H "Content-Type: application/json" -X POST -d "{\"content\":\"🚀 Deployed $REPO_NAME to $ENVIRONMENT at $DEPLOY_TIME\"}" $DISCORD_WEBHOOK_URL
fi

echo -e "${GREEN}Deployment complete!${NC}"
