.PHONY:

build:
	go build -v ./cmd/apiserver

run: build
	./apiserver

tidy:
	go mod tidy

run-db:
	docker run --name todo-db -e POSTGRES_PASSWORD=qwerty -p 5433:5432 -d --rm postgres

migrate-up:
	migrate -path ./migrations -database "postgres://postgres:qwerty@localhost:5433/postgres?sslmode=disable" up

migrate-down:
	migrate -path ./migrations -database "postgres://postgres:qwerty@localhost:5433/postgres?sslmode=disable" down

.DEFAULT_GOAL := run