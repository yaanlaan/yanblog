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

# Install necessary runtime dependencies (if any)
# RUN apk --no-cache add ca-certificates

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/server .

# Copy config directory (assuming config.yaml is present or will be mounted)
COPY --from=builder /app/config ./config
# Create uploads directory
RUN mkdir -p uploads

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./server"]
