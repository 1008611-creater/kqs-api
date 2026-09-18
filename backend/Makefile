.PHONY: build generate test test-unit test-integration test-e2e

VERSION ?= $(shell tr -d '\r\n' < ./cmd/server/VERSION)
LDFLAGS ?= -s -w -X main.Version=$(VERSION)

build:
	CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" -trimpath -o bin/server ./cmd/server

generate:
	go generate ./ent
	go generate ./cmd/server

test:
	go test ./...
	golangci-lint run ./...

test-unit:
	bash -o pipefail -c 'go test -json -tags=unit ./... -count=1 2>&1 | tee unit-test.json; rc=$${PIPESTATUS[0]}; echo "--- unit test failure summary ---"; python3 -c "import json,sys; [print(f\"{e.get(\\\"Package\\\",\\\"\\\")}: {e.get(\\\"Test\\\",\\\"package\\\")}\\n  {e.get(\\\"Output\\\",\\\"\\\").rstrip()}\") for line in open(sys.argv[1],encoding=\"utf-8\",errors=\"replace\") if (lambda e: e.get(\"Action\") == \"fail\")(json.loads(line))]" unit-test.json || true; exit $${rc}'

test-integration:
	go test -tags=integration ./...

test-e2e:
	./scripts/e2e-test.sh

test-e2e-local:
	go test -tags=e2e -v -timeout=300s ./internal/integration/...
