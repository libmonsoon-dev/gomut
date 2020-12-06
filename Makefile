.PHONY: pre-commit build

clear:
	rm -rf build

all: clean generate build

generate:
	go generate ./src/...

build-dir:
	mkdir -p build

typecheck: build-dir
	go build ./src/...

goreport:
	goreportcard-cli -v -t 100.0

ruleguard:
	ruleguard -c=3 -rules=tools/ruleguard/rules.go -fix ./...

test:
	go test ./src/...

clean:
	rm -rf build

pre-commit: generate mod-tidy typecheck ruleguard goreport test

coverage.out: build-dir
	go test ./src/... -coverprofile=build/coverage.out

coverage.html: coverage.out
	go tool cover -html=build/coverage.out -o build/coverage.html

coverage: coverage.html
	open build/coverage.html

mod-tidy:
	go mod tidy -v

install-tools:
	go list -f '{{range .Imports}}{{.}} {{end}}' tools/tools.go | xargs go get -v
