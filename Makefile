.PHONY: all build test lint bench clean

all: build lint test

build:
	go build -o bin/ ./...

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
