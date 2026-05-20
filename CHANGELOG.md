# Changelog

Alle nennenswerten Änderungen an diesem Projekt werden in dieser Datei
dokumentiert. Format gemäß [Keep a Changelog](https://keepachangelog.com/de/1.1.0/),
Versionierung nach [Semantic Versioning](https://semver.org/lang/de/).

## [Unreleased]

### Added
- Initiales Repo-Skelett: Go-Modul `github.com/jmt-labs/lodestone`,
  Cobra-CLI-Stub `cmd/lodestone`, CI-/Release-Workflows, GoReleaser-
  Konfiguration, Makefile, Linter-Konfiguration.
- Doku-Grundgerüst: `README.md`, `CONTRIBUTING.md`, `CLAUDE.md`,
  `AGENTS.md`, `base/models.yaml`, `docs/lodestone.md`.
- Spec + Phase-1-Plan unter `docs/superpowers/`.
- **T1** — Schemas für `Signal`, `Fingerprint`, `Recommendation`,
  `WorkPackage` (`internal/lodestone/schema/`) inkl. JSON-Roundtrip-
  Tests.
- **T2** — Datei-basierter Store (`internal/lodestone/store/`) mit
  JSONL-Signals, JSON-Fingerprint, JSONL-Recommendations und atomarem
  Replace.
