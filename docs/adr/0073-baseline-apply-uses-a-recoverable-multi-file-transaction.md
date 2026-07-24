# ADR-0073: Baseline apply uses a recoverable multi-file transaction

Status: Accepted

Baseline apply acquires a repository-local exclusive lock, records a
Git-private recovery journal, stages every postimage, revalidates the complete
bounded preimage, replaces files in deterministic order, and verifies every
postimage. Any incomplete apply rolls back in reverse order from exact saved
preimages; an interrupted transaction must be recovered before another apply
can start. Sequential best-effort writes and whole-worktree cleanliness checks
were rejected because neither provides the atomic rollback and unrelated-work
isolation required by the confirmed Change Plan.
