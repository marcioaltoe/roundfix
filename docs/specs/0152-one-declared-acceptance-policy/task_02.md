---
task: task_02
spec: 0152-one-declared-acceptance-policy
status: completed
type: backend
complexity: low
---

# Task 02: Archive calls it

## Overview

Archive holds the policy the repository documented. This Task moves it behind
the one decision without changing what archive does.

## Requirements

1. MUST replace archive's inline verdict judgement with a call to the one
   decision.
2. MUST keep archive accepting exactly the reports it accepts today and refusing
   exactly the reports it refuses today, with the same reasons.
3. MUST stay inside the bounded path `internal/spec/archive.go`.
4. MUST NOT change any other archive precondition.

## Subtasks

- [ ] Call the one decision from archive.
- [ ] Pin archive's accepted and refused shapes with tests.

## Acceptance Criteria

- [ ] A `pass` archives, including one carrying an environment-blocked row.
- [ ] A qualifying `partial` archives, and every other shape refuses with
      today's reason.
- [ ] No archive precondition other than the verdict judgement is touched.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/spec/archive.go`

## Verification

- `out="$(go test -count=1 -v -run "^TestArchiveAppliesTheOneEligibilityDecision" ./internal/spec 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestArchiveAppliesTheOneEligibilityDecision"` — expected: exit 0; before this Task the case does not exist, so the command fails.

## References

- [_authorization.md](_authorization.md) — Approved bounded mutation

## Result

Implementation:

- `archiveUnprovenActions` now delegates QA Report acceptance and refusal
  reasons to `QAReportEligibility`. After an eligible `partial`, archive still
  reads the unreachable declarations to preserve the existing `unproven`
  metadata; every other archive precondition is unchanged.
- `TestArchiveAppliesTheOneEligibilityDecision` pins plain and
  environment-blocked `pass` reports, a qualifying `partial` and its unproven
  actions, plus the existing reasons for `fail`, an unsupported verdict, and
  every disqualifying `partial` shape.

Focused checks:

- `rtk go test -count=1 -run '^(TestQAReportEligibility|TestSpec0058ReplayArchivesDeclaredUnreachableRelease|TestSpec0058ReplayReportsWronglyDeclaredReachableRow|TestSpec0058ReplayRefusesUnmatchedBlockedRow|TestArchivedPassCorpusRemainsArchiveEligible)$' ./internal/spec` — passed; 15 focused tests exercised the shared decision and the pre-existing archive acceptance/refusal paths. The first sandboxed attempt could not access the external Go build cache and produced no code verdict; the permitted rerun passed.
- `rtk rg -n 'QAReportEligibility\(specDir, report\)|rows_blocked_finding is|rows_blocked_environment is|rows_blocked_declared is' internal/spec/archive.go` — found only the shared-decision call at line 143; archive no longer carries the duplicated refusal clauses.
- `rtk git diff --check` — passed.

Acceptance evidence:

- A plain `pass` and a `pass` carrying an environment-blocked row are pinned as
  accepted by the archive precondition table; the archived-pass corpus check
  also passed.
- A qualifying `partial` remains accepted with its declared actions recorded;
  the Spec 0058 replay archive check passed. Every refused shape is pinned to
  the exact reason returned before this refactor.
- The production diff changes only archive's verdict judgement inside
  `archiveUnprovenActions`; Task-status checks, report reading, destination
  checks, metadata stamping, and movement are untouched.

Not run:

- The Task's declared `## Verification` command — reserved for the Daemon by
  the assigned execution contract.
