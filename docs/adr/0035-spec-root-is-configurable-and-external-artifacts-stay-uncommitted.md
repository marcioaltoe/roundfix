---
status: accepted
created_at: 2026-07-07T15:59:11Z
updated_at: 2026-08-11T21:35:56Z
deprecated_at: null
superseded_by: null
---

# Spec Root is configurable and external Spec artifacts stay uncommitted

Spec artifacts may live outside the repository working tree — typically a
knowledge workspace repository nested behind a `docs/specs` symlink — so the
Spec Root becomes explicit configuration (`specs.root`, Project > User >
built-in `docs/specs`), resolved once against the user's checkout and carried
into Runs so Worktrees never re-resolve a relative root. When the resolved
root is external, Daemon commits carry only repository changes: task files and
QA Reports settle on disk in the external repository and are never staged into
the code repository, whose history must not absorb another repository's
artifacts; committing the knowledge workspace stays a user action. Staging
also drops any path that crosses a symbolic link (with a journaled warning)
regardless of configuration, because git refuses such pathspecs and a failed
Task commit is strictly worse than an uncommitted artifact.
