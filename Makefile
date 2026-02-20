.PHONY: build build-all run clean

build:
	mkdir -p bin
	go build -o bin/cryptoview ./cmd/cryptoview

build-all:
	mkdir -p bin
	GOOS=windows GOARCH=amd64 go build -o bin/cryptoview-windows-amd64.exe ./cmd/cryptoview
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags ci -o bin/cryptoview-linux-amd64 ./cmd/cryptoview
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o bin/cryptoview-darwin-amd64 ./cmd/cryptoview

run:
	go run ./cmd/cryptoview

clean:
	rm -rf bin
