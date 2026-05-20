# ADR-0007 — Skill-Embed: `flavors/` kanonisch, `internal/.../skills/data/` als Copy

## Status

Accepted, 2026-05-20.

## Kontext

Phase 2 liefert vier Claude-Skills (`lodestone-scout`,
`-recommend`, `-plan`, `-review-trends`). `lodestone init` muss sie
ins User-Projekt nach `.claude/skills/` kopieren — auch wenn der User
nur das vorkompilierte Binary hat. Go braucht für `go:embed` die
Quelldateien im Paket-Verzeichnis. `flavors/lodestone/skills/` ist
der natürliche Ort, ist aber außerhalb des Embed-Pakets.

## Entscheidung

Zwei-Wege-Layout mit klarer Quelle:

- **Kanonische Quelle:** `flavors/lodestone/skills/*.md` — editierbar.
- **Embed-Copy:** `internal/lodestone/skills/data/*.md` — wird via
  `go:embed` ins Binary gezogen. **Nicht direkt bearbeiten.**
- **Sync-Mechanismus:** Make-Target `make skills-coverage` und
  CI-Check verifizieren, dass beide Verzeichnisse byte-identische
  Inhalte haben.

`docs/user/skills.md` ist der User-facing Index und verlinkt
ausschließlich auf `flavors/lodestone/skills/*.md` (keine
Inhalts-Duplikation).

## Konsequenzen

- **Plus:** `go:embed` funktioniert ohne Pfad-Akrobatik.
- **Plus:** Skill-Updates haben einen klaren Quell-Pfad
  (`flavors/`); CI-Check fängt vergessene Kopien.
- **Plus:** User-Doku referenziert die kanonische Quelle, nicht den
  Embed-Mirror.
- **Minus:** Zwei Pfade müssen synchron bleiben. Mitigiert durch
  Make-Target + CI-Job.

## Alternativen

- **Skills nur unter `internal/...skills/data/`.** Verworfen — User,
  die das Repo browsen, würden Skills im Versteck unter `internal/`
  suchen müssen.
- **Skills nur in `flavors/`, `go:embed` mit `../../flavors`.**
  Verworfen — `go:embed` erlaubt nur Pfade innerhalb des Pakets.

## Quelle

[Phase-2-Design](../../superpowers/specs/2026-05-20-lodestone-phase2-design.md),
`internal/lodestone/skills/`.
