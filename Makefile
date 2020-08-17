.PHONY: test verbose_test coverage
dev:
	clear && go run main.go

test:
	clear && go test ./...

verbose_test:
	clear && DATABASE_NAME=uptime_server_service_test go test ./... -v

build:
	clear && go build

coverage:
	clear && go test ./... -coverprofile=coverage.txt -covermode=atomic

test_build: build
	rm -rf uptime-server-microservice