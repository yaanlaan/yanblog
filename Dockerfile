# Stage 1: Build Go Backend
FROM golang:1.24-alpine AS backend-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server .

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

# Install basic dependencies
RUN apk add --no-cache ca-certificates tzdata netcat-openbsd

WORKDIR /app

# 1. Setup Backend
COPY --from=backend-builder /app/server .
COPY entrypoint.sh .
# Fix line endings (CRLF -> LF) for Windows hosts
RUN sed -i 's/\r$//' entrypoint.sh
RUN chmod +x entrypoint.sh
RUN mkdir -p config uploads

# 2. Setup Frontend Files
RUN mkdir -p /usr/share/nginx/html/web
COPY --from=frontend-builder /app/dist /usr/share/nginx/html/web
# COPY --from=frontend-builder /app/public/static /usr/share/nginx/html/web/static
# COPY --from=frontend-builder /app/public/iconfont /usr/share/nginx/html/web/iconfont
# COPY --from=frontend-builder /app/public/favicon.ico /usr/share/nginx/html/web/favicon.ico

# Ensure permissions are correct for Nginx
RUN chmod -R 755 /usr/share/nginx/html

# 3. Setup Admin Files
RUN mkdir -p /usr/share/nginx/html/admin
COPY --from=admin-builder /app/dist /usr/share/nginx/html/admin

# 4. Link Paths for Go Backend Compatibility
# Go code expects: "./web/frontend/public/static/about.md"
# Real location: "/usr/share/nginx/html/web/static/about.md"
RUN mkdir -p /app/web/frontend/public && \
    ln -s /usr/share/nginx/html/web/static /app/web/frontend/public/static && \
    ln -s /usr/share/nginx/html/web/iconfont /app/web/frontend/public/iconfont && \
    ln -s /usr/share/nginx/html/web/favicon.ico /app/web/frontend/public/favicon.ico

# 5. Bake in Configuration (Avoids Host Mount Issues)
COPY config/config.yaml /app/config/config.yaml
# 直接覆盖前端静态目录中的配置，避免 Nginx alias 权限/路径问题
COPY web/frontend/public/config.yaml /usr/share/nginx/html/web/config.yaml

# 创建软链接：让后端(操作 /app/config/frontend_config.yaml) 实际上修改的是 Nginx 的文件
# 这样就实现了后端修改，前端生效，且共享同一份配置
RUN ln -sf /usr/share/nginx/html/web/config.yaml /app/config/frontend_config.yaml

# 确保配置文件存在
RUN test -f /app/config/config.yaml || (echo "Error: config.yaml not found!" && exit 1) && \
    test -L /app/config/frontend_config.yaml || (echo "Error: frontend_config.yaml symlink failed!" && exit 1)

# 确保所有文件对 Nginx 用户(通常是 nginx:nginx)可读
# 赋予所有目录 +x 权限，所有文件 +r 权限
RUN chown -R nginx:nginx /usr/share/nginx/html && \
    chmod -R 755 /usr/share/nginx/html && \
    chmod -R 755 /app

# 6. Setup Nginx Config
COPY nginx.conf.unified /etc/nginx/conf.d/default.conf

# Expose ports
EXPOSE 80 81

CMD ["/app/entrypoint.sh"]
