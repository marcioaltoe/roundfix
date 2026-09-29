---
task: task_02
spec: 0184-baseline-plans-that-show-what-a-history-move-breaks
status: pending
type: backend
complexity: medium
---

# Task 02: Every Baseline Plan with History Relocations reports its Relocation Citations

## Overview

This Task calls the scan from Task 01 where `planHistoryMoves` has built the
ordered History Relocation ledger, and appends the returned findings to that
function's warnings. Every Baseline Plan then carries them: `baseline plan`,
`baseline update` and the interactive review all build their plans through
`BuildPlan`. The Plan Digest binds them through the existing hashing, and
nothing changes for a plan without relocations. The plan document goes to the
reviewer and to `baseline apply`. Apply must therefore write, move and verify
exactly what it did before, and a digest confirmed before a citing file
changed must no longer apply.

## Requirements

1. MUST call `relocationCitationFindings(ctx, root, moves, refused)`, with `refused` built from the `historyDestinationOccupied` collisions `discoverHistoryLayout` reported, in
   `planHistoryMoves` in `internal/baseline/plan.go`. The call happens after the
   ledger is sorted and ordinals are assigned, and only when the ledger is
   non-empty. The result MUST be appended after the existing retained-review
   and collision warnings. A scan error MUST be returned as a planning error
   that names the Relocation Citation scan.
2. MUST NOT add citing files to preimages, postimages, the managed-entry
   ledger, file changes or History Relocations. It MUST NOT change
   `computePlanDigest`, `ValidatePlanDocument`, the plan or result schema
   versions, any flag or exit code, or any renderer.
3. MUST put plan-level tests in a new
   `internal/baseline/history_citation_plan_test.go`. The tests MUST use the
   package's existing plan fixtures (`newPlanRepository`,
   `writeInspectionFile`, `commitInspectionRepository`, `buildTestPlan`) and
   `ApplyPlan`, without editing any existing test file.
4. MUST put CLI tests in a new `internal/cli/baseline_history_citation_test.go`,
   reusing the package's existing Baseline Plan and update test helpers without
   editing them. It covers:
   - `baseline plan --format=text` prints
     `Warning: baseline.history.citation: <citing path>: line <L> cites …`;
   - `--format=json` carries the same entry in `warnings`;
   - `baseline update --confirm-plan <digest>` with a digest taken before a
     citing file changed is refused through the existing digest-mismatch path,
     and leaves the repository bytes unchanged.
5. MUST keep every existing Baseline test green without renaming any.

## Subtasks

- [ ] Call the scan from `planHistoryMoves` and append its findings.
- [ ] Prove that the warnings are in the plan and that the digest binds them.
- [ ] Prove that apply is unchanged and that a plan without relocations is
      unchanged.
- [ ] Prove the text, JSON and stale-confirmation behavior through the CLI.

## Acceptance Criteria

- [ ] A plan whose relocated ADR is cited carries a `baseline.history.citation`
      warning for the citing file, while a current-layout plan carries none and
      performs no scan.
- [ ] Two repositories that differ only in a citing file have different Plan
      Digests, but identical preimages, postimages, History Relocations and file
      changes.
- [ ] Applying both plans yields identical trees apart from the citing file.
- [ ] `baseline plan` prints the warning line in text and the warning entry in
      JSON.
- [ ] `baseline update --confirm-plan` with a pre-edit digest is refused, and
      the repository bytes are unchanged.

## Context

- interface: `internal/baseline/plan.go`
- creates: `internal/baseline/history_citation_plan_test.go`
- creates: `internal/cli/baseline_history_citation_test.go`
- instruction: `internal/cli/baseline_update.go`
- instruction: `internal/cli/baseline_profile.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestBaselinePlanReportsRelocationCitations|TestRelocationCitationsBindThePlanDigestOnly|TestRelocationCitationsLeaveApplyUnchanged)$" ./internal/baseline 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestBaselinePlanReportsRelocationCitations TestRelocationCitationsBindThePlanDigestOnly TestRelocationCitationsLeaveApplyUnchanged; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf "missing PASS for %s\\n" "$name" >&2; exit 1; }; done` — expected: exit 0. Before this Task none of the named tests exists, so the command fails.
- `out="$(go test -count=1 -v -run "^(TestBaselinePlanPrintsRelocationCitationWarnings|TestBaselineUpdateRefusesADigestWhoseCitationsChanged)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestBaselinePlanPrintsRelocationCitationWarnings TestBaselineUpdateRefusesADigestWhoseCitationsChanged; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf "missing PASS for %s\\n" "$name" >&2; exit 1; }; done` — expected: exit 0. Before this Task none of the named tests exists, so the command fails.

## References

- [_prd.md](_prd.md) — Goals 1, 3 and 4; User Stories 1–3; Core Features 1, 4
  and 6; Success Metric 3
- [_techspec.md](_techspec.md) — System Architecture; API Contract 1;
  API Contract 2; Testing Approach 2–3; Build Order 2
- ADR-0173; ADR-0071; ADR-0073; ADR-0103; ADR-0068
