# `lodestone init`

Setzt Lodestone in einem Repository auf: legt `.lodestone.yaml` mit
Defaults an, hängt einen `.gitignore`-Snippet an und installiert die
vier Lodestone-Skills nach `.claude/skills/`.

## Synopsis

```sh
lodestone init [--config|--no-config] [--gitignore|--no-gitignore] [--skills|--no-skills] [--skills-dir <pfad>] [--force]
```

## Flags

| Flag | Default | Zweck |
|---|---|---|
| `--config` | true | `.lodestone.yaml` mit Defaults anlegen |
| `--gitignore` | true | `.gitignore`-Snippet für `.lodestone/` anhängen |
| `--skills` | true | Vier Lodestone-Skills installieren |
| `--skills-dir` | `.claude/skills` | Ziel-Verzeichnis für die Skills |
| `--force` | false | Vorhandene Dateien überschreiben |

## Verhalten

Im Default-Lauf werden drei Dinge angelegt oder erweitert:

1. **`.lodestone.yaml`** mit den konservativen Anti-Hype-Defaults
   (siehe [ADR-0005](../../internals/adr/0005-anti-hype-defaults.md)
   und [Konfiguration](../configuration.md)).
2. **`.gitignore`-Snippet:**

   ```
   .lodestone/
   !.lodestone/decisions.log
   ```

   `.lodestone/` ist ignoriert, der Audit-Trail bleibt committable.
3. **Vier Skills** unter `.claude/skills/`:
   `lodestone-scout`, `lodestone-recommend`, `lodestone-plan`,
   `lodestone-review-trends`. Quelle ist die go:embed-Kopie aus dem
   Binary; kanonisch unter `flavors/lodestone/skills/` (siehe
   [ADR-0007](../../internals/adr/0007-skill-embed-strategie.md)).

Ohne `--force` werden vorhandene Dateien nicht angefasst.

## Beispiele

```sh
# Voller Bootstrap
cd dein-projekt && lodestone init

# Nur Skills, ohne Config-Datei
lodestone init --no-config --no-gitignore

# Skills nach einem Custom-Verzeichnis
lodestone init --skills-dir .agents/skills

# Vorhandene Dateien überschreiben
lodestone init --force
```

## Exit-Codes & Fehler

| Code | Bedeutung |
|---|---|
| 0 | Erfolg |
| ≠0 | Schreibfehler, ungültiger `--skills-dir`, Konflikt ohne `--force` |

## Verwandt

- [Konfiguration](../configuration.md) — `.lodestone.yaml`-Felder.
- [Skills](../skills.md) — Übersicht und Verlinkung zu
  `flavors/lodestone/skills/`.
- [Quickstart](../quickstart.md) — empfohlene Erst-Befehlsfolge.
