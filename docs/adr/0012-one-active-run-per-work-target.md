---
status: accepted
created_at: 2026-07-05T22:17:04Z
updated_at: 2026-08-11T21:35:56Z
deprecated_at: null
superseded_by: null
---

# One active run per work target

A Run now targets either an Open Pull Request's Review Issues or a Spec's Tasks, so PR Head Branch identity can no longer be the global concurrency boundary. The Run Database rejects a second Active Run for the same work target — (Head Repository, PR Head Branch) for review Runs, (repository, spec slug) for spec Runs — and, until worktree-per-task lands, Preflight Validation also rejects a second Active Run of any kind in the same working tree, because concurrent Runs would mutate one checkout. Supersedes ADR 0005.
