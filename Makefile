.PHONY: test verbose_test coverage
dev:
	clear && go run main.go

test:
	clear && GIN_MODE=release go test -p 1 ./...

verbose_test:
	clear && DATABASE_NAME=uptime_server_service_test go test ./... -v

build:
	clear && go build

coverage:
	clear && GIN_MODE=release go test -p 1 ./... -coverprofile=coverage.txt -covermode=atomic

test_build: build
	rm -rf uptime-server-microservice