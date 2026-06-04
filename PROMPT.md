# Go CLI 项目模板提示词

创建一个名为 `<project-name>` 的 Go CLI 项目，遵循以下结构和约定。

## 项目结构

```
<project-name>/
├── .github/
│   └── workflows/
│       ├── ci.yml              # CI: build + vet + test
│       └── release.yml         # 打 tag v* 时自动构建多平台二进制 + GitHub Release
├── cmd/
│   └── root.go                 # cobra 根命令
├── <pkg>/                      # 核心库包
├── testdata/                   # 测试数据
├── tmp/                        # 临时文件（gitignored）
├── .gitignore
├── AGENTS.md                   # AI 助手指南
├── go.mod
├── LICENSE
├── main.go                     # CLI 入口: cmd.Execute()
└── Makefile
```

## Makefile

```makefile
.PHONY: build test clean build-all build-linux build-darwin build-windows lint fmt vet

BINARY_NAME=<name>
BIN_DIR=bin
DIST_DIR=dist
CMD_DIR=.

build: build-all

build-all: build-linux build-darwin build-windows

build-linux:
	mkdir -p $(BIN_DIR) $(DIST_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(BIN_DIR)/$(BINARY_NAME)-linux-amd64 $(CMD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o $(BIN_DIR)/$(BINARY_NAME)-linux-arm64 $(CMD_DIR)
	tar czf $(DIST_DIR)/$(BINARY_NAME)-linux-amd64.tar.gz -C $(BIN_DIR) $(BINARY_NAME)-linux-amd64
	tar czf $(DIST_DIR)/$(BINARY_NAME)-linux-arm64.tar.gz -C $(BIN_DIR) $(BINARY_NAME)-linux-arm64

build-darwin:
	mkdir -p $(BIN_DIR) $(DIST_DIR)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o $(BIN_DIR)/$(BINARY_NAME)-darwin-amd64 $(CMD_DIR)
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o $(BIN_DIR)/$(BINARY_NAME)-darwin-arm64 $(CMD_DIR)
	tar czf $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64.tar.gz -C $(BIN_DIR) $(BINARY_NAME)-darwin-amd64
	tar czf $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64.tar.gz -C $(BIN_DIR) $(BINARY_NAME)-darwin-arm64

build-windows:
	mkdir -p $(BIN_DIR) $(DIST_DIR)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o $(BIN_DIR)/$(BINARY_NAME)-windows-amd64.exe $(CMD_DIR)
	CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -ldflags="-s -w" -o $(BIN_DIR)/$(BINARY_NAME)-windows-arm64.exe $(CMD_DIR)
	tar czf $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.tar.gz -C $(BIN_DIR) $(BINARY_NAME)-windows-amd64.exe
	tar czf $(DIST_DIR)/$(BINARY_NAME)-windows-arm64.tar.gz -C $(BIN_DIR) $(BINARY_NAME)-windows-arm64.exe

test:
	go test ./...

clean:
	rm -rf $(BIN_DIR) $(DIST_DIR)

vet:
	go vet ./...

fmt:
	go fmt ./...
```

## GitHub Actions

### `.github/workflows/ci.yml`

```yaml
name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'
      - run: go build ./...
      - run: go vet ./...
      - run: go test ./...
```

### `.github/workflows/release.yml`

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

permissions:
  contents: write

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'
      - run: make build-all
      - uses: softprops/action-gh-release@v2
        with:
          files: dist/*.tar.gz
          generate_release_notes: true
```

## .gitignore

```
*.exe*
*.test
*.out
coverage.*
*.coverprofile
go.work
go.work.sum
.env
bin
dist
tmp
```

## AGENTS.md 模板

```markdown
# AGENTS.md — <project-name>

## 命令

- 测试全部: `go test ./...`
- 测试单个: `go test -run <TestName> ./...`
- 静态检查: `go vet ./...`
- 构建 CLI: `make build` 或 `go build -o bin/<name> .`
- 多平台发布: `make build-all`
- **禁止自动提交**，除非用户明确要求。

## 架构

- `<pkg>/` — 核心库包，包含公开 API
- `main.go` — CLI 入口，调用 `cmd.Execute()`
- `cmd/root.go` — cobra 根命令
- `Makefile` — `build` / `test` / `clean` / `vet` / `fmt`
- <项目特有约定>

## 依赖

- <列出核心依赖>

## 重要约定

- <项目特有说明>
```

## 开始构建

1. 初始化模块: `go mod init <module-path>`
2. 创建上述目录结构和文件
3. 安装依赖: `go mod tidy`
4. 验证: `make build && make vet && make test`
