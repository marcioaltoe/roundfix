---
status: accepted
created_at: 2026-07-24T21:27:41Z
updated_at: 2026-08-11T21:35:56Z
deprecated_at: null
superseded_by: null
---

# ADR-0073: Baseline apply uses a recoverable multi-file transaction

Baseline apply acquires a repository-local exclusive lock, records a
Git-private recovery journal, stages every postimage, revalidates the complete
bounded preimage, replaces files in deterministic order, and verifies every
postimage. Any incomplete apply rolls back in reverse order from exact saved
preimages; an interrupted transaction must be recovered before another apply
can start. Sequential best-effort writes and whole-worktree cleanliness checks
were rejected because neither provides the atomic rollback and unrelated-work
isolation required by the confirmed Change Plan.
