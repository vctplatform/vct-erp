# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o cms-api ./cmd/api

# Production stage
FROM alpine:3.19
RUN apk add --no-cache git ca-certificates
WORKDIR /app
COPY --from=builder /build/cms-api .
EXPOSE 8080
CMD ["./cms-api"]
