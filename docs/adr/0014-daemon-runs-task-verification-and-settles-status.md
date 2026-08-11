---
status: accepted
created_at: 2026-07-05T22:17:04Z
updated_at: 2026-08-11T21:35:56Z
deprecated_at: null
superseded_by: null
---

# Daemon runs task verification and settles task status

The Daemon runs a Task's Verification commands verbatim after the Agent finishes, and only passing verification allows the Task commit (ADR 0001). The Agent writes task status and the Result section while working, but the Daemon settles the final status — a Task marked completed whose verification fails is settled failed — mirroring how the review path compensates for forgetful agents. A task file without a Verification section fails Preflight Validation before any Run is created.
