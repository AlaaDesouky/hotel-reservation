build:
	@go build -o bin/api

run: build
	@./bin/api

dev:
	@air

seed:
	@go run scripts/seed.go

test:
	@go test -v ./...