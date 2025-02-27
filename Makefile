BINARY_NAME=bot
DOCKER_IMAGE=discord-bot
VERSION=1.0.0

build:
	go build -o bin/$(BINARY_NAME) cmd/bot/main.go

run: build
	./bin/$(BINARY_NAME)

clean:
	go clean
	rm -rf bin/

test:
	go test -v ./...

lint:
	go vet ./...
	go fmt ./...

docker-build:
	docker build -t $(DOCKER_IMAGE):$(VERSION) .

docker-run:
	docker run --env-file config/.env $(DOCKER_IMAGE):$(VERSION)

docker-push:
	docker push $(DOCKER_IMAGE):$(VERSION)

setup:
	cp config/.env.example config/.env
	go mod download
	go mod tidy

.PHONY: build run clean test lint docker-build docker-run docker-push setup 