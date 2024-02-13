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
