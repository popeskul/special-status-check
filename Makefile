.PHONY: build run test generate install-lint lint lint-ci swagger

DOCKER_IMAGE_NAME=go_app

deps:
	go mod download

build:
	go build -o bin/special-status-check cmd/statuscheckservice/main.go

install-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.49.0

lint:
	golangci-lint run ./...

lint-ci:
	golangci-lint run ./... --out-format=colored-line-number --timeout=5m

swagger:
	swag init -g cmd/statuscheckservice/main.go

test: generate
	go test ./...

generate:
	go generate ./...

run: generate
	go run ./cmd/statuscheckservice/main.go

docker-compose-up:
	docker-compose up --build -d

docker-compose-down:
	docker-compose down
