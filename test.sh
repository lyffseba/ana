#!/bin/bash
# Test script for Go repository
# Reference: https://app.warp.dev/session/b660fd8a-f765-449c-a70c-f8c7b971e3c4?pwd=e9ccd7cb-d8be-494e-a2f2-35469f726896

set -e

echo "Running tests with race detection and coverage..."
go test -race -coverprofile=coverage.txt -covermode=atomic ./...

# Check coverage threshold
COVERAGE=$(go tool cover -func=coverage.txt | grep total | awk '{print $3}' | sed 's/%//')
echo "Total coverage: $COVERAGE%"

# Generate HTML coverage report
go tool cover -html=coverage.txt -o coverage.html

# Fail if coverage is below 80%
if (( $(echo "$COVERAGE < 80" | bc -l) )); then
    echo "Error: Code coverage is below 80%"
    exit 1
fi

echo "All tests passed successfully!"
