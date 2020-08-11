.PHONY: test
dev:
	clear && go run main.go

test:
	clear && DATABASE_NAME=uptime_server_service_test go test ./...

build:
	clear && go build

test_build: build
	rm -rf uptime-server-microservice