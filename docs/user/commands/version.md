# `lodestone version`

Zeigt Version, Commit-Hash und Build-Datum des Binaries an.

## Synopsis

```sh
lodestone version
```

## Verhalten

Gibt drei Felder auf stdout aus:

- **Version** — semantische Versionsnummer (z. B. `v0.1.0-alpha`).
- **Commit** — Git-Commit-Hash der Quelle.
- **Build-Datum** — Build-Zeitpunkt.

Die Werte werden zur Build-Zeit über `-ldflags` injiziert
(`internal/version`).

## Beispiele

```sh
lodestone version
# version: v0.1.0-alpha
# commit:  abc1234
# built:   2026-05-20T08:00:00Z
```

## Exit-Codes & Fehler

| Code | Bedeutung |
|---|---|
| 0 | Erfolg |

## Verwandt

- [Release-Prozess](../../contributor/release-process.md) — wie die
  Version-Strings gesetzt werden.
