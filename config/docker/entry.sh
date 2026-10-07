#!/bin/sh
set -e

echo '[entry] Checking config...'

mkdir -p /app/config/backend /app/config/frontend

if [ ! -f /app/config/backend/config.yaml ] && [ ! -f /app/config/config.yaml ]; then
    cp /app/config/config.yaml.template /app/config/config.yaml
    echo '[entry] Backend config initialized from template'
fi

if [ ! -f /app/config/frontend/config.yaml ]; then
    cp /app/web/frontend/public/config.yaml /app/config/frontend/config.yaml
    echo '[entry] Frontend config initialized from demo'
fi

mkdir -p /app/data
mkdir -p /app/uploads/defaults

for f in /app/defaults/*; do
    fname=$(basename $f)
    if [ ! -f /app/uploads/defaults/$fname ]; then
        cp $f /app/uploads/defaults/$fname
        echo "[entry] Restored default image: $fname"
    fi
done

echo '[entry] Starting nginx...'
nginx

echo '[entry] Starting Go backend...'
exec /app/yanblog