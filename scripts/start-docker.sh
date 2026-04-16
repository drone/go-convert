#!/usr/bin/env bash

set -e

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default configuration
PORT=${PORT:-8090}
LOG_LEVEL=${LOG_LEVEL:-debug}
MAX_BATCH_SIZE=${MAX_BATCH_SIZE:-100}
MAX_YAML_BYTES=${MAX_YAML_BYTES:-1048576}

SERVICE_NAME="go-convert-service"
DOCKER_IMAGE="${SERVICE_NAME}:latest"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

echo -e "${GREEN}Starting ${SERVICE_NAME} in Docker...${NC}"
echo ""
echo "Configuration:"
echo "  PORT:            $PORT"
echo "  LOG_LEVEL:       $LOG_LEVEL"
echo "  MAX_BATCH_SIZE:  $MAX_BATCH_SIZE"
echo "  MAX_YAML_BYTES:  $MAX_YAML_BYTES"
echo ""

# Navigate to project root
cd "$PROJECT_ROOT"

# Stop existing container if running
if docker ps -a --format '{{.Names}}' | grep -q "^${SERVICE_NAME}$"; then
    echo -e "${YELLOW}Stopping existing container...${NC}"
    docker stop "$SERVICE_NAME" >/dev/null 2>&1 || true
    docker rm "$SERVICE_NAME" >/dev/null 2>&1 || true
fi

# Build Docker image
echo -e "${YELLOW}Building Docker image...${NC}"
if docker build -f Dockerfile.service -t "$DOCKER_IMAGE" .; then
    echo -e "${GREEN}✓ Docker image built${NC}"
else
    echo -e "${RED}✗ Docker build failed${NC}"
    exit 1
fi

# Run container
echo -e "${GREEN}Starting Docker container...${NC}"
docker run -d \
    --name "$SERVICE_NAME" \
    -p "$PORT:8090" \
    -e LOG_LEVEL="$LOG_LEVEL" \
    -e MAX_BATCH_SIZE="$MAX_BATCH_SIZE" \
    -e MAX_YAML_BYTES="$MAX_YAML_BYTES" \
    "$DOCKER_IMAGE"

# Wait for service to start
echo -e "${YELLOW}Waiting for service to be ready...${NC}"
sleep 2

# Health check
if curl -s "http://localhost:$PORT/healthz" | grep -q "ok"; then
    echo -e "${GREEN}✓ Service is running and healthy${NC}"
    echo ""
    echo "Service details:"
    echo "  Container: $SERVICE_NAME"
    echo "  Health:    http://localhost:$PORT/healthz"
    echo "  API:       http://localhost:$PORT/api/v1/convert/batch"
    echo ""
    echo "View logs with:"
    echo "  docker logs -f $SERVICE_NAME"
    echo ""
    echo "Stop service with:"
    echo "  docker stop $SERVICE_NAME && docker rm $SERVICE_NAME"
else
    echo -e "${RED}✗ Service health check failed${NC}"
    echo "Check logs with: docker logs $SERVICE_NAME"
    exit 1
fi
