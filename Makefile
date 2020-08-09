dev:
	clear && go run main.go

test:
	clear && go test ./...

build:
	clear && go build

test_build: build
	rm -rf uptime-server-microservice