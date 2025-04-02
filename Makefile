build:
	./set_version.sh
	go build ./cmd/pgtester

debug:
	go build -gcflags "all=-N -l" ./cmd/pgtester
	~/go/bin/dlv --headless --listen=:2345 --api-version=2 --accept-multiclient exec ./pgtester ./testdata

run:
	./pgtester testdata/*

fmt:
	gofmt -w .

test: sec lint

sec:
	gosec ./...
lint:
	golangci-lint run

.PHONY: install-go-test-coverage
install-go-test-coverage:
	go install github.com/vladopajic/go-test-coverage/v2@latest

.PHONY: check-coverage
check-coverage: install-go-test-coverage
	go test $$(go list ./... | grep -v /e2e) -coverprofile=./cover.out -covermode=atomic -coverpkg=./...
	${GOBIN}/go-test-coverage --config=./.testcoverage.yml
