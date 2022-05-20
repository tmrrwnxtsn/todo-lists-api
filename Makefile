.PHONY:

build:
	go build -v ./cmd/apiserver

run: build
	./apiserver

tidy:
	go mod tidy

.DEFAULT_GOAL := run