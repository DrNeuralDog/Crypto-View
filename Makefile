POWERSHELL ?= powershell

.PHONY: build build-all run clean test

test:
	$(POWERSHELL) -NoProfile -ExecutionPolicy Bypass -File ./test.ps1

build:
	go build -o bin/cryptoview ./cmd/cryptoview

build-all:
	$(POWERSHELL) -NoProfile -ExecutionPolicy Bypass -File ./build_all_os.ps1

run:
	go run ./cmd/cryptoview

clean:
	go clean ./...
