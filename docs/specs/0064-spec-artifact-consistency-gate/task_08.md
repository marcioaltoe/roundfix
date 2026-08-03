---
task: task_08
spec: 0064-spec-artifact-consistency-gate
status: pending
type: infra
complexity: low
---

# Task 08: Wire the check into the local gate

## Overview

Make the check fail-closed: add a `spec-check` target and append it to the
`verify` gate, so a Spec contradiction stops a commit the way a test failure
does. This is the authorized tooling Task, and it runs last among the
implementation Tasks because task_07 has already brought the active Specs to a
clean report.

## Requirements

1. MUST add a `spec-check` target that runs the built binary's `spec check`
   over the Spec Root and propagates its exit status.
2. MUST append `spec-check` to the `verify` gate after `build`, so the target
   runs against a freshly built binary.
3. MUST NOT let a pipe hide the target's exit status — no pager, no `tail`, no
   pipeline whose last command masks a failure.
4. MUST declare the target in `.PHONY` and give it a `##` help description
   matching the file's existing style.
5. MUST change exactly one file: `Makefile`. This is the complete bounded file
   list from the Tooling authority row in `_prd.md` and `_techspec.md`, granted
   by the 2026-08-02 authorization record that names this Spec. Any other path
   is out of scope — stop rather than widen it.

## Subtasks

- [ ] Add the `spec-check` target with its help description.
- [ ] Append it to `verify` after `build`.
- [ ] Add it to `.PHONY`.

## Acceptance Criteria

- [ ] `make spec-check` exits 0 on the current clean corpus.
- [ ] `make spec-check` exits non-zero when a Spec carries an `error`, proven
      by a temporary edit reverted within the same check.
- [ ] `verify` lists `spec-check` after `build`.
- [ ] `spec-check` appears in `.PHONY` and in `make help`.
- [ ] `Makefile` is the only file this Task changed.

## Context

- instruction: `docs/agents/agent-instructions.md`
- interface: `Makefile`

## Verification

- `make build && make spec-check` — expected: exit 0; the gate target passes
  on the clean corpus.
- `grep -q "^verify:.*spec-check" Makefile` — expected: exit 0; the gate
  includes the target.
- `grep -q "^\.PHONY:.*spec-check" Makefile` — expected: exit 0.
- `make help | grep -q "spec-check"` — expected: exit 0.
- `git diff --name-only HEAD | grep -v "^Makefile$" | grep -v "^docs/specs/0064-spec-artifact-consistency-gate/task_08.md$" | grep -q . && exit 1 || exit 0`
  — expected: exit 0; only the bounded file and this Task file changed.

## References

- `_prd.md` → Project Constraints (tooling authority); Core Feature 1.
- `_techspec.md` → Integration Points; Build Order 8.
- `docs/workflow/authorizations/2026-08-02-queued-spec-tooling.md`.
