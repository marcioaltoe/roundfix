---
task: task_04
spec: 0067-derived-artifact-regeneration-boundary
status: completed
type: backend
complexity: low
---

# Task 04: Say what a human must do when the command cannot

## Overview

`make baseline-digests` prints "Read the failing test output above, fix the
canonical source it validates, then rerun make baseline-digests" — advice that
cannot be followed when the canonical source is already correct and only a
non-sanctioned corpus is stale. This slice makes the diagnostic read the
ownership record and name the actual human action.

## Requirements

1. MUST name, when a failure's remediation lies outside the sanctioned command,
   the exact invocation from the owning record rather than printing a command
   that will not help.
2. MUST name the owning record's path, so the reader can verify the claim.
3. MUST keep the existing diagnostic for failures the sanctioned command *can*
   fix, unchanged.
4. MUST state plainly, for a `frozen` artifact, that nothing regenerates it and
   why, instead of suggesting a regeneration.

## Subtasks

- [ ] Read the ownership record when a failure falls outside the command.
- [ ] Emit the declared invocation and the record path.
- [ ] Keep the in-scope diagnostic unchanged.

## Acceptance Criteria

- [ ] A stale `dedicated` corpus produces a diagnostic naming its exact
      invocation and its record path.
- [ ] A failure the sanctioned command can fix keeps today's diagnostic
      verbatim.
- [ ] A `frozen` artifact's diagnostic states nothing regenerates it and gives
      the recorded reason.
- [ ] No diagnostic suggests a command that cannot resolve the failure it
      reports, asserted across all three owner classes.

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/baseline -count=1 -run 'Diagnostic|Remediation|Ownership' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the diagnostic tests ran and passed.
- `go test ./internal/baseline -count=1` — expected: exit 0.
- `make verify` — expected: exit 0.

## References

- `_prd.md` → Core Feature 4; Goals.
- `_techspec.md` → API Contracts; Build Order 4.

## Result

### Implementation

- Added ownership-backed remediation lookup for derived artifacts. The lookup
  resolves and reads the governing `_ownership.yml`, then formats one action
  for its declared `sanctioned`, `dedicated`, or `frozen` owner.
- Kept the sanctioned remediation text `run 'make baseline-digests'` unchanged.
  Dedicated remediation now includes the record's exact `command` and path;
  frozen remediation says that nothing regenerates the artifact and includes
  the record's reason and path.
- Routed stale plan-characterization failures and frozen parity-corpus identity
  failures through the ownership-backed remediation. The Makefile and ownership
  records remain unchanged in this slice.

### Focused checks

- Before the implementation, `rtk env
  GOCACHE=/Users/marcio/.roundfix/worktrees/wt67run-dee297f2/run_20260805T050602Z_bfee71697528c7d7/.gocache
  go test ./internal/baseline -count=1 -run
  '^TestDerivedOwnershipRemediationDiagnostics$'` exited 1 because
  `derivedArtifactRemediation` was undefined. This was the red starting signal.
- After the final code edit, `rtk env
  GOCACHE=/Users/marcio/.roundfix/worktrees/wt67run-dee297f2/run_20260805T050602Z_bfee71697528c7d7/.gocache
  go test ./internal/baseline -count=1 -run
  '^(TestBaselinePlanCharacterization|TestBaselineCompatibilityCorpus|TestDerivedOwnershipRemediation.*)$'`
  exited 0 (`ok roundfix/internal/baseline`, 1.831s). This exercised the two
  changed failure-owning suites plus owner-class and unowned-artifact
  remediation checks.
- `rtk git -c core.fsmonitor=false diff --check` exited 0 with no output.

### Acceptance evidence

1. `TestDerivedOwnershipRemediationDiagnostics/dedicated_plan_characterization`
   reads the repository's dedicated record and asserts that remediation contains
   both its exact `command` and
   `testdata/plan-characterization/_ownership.yml`. The stale-golden failure
   now uses that lookup.
2. `TestDerivedOwnershipRemediationDiagnostics/sanctioned_setup` asserts exact
   equality with the pre-existing `run 'make baseline-digests'` diagnostic;
   the existing sanctioned call sites continue to read that same constant.
3. `TestDerivedOwnershipRemediationDiagnostics/frozen_parity_corpus` asserts
   that remediation says `nothing regenerates this artifact` and contains both
   `testdata/parity-corpus/_ownership.yml` and its recorded reason. The parity
   identity failure now emits that result instead of the sanctioned command.
4. The owner-class table rejects the sanctioned command for `dedicated` and
   `frozen`, requires exact sanctioned text for `sanctioned`, and
   `TestDerivedOwnershipRemediationRejectsUnownedArtifact` proves an
   unclassified artifact receives no guessed action.

The Task's declared `## Verification` commands were not run; the Daemon owns
that gate and terminal status settlement.
