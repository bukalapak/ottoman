dep: dep-metalinter
dep-metalinter:
	@go get gopkg.in/alecthomas/gometalinter.v1
	@gometalinter.v1 --install > /dev/null
check: check-metalinter
check-metalinter: dep-metalinter
	@(! gometalinter.v1 --fast --config=scripts/metalinter.json ./... | read) || exit 1
test:
	go test -v $$(go list ./... | grep -v /vendor/)
coverage:
	@./scripts/coverage.sh coverage.out
cover:
	go tool cover -html=coverage.out
