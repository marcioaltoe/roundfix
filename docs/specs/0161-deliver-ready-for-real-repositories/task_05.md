---
task: task_05
spec: 0161-deliver-ready-for-real-repositories
status: pending
type: backend
complexity: medium
---

# Task 05: Accept the archive commit the real command makes

## Overview

Corrective Task from the pre-PR review of 2026-09-24. `archiveCommitIsExact` compares the archived folder's tree with the active one, but `spec.Archive` stamps `_prd.md` (`status`, `archived`, `source_slug`, `unproven`) before moving it, so a resume after a crash following a real archive commit still parks as `review-stale`; the test used a bare rename.

## Requirements

1. MUST accept on resume a HEAD whose parent is the reviewed head and whose change is the Spec folder moved into the archive root, with `_prd.md` differing only in the stamped archive keys and its body byte-identical.
2. MUST refuse a commit carrying any other change.
3. MUST build the positive test from the output of the real `spec.Archive`.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] A commit made by the real archive command is accepted.
- [ ] A commit with an extra change, or a changed PRD body, is refused.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/cli/deliver_workflow.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestResumeAcceptsARealArchiveCommit|TestResumeRefusesAnArchiveCommitWithExtraChanges)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestResumeAcceptsARealArchiveCommit TestResumeRefusesAnArchiveCommitWithExtraChanges; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task neither case exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — Park and resume
