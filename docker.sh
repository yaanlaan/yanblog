#!/bin/bash
set -e

# Ensure LF line endings
echo ">>> Starting YanBlog Deployment (Single Container Mode)..."

# --- 0. Auto-Configure Ports ---
COMPOSE_FILE="docker-compose.yaml"
CONFIG_FILE="web/frontend/public/config.yaml"

if [ -f "$COMPOSE_FILE" ]; then
    # Extract Admin Port (mapped to internal 81)
    # Uses sed to capture the digits before :81
    ADMIN_PORT=$(sed -n 's/.*- "\?\([0-9]*\):81"\?.*/\1/p' "$COMPOSE_FILE" | head -n 1 | xargs)
    
    if [ ! -z "$ADMIN_PORT" ]; then
        echo ">>> Detected Admin Port from docker-compose: $ADMIN_PORT"
        
        if [ -f "$CONFIG_FILE" ]; then
            # Update dev_admin_port using sed. 
            # Uses a temp file approach for compatibility with both GNU sed (Linux) and BSD sed (Mac) logic
            sed "s/dev_admin_port: [0-9]*/dev_admin_port: $ADMIN_PORT/" "$CONFIG_FILE" > "${CONFIG_FILE}.tmp" && mv "${CONFIG_FILE}.tmp" "$CONFIG_FILE"
            echo ">>> Updated $CONFIG_FILE with port $ADMIN_PORT"
        fi
    fi
fi
# -----------------------------

# 1. Cleaning up old containers
echo "1. Removing old containers..."
docker compose down

# 2. Build and Start
echo "2. Building and Starting..."
docker compose up -d --build

echo ">>> Deployment Complete!"
echo "    Frontend: http://localhost:3002"
echo "    Admin:    http://localhost:3011"

