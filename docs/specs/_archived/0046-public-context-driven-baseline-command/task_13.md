---
task: task_13
spec: 0046-public-context-driven-baseline-command
status: completed
type: backend
complexity: high
---

# Task 13: Restore Repository Skill Sets through Baseline

## Overview

Port the maintained skill-restoration operation to
`roundfix baseline skills restore`. Users can preview and confirm immutable
external skill recovery without retaining a Python fallback or adding states
to Baseline adoption.

## Requirements

1. MUST preserve the maintained profile, selected-skill, offline source,
   immutable provenance, lock-digest, preview, confirmation, stale-plan,
   transaction, and rollback contracts.
2. MUST expose the TechSpec command, text/JSON output, stdout/stderr, and exit
   categories under the Baseline command family.
3. MUST group restoration by exact provenance, validate source bytes and
   adapter compatibility before mutation, and update the Repository Skill Set
   plus lock evidence atomically.
4. MUST use the recoverable transaction and exact Plan Digest confirmation for
   every non-empty restoration.
5. MUST reproduce every exact or semantic restoration fixture and designed
   delta recorded in the compatibility matrix.

## Subtasks

- [x] Port restoration planning and immutable provenance validation.
- [x] Implement online and declared offline source acquisition boundaries.
- [x] Integrate exact confirmation, transaction, and rollback.
- [x] Expose stable text/JSON command results.
- [x] Add compatibility, stale, unsafe-target, and failure-matrix tests.

## Acceptance Criteria

- [x] Preview performs zero mutation and returns the exact confirmation digest.
- [x] Exact offline sources restore without network access.
- [x] Provenance, source, adapter, target, or lock mismatch fails before mutation.
- [x] Partial write or lock-update failure restores the complete preimage.
- [x] Empty restoration is idempotent.
- [x] Maintained Python restoration cases have Go parity evidence.

## Context

- instruction: `docs/adr/0072-baseline-go-cutover-preserves-python-contracts.md`
- interface: `.agents/skills/setup-context-driven/scripts/context_setup.py`
- interface: `skills-lock.json`

## Verification

- `rtk go test -count=1 ./internal/baseline ./internal/cli -run 'TestBaselineSkillsRestore|TestSkillsRestoreOffline|TestSkillsRestoreProvenance|TestSkillsRestoreStalePlan|TestSkillsRestoreRollback|TestSkillsRestoreCompatibility'` — expected: public command, proof, confirmation, idempotence, rollback, and parity cases pass.
- `rtk go run -buildvcs=false ./cmd/roundfix baseline skills restore --help` — expected: help matches the approved repository, profile, skill, source, confirmation, and format contract.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → Goal 5; Core Features 19–20; Success Metrics.
- `_techspec.md` → API Contracts: maintenance operations; Testing Approach; Build Order 8.
- ADR-0072 → Go destination for maintained restore-skills behavior.

## Result

Implemented `roundfix baseline skills restore` as a Go-owned maintenance
operation. The command resolves the maintained profile and optional selected
skills, acquires immutable Git source trees from either HTTPS or the declared
offline source, groups work by exact provenance, validates the source tree,
target, adapter, and ordered lock document, and emits a deterministic
`restore-v1` preview with the exact confirmation digest.

Every non-empty confirmed restoration runs through the recoverable Baseline
file transaction. Skill content and `skills-lock.json` therefore commit
atomically, verify against exact postimages, and restore the complete preimage
after a partial write, lock replacement, or post-write verification failure.
The public command preserves the approved text/JSON, stdout/stderr, help, and
exit-category contracts without retaining a Python fallback or changing
Baseline adoption state.

### Acceptance evidence

1. `TestSkillsRestoreOfflinePreviewApplyAndIdempotence` snapshots the complete
   repository before and after preview, asserts zero mutation and an exact
   64-character Plan Digest, then confirms with that digest.
2. The same test uses a real local Git repository through `--source-dir`; its
   injected acquisition boundary records one offline acquisition and no
   network operation.
3. `TestSkillsRestoreProvenanceAndPreMutationRefusals` covers exact provenance
   grouping plus commit, tree digest, lock adapter, unsafe target, and malformed
   lock refusals, asserting the repository remains byte-identical.
4. `TestSkillsRestoreRollbackRestoresSkillAndLockPreimage` injects failures
   after a partial skill write, during lock replacement, and during post-write
   verification; every case restores all original skill and lock bytes.
5. `TestSkillsRestoreOfflinePreviewApplyAndIdempotence` repeats restoration
   after the confirmed apply and receives the deterministic empty `noop`
   result without mutation.
6. `TestSkillsRestoreCompatibilityMatchesMaintainedPythonShape` proves the Go
   `restore-v1` payload, portable tree digest, planned file changes, ordered
   lock entry, and lock mode match the maintained Python contract.
7. `TestSkillsRestoreStalePlanDoesNotMutate` proves an incorrect or stale
   confirmation digest is action-required and cannot mutate the Repository
   Skill Set.
8. `TestBaselineSkillsRestoreCommand` exercises public help, JSON preview, text
   success, stdout/stderr separation, and exit codes through the CLI boundary.

### Verification

- `rtk go test -count=1 ./internal/baseline ./internal/cli -run 'TestBaselineSkillsRestore|TestSkillsRestoreOffline|TestSkillsRestoreProvenance|TestSkillsRestoreStalePlan|TestSkillsRestoreRollback|TestSkillsRestoreCompatibility'`
  passed with 18 tests across 2 packages.
- `rtk go run -buildvcs=false ./cmd/roundfix baseline skills restore --help`
  passed and printed the repository, profile, repeatable skill, offline source,
  exact confirmation, format, and exit-category contract.
- `rtk make verify` passed: 1,976 Go tests, both 256-test maintained Python
  suites, asset catalog validation, Roundfix skill checks, and the binary
  build.
- `rtk git diff --check` passed.

### Follow-up

- Outside this Task's slice, the first full gate exposed that
  `TestRepositoryInspectionNoMutation` inherits the user's global
  `core.fsmonitor=true` and can snapshot the mutable
  `.git/fsmonitor--daemon*` state. The isolated test passed 20 consecutive
  runs and the complete gate rerun passed; no unrelated test change was made.
- The later documentation and skill-sync Task must publish the maintenance
  command before any CLI-changing pull request is opened.
