---
task: task_07
spec: 0057-baseline-capability-evidence-and-retention
status: pending
type: backend
complexity: medium
---

# Task 07: Add the remediate-and-re-run outcome

## Overview

Facing a mixed divergence set, a maintainer who wants to fix things in the
repository and come back has only decline, which records a refusal of the
Baseline. Pausing for repository work and refusing the Baseline are different
intentions and should not share a journal record. This Task gives the
divergence prompt a fourth outcome.

## Requirements

1. MUST add a fourth prompt outcome that exits without writing, prints
   per-divergence remediation, and names the read-only re-check command to run
   afterwards.
2. MUST journal that outcome distinctly from decline, so pausing for repository
   work is never recorded as refusing the Baseline.
3. MUST write nothing when that outcome is chosen.
4. MUST document, on the adaptation option, its removal-only constraint, so the
   two options are not confused.
5. MUST leave the three existing outcomes' behavior, journal records, and exit
   codes unchanged.

## Subtasks

- [ ] Add the fourth outcome to the divergence prompt.
- [ ] Print per-divergence remediation and the re-check command.
- [ ] Journal it distinctly from decline.
- [ ] State the adaptation option's removal-only constraint.

## Acceptance Criteria

- [ ] Choosing the new outcome exits without writing anything.
- [ ] It prints remediation for each unsatisfied divergence and names the
      re-check command.
- [ ] Its journal record is distinguishable from a decline record.
- [ ] The adaptation option states it can only remove.
- [ ] Decline, adapt, and proceed keep their existing behavior and records.
- [ ] `git status --porcelain` shows no path outside `internal/baseline/`,
      `internal/cli/`, and this task file.

## Context

- interface: `internal/cli/cli.go`
- interface: `internal/baseline/profile_alignment.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/cli -run '^TestDivergencePromptRemediateOutcome$' -count=1 -v | grep -q -- "--- PASS: TestDivergencePromptRemediateOutcome"` —
  expected: exit 0; the outcome exits without writing and names the re-check.
- `go test ./internal/cli -run '^TestDivergencePromptJournalsDistinctly$' -count=1 -v | grep -q -- "--- PASS: TestDivergencePromptJournalsDistinctly"` —
  expected: exit 0; its record differs from decline.
- `go test ./internal/baseline -run '^TestBaselinePlanCharacterization$' -count=1 -v | grep -q -- "--- PASS: TestBaselinePlanCharacterization"` —
  expected: exit 0.
- `go test ./internal/baseline ./internal/cli -count=1` — expected: exit 0.

## References

- `_prd.md` → User Story 5; Core Features 7.
- `_techspec.md` → API Contracts; Build Order 10.
- ADR-0075.

