# Release-Prozess

> Stand: 2026-05-20, gilt ab Phase-1-Abschluss.

Lodestone wird tag-getrieben mit [GoReleaser](https://goreleaser.com)
released. Der erste geplante Tag ist `v0.1.0-alpha`.

## Ablauf

1. **Vorprüfung** auf `main`:
   ```sh
   go vet ./...
   go test -race -count=1 ./...
   golangci-lint run
   govulncheck ./...
   make e2e
   ```
   Alle fünf müssen grün sein. Bei Abweichungen: Fix-PR statt
   Release-Tag.

2. **CHANGELOG** auf `main` finalisieren: `## [Unreleased]`-Block in
   `## [vX.Y.Z] – YYYY-MM-DD` umbenennen, neuen leeren
   `## [Unreleased]`-Block einfügen.

3. **Tag setzen und pushen**:
   ```sh
   git tag -a vX.Y.Z -m "lodestone vX.Y.Z"
   git push origin vX.Y.Z
   ```
   Das löst den `release.yml`-Workflow aus, der GoReleaser im CI
   ausführt.

4. **GoReleaser-Lokal-Verifikation** (optional, vor `git push --tags`):
   ```sh
   goreleaser check                                # Config validieren
   goreleaser release --snapshot --clean           # Lokaler Build
   ls dist/                                        # 6 Archives + checksums.txt erwartet
   rm -rf dist/                                    # vor echtem Tag aufräumen
   ```

## Erwartete Artefakte

Pro Release sechs Archive (linux/darwin/windows × amd64/arm64) plus
`checksums.txt`:

```
lodestone_<version>_linux_amd64.tar.gz
lodestone_<version>_linux_arm64.tar.gz
lodestone_<version>_darwin_amd64.tar.gz
lodestone_<version>_darwin_arm64.tar.gz
lodestone_<version>_windows_amd64.zip
lodestone_<version>_windows_arm64.zip
checksums.txt
```

Jedes Archive enthält das `lodestone`-Binary, `LICENSE`, `README.md`,
`CHANGELOG.md`.

## Was noch nicht aktiv ist

- **Homebrew-Tap** und **nfpms (apt/deb)** in `.goreleaser.yaml`
  bewusst auskommentiert, bis Tap-Repo und Release-Token in der
  jmt-labs-Org konfiguriert sind. Aktivierung erfolgt vor dem ersten
  stabilen Release (Phase 2 oder später).

## Versionierungs-Schema

Semantic Versioning. Phase-Tagging-Konvention:

- Phase 1 MVP: `v0.1.0-alpha`
- Erstes Beta: `v0.2.0-beta`
- Erste stabile API (nach Phase 2 abgeschlossen): `v1.0.0`

Pre-Release-Suffixe (`-alpha`, `-beta`, `-rc.N`) bleiben so lange in
Verwendung, wie die `.lodestone/`-Datei-Layouts noch brechen dürfen.
