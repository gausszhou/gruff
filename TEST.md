# Benchmark

## Commands

```bash
# All benchmarks (gruff package)
go test -run '^$' -bench=BenchmarkRender -benchmem -timeout 120s .

# All benchmarks (benchmark package — includes large/medium/100×10k)
go test -run '^$' -bench=BenchmarkRender -benchmem -timeout 120s ./benchmark/

# Gruff vs Glamour
go test -run '^$' -bench=BenchmarkGruff,BenchmarkGlamour -benchmem -timeout 120s ./benchmark/

# Single large markdown test (5M chars)
go test -run TestRenderMarkdown5MChars -v -timeout 120s ./benchmark/

# All tests
go test -timeout 120s ./...
```

## Build Commands (Windows)

```pwsh
go build -o bin/examples/api.exe               ./examples/api/
go build -o bin/examples/basic.exe             ./examples/basic/
go build -o bin/examples/codeblock.exe         ./examples/codeblock/
go build -o bin/examples/compare-benchmark.exe ./examples/compare-benchmark/
go build -o bin/examples/compare-glamour.exe   ./examples/compare-glamour/
go build -o bin/examples/compare-theme.exe     ./examples/compare-theme/
go build -o bin/examples/custom-theme.exe      ./examples/custom-theme/
go build -o bin/examples/table.exe             ./examples/table/
```

## Build Commands (Unix)

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
