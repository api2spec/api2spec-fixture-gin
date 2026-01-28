FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go mod files first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build
RUN CGO_ENABLED=0 go build -o server ./cmd/server/main.go

FROM alpine:3.20

WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 3000
CMD ["./server"]
