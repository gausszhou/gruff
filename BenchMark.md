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
