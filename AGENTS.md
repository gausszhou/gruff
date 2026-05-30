# AGENTS.md — gruff

## Commands

- Test all: `go test ./...`
- Single test: `go test -run TestRender_Heading/h1 ./...`
- Lint: `go vet ./...`
- Build CLI: `make build` (outputs to `dist/gruff`) or `go build -o bin/gruff ./cmd/gruff/`
- Release (multi-platform): `make release`
- Benchmark: `go test -bench=. -benchmem ./benchmark/`
- **NEVER commit** unless explicitly asked. Do NOT stage or commit any changes without user instruction.

## Architecture

- `gruff.go` — public API (`Render`, `RenderBytes`, `WithDark`, `WithLight`, `WithWordWrap`)
- `gruff_renderer.go` — goldmark AST walker → ANSI output; `renderListItem` uses `renderChildren` (no manual child loop, no trailing `\n`)
- `gruff_ansi.go` — SGR codes, `Style`, `Color`, `Theme` types, helpers (`displayWidth`, `stripANSI`)
- `cmd/gruff/main.go` — CLI binary (reads file path arg, auto-detects terminal width)
- `Makefile` — `build` / `release` / `clean` targets (`dist/` output)
- `.github/workflows/release.yml` — tag `v*` triggers `make release` + uploads archives
- Benchmark input: `testdata/benchmark.md`

## Important

- Do NOT use `\x1b[0m` (full reset) in inline styles — use specific undo codes (`\x1b[22m` noBold, `\x1b[39m` default fg, `\x1b[49m` default bg) to preserve outer style state during nesting
- goldmark `Emphasis{Level: 2}` = bold; `Emphasis{Level: 1}` = italic; `***both***` nests Level 1 wrapping Level 2. `ast.Strong` does NOT exist
- `name` and `Theme` struct fields in `ansi.go` must be kept in sync
- `displayWidth` uses `go-runewidth` (not hand-rolled ranges); handles U+FE0F by forcing preceding char to width 2
- `wrapText` fills each line to `width` with bg color (no `\x1b[K`); hard `\n` preserved
- `defaultWordWrap` = 80 in `gruff.go`; tables inherit this for column width capping
- `Document.Padding` = 1 in both dark/light themes; `wrapText` respects padding
- `*ast.TextBlock` case exists in `renderNode` (tight list items); writes `\n` after children
- `isBlockLevel` helper used by Document to insert `\n` between block siblings
- `ansi.go` constants include `ansiDefaultBg` (`\x1b[49m`); used at end of `Render` output
- `strings.TrimSpace` applied to raw output before `wrapText`
- Put temporary/verification files in `tmp/` (gitignored)

## Dependencies

- **Core library runtime:** `goldmark` (parser), `go-runewidth` (Unicode width)
- **CLI only:** `golang.org/x/term` (terminal size detection)
- **Examples only:** glamour, bubbletea, bubbles, lipgloss, etc. — not part of library builds
