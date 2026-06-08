# Stage 1: Build Go Backend
FROM golang:1.24-alpine AS backend-builder
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -o server .

# Stage 2: Build Frontend (Public)
FROM node:20-alpine AS frontend-builder
WORKDIR /app
COPY web/frontend/package*.json ./
RUN npm install
COPY web/frontend ./
RUN npm run build

# Stage 3: Build Admin (Backend UI)
FROM node:20-alpine AS admin-builder
WORKDIR /app
COPY web/backend/package*.json ./
RUN npm install
COPY web/backend ./
RUN npm run build

# Stage 4: Final Unified Image
FROM nginx:alpine

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# --- Backend ---
COPY --from=backend-builder /app/server .
RUN mkdir -p /app/config /app/data /app/uploads
# Default images for recovery
COPY uploads/defaults /app/uploads/defaults

# Copy backend config template (used on first run if no host config exists)
COPY config/config_template.yaml /app/config/config.yaml.template
# Frontend config is initialized at runtime from demo if missing

# --- Frontend ---
RUN mkdir -p /usr/share/nginx/html/web
COPY --from=frontend-builder /app/dist /usr/share/nginx/html/web

# --- Admin ---
RUN mkdir -p /usr/share/nginx/html/admin
COPY --from=admin-builder /app/dist /usr/share/nginx/html/admin

# --- Static files for Go backend routes ---
# Go serves /static, /iconfont, /favicon.ico from ./web/frontend/public/
RUN mkdir -p /app/web/frontend/public
COPY web/frontend/public/static /app/web/frontend/public/static
COPY web/frontend/public/iconfont /app/web/frontend/public/iconfont
COPY web/frontend/public/favicon.ico /app/web/frontend/public/favicon.ico
COPY web/frontend/public/favicon.svg /app/web/frontend/public/favicon.svg
# Demo config (used as template on first run if no user config exists)
COPY web/frontend/public/config.yaml /app/web/frontend/public/config.yaml

# --- Permissions ---
RUN chmod -R 755 /usr/share/nginx/html && chmod -R 755 /app

# --- Nginx Config ---
COPY nginx.conf.unified /etc/nginx/conf.d/default.conf

EXPOSE 80 81

# Startup: initialize config on first run, then start services
# Config is mounted from host at /app/config — changes persist across restarts
CMD ["/bin/sh", "-c", "\
  echo '[entry] Checking config...' && \
  if [ ! -f /app/config/backend/config.yaml ] && [ ! -f /app/config/config.yaml ]; then \
    cp /app/config/config.yaml.template /app/config/config.yaml; \
    echo '[entry] Backend config initialized from template'; \
  fi && \
  if [ ! -f /app/config/frontend/config.yaml ]; then \
    mkdir -p /app/config/frontend; \
    cp /app/web/frontend/public/config.yaml /app/config/frontend/config.yaml; \
    echo '[entry] Frontend config initialized from demo'; \
  fi && \
  mkdir -p /app/data && \
  echo '[entry] Starting nginx...' && \
  nginx && \
  echo '[entry] Starting Go backend...' && \
  exec /app/server \
"]
