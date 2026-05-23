.PHONY: all build test lint bench clean

all: build lint test

build:
	CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o bin/ ./...

test:
	go test ./...

lint:
	go vet ./...

bench:
	go test -bench=. -benchmem ./...

bench-compare:
	go test -bench=BenchmarkGruff,BenchmarkGlamour -benchmem ./benchmark/

clean:
	rm -rf bin/
