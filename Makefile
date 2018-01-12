SHELL := /bin/bash

dep: dep-tool dep-metalinter
dep-tool:
	@go get github.com/golang/dep/cmd/dep
dep-metalinter:
	@go get github.com/alecthomas/gometalinter
	@gometalinter --install > /dev/null
check: check-metalinter
check-metalinter: dep-metalinter
	@gometalinter --fast --config=metalinter.json ./...
check-all: check-all-metalinter
check-all-metalinter: dep-metalinter
	@gometalinter --config=metalinter.json ./...
redis-cluster-start:
	./scripts/redis-cluster.sh start
redis-cluster-stop:
	./scripts/redis-cluster.sh stop
update-x-http:
	./scripts/x-http-updater.sh
test:
	go test -race -v $$(go list ./... | grep -Ev "vendor|qtest")
coverage:
	@./scripts/coverage.sh coverage.out
cover:
	go tool cover -html=coverage.out
