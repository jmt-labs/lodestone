SHELL := /usr/bin/env bash
.SHELLFLAGS := -eu -o pipefail -c

BINARY := lodestone
BIN_DIR := bin
PKG := ./...
LODESTONE_BIN ?= $(CURDIR)/$(BIN_DIR)/$(BINARY)

.PHONY: build test lint vuln e2e clean release-check release-dry tidy help

help:
	@echo "Targets:"
	@echo "  build        - Baut cmd/lodestone nach $(BIN_DIR)/$(BINARY)"
	@echo "  test         - go test $(PKG)"
	@echo "  lint         - golangci-lint run"
	@echo "  vuln         - govulncheck $(PKG)"
	@echo "  e2e          - Phase-1 E2E-Smoke-Test (ab T9)"
	@echo "  release-check- goreleaser check (Config validieren)"
	@echo "  release-dry  - goreleaser release --snapshot --clean (dry-run)"
	@echo "  tidy         - go mod tidy"
	@echo "  clean        - Entfernt $(BIN_DIR)/ und dist/"

build:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(BINARY) ./cmd/lodestone

test:
	go test $(PKG)

lint:
	golangci-lint run

vuln:
	govulncheck $(PKG)

e2e: build
	@if [ -x e2e/lodestone_test.sh ]; then \
		LODESTONE_BIN=$(LODESTONE_BIN) bash e2e/lodestone_test.sh; \
	else \
		echo "e2e/lodestone_test.sh nicht vorhanden (kommt in T9)"; \
	fi

release-check:
	goreleaser check

release-dry:
	goreleaser release --snapshot --clean

tidy:
	go mod tidy

clean:
	rm -rf $(BIN_DIR) dist
