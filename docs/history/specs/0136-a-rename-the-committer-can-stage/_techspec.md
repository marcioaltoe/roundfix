---
spec: 0136-a-rename-the-committer-can-stage
status: active
created: 2026-09-14
---

# TechSpec — A rename the committer can stage

## Executive Summary

One changed-path set feeds two stages with different needs. Classification wants
every path the turn touched, including the source a rename removed. Staging
wants only paths `git add` can match. Spec 0135 gave the set what
classification needed and broke what staging requires.

The repair keeps one source of truth and filters at the point of use: the
staging filter already exists and already drops paths it cannot stage, so it
gains one more reason. Two adjacent stages are corrected to ask the change
rather than the record.

## Project Constraints

- Identifier strategy: applicable — preserve the existing refusal codes, reason codes and record field names; no new identifier is introduced. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no credential, network or HTTP surface is touched. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0130 keeps the Governed Path set monotonic and is not narrowed here; ADR-0057 keeps the Daemon the exclusive writer of Task status. Source: `docs/agents/spec-routing.md`, `docs/agents/domain.md`.
- Tooling authority: applicable — no protected tooling mutation proposed or authorized. The repair is ordinary source in `internal/daemon`, `internal/cli` and `internal/spec`. Source: `docs/agents/agent-instructions.md`.

## System Architecture

| Component | Responsibility | Change |
| --- | --- | --- |
| Stageable-path filter | Reduce the changed set to paths the committer can stage | Drop a path absent from both the worktree and the index, with a recorded reason |
| Governed-mutation classifier | Decide what authority the turn requires | None; it reads the unfiltered snapshot pair |
| Run final push | Decide whether push authority is required | Derive governed mutation from the Run's changed paths |
| Record location resolver | Report an unresolved record | Carry the derived record path into the unresolved result |
| Prior-changed-path reader | Report the paths a Run changed since its initial head | Stop collapsing a rename into its destination |

## Implementation Design

### Stage only what Git can match

`git add -f -- <path>` matches against the worktree and the index together. A
deletion is matchable because the path is still in the index until the deletion
is recorded. An already-staged rename source is matchable in neither, because
`git mv` has already written the rename to the index, so Git answers
`pathspec did not match any files` and fails the whole commit.

The filter that builds the stageable list already drops paths for reasons it
records — external to the repository, crossing a symbolic link, an executable
mode. Absence from both the worktree and the index is one more such reason, and
recording it keeps the drop visible rather than silent.

Classification is unaffected because it reads the snapshot pair directly rather
than the staged list. That separation already has a named regression lock in the
suite, which this Spec keeps.

### Ask the change, not the record

The Run's final push passes governed mutation as "an authorization record was
granted", which is a statement about the record rather than about the work. A
Spec whose record grants implement and commit but not push, and whose Run
changes only ordinary files, then fails a check for authority its change never
needed. The Run already knows which paths it changed; the decision must come
from those.

The true refusal is preserved: a Run that did change a Governed Path still
requires push authority, which is the case the regression lock covers.

### Name the record a reader can open

The unresolved result labels the record with the Spec-relative fragment rather
than the path under the configured Spec Root, and the refusal message prints
that label. Carry the path the resolver derived; when resolution failed before
deriving one, compose it from the configured Spec Root so the label is still a
path that can be looked at.

### A rename survives the diff too

The QA gate proved a third reader collapses the pair. After integration the Run
asks for the paths changed between its initial head and HEAD with
`git diff --name-only`, and Git's own rename detection answers with the
destination alone. The governed source disappears before the final-push
decision reads it, so a governed rename auto-pushed under a record that grants
no `push`, while a direct governed edit refused correctly.

Ask the diff not to detect renames. The same reader feeds the changed-path
audit against the grant and the Spec Context Bundle, and both are more correct
with both sides present: the audit then sees the governed path a rename
removed, and the bundle names both paths the turn touched.

## Coverage Map

- PRD Goal 1 → Stageable-path filter, Governed-mutation classifier.
- PRD Goal 2 → Run final push, Prior-changed-path reader.
- PRD Goal 3 → Record location resolver.
- Core Features 1-3 → Stageable-path filter, Governed-mutation classifier.
- Core Feature 4 → Run final push.
- Core Feature 5 → Record location resolver.

## Testing Approach

Focused tests at the three existing seams; no new seam is needed.

1. A Task that renames a file with `git mv` commits, in a repository the test
   creates.
2. A Task that renames a file without staging it also commits, and the
   deletion is recorded.
3. A Task that deletes a file still stages the deletion.
4. A governed rename still refuses without its operation, so staging and
   classification stay independent.
5. An ordinary Run under a record without push authority is not refused at the
   final push, and a governed Run without push authority still is.
6. An unresolved record's reported path resolves under the configured Spec
   Root, and the existing reason codes and fields are unchanged.
7. A governed rename is refused at the final push, and the prior-changed reader
   reports both sides of a rename rather than the destination alone.

Observation 1 rests on evidence this Spec did not author: Git's own pathspec
matching, measured directly — `git add -f` on a rename source exits non-zero
after `git mv` and zero after a plain `mv`. The assertion reads Git's behavior
rather than a rule this Spec composed. If Git cannot be run, the row records
blocked with that reason.

## Build Order

1. Stage only paths Git can match, proven by committing a real `git mv` rename
   (depends on: none).
2. Derive the final push's governed mutation from the Run's changed paths
   (depends on: none).
3. Carry a resolvable record path into the unresolved result (depends on: none).
4. Report both sides of a rename in the prior-changed reader, so the final
   push classifies a governed rename (depends on: 1, 2).
5. Terminal QA (depends on: 1, 2, 3, 4).

## Risks & Considerations

- Dropping a path from staging must never drop it from classification. The
  existing lock that classification reads the unfiltered snapshot is kept, and
  the governed-rename refusal is asserted alongside the successful ordinary
  rename.
- A recorded drop is visible in the commit's dropped-path report, so an
  unexpected drop is observable rather than silent.

## Decisions

- Filter at the point of use rather than narrowing the changed set, so the set
  keeps one meaning: every path the turn touched.
- Measure Git rather than reason about it. This defect existed because the
  previous test asserted staging on a tree with no rename in it.

## Vocabulary Contract

No token is coined. The repair reuses existing refusal tokens, reason codes and
drop reasons, and their glossary owners are unchanged.
