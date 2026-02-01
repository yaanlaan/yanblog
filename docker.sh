#!/bin/bash

# Ensure LF line endings
echo ">>> Starting YanBlog Deployment (Single Container Mode)..."

# 1. Cleaning up old containers
echo "1. Removing old containers..."
docker-compose down

# 2. Build and Start
echo "2. Building and Starting..."
docker-compose up -d --build

echo ">>> Deployment Complete!"
echo "    Frontend: http://localhost:3002"
echo "    Admin:    http://localhost:3011"

