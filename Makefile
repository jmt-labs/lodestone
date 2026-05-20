SHELL := /usr/bin/env bash
.SHELLFLAGS := -eu -o pipefail -c

BINARY := lodestone
BIN_DIR := bin
PKG := ./...
LODESTONE_BIN ?= $(CURDIR)/$(BIN_DIR)/$(BINARY)

.PHONY: build test lint vuln e2e clean release-check release-dry tidy help \
	docs docs-status-check docs-cmd-coverage skills-coverage docs-links

help:
	@echo "Targets:"
	@echo "  build              - Baut cmd/lodestone nach $(BIN_DIR)/$(BINARY)"
	@echo "  test               - go test $(PKG)"
	@echo "  lint               - golangci-lint run"
	@echo "  vuln               - govulncheck $(PKG)"
	@echo "  e2e                - Phase-1 E2E-Smoke-Test (ab T9)"
	@echo "  docs               - alle Doku-Coverage-Checks (status, cmd, skills, links)"
	@echo "  docs-status-check  - Phasen-Status konsistent in README, roadmap, CHANGELOG"
	@echo "  docs-cmd-coverage  - jedes Cobra-Verb hat eine docs/user/commands/<verb>.md"
	@echo "  skills-coverage    - jeder flavors-Skill ist in docs/user/skills.md verlinkt + Embed-Mirror byte-identisch"
	@echo "  docs-links         - relative Markdown-Links innerhalb docs/ existieren"
	@echo "  release-check      - goreleaser check (Config validieren)"
	@echo "  release-dry        - goreleaser release --snapshot --clean (dry-run)"
	@echo "  tidy               - go mod tidy"
	@echo "  clean              - Entfernt $(BIN_DIR)/ und dist/"

build:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(BINARY) ./cmd/lodestone
	go build -o $(BIN_DIR)/$(BINARY)-mcp ./cmd/lodestone-mcp

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

docs: docs-status-check docs-cmd-coverage skills-coverage docs-links

docs-status-check:
	@missing=0; \
	for f in README.md docs/internals/roadmap.md CHANGELOG.md; do \
		if ! LC_ALL=C.UTF-8 grep -Fq 'Phasen 1–4' $$f; then \
			echo "::error::$$f enthält nicht das aktuelle Phasen-Statement (Phasen 1–4)"; \
			missing=1; \
		fi; \
	done; \
	exit $$missing

docs-cmd-coverage:
	@missing=0; \
	for verb in $$(grep -hE '^\s+Use:\s+"[a-z]+' cmd/lodestone/*.go | sed -E 's/.*Use:\s+"([a-z]+).*/\1/' | sort -u); do \
		case "$$verb" in \
			lodestone) continue ;; \
		esac; \
		f="docs/user/commands/$$verb.md"; \
		if [ ! -f "$$f" ]; then \
			echo "::error::Verb '$$verb' hat keine Doku unter $$f"; \
			missing=1; \
		fi; \
	done; \
	exit $$missing

skills-coverage:
	@missing=0; \
	for skill in flavors/lodestone/skills/*.md; do \
		name=$$(basename $$skill .md); \
		if ! grep -q "$$name" docs/user/skills.md; then \
			echo "::error::Skill '$$name' fehlt in docs/user/skills.md"; \
			missing=1; \
		fi; \
		mirror="internal/lodestone/skills/data/$$(basename $$skill)"; \
		if [ ! -f "$$mirror" ]; then \
			echo "::error::Embed-Mirror fehlt: $$mirror"; \
			missing=1; \
		elif ! cmp -s "$$skill" "$$mirror"; then \
			echo "::error::Embed-Mirror $$mirror weicht von $$skill ab (sync nötig)"; \
			missing=1; \
		fi; \
	done; \
	exit $$missing

docs-links:
	@bash scripts/check-docs-links.sh
