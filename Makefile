.PHONY: pre-commit build

clear:
	rm -rf build

all: clean generate build

generate:
	go generate ./src/...

fmt-staged:
	./fmt-staged.sh

build-dir:
	mkdir -p build

build: build-dir
	cd build && go build ../src/...

goreport:
	goreportcard-cli -v -t 100.0

# TODO: run this in Docker
ruleguard:
	ruleguard -c=3 -rules=rules.go -fix ./...

test:
	go test -v ./src/...

clean:
	rm -rf build

pre-commit: generate fmt-staged build ruleguard goreport test

coverage.out: build-dir
	go test ./src/... -coverprofile=build/coverage.out

coverage.html: coverage.out
	go tool cover -html=build/coverage.out -o build/coverage.html

coverage: coverage.html
	open build/coverage.html
