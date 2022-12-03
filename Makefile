# ==============================================================================
# Main

http:
	export config=local && go run ./cmd/api/http/main.go

grpc:
	export config=local && go run ./cmd/api/grpc/main.go

redis-pubsub:
	export config=local && go run ./cmd/elasticsearch-indexer-redis/main.go

build:
	export config=local &&\
	go env -w CGO_ENABLED=1 &&\
		go build -o ./bin/api-http ./cmd/api/http/main.go &&\
		go build -o ./bin/api-grpc ./cmd/api/grpc/main.go &&\
		go build -o ./bin/elasticsearch-indexer-redis ./cmd/elasticsearch-indexer-redis/main.go

# ==============================================================================
# Modules support

deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

tidy:
	go mod tidy
	go mod vendor

deps-upgrade:
	# go get $(go list -f '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}' -m all)
	go get -u -t -d -v ./...
	go mod tidy
	go mod vendor

deps-cleancache:
	go clean -modcache

# ==============================================================================
# Docker compose commands

develop:
	echo "Starting docker environment"
	export config=docker
	config=docker && docker compose -f docker-compose.yml up --build -d

kill_dev:
	echo "Killing docker environment"
	docker compose -f docker-compose.yml kill

stop_dev:
	echo "Stopping docker environment"
	docker compose -f docker-compose.yml stop

remove_dev:
	echo "Removing docker environment"
	docker compose -f docker-compose.yml down -v

local:
	echo "Starting local environment"
	export config=local
	docker compose -f docker-compose.local.yml up --build -d

kill_local:
	echo "Killing docker environment"
	docker compose -f docker-compose.local.yml kill

stop_local:
	echo "Stopping docker environment"
	docker compose -f docker-compose.local.yml stop

remove_local:
	echo "Removing docker environment"
	docker compose -f docker-compose.local.yml down -v

# ==============================================================================
# Tools

proto_update:
	buf mod update ./pkg/grpc/v1

proto_gen:
	buf generate ./pkg/grpc/v1

proto_lint:
	buf lint ./pkg/grpc/v1

swagger_init:
	swag init -g ./cmd/api/http/main.go