DSN ?= postgres://127.0.0.1/todo_lists_db?sslmode=disable&user=postgres&password=qwerty

.PHONY: tidy
tidy: ## check go.mod
	go mod tidy

.PHONY: test
test: ## run unit tests
	@go test -cover -covermode=count ./...

.PHONY: test-cover
test-cover: ## show test coverage information
	go test -cover -coverprofile=cover.out ./... && go tool cover -html=cover.out -o cover.html

.PHONY: build
build:  ## build the API server binary
	go build -o apiserver -a ./cmd/apiserver

.PHONY: run-binary
run-binary: build  ## run the API server binary
	./apiserver

.PHONY: run
run: ## run the API server
	go run cmd/apiserver/main.go

.PHONY: migrate-up
migrate-up: ## run all new local database migrations
	@echo "Running all new local database migrations..."
	@migrate -path ./migrations -database "$(DSN)" up

.PHONY: migrate-down
migrate-down: ## revert local database to the last migration step
	@echo "Reverting local database to the last migration step..."
	@migrate -path ./migrations -database "$(DSN)" down 1

.PHONY: compose-up
compose-up: ## builds the images if the images do not exist and starts the containers (API server and postgres DB)
	docker-compose up

.DEFAULT_GOAL := compose-up