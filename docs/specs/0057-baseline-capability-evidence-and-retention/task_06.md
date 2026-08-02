---
task: task_06
spec: 0057-baseline-capability-evidence-and-retention
status: pending
type: backend
complexity: high
---

# Task 06: Offer a read-only capability re-check

## Overview

After remediating a blocking divergence there is no way to ask whether it is
fixed without resolving decisions and driving a full plan, so the
remediate-and-re-check loop does not exist. This Task adds a read-only re-check
that needs no decisions, writes nothing, and produces the same capability
outcomes a full plan would.

## Requirements

1. MUST provide a read-only capability re-check that requires no decisions to
   be supplied and resolves none.
2. MUST write nothing: no file, no journal mutation, no configuration change.
3. MUST produce capability outcomes identical to those a full plan produces for
   the same repository, by sharing the evaluation path rather than
   reimplementing it.
4. MUST render the same probe evidence and requirement grouping the full plan
   renders.
5. MUST report clearly when the repository has no resolvable Profile, rather
   than failing obscurely.
6. MUST leave the full plan's behavior and output unchanged.

## Subtasks

- [ ] Add the read-only re-check entry point over the shared evaluation.
- [ ] Require no decisions and resolve none.
- [ ] Render probes and grouping identically to the full plan.
- [ ] Confirm nothing is written.

## Acceptance Criteria

- [ ] The re-check completes with zero decisions supplied.
- [ ] Its capability outcomes match a full plan's for the same repository,
      asserted field by field.
- [ ] It renders the same probe evidence and grouping as the full plan.
- [ ] Running it leaves the repository byte-identical and writes no journal
      entry.
- [ ] A repository with no resolvable Profile produces a named error rather
      than a panic or an empty result.
- [ ] The full plan's output is unchanged, proven by the characterization
      corpus.
- [ ] `git status --porcelain` shows no path outside `internal/baseline/`,
      `internal/cli/`, and this task file.

## Context

- interface: `internal/baseline/profile_alignment.go`
- interface: `internal/cli/cli.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/baseline ./internal/cli -run '^TestCapabilityRecheck$' -count=1 -v | grep -q -- "--- PASS: TestCapabilityRecheck"`
  — expected: exit 0; the re-check needs no decisions and writes nothing.
- `go test ./internal/baseline -run '^TestCapabilityRecheckMatchesFullPlan$' -count=1 -v | grep -q -- "--- PASS: TestCapabilityRecheckMatchesFullPlan"`
  — expected: exit 0; outcomes match field by field.
- `go test ./internal/baseline -run '^TestBaselinePlanCharacterization$' -count=1 -v | grep -q -- "--- PASS: TestBaselinePlanCharacterization"` —
  expected: exit 0.
- `go test ./internal/baseline ./internal/cli -count=1` — expected: exit 0.

## References

- `_prd.md` → User Story 4; Core Features 6; Success Metrics (obtainable with
  zero decisions and matching full-plan outcomes).
- `_techspec.md` → API Contracts; Build Order 5.

