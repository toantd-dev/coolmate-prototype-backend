# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git make

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Run tests
RUN go test -v ./... -timeout 60s

# Build application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o coolmate-backend ./cmd/api

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/coolmate-backend .

# Copy environment template
COPY .env.example .env

EXPOSE 8080

CMD ["./coolmate-backend"]
