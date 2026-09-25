---
type: fix
status: open
created: 2026-09-25
spec: null
reason: null
---

# `roundfix events` aborts for the whole Run after a vacuous pre-work Verification

## Symptom

The vacuous pre-work Verification event is published with `probed_commands` while the projection requires `commands`, so the Supervisor event stream aborts with `missing payload field "commands"`.

## Where

`internal/daemon/task_engine.go` (publisher), `internal/runevent/stream.go` (projection), `internal/cli/events.go`.

## Expected

Align key and shape, add a publisher→projection contract test, and degrade a malformed record to a warning.

## Evidence

secondbrain `inbox/roundfix/2026-09-17-evento-de-verificacao-vacua-usa-chave-que-a-projecao-nao-le.md` (also defect 1 of `2026-09-10-events-aborta-e-carry-forward-quebra-no-worktree.md`).
