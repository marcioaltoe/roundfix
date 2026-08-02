---
task: task_06
spec: 0057-baseline-capability-evidence-and-retention
status: completed
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
- `go test ./internal/baseline ./internal/cli -run TestCapabilityRecheck -count=1`
  — expected: exit 0; the re-check needs no decisions and writes nothing.
- `go test ./internal/baseline -run TestCapabilityRecheckMatchesFullPlan -count=1`
  — expected: exit 0; outcomes match field by field.
- `go test ./internal/baseline -run TestBaselinePlanCharacterization -count=1` —
  expected: exit 0.
- `go test ./internal/baseline ./internal/cli -count=1` — expected: exit 0.

## References

- `_prd.md` → User Story 4; Core Features 6; Success Metrics (obtainable with
  zero decisions and matching full-plan outcomes).
- `_techspec.md` → API Contracts; Build Order 5.

## Result

Implementation did not start because this Task's required Task 04 dependency
is marked completed without its implementation. Task 04's `## Result` records
that probe rendering and blocking, advisory, and informational grouping were
not implemented because its changed-path criterion excluded the sole renderer
in `internal/cli/baseline_human.go`. The current full-plan renderer still
prints only each divergence's message and next action, derives only blocking or
advisory labels, and never reads `ProfileDivergence.Probe`, `.Evidence`, or
`.Requirement` for the required presentation.

Task 06 therefore cannot both render the same probe evidence and requirement
grouping as the full plan and leave the full plan's behavior and output
unchanged. Implementing the missing renderer here would absorb Task 04's slice,
contrary to the bounded Task contract.

Focused inspection evidence:

- `rtk git -c core.fsmonitor=false status --short` showed only the Daemon-owned
  `pending` to `in_progress` edit in this task file before inspection.
- `rtk git -c core.fsmonitor=false show --stat --oneline 6192029` showed that
  the Task 04 commit changed only `task_04.md`.
- `rtk rg -n
  "TestDivergenceRendersProbe|TestDivergenceGroupsByRequirement|DivergenceGroup|informational"
  internal` found no Task 04 implementation or regression test.
- Reading `renderBaselineProfileAlignment` confirmed it has no probe renderer
  and no informational requirement group.
- Repository search found no other shared alignment renderer that Task 06 can
  reuse.

Acceptance evidence:

- The zero-decision re-check, field-by-field equivalence, read-only behavior,
  named missing-Profile error, and CLI rendering were not implemented or
  claimed because the required presentation dependency is absent.
- Full-plan behavior and characterization artifacts remain unchanged because
  no implementation code, test, or golden file was edited.
- The changed-path postflight contains no path outside this task file.

Required follow-up: reopen Task 04 with `internal/cli/baseline_human.go` and its
canonical tests in scope, implement and verify its declared presentation
contract, then run Task 06 fresh against that shared renderer.

Daemon Verification was not run in this Agent turn.
