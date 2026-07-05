---
task: task_07
spec: 0003-dogfood-polish
status: pending
type: backend
complexity: low
---

# Task 07: Place spec-Run agent logs under the Artifact Directory

## Overview

Review Runs log agent output under the configured Artifact Directory; spec
Runs invented a `<workdir>/.roundfix/runs/<run-id>/` location because
`TaskPlan` carried no Artifact Directory. Unify: one setting governs all
agent logs. Verifiable through engine tests asserting the log path.

## Requirements

1. MUST add the Artifact Directory to the spec-Run plan and derive agent log
   paths through the same helper the review path uses (per-Run, per-Batch
   naming consistent with review logs).
2. MUST resolve the directory exactly like review Runs do (config
   `defaults.artifact_dir`, `~` and repo-relative handling) and create it as
   needed.
3. MUST stop writing anything under `<workdir>/.roundfix/` — no new dotdir in
   user repositories; the QA step's log follows the same rule.
4. MUST print the agent log path on stderr exactly as today, just pointing at
   the new location.

## Subtasks

- [ ] Artifact Directory on the spec-Run plan and CLI wiring
- [ ] Log-path derivation shared with the review path
- [ ] QA step log alignment
- [ ] Engine and CLI tests asserting the new location and the absence of
      `.roundfix/` in the worktree

## Acceptance Criteria

- [ ] A spec Run with a configured artifact dir writes `batch-NNN` agent logs
      under it and creates no `.roundfix/` directory in the repo.
- [ ] The stderr diagnostics name the real log path.
- [ ] Review-path log behavior is byte-identical.

## Verification

- `rtk go test ./internal/daemon/ ./internal/cli/` — expected: all tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 8; Core Feature 8; Decisions. `_techspec.md` →
Interfaces (TaskPlan.ArtifactDir), Build Order 7. Dogfood finding 8.
