#!/bin/bash

set -e

CONFIG_DIR="config"
BACKEND_CONFIG="$CONFIG_DIR/backend/config.yaml"
FRONTEND_CONFIG="$CONFIG_DIR/frontend/config.yaml"
BACKEND_TEMPLATE="$CONFIG_DIR/config_template.yaml"
FRONTEND_TEMPLATE="$CONFIG_DIR/frontend/config_template.yaml"
OLD_CONFIG="config/config.yaml"
OLD_FRONTEND_CONFIG="web/frontend/public/config.yaml"

echo "=== YanBlog Configuration Initialization ==="

mkdir -p "$CONFIG_DIR/backend" "$CONFIG_DIR/frontend"

if [ -f "$OLD_CONFIG" ] && [ ! -f "$BACKEND_CONFIG" ]; then
    echo "Migrating old backend config..."
    cp "$OLD_CONFIG" "$BACKEND_CONFIG"
    echo "Old config migrated to $BACKEND_CONFIG"
elif [ -f "$BACKEND_TEMPLATE" ] && [ ! -f "$BACKEND_CONFIG" ]; then
    cp "$BACKEND_TEMPLATE" "$BACKEND_CONFIG"
    echo "$BACKEND_TEMPLATE 已复制为 $BACKEND_CONFIG"
elif [ ! -f "$BACKEND_CONFIG" ]; then
    cat > "$BACKEND_CONFIG" << 'EOF'
server:
  AppMode: debug
  HttpPort: :8080
  SiteUrl: ""

database:
  Db: MYSQL
  DbHost: db
  DbPort: 3306
  DbUser: root
  DbPassWord: rootpassword
  DbName: yanblog

JwtKey: ""

weather:
  Provider: openweathermap
  ApiKey: ""
  DefaultCity: Hefei

FrontEndConfigPath: config/frontend/config.yaml
EOF
    echo "Created default backend config at $BACKEND_CONFIG"
fi

if [ -f "$OLD_FRONTEND_CONFIG" ] && [ ! -f "$FRONTEND_CONFIG" ]; then
    echo "Migrating old frontend config..."
    cp "$OLD_FRONTEND_CONFIG" "$FRONTEND_CONFIG"
    echo "Old frontend config migrated to $FRONTEND_CONFIG"
elif [ -f "$FRONTEND_TEMPLATE" ] && [ ! -f "$FRONTEND_CONFIG" ]; then
    cp "$FRONTEND_TEMPLATE" "$FRONTEND_CONFIG"
    echo "$FRONTEND_TEMPLATE 已复制为 $FRONTEND_CONFIG"
elif [ ! -f "$FRONTEND_CONFIG" ]; then
    echo "Error: No frontend config template found!"
    exit 1
fi

echo "=== Configuration Initialization Complete ==="
echo ""
echo "Please edit the following files to configure your blog:"
echo "- Backend: $BACKEND_CONFIG"
echo "- Frontend: $FRONTEND_CONFIG"
echo ""
echo "Important:"
echo "- Set JwtKey in backend config (run: openssl rand -hex 32)"
echo "- Configure database connection if not using Docker"
echo "- Customize frontend settings (blog name, author info, etc.)"