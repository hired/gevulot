OUT := build/gevulot

default: build

test:
	go test ./src/...

build: clean
	go build -o $(OUT) ./src

lint:
	golint -set_exit_status=1 ./...

run: build
	$(OUT)

clean:
	rm -f $(OUT)
