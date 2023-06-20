FROM golang:1.19 AS builder

WORKDIR /src/
COPY . .
RUN go mod tidy &&\
    GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -a -o bin/http cmd/api/http/main.go &&\
    GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -a -o bin/grpc cmd/api/grpc/main.go &&\
    GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -a -o bin/redis-ce cmd/elasticsearch-indexer-redis/main.go

FROM debian:buster-slim AS api-http
WORKDIR /api/
COPY --from=builder ["/src/bin/http", "/src/wait-for-it.sh", "/api/"]
RUN chmod +x /api/wait-for-it.sh

FROM debian:buster-slim AS api-grpc
WORKDIR /api/
COPY --from=builder ["/src/bin/grpc", "/src/wait-for-it.sh", "/api/"]
RUN chmod +x /api/wait-for-it.sh

FROM debian:buster-slim AS redis-pubsub
WORKDIR /app/
COPY --from=builder ["/src/bin/redis-ce", "/app/"]
