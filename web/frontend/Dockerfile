# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .
RUN npm run build

# Production stage
FROM nginx:alpine

COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Use config.docker.yaml as the main config.yaml
RUN cp /usr/share/nginx/html/config.docker.yaml /usr/share/nginx/html/config.yaml

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
