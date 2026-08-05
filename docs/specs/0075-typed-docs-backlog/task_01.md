---
task: task_01
spec: 0075-typed-docs-backlog
status: pending
type: backend
complexity: medium
---

# Task 01: Give the layout a home for typed intent

## Overview

The documentation layout gives observations a home — `docs/findings/`, dated and
evidence-backed — and gives intent none. Suggestions without observed evidence
were either forced into findings, diluting what a finding means, or lost.

This slice adds `docs/backlog/` to the layout contract as module clauses, which
is the unit adopting repositories actually receive.

## Requirements

1. MUST extend the one-job clause so `docs/backlog/` has a stated single job,
   beside the existing `docs/findings/` entry.
2. MUST add the Backlog Operational Contract clause carrying the frontmatter
   shape — `type`, `status`, `created`, `spec`, `reason` — and the body template
   for each of the four types.
3. MUST use the Conventional Commits intent vocabulary verbatim: `feat`, `fix`,
   `perf`, `refactor`. `refactor` is the canonical token, never an
   abbreviation.
4. MUST add the boundary clause stating both directions: a finding records what
   happened and is never a commitment; a backlog entry records what to do next
   and is never evidence. A finding may spawn an entry; an entry needs no
   finding.
5. MUST state that a `feat` entry is upstream raw material the spec pipeline may
   consume, never the `write-idea` artifact itself.
6. MUST state the extension rule in the contract itself: the type set is open,
   and a new type must be a Conventional Commits type that expresses intent.
   Extension is a contract change with a corpus re-record, never an informal
   addition — and the contract must say so where a reader adding a type will
   find it, not only in the Spec.
7. MUST leave the findings contract clause byte-identical. This Spec adds beside
   it and never edits it.
8. MUST NOT change any binary behaviour: no command reads or validates the
   backlog, and validation stays editorial exactly as findings work today.

## Subtasks

- [ ] Extend the one-job clause with `docs/backlog/`.
- [ ] Add the Backlog Operational Contract clause with all four templates.
- [ ] Add the boundary clause beside the findings contract.
- [ ] Confirm the findings clause is byte-identical.

## Acceptance Criteria

- [ ] The layout guide, when generated, documents `docs/backlog/` with its one
      job.
- [ ] The contract carries the frontmatter shape and a body template for each
      of `feat`, `fix`, `perf` and `refactor`.
- [ ] The boundary clause states both directions and the `write-idea`
      distinction.
- [ ] The contract states the extension rule — a new type must be a
      Conventional Commits type expressing intent, added as a contract change
      with a corpus re-record.
- [ ] The findings contract clause is byte-identical, asserted by diff rather
      than by reading.
- [ ] No Go source file changed.

## Context

- instruction: `docs/agents/docs-layout.md`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/baseline -count=1 -run 'Guide|Layout|Clause|Module' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the guide-generation tests ran and passed.
- `grep -q "docs/backlog" internal/baseline/assets/modules/context-workflow.json`
  — expected: exit 0; the layout contract carries the directory.
- `for t in feat fix perf refactor; do grep -q "$t" internal/baseline/assets/modules/context-workflow.json || exit 1; done`
  — expected: exit 0; all four types are present.
- `grep -q "Conventional Commits" internal/baseline/assets/modules/context-workflow.json`
  — expected: exit 0; the contract carries the extension rule a reader needs.
- `if grep -rn "refact\b" internal/baseline/assets/modules/context-workflow.json | grep -q .; then exit 1; fi`
  — expected: exit 0; the canonical token is never abbreviated.
- `git diff --name-only HEAD | grep -E "\.go$" | grep -q . && exit 1 || exit 0`
  — expected: exit 0; no Go source changed.

## References

- `_prd.md` → Core Features 1-5 and 7; Decisions.
- `_techspec.md` → Implementation Design; Build Order 1.
- ADR-0092.
