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
RUN mkdir -p /app/config /app/uploads

# Copy config template (host mount overrides at runtime)
COPY config/config_template.yaml /app/config/config.yaml.template
# Copy frontend config (Go expects config/frontend/config.yaml)
RUN mkdir -p /app/config/frontend
COPY config/frontend/config.yaml /app/config/frontend/config.yaml

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
# Frontend config (for Go to serve at /config.yaml)
COPY config/frontend/config.yaml /app/web/frontend/public/config.yaml

# --- Permissions ---
RUN chmod -R 755 /usr/share/nginx/html && chmod -R 755 /app

# --- Nginx Config ---
COPY nginx.conf.unified /etc/nginx/conf.d/default.conf

EXPOSE 80 81

# Startup: use host config if mounted, otherwise fall back to template
# Nginx runs in background, Go server runs in foreground (so logs/errors are visible)
CMD ["/bin/sh", "-c", "\
  if [ -f /app/host-config/backend/config.yaml ]; then \
    cp /app/host-config/backend/config.yaml /app/config/config.yaml; \
    echo '[entry] Using host backend config'; \
  else \
    cp /app/config/config.yaml.template /app/config/config.yaml; \
    echo '[entry] Using default config (template)'; \
  fi && \
  if [ -f /app/host-config/frontend/config.yaml ]; then \
    cp /app/host-config/frontend/config.yaml /app/config/frontend/config.yaml; \
    cp /app/host-config/frontend/config.yaml /app/web/frontend/public/config.yaml; \
    echo '[entry] Using host frontend config'; \
  fi && \
  echo '[entry] Starting nginx...' && \
  nginx && \
  echo '[entry] Starting Go backend...' && \
  exec /app/server \
"]
