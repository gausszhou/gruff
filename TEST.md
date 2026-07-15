# Test & Benchmark

## All Tests

```bash
go test -timeout 120s ./...
```

## Benchmarks

```bash
# All benchmarks
go test -run '^$' -bench='BenchmarkGruff|BenchmarkGlamour' -benchmem -timeout 120s ./benchmark/
```

## Build Examples

```sh
go build -o bin/examples/api               ./examples/api/
go build -o bin/examples/basic             ./examples/basic/
go build -o bin/examples/codeblock         ./examples/codeblock/
go build -o bin/examples/compare-benchmark ./examples/compare-benchmark/
go build -o bin/examples/compare-glamour   ./examples/compare-glamour/
go build -o bin/examples/compare-theme     ./examples/compare-theme/
go build -o bin/examples/custom-theme      ./examples/custom-theme/
go build -o bin/examples/table             ./examples/table/
```
