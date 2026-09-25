---
type: fix
status: open
created: 2026-09-25
spec: null
reason: null
---

# Read-only commands refuse a Run Database of another schema version and advise a migration that cannot run

## Symptom

`roundfix runs list` (and events, attach, doctor, gc, reconcile inspection, deliver status) refuses when the database schema differs from the binary's and says to run an operational command. There is no migrate command; an operational command needs a PR or Spec; and the same advice is printed when the database is newer, where no migration exists. On 2026-09-24 this silently stopped an operator's monitoring for 13 hours.

## Where

`internal/store/store.go` `openReader` (version check and message).

## Expected

An explicit `roundfix migrate` (or `doctor --migrate`), and a direction-aware message: older database → run migrate; newer database → use the newer binary.

## Evidence

Session 2026-09-25: `runs list` → `has schema version 13, but this binary supports schema version 18`; installed 0.15.0 binary refused a v18 database.
