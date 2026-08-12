---
status: accepted
created_at: 2026-06-09T18:18:38Z
updated_at: 2026-08-11T21:35:56Z
deprecated_at: null
superseded_by: null
---

# Central run database

Roundfix stores Run state in a global SQLite database at `~/.roundfix/roundfix.db` instead of inside each repository Artifact Directory. Centralizing Run state gives the Daemon one place to track progress across multiple repositories and future concurrent Runs, while repositories may contain only markdown review artifacts when users choose that layout.

The Run Database also holds the Run Event Journal: an append-only history of Run Events ordered by a per-Run monotonic cursor. Keeping the journal in the central database, instead of in per-repository files or a separate broker, lets attach and replay surfaces read any Run's history through one durable store that survives terminal disconnects and process restarts. Run Event payload rules are recorded in ADR 0008.
