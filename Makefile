OUT := build/gevulot

# List of flags to pass to `go test` command
TEST_FLAGS :=

# If VERBOSE env var is set add -v flag to the TEST_FLAGS
ifdef VERBOSE
	TEST_FLAGS += -v
endif

default: build

test:
	go test $(TEST_FLAGS) ./pkg/...

build: clean
	go build -o $(OUT) ./pkg/cmd/gevulot

lint:
	golint -set_exit_status=1 ./pkg/...

run: build
	$(OUT)

clean:
	rm -f $(OUT)
