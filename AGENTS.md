# AGENTS.md — gruff

## Commands

- Test all: `go test ./...`
- Benchmark all: `go test -bench=. -benchmem ./...`
- Benchmark gruff vs glamour: `go test -bench=BenchmarkGruff,BenchmarkGlamour -benchmem ./benchmark/`
- Lint: `go vet ./...`
- Single test: `go test -run TestRender_Heading/h1 ./...`
- Build all: `go build -o bin/ ./...` (all build artifacts go to `bin/`)

## Architecture

- `gruff.go` — public API (`Render`, `RenderBytes`, `WithDark`, `WithLight`, `WithWordWrap`)
- `renderer.go` — goldmark AST walker → ANSI output
- `ansi.go` — SGR codes, `Style`, `Color`, `Theme` types, helpers (`displayWidth`, `stripANSI`)
- `gruff_test.go` — table-driven tests per syntax element
- `benchmark/` — benchmark comparing gruff vs glamour

## Table rendering

Two-pass: collect all cell content + calculate column widths (Pass 1), then render with UTF-8 box‑drawing column separators (Pass 2). No outer borders. Columns auto-expand to content, capped per column at `(wordWrap - overhead) / numCols`. Word wrap inside cells (ANSI-aware, at space boundaries). Alignment per cell via `TableCell.Alignment`.

## Important

- Do NOT use `\x1b[0m` (full reset) in inline styles — use specific undo codes (`\x1b[22m` noBold, `\x1b[39m` default fg, `\x1b[49m` default bg) to preserve outer style state during nesting
- goldmark `Emphasis{Level: 2}` = bold; `Emphasis{Level: 1}` = italic; `***both***` nests Level 1 wrapping Level 2
- `ast.Strong` does NOT exist in goldmark — use `ast.Emphasis` with Level check
- `name` and `Theme` struct fields in `ansi.go` must be updated together
- `testdata/_data.md` is the single benchmark input file
- `WithWordWrap` sets document width (default 120); tables inherit this for column width capping
- `displayWidth` uses `go-runewidth` (not hand-rolled ranges) for correct Unicode width
- Do NOT create `cmd/` directory for diagnostics — use inline test or remove after use
- Put temporary/verification files in `tmp/` directory (already gitignored)

## Constraints

- `goldmark` is the only runtime dependency (parser)
- `glamour` is a benchmark-only dependency (not in production build)
- `go-runewidth` is a runtime dependency (Unicode display width)
