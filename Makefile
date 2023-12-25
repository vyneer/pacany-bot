.PHONY: FORCE

DIST_FOLDER = ./dist
BINARY = pacany-bot
GOPATH = ${shell go env GOPATH}
BIN_PATH = ${GOPATH}/bin
GOLANGCI_LINT = ${BIN_PATH}/golangci-lint

GOFLAGS := -tags netgo -ldflags -extldflags="-static"

# Build

build: ${DIST_FOLDER}/${BINARY}
.PHONY: build

run: build
	${DIST_FOLDER}/${BINARY}
.PHONY: run

clean:
	go clean
	rm -rf ${DIST_FOLDER}
.PHONY: clean

release: clean
	goreleaser release
.PHONY: release

# Test

test:
	go test
.PHONY: test

# Lint

lint: ${GOLANGCI_LINT}
	${GOLANGCI_LINT} run
.PHONY: lint

tidy: 
	go mod tidy
.PHONY: tidy

fix: ${GOLANGCI_LINT} tidy
	${GOLANGCI_LINT} run --fix
.PHONY: fix

# Setup development environment

download:
	go mod download
.PHONY: download

verify:
	go mod verify
.PHONY: verify

setup-linter: ${GOLANGCI_LINT}
.PHONY: setup-linter

setup-hooks: .git/hooks/pre-commit
.PHONY: setup-hooks

setup: download verify setup-linter setup-hooks
.PHONY: setup

# Non-PHONY targets

${DIST_FOLDER}/${BINARY}: FORCE
	CGO_ENABLED=0 go build ${GOFLAGS} -v -o $@

${GOLANGCI_LINT}:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ${BIN_PATH}

.git/hooks/pre-commit:
	cp .githooks/pre-commit .git/hooks/pre-commit
	chmod 755 .git/hooks/pre-commit