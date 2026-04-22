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
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

echo -e "${GREEN}Starting ${SERVICE_NAME}...${NC}"
echo ""
echo "Configuration:"
echo "  PORT:            $PORT"
echo "  LOG_LEVEL:       $LOG_LEVEL"
echo "  MAX_BATCH_SIZE:  $MAX_BATCH_SIZE"
echo "  MAX_YAML_BYTES:  $MAX_YAML_BYTES"
echo ""

# Navigate to project root
cd "$PROJECT_ROOT"

# Build the service
echo -e "${YELLOW}Building service...${NC}"
if go build -o "$SERVICE_NAME" ./cmd/server; then
    echo -e "${GREEN}✓ Build successful${NC}"
else
    echo -e "${RED}✗ Build failed${NC}"
    exit 1
fi

# Check if port is already in use
if lsof -Pi :$PORT -sTCP:LISTEN -t >/dev/null 2>&1 ; then
    echo -e "${RED}✗ Port $PORT is already in use${NC}"
    echo "Please stop the existing process or use a different port:"
    echo "  PORT=9000 ./scripts/start-service.sh"
    exit 1
fi

# Start the service
echo -e "${GREEN}Starting service on port $PORT...${NC}"
echo ""
echo -e "${YELLOW}Service logs:${NC}"
echo "============================================"

export PORT LOG_LEVEL MAX_BATCH_SIZE MAX_YAML_BYTES

./"$SERVICE_NAME"
