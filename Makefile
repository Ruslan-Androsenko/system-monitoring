BIN_CLIENT := "./bin/client-monitoring"
BIN_SERVER := "./bin/system-monitoring"
CONTAINER_NAME="system-monitoring"
DOCKER_IMG="system-monitoring:develop"
LINTER_PATH=/tmp/bin
LINTER_BIN=/tmp/bin/golangci-lint

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build-client:
	go build -v -o $(BIN_CLIENT) -ldflags "$(LDFLAGS)" ./cmd/client

build-server:
	go build -v -o $(BIN_SERVER) -ldflags "$(LDFLAGS)" ./cmd/server

run: build-server
	$(BIN_SERVER) -config ./configs/config.toml -port 8070

client: build-client
	$(BIN_CLIENT) -config ./configs/config-client.toml -messages 100

server: build-server
	$(BIN_SERVER) -config ./configs/config-client.toml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run --rm \
		--env="SERVER_PORT=8070" \
		--name=$(CONTAINER_NAME) \
		--publish="8070:8070" \
		$(DOCKER_IMG)

up:
	docker compose up -d --build

down:
	docker compose down

restart: down up

version: build-server
	$(BIN_SERVER) version

test:
	go test -race -v -count 100 -timeout=120m ./internal/...

integration-test:
	go test -race -v -count 100 -timeout=120m ./cmd/server/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(LINTER_PATH) v1.55.2

lint: install-lint-deps
	$(LINTER_BIN) run ./...

.PHONY: build-client build-server run server client build-img run-img up down restart version test integration-test lint
