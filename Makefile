SHELL := /bin/bash

export PATH := $(shell pwd)/bin:${PATH}

tool-metalinter:
	@./scripts/install-metalinter.sh
check: check-metalinter
check-metalinter: tool-metalinter
	./bin/gometalinter --fast --config=metalinter.json ./...
check-all: check-all-metalinter
check-all-metalinter: tool-metalinter
	./bin/gometalinter --config=metalinter.json ./...
update-x-http:
	./.scripts/x-http-updater.sh
test:
	go test -race -v $$(go list ./...)
coverage:
	go test -race -v -cover -coverprofile=coverage.out $$(go list ./...)
cover:
	go tool cover -html=coverage.out
