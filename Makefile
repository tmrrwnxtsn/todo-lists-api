.PHONY:

build:
	docker-compose build

run:
	docker-compose up

tidy:
	go mod tidy

migrate-up:
	migrate -path ./migrations -database "postgres://user:pass@0.0.0.0:5433/postgres?sslmode=disable" up

migrate-down:
	migrate -path ./migrations -database "postgres://user:pass@0.0.0.0:5433/postgres?sslmode=disable" down

.DEFAULT_GOAL := run