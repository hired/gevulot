OUT := build/gevulot

default: build

test:
	go test ./src/...

build: clean
	go build -o $(OUT) ./src

run: build
	$(OUT)

clean:
	rm -f $(OUT)
