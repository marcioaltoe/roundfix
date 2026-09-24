---
task: task_06
spec: 0157-reconciliation-and-one-repository-identity
status: completed
type: backend
complexity: low
---

# Task 06: Same report means same content

## Overview

Corrective Task from the pre-PR review of 2026-09-24. `supersedingQAReport` now treats a Run-side report and a target-side report with the same file name as the same report without comparing their content, so a Run-side report with a different verdict can be declared superseded and its Run Branch released.

## Requirements

1. MUST treat a Run-side report and a target-side report with the same file name as the same report only when their blobs are identical.
2. MUST keep recognising an archived copy whose content is identical to the Run-side report.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] Same name with different content is not proven superseded.
- [ ] Same name with identical content across the active and archived roots is recognised.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/worktree/worktree.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestSupersedingQAReportRequiresTheSameContent|TestSupersedingQAReportRecognisesAnArchivedCopy)$" ./internal/worktree 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestSupersedingQAReportRequiresTheSameContent TestSupersedingQAReportRecognisesAnArchivedCopy; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task the first case does not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — Report identity

## Result

Implemented content-backed QA Report identity for the same-name,
different-root case. `supersedingQAReport` now resolves the Run-side and
target-side Git blob IDs and refuses supersession unless both IDs are present
and equal. Differently named reports retain their existing newest-report
ordering.

Focused-check evidence:

- Before the production change,
  `rtk env GOCACHE=/private/tmp/roundfix-task06-gocache go test -count=1 -run '^TestSupersedingQAReportRequiresTheSameContent$' ./internal/worktree`
  failed because an archived `verdict: pass` report superseded a same-named
  active `verdict: fail` report. After the change, the command exited 0. This
  covers the criterion that same name with different content is not proven
  superseded.
- `rtk env GOCACHE=/private/tmp/roundfix-task06-gocache go test -count=1 -run '^TestSupersedingQAReportRecognisesAnArchivedCopy$' ./internal/worktree`
  exited 0 after the change. This covers the criterion that byte-identical
  content under the active and archived QA roots is recognised.
- `rtk env GOCACHE=/private/tmp/roundfix-task06-gocache go test -count=1 -run '^TestSupersedingQAReport' ./internal/worktree`
  exited 0, covering both identity cases and the existing later-date and
  later-sequence behavior.
- `rtk env GOCACHE=/private/tmp/roundfix-task06-gocache go test -count=1 ./internal/worktree`
  exited 0 after the final Go edits, and `rtk git diff --check` exited 0.

The first pre-change attempt used Go's default build cache and could not open
the sandbox-external cache, so all recorded Go evidence uses the task-scoped
cache above. The Daemon-owned `## Verification` command was not run.

## Carry-forward provenance

- Source Run: `run_20260924T173241Z_9612fd1f6966c280`
- Source commit: `e1a3dff483c45fcc565cbb23bab9071564cd134f`
