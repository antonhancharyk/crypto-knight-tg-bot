.PHONY: build test lint run

build:
	go build -o bin/tgbot ./cmd/tgbot

test:
	go test -v -race -count=1 ./...

lint:
	golangci-lint run ./...

run: build
	./bin/tgbot
