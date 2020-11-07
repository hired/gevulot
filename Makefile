# Result binary path
OUT := build/gevulot

# List of flags to pass to `go test` command
TEST_FLAGS := -timeout 5s -race

# List of LDFLAGS to pass to `go build` command
LDFLAGS :=

# If VERBOSE env var is set add -v flag to the TEST_FLAGS
ifdef VERBOSE
	TEST_FLAGS += -v
endif

# Test database name
PGDATABASE ?= gevulot_test

default: build

prepare-test-db:
	psql $(PGDATABASE) -f scripts/gevulot_test_schema.sql

test:
	go test $(TEST_FLAGS) ./pkg/...

build: clean
	go build -ldflags "$(LDFLAGS)" -o $(OUT) .

# Lazily get build information
build: VERSION ?= $(shell git describe --tags --exact-match 2> /dev/null || git symbolic-ref -q --short HEAD)
build: COMMIT_HASH ?= $(shell git rev-parse --short HEAD 2> /dev/null)
build: BUILD_DATE ?= $(shell date +%FT%T%z)

# Set LDFLAGS for build action
build: LDFLAGS += -X github.com/hired/gevulot/pkg/cli.version=$(VERSION) -X github.com/hired/gevulot/pkg/cli.commitHash=$(COMMIT_HASH) -X github.com/hired/gevulot/pkg/cli.buildDate=$(BUILD_DATE)

lint:
	golangci-lint run ./pkg/...

run: build
	$(OUT) --verbose

clean:
	rm -f $(OUT)
