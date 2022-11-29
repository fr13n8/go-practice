from golang:1.19 AS builder

WORKDIR /src/
COPY . .
# go fmt $(go list ./... | grep -v /vendor/) &&\
# go vet $(go list ./... | grep -v /vendor/) &&\
RUN go mod tidy &&\
    GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -a -o bin/http cmd/api/http/main.go

FROM debian:buster-slim
WORKDIR /api/
COPY --from=builder ["/src/bin/http", "/src/database.db", "/api/"]

EXPOSE 80
CMD ["./http"]
