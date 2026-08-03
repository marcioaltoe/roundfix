---
status: accepted
created_at: 2026-08-03T00:00:00Z
updated_at: 2026-08-03T00:00:00Z
deprecated_at: null
superseded_by: null
---

# Repository facts are read in batches, never cached across mutations

Roundfix issued one git subprocess per repository fact and, in two loops, one
per file — about six thousand spawns per verification run and the same
pattern against every user repository — so reads that share an immutable
scope now go through git's own batch interfaces (`cat-file --batch`, and
multi-query `rev-parse`), making subprocess count proportional to operations
rather than to files. Caching facts across mutation boundaries was rejected:
a stale repository fact corrupts a Run's decisions, while a subprocess only
slows one, so batching is bounded by scopes where the repository state
provably cannot change mid-read. Extracting a shared git client was also
rejected — the cost was process count, not code duplication, and the
per-package runners keep their distinct error vocabularies.
