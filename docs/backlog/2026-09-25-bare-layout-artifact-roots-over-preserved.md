---
type: fix
status: open
created: 2026-09-25
spec: null
reason: null
---

# `gc --sanitize` preserves pre-upgrade artifact roots in bare-repository layouts

## Symptom

In bare-repository layouts the default Artifact Root moved to the common Git directory; `gc --sanitize` classifies a pre-upgrade per-worktree root `overridden` and never reclaims it.

## Where

`internal/cli/gc.go` `classifyGCSanitationRoot`; `internal/config/config.go`.

## Expected

Accept a root equal to the default derived from either the recorded key or the recorded checkout.

## Evidence

Recorded limit in `docs/history/specs/0162-a-durable-repository-key-per-run/_prd.md`.
