---
status: accepted
created_at: 2026-10-27T01:30:00Z
updated_at: 2026-10-27T01:30:00Z
deprecated_at: null
superseded_by: null
---

# An untracked executable rides a Task commit only when the Spec's authorization names it

ADR-0157 kept the Daemon's refusal of every untracked executable file and
accepted the cost that a Task which deliberately produces one must add it to
the repository first. A Run cannot do that: the Agent's work only reaches the
repository through the Task commit, so a Spec that asks for an executable file
loses its Run. On 2026-10-26 Quaestor's Spec for launchd scripts lost its Task
this way after the Verification passed, and the maintainer chose that the
Spec's `_authorization.md` names the executables it may commit.

## Decision

A Spec's `_authorization.md` may carry `executables:`, a list of exact
repository-relative file paths. A Task commit keeps an untracked regular file
with an execute bit when its path is on that list of a granted record, and Git
records the mode it reads: `100755` under its default `core.fileMode: true`. Every other untracked executable stays refused as
ADR-0157 decided, and the refusal of a Task commit says how to name one. The
list is read from the authorization the Run resolved when it started, so no
Task or Agent can add to it, and it bounds no Governed Path: a listed path that
is also governed still needs `paths:`.

The stage filter judges files, never a directory entry: a Task commit and a
Batch commit expand an untracked directory into its files before the filter,
so an executable inside a new directory is refused like one beside it. A QA
Report commit and a Batch commit never read the list: a QA Report and its
evidence are text the gate writes, and a Batch commit belongs to a review Run,
which has no Spec authorization.

## Alternatives rejected

- Staging every executable, the alternative ADR-0157 rejected, stays rejected:
  it lets compiled output into the tree.
- Naming executables in the Task file: an Agent can edit its Task file during
  the Run, so the grant would widen from inside the work it bounds.
- Globs or directories in the list: one exact path per file keeps the grant
  reviewable and matches how `paths:` bounds Governed Paths.
