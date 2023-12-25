.PHONY: FORCE

TARGET_FOLDER = ./target
BINARY = pacani-bot
GOPATH = ${shell go env GOPATH}
BIN_PATH = ${GOPATH}/bin
GOLANGCI_LINT = ${BIN_PATH}/golangci-lint

GOFLAGS := -tags netgo

DEBUG ?= 1
ifeq ($(DEBUG), 1)
	LDFLAGS := '-extldflags="-static"'
else
	GOFLAGS += -trimpath
	LDFLAGS := '-s -w -extldflags="-static"'
endif

GOFLAGS += -ldflags ${LDFLAGS}

# Build

build: ${TARGET_FOLDER}/${BINARY}
.PHONY: build

run: build
	${TARGET_FOLDER}/${BINARY}
.PHONY: run

clean:
	go clean
	rm -rf ${TARGET_FOLDER}
.PHONY: clean

ko:
	ko build --local --base-import-paths
.PHONY: ko

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

${TARGET_FOLDER}/${BINARY}: FORCE
	CGO_ENABLED=0 go build ${GOFLAGS} -v -o $@

${GOLANGCI_LINT}:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ${BIN_PATH}

.git/hooks/pre-commit:
	cp .githooks/pre-commit .git/hooks/pre-commit
	chmod 755 .git/hooks/pre-commit