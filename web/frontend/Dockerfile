# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .
RUN npm run build

# Production stage
FROM nginx:alpine

# 1. 复制 Vite 构建产物
COPY --from=builder /app/dist /usr/share/nginx/html

# 2. [强制修复] 显式复制 static 目录，防止 Vite 构建时遗漏
# 注意：前端 build context 是 web/frontend
COPY public/static /usr/share/nginx/html/static

COPY nginx.conf /etc/nginx/conf.d/default.conf

# 3. 权限修复
RUN chmod -R 755 /usr/share/nginx/html

EXPOSE 80



CMD ["nginx", "-g", "daemon off;"]
