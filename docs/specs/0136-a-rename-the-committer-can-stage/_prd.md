---
spec: 0136-a-rename-the-committer-can-stage
status: active
created: 2026-09-14
surfaces: [backend, cli]
---

# A rename the committer can stage

Spec 0135 made the changed-path reader carry a rename's source, so the governed
side of a rename reaches the classifier. That set is read by two stages, and
only one of them wanted the extra path. The commit stage stages explicit paths
with `git add -f`, and after `git mv` the source exists in neither the worktree
nor the index, because the index already records the rename. Git answers
`pathspec did not match any files` and exits non-zero, so a Task that renames a
file with `git mv` now fails to commit. Measured directly against Git: the same
command succeeds for an unstaged rename and fails for a staged one.

Two further reports in the same candidate are about a stage asking the wrong
question. The Run's final push derives whether a governed mutation occurred from
whether an authorization record was granted at all, so an ordinary Run under a
record that grants implement and commit but not push is refused for missing
push authority it never needed. And when a record's location cannot be
resolved, the unresolved result names the record as `<slug>/_authorization.md`,
so the refusal a reader sees points at a path that does not exist.

All three were found by the pre-Pull-Request review rounds on 2026-09-14, the
last round permitted for this candidate.

## Project Constraints

- Identifier strategy: applicable — preserve the existing refusal codes, reason codes and record field names; no new identifier is introduced. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no credential, network or HTTP surface is touched. The push in question is an existing Git operation whose authority check is being corrected, not added. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0130 keeps the Governed Path set monotonic and is not narrowed here; ADR-0057 keeps the Daemon the exclusive writer of Task status. Source: `docs/agents/spec-routing.md`, `docs/agents/domain.md`.
- Tooling authority: applicable — no protected tooling mutation proposed or authorized. The repair is ordinary source in `internal/daemon`, `internal/cli` and `internal/spec`, and their test files are ordinary too. Source: `docs/agents/agent-instructions.md`.

## Goals

- A Task that renames a file commits, whether the rename was staged by `git mv`
  or left unstaged, while a governed rename still requires its operation.
- The final push asks the change what it mutated, not the record whether it
  exists, so an ordinary Run is never refused for authority it does not need.
- An unresolved record names the path a reader can actually look at.

## Core Features

1. The commit stage stages only paths Git can match: a path absent from both
   the worktree and the index is dropped with a recorded reason, so an
   already-staged rename source is not offered to `git add`.
2. A deletion is still staged, because a deleted path remains in the index
   until the deletion is recorded.
3. Governed-mutation classification keeps reading the unfiltered snapshot, so
   dropping a path from staging never changes what authority is required.
4. The final push derives governed mutation from the paths the Run actually
   changed.
5. An unresolved record carries the record path the resolver derived, or the
   repository-relative path under the configured Spec Root when resolution
   failed before deriving one.

## Non-Goals / Out of Scope

- The changed-path reader's behavior, which 0135 settled: it keeps emitting
  both sides of a rename.
- External and symlinked Spec Root identity, including the external Task Graph
  authority gap the same review round reported, which stays promoted to its own
  Spec.
- Any change to the Governed Path set, to what a grant permits, or to the
  operation vocabulary.

## Declared intentional breaks

1. A path absent from both the worktree and the index stops being offered to
   `git add` and is reported as dropped, where it previously failed the commit.
2. An ordinary Run under a record without `push` stops being refused at the
   final push, where it was refused before.
3. An unresolved record reports a longer, resolvable path where it reported a
   bare slug-relative one.

Everything else keeps behaving as it does: refusal tokens, the operation
vocabulary, the Governed Path set, what authority a governed change requires,
and Daemon ownership of Task status.

## Regression locks

- The rename repair is proven by committing a real `git mv` rename in a
  repository the test creates, because Git's own answer is what this Spec
  measured and the previous test did not perform a rename at all.
- A governed rename still refuses without its operation, so dropping the source
  from staging does not drop it from classification.
- A governed Run still requires push authority, so correcting the false refusal
  does not remove the true one.
