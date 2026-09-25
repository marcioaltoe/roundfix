---
task: task_02
spec: 0159-archive-override-and-authoring-rules
status: completed
type: backend
complexity: medium
---

# Task 02: Claimed ADR ordinals

## Overview

A Task names the ADR it creates with `creates:`, and nothing checks that the number is free, so two Specs authored in parallel can each plan the same ADR number.

## Requirements

1. MUST add `SC-ORDINAL-CLAIMED` as an error at the Tasks stage, registered with the other detectors.
2. MUST report it when a Task's `creates:` path `docs/adr/NNNN-*.md` uses a number another file on the tree holds under a different name.
3. MUST report it when Tasks of two different active Specs create paths with the same number.
4. MUST accept a claim whose exact path already exists, and distinct numbers.
5. MUST update the corpus golden only if the active corpus counts change, recording why.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] A number held on the tree and a number claimed by two active Specs each fail with `SC-ORDINAL-CLAIMED`.
- [ ] Distinct numbers and a fulfilled claim pass.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/speccheck/coherence.go`
- creates: `internal/speccheck/ordinal.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestOrdinalClaimedOnTree|TestOrdinalClaimedByAnotherActiveSpec|TestDistinctOrdinalsAreAccepted)$" ./internal/speccheck 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestOrdinalClaimedOnTree TestOrdinalClaimedByAnotherActiveSpec TestDistinctOrdinalsAreAccepted; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task none of the three cases exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — Claimed ordinals

## Result

Implemented `SC-ORDINAL-CLAIMED` as a Tasks-stage error. The detector checks
ADR `creates:` claims from the Spec under review against differently named
files in `docs/adr/` and against claims from Tasks in other active Specs. An
existing file at the exact claimed path remains a fulfilled claim.

Focused checks:

- Red signal: `GOCACHE=/tmp/roundfix-task-02-gocache go test -count=1 -run '^TestOrdinalClaimedOnTree$' ./internal/speccheck` failed to compile because `CodeOrdinalClaimed` did not exist before the implementation.
- `GOCACHE=/tmp/roundfix-task-02-gocache go test -count=1 -run 'TestOrdinal|TestStageScope' ./internal/speccheck` passed.
- `GOCACHE=/tmp/roundfix-task-02-gocache go test -count=1 ./internal/speccheck` passed.
- `GOCACHE=/tmp/roundfix-task-02-gocache go test -count=1 -tags docscontract -run '^TestCheckCorpusGolden$' ./internal/docscontract` passed. Active corpus counts did not change, so the corpus golden was not updated.

Acceptance evidence:

- `TestOrdinalClaimedOnTree` and `TestOrdinalClaimedByAnotherActiveSpec` cover the tree-held and cross-active-Spec collision cases and assert `SC-ORDINAL-CLAIMED` error findings.
- `TestDistinctOrdinalsAreAccepted` covers both distinct active-Spec ordinals and an exact-path fulfilled claim without an ordinal finding.

Follow-ups: none.
