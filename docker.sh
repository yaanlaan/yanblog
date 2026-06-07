#!/bin/bash
set -e

echo ">>> Starting YanBlog Deployment (Single Container Mode)..."

COMPOSE_FILE="docker-compose.yaml"
CONFIG_DIR="config"
BACKEND_CONFIG="$CONFIG_DIR/backend/config.yaml"
FRONTEND_CONFIG="$CONFIG_DIR/frontend/config.yaml"

if [ ! -d "$CONFIG_DIR/backend" ] || [ ! -f "$BACKEND_CONFIG" ]; then
    echo ">>> Initializing configuration..."
    ./init.sh
fi

if [ -f "$COMPOSE_FILE" ]; then
    ADMIN_PORT=$(sed -n 's/.*- "\?\([0-9]*\):81"\?.*/\1/p' "$COMPOSE_FILE" | head -n 1 | xargs)
    
    if [ ! -z "$ADMIN_PORT" ] && [ -f "$FRONTEND_CONFIG" ]; then
        echo ">>> Detected Admin Port from docker-compose: $ADMIN_PORT"
        sed "s/dev_admin_port: [0-9]*/dev_admin_port: $ADMIN_PORT/" "$FRONTEND_CONFIG" > "${FRONTEND_CONFIG}.tmp" && mv "${FRONTEND_CONFIG}.tmp" "$FRONTEND_CONFIG"
        echo ">>> Updated $FRONTEND_CONFIG with port $ADMIN_PORT"
    fi
fi

echo ">>> Removing old containers..."
docker compose down

echo ">>> Building and Starting..."
docker compose up -d --build

echo ">>> Deployment Complete!"
echo "    Frontend: http://localhost:3002"
echo "    Admin:    http://localhost:3011"
echo "    API:      http://localhost:8080"
echo ""
echo ">>> Configuration Files:"
echo "    Backend: $BACKEND_CONFIG"
echo "    Frontend: $FRONTEND_CONFIG"
echo ""
echo ">>> Tips:"
echo "    - Modify config files and restart container to apply changes"
echo "    - Use environment variables to override config values"
echo "    - Run 'docker compose logs -f' to view logs"