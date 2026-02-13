.PHONY: build run clean

build:
	mkdir -p bin
	go build -o bin/cryptoview ./cmd/cryptoview

run:
	go run ./cmd/cryptoview

clean:
	rm -rf bin
