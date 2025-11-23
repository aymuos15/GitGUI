#!/bin/bash
set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}Building gg binaries...${NC}"

# Build Linux AMD64
echo -e "${BLUE}Building linux-amd64...${NC}"
GOOS=linux GOARCH=amd64 go build -o gg-linux-amd64 .
echo -e "${GREEN}✓ gg-linux-amd64 built successfully${NC}"

# Build Linux ARM64
echo -e "${BLUE}Building linux-arm64...${NC}"
GOOS=linux GOARCH=arm64 go build -o gg-linux-arm64 .
echo -e "${GREEN}✓ gg-linux-arm64 built successfully${NC}"

echo -e "${GREEN}All binaries built successfully!${NC}"
ls -lh gg-linux-*
