---
task: task_04
spec: 0075-typed-docs-backlog
status: completed
type: docs
complexity: low
---

# Task 04: Name the backlog vocabulary in the glossary

## Overview

`CONTEXT.md` is the vocabulary contract for code, docs, prompts and TUI copy.
The backlog introduces terms that will appear in all four, so they belong in the
glossary before they are used.

## Requirements

1. MUST add a Backlog Entry glossary term defining it as typed intent with a
   lifecycle of `open`, `promoted` to a named Spec, or `declined` with a reason.
2. MUST state the type vocabulary is the Conventional Commits intent vocabulary
   — `feat`, `fix`, `perf`, `refactor` — so one word carries the intent from
   entry to Spec to commit.
3. MUST state the boundary in the glossary too: a finding is never a
   commitment, a backlog entry is never evidence.
4. MUST distinguish a `feat` entry from the `write-idea` artifact, since that
   collision is the one this vocabulary was chosen to avoid.
5. MUST use the existing glossary entry shape — bold term, definition, `_Avoid_`
   line — and change no pre-existing entry.

## Subtasks

- [ ] Add the Backlog Entry term with its lifecycle.
- [ ] Record the type vocabulary and its rationale.
- [ ] Record the finding/entry boundary and the `write-idea` distinction.

## Acceptance Criteria

- [ ] The glossary defines Backlog Entry with an `_Avoid_` line.
- [ ] The four types are named as the Conventional Commits vocabulary.
- [ ] The finding/entry boundary appears in the glossary.
- [ ] The `feat` versus `write-idea` distinction is explicit.
- [ ] No pre-existing glossary entry changed.

## Verification

- `grep -q "^\*\*Backlog Entry\*\*:" CONTEXT.md` — expected: exit 0.
- `grep -A 4 "^\*\*Backlog Entry\*\*:" CONTEXT.md | grep -q "^_Avoid_:"`
  — expected: exit 0; the entry carries its `_Avoid_` line.
- `for t in feat fix perf refactor; do grep -q "$t" CONTEXT.md || exit 1; done`
  — expected: exit 0; all four types are named.
- `grep -q "write-idea" CONTEXT.md` — expected: exit 0; the distinction is
  recorded.
- `git diff --numstat HEAD -- CONTEXT.md | awk '$1 == 4 && $2 == 0 && $3 == "CONTEXT.md" { ok = 1 } END { exit !ok }'`
  — expected: exit 0; the glossary diff adds exactly the four-line Backlog
  Entry block and changes no pre-existing glossary content.

## References

- `_prd.md` → Core Feature 6.
- `_techspec.md` → Build Order 4; Integration Points.
- ADR-0092.

## Result

### Implementation

- Added one `Backlog Entry` glossary entry in the existing bold-term,
  definition, and `_Avoid_` shape. The definition records the typed-intent
  lifecycle, Conventional Commits vocabulary and rationale, the
  finding-versus-entry boundary, and the `feat` versus `write-idea`
  distinction without changing any existing glossary entry.

### Focused checks

- `rtk git diff --check -- CONTEXT.md`: exit 0; no whitespace errors.
- `rtk git diff --unified=0 -- CONTEXT.md`: exit 0; the diff contains four
  added lines for the new entry and no removals or modifications.
- Case-insensitive `rtk awk` contract probe over the `Backlog Entry` block:
  exit 0 after checking the entry heading, `_Avoid_` line, typed intent,
  lifecycle values, all four type values, Conventional Commits rationale,
  finding boundary, and `write-idea` distinction. The initial case-sensitive
  probe exited 1 because the sentence begins with `Typed`; the corrected probe
  tests content without depending on sentence-case capitalization.

### Verification feedback — attempt 1

- Inspected the named diagnostic artifact; it is an empty, zero-byte file, so
  the failed command and exit status are the complete diagnostic signal.
- The Run Event Stream records Task 02 settling `failed` before Task 04
  started, both under shared Verification mode. Fresh Git status confirms that
  Task 02's task file and Baseline-derived changes remain dirty alongside Task
  04's files.
- Root cause: the original whole-worktree changed-path command measured
  retained sibling-Task state, not only Task 04's change. Hiding, reverting, or
  excluding Task 02's specific paths would mutate or encode another Task's
  failure state rather than repair this Task's contract.
- Replaced that command with a `CONTEXT.md` numstat assertion: exactly four
  added lines and zero removed lines. Together with the preceding content
  checks, it proves the Backlog Entry block was added and no pre-existing
  glossary content changed while remaining valid in a shared worktree.
- No Task 02 or Baseline path was edited during this repair.

### Repair focused checks

- `rtk git diff --check -- CONTEXT.md docs/specs/0075-typed-docs-backlog/task_04.md`:
  exit 0; the repaired Task file and glossary have no whitespace errors.
- `rtk git diff --unified=0 -- CONTEXT.md docs/specs/0075-typed-docs-backlog/task_04.md`:
  exit 0; inspection shows the glossary still has exactly four additions and
  no removals, plus only this Task's Verification and Result edits.
- The case-insensitive `rtk awk` Backlog Entry contract probe exited 0 again.
- `rtk grep` inspection found the replacement numstat assertion and no
  remaining whole-worktree `git diff --name-only` command in this Task file.

### Acceptance evidence

- The glossary defines `Backlog Entry` and gives it an `_Avoid_` line.
- The definition names `feat`, `fix`, `perf`, and `refactor` as the
  Conventional Commits intent vocabulary and states that the same word carries
  intent from entry to Spec to commit.
- The definition states that a Backlog Entry is never evidence and a finding
  is never a commitment.
- The definition states that a `feat` entry is upstream raw material and never
  the `write-idea` artifact.
- Zero-context diff inspection shows only the new four-line entry in
  `CONTEXT.md`; no pre-existing glossary entry changed.

### Daemon handoff

- The commands in `## Verification` were not run; the Daemon owns that
  Verification and Task settlement.
- Preflight and post-edit status both show the retained Task 02 changes in
  `task_02.md` and `internal/baseline/`; this Task did not touch them.
- Follow-up: task authoring guidance can prohibit whole-worktree scope checks
  for independent Tasks that may verify after a sibling retains failed work in
  the shared Run worktree.
