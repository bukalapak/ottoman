SHELL := /bin/bash

dep: dep-metalinter
dep-metalinter:
	@go get gopkg.in/alecthomas/gometalinter.v1
	@gometalinter.v1 --install > /dev/null
check: check-metalinter
check-metalinter: dep-metalinter
	@(! gometalinter.v1 --fast --config=scripts/metalinter.json ./... | read) || exit 1
status: status-metalinter
status-metalinter: dep-metalinter
	@gometalinter.v1 --fast --config=scripts/metalinter.json ./...
report: report-metalinter
report-metalinter: dep-metalinter
	@gometalinter.v1 --config=scripts/metalinter.json ./...
test:
	go test -race -v $$(go list ./... | grep -Ev "vendor|qtest")
coverage:
	@./scripts/coverage.sh coverage.out
cover:
	go tool cover -html=coverage.out
