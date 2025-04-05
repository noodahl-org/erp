FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o main ./internal/api

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy only the built binary and environment file
COPY --from=builder /app/internal/api/.env ./
COPY --from=builder /app/main ./

USER appuser
CMD ["./main"]