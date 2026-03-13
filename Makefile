.PHONY: help dev prod build test clean

help:
	@echo "Available commands:"
	@echo "  make dev    - Start development environment (docker)"
	@echo "  make prod   - Start production environment (docker)"
	@echo "  make build  - Build backend binary"
	@echo "  make test   - Run tests"
	@echo "  make clean  - Clean up"

dev:
	docker compose --profile dev up --build

prod:
	docker compose --profile prod up --build

build:
	cd backend && go build -o bin/server ./cmd/server

test:
	cd backend && go test -v ./...

clean:
	cd backend && rm -f bin/server
	docker compose down -v
