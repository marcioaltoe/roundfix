---
task: task_04
spec: 0075-typed-docs-backlog
status: pending
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
- `git diff --name-only HEAD | grep -vE "^(CONTEXT\.md$|docs/specs/0075-typed-docs-backlog/task_04\.md$)" | grep -q . && exit 1 || exit 0`
  — expected: exit 0; only the glossary and this Task file changed.

## References

- `_prd.md` → Core Feature 6.
- `_techspec.md` → Build Order 4; Integration Points.
- ADR-0092.
