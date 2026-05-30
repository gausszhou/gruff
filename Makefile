.PHONY: build release clean

BINARY ?= gruff
DIST ?= dist

build:
	@mkdir -p $(DIST)
	go build -o $(DIST)/$(BINARY) ./cmd/gruff/

release:
	@mkdir -p $(DIST)
	for os in linux darwin windows; do \
		for arch in amd64 arm64; do \
			ext=""; \
			[ "$$os" = "windows" ] && ext=".exe"; \
			GOOS=$$os GOARCH=$$arch go build -o "$(DIST)/$(BINARY)_$${os}_$${arch}$${ext}" ./cmd/gruff/; \
			file="$(DIST)/$(BINARY)_$${os}_$${arch}$${ext}"; \
			if [ "$$os" = "windows" ]; then \
				cd $(DIST) && zip "$(BINARY)_$${os}_$${arch}.zip" "$(BINARY)_$${os}_$${arch}$${ext}" && cd ..; \
			else \
				cd $(DIST) && tar czf "$(BINARY)_$${os}_$${arch}.tar.gz" "$(BINARY)_$${os}_$${arch}$${ext}" && cd ..; \
			fi; \
		done; \
	done

clean:
	rm -rf $(DIST)
