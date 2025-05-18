#!/bin/bash
# Update and push all changes
# Reference: https://app.warp.dev/session/b660fd8a-f765-449c-a70c-f8c7b971e3c4?pwd=e9ccd7cb-d8be-494e-a2f2-35469f726896

set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

# Warp session reference
WARP_SESSION="https://app.warp.dev/session/b660fd8a-f765-449c-a70c-f8c7b971e3c4?pwd=e9ccd7cb-d8be-494e-a2f2-35469f726896"

echo -e "${BLUE}Updating documentation with session reference...${NC}"

# Update all .md files with session reference
for file in $(find . -name "*.md"); do
    if ! grep -q "$WARP_SESSION" "$file"; then
        echo -e "\n## Development Session\nReference: $WARP_SESSION\nLast Updated: $(date)" >> "$file"
        echo -e "${GREEN}Updated${NC} $file"
    fi
done

# Update all .go files with session reference
for file in $(find . -name "*.go"); do
    if ! grep -q "$WARP_SESSION" "$file"; then
        sed -i "1i// Reference: $WARP_SESSION\n// Last Updated: $(date)\n" "$file"
        echo -e "${GREEN}Updated${NC} $file"
    fi
done

# Update main README if it exists
if [ -f "README.md" ]; then
    if ! grep -q "$WARP_SESSION" "README.md"; then
        sed -i "1i# ANA Project\n\n## Development Session\nReference: $WARP_SESSION\nLast Updated: $(date)\n" "README.md"
        echo -e "${GREEN}Updated${NC} README.md"
    fi
fi

echo -e "${BLUE}Checking git status...${NC}"
git status

echo -e "${BLUE}Adding all changes...${NC}"
git add .

echo -e "${BLUE}Committing changes...${NC}"
git commit -m "docs: update all files with session reference

Session: $WARP_SESSION"

echo -e "${BLUE}Pushing changes...${NC}"
git push

echo -e "${GREEN}All changes have been updated and pushed!${NC}"
