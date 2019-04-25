SHELL := /bin/bash

export PATH := $(shell pwd)/bin:${PATH}

tool-lint:
	@curl -sfL https://raw.githubusercontent.com/bukalapak/toolkit-installer/master/golangci-lint.sh | sh
check: check-lint
check-lint: tool-lint
	./bin/golangci-lint run
update-x-http:
	./scripts/x-http-updater.sh
test:
	go test -race -v $$(go list ./...)
coverage:
	go test -race -v -cover -coverprofile=coverage.out $$(go list ./...)
cover:
	go tool cover -html=coverage.out
