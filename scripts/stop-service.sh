#!/usr/bin/env bash

set -e

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

SERVICE_NAME="go-convert-service"
PORT=${PORT:-8090}

echo -e "${YELLOW}Stopping ${SERVICE_NAME}...${NC}"

# Stop local process
if lsof -Pi :$PORT -sTCP:LISTEN -t >/dev/null 2>&1; then
    PID=$(lsof -Pi :$PORT -sTCP:LISTEN -t)
    echo "Killing process on port $PORT (PID: $PID)..."
    kill -9 "$PID" 2>/dev/null || true
    echo -e "${GREEN}✓ Local service stopped${NC}"
fi

# Stop Docker container
if docker ps --format '{{.Names}}' | grep -q "^${SERVICE_NAME}$"; then
    echo "Stopping Docker container..."
    docker stop "$SERVICE_NAME" >/dev/null 2>&1 || true
    docker rm "$SERVICE_NAME" >/dev/null 2>&1 || true
    echo -e "${GREEN}✓ Docker container stopped${NC}"
fi

echo -e "${GREEN}Service stopped${NC}"
