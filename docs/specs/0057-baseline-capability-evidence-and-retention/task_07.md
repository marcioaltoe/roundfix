---
task: task_07
spec: 0057-baseline-capability-evidence-and-retention
status: completed
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

## Result

Implemented the fourth Profile-divergence prompt outcome. Choosing repository
remediation now returns the existing action-required exit contract without
mutation, renders one remediation line for every unsatisfied divergence, and
prints an exact shell-quoted `roundfix baseline capabilities check` command for
the selected Profile and repository. Its machine result uses the distinct
`remediation` category, while decline retains its existing `decision` category,
message, next action, and exit. The adaptation prompt now labels the path as
`removal-only`.

Focused implementation evidence:

- Pre-change inspection found only the three existing Profile-divergence
  choices and no remediate branch; the Task 06 re-check command was present and
  callable through `roundfix baseline capabilities check`.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task07-gocache rtk go test
  ./internal/cli -run '^(TestDivergencePrompt.*|TestBaselineHumanProfileAdaptation|TestHumanBaselineAdoption|TestHumanBaselineUpdate)$'
  -count=1` passed 7 tests in 1 package.
- `rtk git diff --check` exited 0 with no diagnostics.
- `rtk git -c core.fsmonitor=false status --short` listed only this task file,
  `internal/cli/baseline_human.go`, and
  `internal/cli/baseline_human_test.go`.

Acceptance evidence:

- `TestDivergencePromptRemediateOutcome` snapshots the repository tree before
  selection and proves it is byte-identical afterwards; it also proves the new
  outcome preserves the action-required exit code.
- The same test resolves a mixed alignment fixture first, then requires the
  remediation output to contain every returned divergence ID with its exact
  next action and the exact Profile/repository capability re-check command.
- `TestDivergencePromptJournalsDistinctly` reopens the JSON records and proves
  remediation journals as category `remediation` while decline remains
  category `decision`; the serialized records differ.
- The prompt assertion requires the adaptation label to contain
  `removal-only`.
- The existing Profile adaptation test still exercises reviewed adaptation
  through a ready Plan and the unchanged decline result/no-write behavior;
  focused human adoption and update tests also remain green. The profile-change
  branch and adaptation selection indices remain unchanged.
- The changed-path inspection contains no path outside `internal/cli/` and
  this task file, which is narrower than the Task's allowed path set.

The Daemon-owned commands under `## Verification` were not run in this Agent
turn.
