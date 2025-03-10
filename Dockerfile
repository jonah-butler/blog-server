# ---- Step 1: Build Stage ----
  FROM golang:1.22.5-alpine AS builder

  # Set environment variables for Go optimization
  ENV CGO_ENABLED=0 \
      GOOS=linux \
      GOARCH=amd64
  
  # Install dependencies
  RUN apk --no-cache add git
  
  WORKDIR /app
  
  # Cache dependencies first for better build efficiency
  COPY go.mod go.sum ./
  RUN go mod download
  
  # Copy all project files
  COPY . .
  
  # Build the Go application (adjusted path)
  RUN go build -o /blog_api ./cmd
  
  # ---- Step 2: Final Stage ----
  FROM alpine:latest
  
  # Create a non-root user
  RUN addgroup -S appgroup && adduser -S appuser -G appgroup
  
  WORKDIR /app
  
  # Copy the built binary from the builder stage
  COPY --from=builder /blog_api /app/api

  # COPY .env file from the host machine into the final container
COPY .env /app/.env
  
  # Grant execution permission
  RUN chmod +x /app/api
  
  # Expose the API port
  EXPOSE 8080
  
  # Switch to non-root user
  USER appuser

  ENV RUNNING_IN_DOCKER=true

  # Run the API
  CMD ["/app/api"]
  