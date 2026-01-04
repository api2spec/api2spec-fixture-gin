.PHONY: run build test clean tidy lint

run:
	go run ./cmd/server/main.go

build:
	go build -o bin/server ./cmd/server/main.go

test:
	go test -v ./...

clean:
	rm -rf bin/

tidy:
	go mod tidy

lint:
	golangci-lint run

validate:
	api2spec check --ci
