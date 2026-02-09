#!/bin/bash

# Integration test script for githubCompare
# This script tests the application with a real public repository

set -e

echo "=== githubCompare Integration Test ==="
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Check if binary exists
BINARY="./dist/githubCompare"
if [ ! -f "$BINARY" ]; then
    echo -e "${RED}Error: Binary not found. Please run 'make build' first.${NC}"
    exit 1
fi

echo "✓ Binary found: $BINARY"
echo ""

# Test 1: Help command
echo "Test 1: Testing --help command..."
if $BINARY --help > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Help command works${NC}"
else
    echo -e "${RED}✗ Help command failed${NC}"
    exit 1
fi
echo ""

# Test 2: Missing required flag
echo "Test 2: Testing missing --repo flag..."
if $BINARY 2>&1 | grep -q "required"; then
    echo -e "${GREEN}✓ Correctly requires --repo flag${NC}"
else
    echo -e "${RED}✗ Should require --repo flag${NC}"
    exit 1
fi
echo ""

# Test 3: Invalid repository URL
echo "Test 3: Testing invalid repository URL..."
if $BINARY --repo "not-a-valid-url" 2>&1 | grep -q "Error"; then
    echo -e "${GREEN}✓ Correctly handles invalid URL${NC}"
else
    echo -e "${RED}✗ Should handle invalid URL${NC}"
    exit 1
fi
echo ""

echo -e "${GREEN}=== All basic tests passed! ===${NC}"
echo ""
echo "Note: To test with a real repository, run:"
echo "  $BINARY --repo https://github.com/owner/repo"
echo ""
