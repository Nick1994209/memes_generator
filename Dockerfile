# Build React frontend
FROM node:16-alpine AS frontend-build

WORKDIR /app
COPY web/package*.json ./
RUN npm install

COPY web/ .
RUN npm run build

# Build Golang backend
FROM golang:1.24-alpine AS backend-build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o memes-generator cmd/web/main.go
RUN go build -o generate-meme cmd/generate/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates and required packages for user management
RUN apk --no-cache add ca-certificates shadow

WORKDIR /app

# Copy built frontend
COPY --from=frontend-build /app/build ./web/build

# Copy built backend
COPY --from=backend-build /app/memes-generator .
COPY --from=backend-build /app/generate-meme .

# Create data directory
RUN mkdir -p ./data/memes ./data/templates

# Non-root user stage
# Create user and set ownership
RUN groupadd -g 1000 appuser && \
    useradd -u 1000 -g appuser -s /bin/bash -m appuser && \
    chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Run the application
CMD ["./memes-generator"]
