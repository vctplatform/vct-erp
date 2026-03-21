FROM golang:1.26.0 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/vct-finance-api ./cmd/api

FROM debian:bookworm-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates tzdata \
    && rm -rf /var/lib/apt/lists/* \
    && useradd --system --uid 10001 --create-home appuser

WORKDIR /app

COPY --from=builder /out/vct-finance-api /usr/local/bin/vct-finance-api
COPY --from=builder /app/migrations /app/migrations
COPY --from=builder /app/schema.sql /app/schema.sql

EXPOSE 8080

USER appuser

ENTRYPOINT ["vct-finance-api"]
