# `lodestone apply <rec-id>`

Wandelt eine Recommendation in einen **Draft-PR** um. Vier Safety-Gates
müssen alle passen — siehe
[ADR-0008](../../internals/adr/0008-apply-safety-gates.md).

## Synopsis

```sh
lodestone apply <rec-id> [--model <id>] [--root <pfad>]
```

## Flags

| Flag | Default | Zweck |
|---|---|---|
| `<rec-id>` | — | Recommendation-ID aus `lodestone signals --top N` |
| `--model` | aus `base/models.yaml` | Modell-Override für die Spec/Plan-Generierung |
| `--root` | `$PWD` | Projekt-Wurzel |

## Safety-Gates

Alle vier müssen passen, sonst wird der Apply abgebrochen:

1. `recommendation.risk == "low"`.
2. `recommendation.effort == "XS"`.
3. `recommendation.compatibility >= 0.85`.
4. Kein Apply in den letzten 24 h (Cooldown via
   `.lodestone/applies.jsonl`).

Zusätzlich muss `git status` sauber sein — kein staged oder unstaged
Diff. Branch-Name immer `lodestone/apply-<rec-suffix>-<date>`. Der PR
ist **immer Draft**, **immer gegen `main`**, ohne Auto-Merge.

## Verhalten

1. Recommendation und Fingerprint laden, Safety-Gates prüfen.
2. Branch anlegen (`git checkout -b lodestone/apply-…`).
3. Spec + Plan über die Planning-Engine generieren (wie
   [`plan`](plan.md), aber inline).
4. Commit erstellen.
5. Push, Draft-PR via `gh pr create --draft` öffnen.
6. Eintrag in `.lodestone/applies.jsonl` mit Status `draft_open`.
7. Audit-Eintrag in `.lodestone/decisions.log`.

Pluggable Runner: `GitRunner` und `PRRunner` haben Real- und
Fake-Implementierungen. Tests laufen gegen die Fakes, der echte Pfad
braucht `git` und `gh` im `$PATH`.

## Beispiele

```sh
# Top-Recommendation finden und applyen
lodestone signals --top 5 --json | jq '.[0]'
lodestone apply sha256:abc123…

# Mit Sonnet statt Opus für die Spec
lodestone apply sha256:abc123… --model claude-sonnet-4-6

# Erfolgs-Statistik ansehen
lodestone stats
```

## Exit-Codes & Fehler

| Code | Bedeutung |
|---|---|
| 0 | Erfolg (Draft-PR geöffnet) |
| ≠0 | Safety-Gate failed, `git status` nicht sauber, `git`/`gh` nicht im PATH, Cooldown noch aktiv |

Typische Abbrüche (alle ausführlich auf stderr begründet):

- `apply rejected: risk != low` — Recommendation ist nicht risikoarm
  genug für Auto-Apply.
- `apply rejected: cooldown active (last apply Xh ago)` — 24h-Sperre
  noch aktiv.
- `apply rejected: git status not clean` — `git stash` oder commit
  zuerst.

## Verwandt

- [`undo`](undo.md) — PR schließen + Branch löschen.
- [`stats`](stats.md) — Apply-Erfolgs-Statistik.
- [ADR-0008](../../internals/adr/0008-apply-safety-gates.md) — Begründung
  der Safety-Gates.
