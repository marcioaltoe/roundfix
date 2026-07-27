---
task: task_05
spec: 0038-terminal-run-worktree-reconciliation
status: pending
type: backend
complexity: medium
---

# Task 05: Surface retained terminal Run Worktrees in Runs List

## Overview

Make retained terminal Run Worktrees discoverable without changing the stable
Runs List stdout rows. Each view reports the exact relevant count on stderr and
points to the Reconcile Command as the sole classification surface.

## Requirements

1. MUST preserve existing Runs List stdout row shape byte-for-byte.
2. MUST count terminal spec Runs that retain a recorded worktree path or Run
   Branch.
3. MUST report hidden retained residue when the default Active view has no
   matching terminal rows.
4. MUST report retained residue represented in terminal and all-state views.
5. MUST write the bounded note to stderr only.
6. MUST point to `roundfix reconcile` without classifying safety in Runs List.
7. MUST omit the note when no retained terminal surface exists.

## Subtasks

- [ ] Derive repository-scoped retained-worktree counts.
- [ ] Add the exact stderr guidance.
- [ ] Preserve stdout row serialization.
- [ ] Cover Active, terminal, and all-state views.
- [ ] Cover zero-retention and mixed-repository cases.

## Acceptance Criteria

- [ ] Active view reports the exact hidden terminal residue count on stderr.
- [ ] Terminal and all-state views report only retained Runs relevant to their
      repository scope.
- [ ] Existing stdout rows remain byte-identical.
- [ ] The note contains the exact `roundfix reconcile` pointer.
- [ ] Runs List never labels a retained entry safe or unsafe.
- [ ] No note is printed when neither worktree nor Run Branch remains.

## Context

- instruction: `.agents/skills/agentic-cli-design/SKILL.md`
- instruction: `.agents/skills/golang-cli/SKILL.md`
- instruction: `.agents/skills/golang-testing/SKILL.md`
- interface: `internal/cli/runs.go`
- interface: `internal/cli/cli_test.go`
- interface: `internal/store/store.go`
- interface: `internal/worktree/worktree.go`

## Verification

- `rtk go test ./internal/cli -run 'TestRun.*List.*Retained.*Worktree' -count=1`
  — expected: all views report exact stderr counts while stdout stays
  byte-stable.
- `rtk go test ./internal/cli -run 'TestRun.*List.*(Active|Terminal|All)' -count=1`
  — expected: existing Runs List view contracts remain passing.

## References

- `_prd.md` → Goal 3; User Story 3; Core Feature 8; Success Metrics.
- `_techspec.md` → API Contracts: Run listing; Testing Approach; Build Order 5.
- `CONTEXT.md` → Runs List and Reconcile Command vocabulary.
