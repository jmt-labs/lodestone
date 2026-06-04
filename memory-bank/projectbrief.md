# Project Brief — lodestone

lodestone sammelt AI-Ökosystem-Signale (GitHub Trending, HackerNews, ArXiv, npm,
Anthropic/OpenAI-Changelogs), scored sie deterministisch gegen einen Repo-Fingerprint
und liefert reproduzierbare Empfehlungen als CLI + MCP-Server.

## Was ist dieses Projekt?

Ein lokales CLI-Tool, das Rauschen im AI-Ökosystem filtert und gezielt die Signale
herausfiltert, die für ein konkretes Projekt relevant sind — basierend auf einem
Repo-Fingerprint (Goals, TechInterests, Schwellwerte).

## Ziele

- Reproduzierbare, deterministische Empfehlungen (zwei Läufe = byte-identisches Ergebnis)
- Signal-Rausch-Filterung für AI-Ökosystem-News
- Dateibasierter State, kein Server, kein Daemon
- CLI-First; MCP-Server als Erweiterung (Phase 2)

## Nicht-Ziele

- Kein Frontend
- Kein Daemon / kein Hintergrundprozess
- Keine Telemetrie (lokale Artefakte bleiben lokal)
- Kein auto-push auf main
- Kein ORM, keine Datenbank
- Kein auto-merge (Auto-PRs erst Phase 4, immer Draft + Feature-Branch)

## Phasen-Übersicht

| Phase | Status | Inhalt |
|---|---|---|
| Phase 1 | Abgeschlossen | CLI-Grundgerüst, deterministisches Scoring, Konfiguration |
| Phase 2 | Aktiv | LLM-Integration (`lodestone plan`, `lodestone apply`), MCP-Server |
| Phase 3 | Geplant | Weitere Signal-Quellen |
| Phase 4 | Geplant | Auto-PRs (Draft, Feature-Branch, harte Schranken) |
