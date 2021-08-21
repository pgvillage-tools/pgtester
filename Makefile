build:
	./set_version.sh
	go build ./cmd/pgtester

debug:
	~/go/bin/dlv debug --headless --listen=:2345 --api-version=2 --accept-multiclient ./cmd/pgtester

run:
	./pgtester -f tests.yaml

fmt:
	gofmt -w .

test: sec lint

sec:
	gosec ./...
lint:
	golangci-lint run
