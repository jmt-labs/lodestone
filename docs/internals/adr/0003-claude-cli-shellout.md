# ADR-0003 — Planning via `claude --print`, kein SDK

## Status

Accepted, 2026-05-20.

## Kontext

Phase 2 braucht LLM-Aufrufe für `lodestone plan`. Optionen sind:

1. Anthropic-Go-SDK direkt einbetten.
2. HTTP-Client gegen die Anthropic-API mit eigener Auth.
3. Shell-out an die lokale `claude`-CLI.

## Entscheidung

Shell-out an `claude --print --model <id>` über `os/exec`.

```go
type Runner interface {
    Run(ctx context.Context, prompt string, model string) (string, error)
}

type ClaudeRunner struct { /* exec.Command-Wrapper */ }
type FakeRunner   struct { /* deterministischer Test-Stub */ }
```

## Konsequenzen

- **Plus:** Null neue Go-Dependencies (die `anthropic-sdk-go` würde
  fünf transitive Deps mitbringen).
- **Plus:** Auth-Handling liegt bei der `claude`-CLI — Lodestone sieht
  keine API-Keys, kein Token-Refresh.
- **Plus:** Modell-Routing über `--model` ist eine reine
  Argument-Frage; `base/models.yaml` definiert die Mapping-Tabelle.
- **Plus:** `Runner`-Interface erlaubt `FakeRunner` für Tests ohne
  Netzwerk- oder API-Kosten.
- **Minus:** Setzt voraus, dass die `claude`-CLI installiert ist —
  expliziter Hinweis in der [FAQ](../../user/faq.md) und im
  [Plan-Befehl](../../user/commands/plan.md).
- **Minus:** Stdout-Parsing ist fragiler als ein typisierter
  SDK-Aufruf. Mitigiert durch `===SPEC===` / `===PLAN===`-Marker im
  Prompt-Template (siehe `internal/lodestone/planning/`).

## Alternativen

- **Anthropic-SDK direkt.** Bricht Dependency-Budget aus
  [ADR-0005](0005-anti-hype-defaults.md) und zwingt zur Auth-Mechanik.
  Verworfen.
- **Raw-HTTP gegen Anthropic-API.** Doppelte Implementierung von
  Auth-Logik, die `claude`-CLI schon hat. Verworfen.

## Quelle

[Phase-2-Design](../../superpowers/specs/2026-05-20-lodestone-phase2-design.md),
`internal/lodestone/planning/`.
