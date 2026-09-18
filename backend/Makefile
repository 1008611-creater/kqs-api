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
	bash -o pipefail -c 'set +e; TLSFINGERPRINT_NETWORK_TESTS=0 go test -short -v -tags=unit ./... -count=1 2>&1 | tee unit-test.log; rc=$${PIPESTATUS[0]}; echo "--- unit test exit code: $${rc} ---"; if [ "$${rc}" -ne 0 ]; then echo "--- unit test failure lines ---"; grep -E "^(--- FAIL:|FAIL[[:space:]]|panic:|.*Error:|.*undefined:|.*cannot |.*no required module)" unit-test.log || true; fi; exit $${rc}'

test-integration:
	go test -tags=integration ./...

test-e2e:
	./scripts/e2e-test.sh

test-e2e-local:
	go test -tags=e2e -v -timeout=300s ./internal/integration/...
