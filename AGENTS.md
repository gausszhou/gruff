# AGENTS.md — gruff

## Commands

- Test all: `go test ./...`
- Benchmark all: `go test -bench=. -benchmem ./...`
- Benchmark gruff vs glamour: `go test -bench=BenchmarkGruff,BenchmarkGlamour -benchmem ./benchmark/`
- Lint: `go vet ./...`
- Single test: `go test -run TestRender_Heading/h1 ./...`

## Architecture

- `gruff.go` — public API (`Render`, `RenderBytes`, `WithDark`, `WithLight`, `WithWordWrap`)
- `renderer.go` — goldmark AST walker → ANSI output
- `ansi.go` — SGR codes, `style` struct, `theme` (dark/light), helpers (`displayWidth`, `stripANSI`)
- `gruff_test.go` — table-driven tests per syntax element
- `benchmark/` — embedded markdown benchmark comparing gruff vs glamour

## Important

- Do NOT use `\x1b[0m` (full reset) in inline styles — use specific undo codes (`\x1b[22m` noBold, `\x1b[39m` default fg, `\x1b[49m` default bg) to preserve outer style state during nesting
- goldmark `Emphasis{Level: 2}` = bold; `Emphasis{Level: 1}` = italic; `***both***` nests Level 1 wrapping Level 2
- `ast.Strong` does NOT exist in goldmark — use `ast.Emphasis` with Level check
- `name` and `Theme` struct fields in `ansi.go` must be updated together
- `benchmark/testdata/sample.md` must mirror `testdata/sample.md`
- `testdata/sample.md` is the single source of truth for example content; `examples/basic/main.go` and `benchmark/benchmark_test.go` read it at runtime via relative path
- Do NOT create `cmd/` directory for diagnostics — use inline test or remove after use

## Constraints

- `goldmark` is the only runtime dependency (parser)
- `glamour` is a benchmark-only dependency (not in production build)
