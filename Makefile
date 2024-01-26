BIN := "./bin/system-monitoring"
CONTAINER_NAME="system-monitoring"
DOCKER_IMG="system-monitoring:develop"
LINTER_PATH=/tmp/bin
LINTER_BIN=/tmp/bin/golangci-lint

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd

run: build
	$(BIN) -config ./configs/config.toml -port 8090

server: build
	$(BIN) -config ./configs/config-client.toml

client: build
	$(BIN) -config ./configs/config-client.toml -messages 100 grpc-client

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run --rm --name=$(CONTAINER_NAME) --network="host" $(DOCKER_IMG)

up:
	docker compose up -d --build

down:
	docker compose down

restart: down up

version: build
	$(BIN) version

test:
	go test -race -v -count 100 -timeout=20m ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(LINTER_PATH) v1.55.2

lint: install-lint-deps
	$(LINTER_BIN) run ./...

.PHONY: build run build-img run-img version test lint
