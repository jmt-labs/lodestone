# ADR-0004 — Cobra als CLI-Framework

## Status

Accepted, 2026-05-20.

## Kontext

Lodestone braucht eine Subkommando-CLI mit elf Verben, Flag-Parsing,
Help-Output und shell-completion. Go-Optionen: `flag` (stdlib),
`urfave/cli`, `spf13/cobra`.

## Entscheidung

`github.com/spf13/cobra` als einziges CLI-Framework. Zusammen mit
`gopkg.in/yaml.v3` für die Konfig ist Cobra die einzige nicht-stdlib-
Dependency, die Phase 1 zugelassen hat — siehe
[ADR-0005](0005-anti-hype-defaults.md).

## Konsequenzen

- **Plus:** Subkommando-Tree, Auto-Generation von Help und Bash-/
  Zsh-/Fish-Completions sind frei mitgegeben.
- **Plus:** De-facto-Standard im Go-Ökosystem — Kubernetes, Helm,
  GitHub-CLI, GoReleaser nutzen Cobra. Vertrautes Pattern für
  Contributors.
- **Plus:** Aktiv gewartet, stabile API seit v1.0 (2020).
- **Minus:** ~200 kB Binary-Wachstum gegenüber `flag`.
  Akzeptabel für eine CLI mit elf Verben.
- **Minus:** Subkommando-Definition ist verbose
  (`cobra.Command{Use:…, Run:…}`). Wird kompensiert durch konsistente
  Datei-Aufteilung pro Verb unter `cmd/lodestone/<verb>.go`.

## Alternativen

- **stdlib `flag`.** Kein Subkommando-Konzept; jeder Verb-Block müsste
  Flag-Parsing manuell sequenzieren. Verworfen für ≥ 4 Verben.
- **urfave/cli v3.** Aktiv, weniger verbose, aber kleinere Community
  und fragmentiertere Migrations-Story (v1 → v2 → v3). Verworfen.

## Quelle

[Phase-1-Design](../../superpowers/specs/2026-05-20-lodestone-design.md),
`cmd/lodestone/main.go`.
