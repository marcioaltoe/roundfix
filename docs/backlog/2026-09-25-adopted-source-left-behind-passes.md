---
type: fix
status: open
created: 2026-09-25
spec: null
reason: null
---

# Spec check accepts an adopted Finding or Backlog Entry whose original is still in place

## Symptom

Adoption is one move (ADR-0083), but the checker only proves the destination exists; a leftover original in `docs/findings/` or `docs/backlog/` passes `--strict`.

## Where

`internal/speccheck/citations.go` indexed references; `internal/speccheck/backlog.go`.

## Expected

Resolve the index's source path and refuse when the original is still present.

## Evidence

secondbrain `inbox/roundfix/2026-09-08-adocao-indexada-aceita-origem-duplicada.md`.
