.PHONY: pre-commit build

all: clean generate build

generate:
	go generate ./src/...

build:
	mkdir -p build && cd build && go build ../src/...

goreport:
	goreportcard-cli -v -t 100.0

test:
	go test -v ./src/...

clean:
	rm -rf build

pre-commit: generate build goreport test
