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

# Create config directory and copy default config if it doesn't exist
RUN mkdir -p config
COPY --from=builder /app/config/config.yaml ./config/config.yaml

# Create uploads directory
RUN mkdir -p uploads

# Create the directory structure for frontend static files and copy them
RUN mkdir -p web/frontend/public
COPY --from=builder /app/web/frontend/public/static/. ./web/frontend/public/static/

# Create the frontend config path and make sure both paths point to the same file
RUN mkdir -p web/frontend/public
# Copy the initial config from the docker_field to the shared location
COPY --from=builder /app/docker_field/frontend/config.yaml ./frontend_config.yaml
# Create a symbolic link so both paths point to the same file
RUN ln -sf /app/frontend_config.yaml /app/web/frontend/public/config.yaml

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./server"]