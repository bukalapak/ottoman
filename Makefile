SHELL := /bin/bash

dep: dep-metalinter
dep-metalinter:
	@go get gopkg.in/alecthomas/gometalinter.v1
	@gometalinter.v1 --install > /dev/null
check: check-metalinter
check-metalinter: dep-metalinter
	@gometalinter.v1 --fast --config=scripts/metalinter.json ./...
check-all: check-all-metalinter
check-all-metalinter: dep-metalinter
	@gometalinter.v1 --config=scripts/metalinter.json ./...
redis-cluster-start:
	./scripts/redis-cluster.sh start
redis-cluster-stop:
	./scripts/redis-cluster.sh stop
test:
	go test -race -v $$(go list ./... | grep -Ev "vendor|qtest")
coverage:
	@./scripts/coverage.sh coverage.out
cover:
	go tool cover -html=coverage.out
