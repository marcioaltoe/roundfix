---
task: task_07
spec: 0003-dogfood-polish
status: completed
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

- [x] Artifact Directory on the spec-Run plan and CLI wiring
- [x] Log-path derivation shared with the review path
- [x] QA step log alignment
- [x] Engine and CLI tests asserting the new location and the absence of
      `.roundfix/` in the worktree

## Acceptance Criteria

- [x] A spec Run with a configured artifact dir writes `batch-NNN` agent logs
      under it and creates no `.roundfix/` directory in the repo.
- [x] The stderr diagnostics name the real log path.
- [x] Review-path log behavior is byte-identical.

## Verification

- `rtk go test ./internal/daemon/ ./internal/cli/` — expected: all tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 8; Core Feature 8; Decisions. `_techspec.md` →
Interfaces (TaskPlan.ArtifactDir), Build Order 7. Dogfood finding 8.

## Result

Spec Run agent logs now use the same Artifact Directory contract as review
Runs. `TaskPlan` carries `ArtifactDir`, task and QA batches derive their log
paths through `agent.LogPath(plan.ArtifactDir, runID, batch)`, and the
implement command resolves `defaults.artifact_dir` through
`config.ValidateArtifactDirectory` before creating the Run. The existing
stderr line still says `Agent log: <path>`; only the path changes to the
configured Artifact Directory.

Evidence:

- Red signal: after adding the regression tests, `rtk go test ./internal/daemon/ ./internal/cli/`
  failed because `TaskPlan.ArtifactDir` did not exist and implement still
  produced `<workdir>/.roundfix/runs/...` log paths.
- `TestTaskCycleExecutesAgentVerifySettleCommitContract` now writes a fake
  task Agent log at `agent.LogPath(fixture.artifactDir, runID, 1)` and asserts
  no repo `.roundfix/` directory exists.
- `TestTaskCycleQAVerdictMatrixSettlesRunAndCommitsReport` asserts the QA
  batch log uses `agent.LogPath(fixture.artifactDir, runID, 2)` and creates no
  repo `.roundfix/` directory.
- `TestRunImplementUsesConfiguredArtifactDirectoryForAgentLogs` covers both
  repo-relative and `~` configured Artifact Directory values, asserts the fake
  Agent log file exists under the resolved directory, and checks stderr names
  that real path.
- Review-path log code was not edited; `rtk go test ./internal/daemon/ ./internal/cli/`
  passed with existing review-path tests unchanged.

Verification:

- `rtk go test ./internal/daemon/ ./internal/cli/` passed: 203 tests in 2 packages.
- `rtk go test ./...` passed: 478 tests in 16 packages.
- `rtk make verify` passed: full Go suite, Roundfix skill check, and build.
