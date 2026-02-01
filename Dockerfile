# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o server .

# Final stage
FROM alpine:latest

WORKDIR /app

# Install necessary runtime dependencies
RUN apk --no-cache add ca-certificates

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/server .

# Create config directory
RUN mkdir -p config
# Copy the config from the build context (which was modified by docker.sh if needed)
COPY config/config.yaml ./config/config.yaml

# Create uploads directory
RUN mkdir -p uploads

# Create the directory structure for frontend static files and copy them
# ensuring parent directories exist first
RUN mkdir -p web/frontend/public/static
COPY --from=builder /app/web/frontend/public/static/ ./web/frontend/public/static/

# Frontend config handling
RUN mkdir -p web/frontend/public
# Copy the frontend config (which might have been modified by docker.sh)
COPY web/frontend/public/config.yaml ./frontend_config.yaml
# Create a symbolic link so both paths point to the same file
RUN ln -sf /app/frontend_config.yaml /app/web/frontend/public/config.yaml

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./server"]