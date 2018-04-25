SHELL := /bin/bash

dep-metalinter:
	@go get -u gopkg.in/alecthomas/gometalinter.v2
	@gometalinter.v2 --install > /dev/null
check: check-metalinter
check-metalinter: dep-metalinter
	@gometalinter.v2 --fast --config=metalinter.json ./...
check-all: check-all-metalinter
check-all-metalinter: dep-metalinter
	@gometalinter.v2 --config=metalinter.json ./...
update-x-http:
	./scripts/x-http-updater.sh
test:
	go test -race -v $$(go list ./...)
coverage:
	go test -race -v -cover -coverprofile=coverage.out $$(go list ./...)
cover:
	go tool cover -html=coverage.out
