BINARY = pacany-bot
OUTPUT = ./dist
GOPATH = ${shell go env GOPATH}
BIN_PATH = ${GOPATH}/bin

GOLANGCI_LINT = ${BIN_PATH}/golangci-lint

GOFLAGS := -trimpath -tags netgo -ldflags -extldflags="-static"

.PHONY: FORCE build run clean release prerelease test lint tidy fix download verify setup-linter setup-hooks setup

# Help

help: ## Output this help screen
	@grep -E '(^[a-zA-Z0-9_-]+:.*?##.*$$)|(^##)' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}{printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}' | sed -e 's/\[32m##/[33m/'

# Build

build: ${OUTPUT}/${BINARY} ## Build the binary

run: build ## Build and run the binary
	${OUTPUT}/${BINARY}

clean: ## Clean the project directory and remove the previously built binary
	go clean
	rm -rf ${OUTPUT}

.ONESHELL:
release: clean ## Create a new tag and make a release on GitHub
	@read -p "Enter tag name: " tag_name
	@read -p "Selected tag name is $$tag_name, is this correct? [y/n]: " -n 1 -r confirm
	echo
	if [[ $$confirm =~ ^[Yy]$$ ]]; then \
		git tag -a $$tag_name -m "$$tag_name"; \
		git push origin $$tag_name; \
		goreleaser release; \
	fi

.ONESHELL:
prerelease: clean ## Create a new tag and make a prerelease on GitHub
	@read -p "Enter tag name: " tag_name
	@read -p "Enter prerelease iteration: " count
	tag_name+="-rc.$$count"
	@read -p "Selected tag name is $$tag_name, is this correct? [y/n]: " -n 1 -r confirm
	echo
	if [[ $$confirm =~ ^[Yy]$$ ]]; then \
		git tag -a $$tag_name -m "$$tag_name"; \
		git push origin $$tag_name; \
		goreleaser release -f .goreleaser.prerelease.yaml; \
	fi

# Test

test: ## Run tests
	go test ./...

# Lint

lint: ${GOLANGCI_LINT} ## Lint the project
	${GOLANGCI_LINT} run

tidy: ## Tidy the `go.mod` file
	go mod tidy

fix: ${GOLANGCI_LINT} tidy ## Lint the project and apply the `gofumpt` code style
	${GOLANGCI_LINT} run --fix

# Setup development environment

download: ## Download the project's dependencies
	go mod download

verify: ## Verify the project's dependencies
	go mod verify

setup-linter: ${GOLANGCI_LINT} ## Install `golanci-lint` into `GOPATH` (if it's not installed)

setup-hooks: .git/hooks/pre-commit ## Install the git hooks (it they're not installed)

setup: download verify setup-linter setup-hooks ## Set up the development environment (runs the `download`, `verify`, `setup-linter` and `setup-hooks` recipes)

# Non-PHONY targets

${OUTPUT}/${BINARY}: FORCE
	CGO_ENABLED=0 go build ${GOFLAGS} -v -o $@

${GOLANGCI_LINT}:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ${BIN_PATH}

.git/hooks/pre-commit:
	cp .githooks/pre-commit .git/hooks/pre-commit
	chmod 755 .git/hooks/pre-commit
