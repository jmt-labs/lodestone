# `lodestone fingerprint`

Scannt das aktuelle Repo und erzeugt ein strukturiertes Profil unter
`.lodestone/fingerprint.json`. Erste Pipeline-Stufe — `score` braucht
den Fingerprint als Input.

## Synopsis

```sh
lodestone fingerprint [--root <pfad>]
```

## Flags

| Flag | Default | Zweck |
|---|---|---|
| `--root` | `$PWD` | Projekt-Wurzel, die gescannt wird |

## Verhalten

Der Walker scannt rekursiv, überspringt aber `.git`, `vendor`,
`node_modules`, `dist`, `build`. Erkannt werden:

- **Sprachen** über Dateierweiterungen (`.go`, `.js`/`.jsx`/`.ts`/`.tsx`).
- **Frameworks** über `go.mod` und `package.json`: `cobra`,
  `anthropic-sdk`, `gin`, `echo`, `chi`, `react`, `vue`, `next`,
  `svelte`, `mcp-sdk`, `@anthropic-ai/sdk`.
- **Dependencies** aus `go.mod` (Regex-Parse für inline + Block-Require)
  und `package.json` (`dependencies` + `devDependencies`).
- **LOC pro Sprache.**
- **Test-Ratio** = Test-LOC / Non-Test-LOC.
- **CI-Provider:** `github_actions` (`.github/workflows/`),
  `gitlab_ci` (`.gitlab-ci.yml`), `circleci` (`.circleci/config.yml`).
- **MCP-Konfiguration** über Existenz von `.mcp.json`.
- **Goals und TechInterests** aus `.lodestone.yaml`.

Output ist atomar via `tmp + rename` geschrieben. Schema-Details:
[Datenmodell § Fingerprint](../../internals/data-model.md#fingerprint).

## Beispiele

```sh
# Aktuelles Verzeichnis scannen
lodestone fingerprint

# Anderen Pfad scannen
lodestone fingerprint --root /pfad/zum/repo

# Ergebnis ansehen
cat .lodestone/fingerprint.json | jq .
```

## Exit-Codes & Fehler

| Code | Bedeutung |
|---|---|
| 0 | Erfolg |
| ≠0 | Walker-Fehler (z. B. Lesefehler), Schreibfehler auf `.lodestone/` |

## Verwandt

- [`ingest`](ingest.md) — externe Signale parallel holen.
- [`score`](score.md) — Fingerprint × Signals scoren.
- [Konfiguration](../configuration.md#lodestoneyaml) — Goals und
  TechInterests setzen.
