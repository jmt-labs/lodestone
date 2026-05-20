# Installation

> Noch nicht released. Bis `v0.1.0-alpha` getagged ist
> ([Roadmap](../internals/roadmap.md)), nur `go install` oder lokaler
> Build möglich.

## Go install

Immer verfügbar, baut beide Binaries aus dem aktuellen `main`:

```sh
go install github.com/jmt-labs/lodestone/cmd/lodestone@latest
go install github.com/jmt-labs/lodestone/cmd/lodestone-mcp@latest
```

`$GOBIN` (bzw. `$GOPATH/bin`) muss im `$PATH` sein. Verifikation:

```sh
lodestone version
lodestone-mcp --version
```

## Lokal bauen

```sh
git clone https://github.com/jmt-labs/lodestone.git
cd lodestone
make build      # legt bin/lodestone und bin/lodestone-mcp ab
```

Voraussetzung: Go ≥ 1.24.

## Geplant ab `v0.1.0`

| Distributor | Status |
|---|---|
| GoReleaser-Archives (linux/darwin/windows × amd64/arm64) | aktiv nach erstem Tag |
| Homebrew-Tap `jmt-labs/tap/lodestone` | konfiguriert, ausgeschaltet bis Tap-Repo bereitsteht |
| apt/deb via nfpms | konfiguriert, ausgeschaltet bis Release-Token konfiguriert ist |

Details: [Release-Prozess](../contributor/release-process.md).

## Verwandt

- [Quickstart](quickstart.md) — die ersten Schritte nach der Installation.
- [`init`-Befehl](commands/init.md) — Projekt bootstrappen.
