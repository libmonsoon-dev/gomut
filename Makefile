.PHONY: pre-commit build

all: clean generate build

generate:
	go generate ./src/...

fmt-staged:
	./fmt-staged.sh

build:
	mkdir -p build && cd build && go build ../src/...

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
