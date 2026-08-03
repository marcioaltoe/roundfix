---
task: task_06
spec: 0072-qa-is-a-task-not-a-flag
status: pending
type: docs
complexity: low
---

# Task 06: Align the agent guides with the authored gate

## Overview

`docs/agents/autonomous-work.md` still teaches the flag: it shows
`roundfix implement --spec <slug> --qa --detach`, tells the Supervisor when
to pass `--qa`, and reasons about re-requesting it per cycle. All of that
lives in the repository-owned section outside the setup markers (the
baseline assets carry no `--qa` reference, verified during design). The
guidance moves to the authored-gate contract: the gate is declared at
decomposition, the graph runs it, and a gate that reported on a graph that
later grew is invalidated by load-time validation.

## Requirements

1. MUST replace every `--qa` reference in `docs/agents/autonomous-work.md`
   with the authored-gate contract, keeping the section's existing rules
   about corrective-Task caps and serial-chain warnings intact.
2. MUST NOT edit inside setup-owned marker blocks; the changes live in the
   repository-owned section only.
3. MUST leave no `--qa` reference anywhere under `docs/agents/`.
4. MUST keep the guidance consistent with what task_03 shipped: the
   parameter does not exist, and an unsettled authored gate makes an
   all-completed graph runnable.

## Subtasks

- [ ] Rewrite the Supervisor guidance around the authored gate.
- [ ] Sweep `docs/agents/` for residual flag references.

## Acceptance Criteria

- [ ] `docs/agents/autonomous-work.md` describes declaring the gate at
      decomposition and never mentions the flag.
- [ ] No file under `docs/agents/` contains `--qa`.
- [ ] Setup-owned marker blocks are byte-identical.
- [ ] `git status --porcelain` shows no path outside `docs/agents/` and
      this task file.

## Verification

- `grep -rL -- "--qa" docs/agents/*.md | grep -q autonomous-work` —
  expected: exit 0; the guide no longer references the flag.
- `grep -rq -- "--qa" docs/agents/` — expected: exit 1; nothing under the
  guides references the flag.
- `go build -buildvcs=false ./...` — expected: exit 0.

## References

- `_prd.md` → Goals (declared once, not per invocation).
- `_techspec.md` → Build Order 6.
