.PHONY: build test clean build-all build-linux build-darwin build-windows lint fmt vet ci install-tools test-race vulncheck osvscanner nilaway

BINARY_NAME=gruff
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

test-race:
	go test -race ./...

clean:
	rm -rf $(BIN_DIR)
	rm -rf $(DIST_DIR)
	rm -f coverage.out

GOPATH_BIN := $(shell go env GOPATH)/bin

lint:
	@test -x "$(GOPATH_BIN)/golangci-lint" || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GOPATH_BIN)/golangci-lint run ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install github.com/google/osv-scanner/cmd/osv-scanner@latest
	go install go.uber.org/nilaway/cmd/nilaway@latest

vulncheck:
	@test -x "$(GOPATH_BIN)/govulncheck" || go install golang.org/x/vuln/cmd/govulncheck@latest
	$(GOPATH_BIN)/govulncheck ./...

osvscanner:
	@test -x "$(GOPATH_BIN)/osv-scanner" || go install github.com/google/osv-scanner/cmd/osv-scanner@latest
	$(GOPATH_BIN)/osv-scanner -r .

nilaway:
	@test -x "$(GOPATH_BIN)/nilaway" || go install go.uber.org/nilaway/cmd/nilaway@latest
	$(GOPATH_BIN)/nilaway -include-pkgs "github.com/gausszhou/gruff/..." ./...

ci: vet test-race lint vulncheck osvscanner nilaway
