---
task: task_02
spec: 0065-loop-order-and-verification-honesty
status: completed
type: backend
complexity: high
---

# Task 02: Refuse a Verification that cannot fail

## Overview

Spec 0060's `task_03` existed to prove an instruction-level gate fires. Its
Verification was `make verify` plus a clean `git status` — both of which pass
most easily when no work happened. It settled `completed` having run none of
its four cases.

This slice makes that shape refusable at authoring time, as a mechanical
`SC-VERIFY-WORK-INDEPENDENT` rule inside `roundfix spec check`, which already
fails `make verify`. A skill instruction would be advice to an Agent free to
ignore it, which is exactly how the defect happened.

## Requirements

1. MUST add `SC-VERIFY-WORK-INDEPENDENT`, reported through the existing
   `spec check` finding surface with its file, line, and fix line.
2. MUST decide the property from the Task's **declared** Verification commands,
   never from prose or intent, keeping ADR-0093's detection boundary.
3. MUST refuse a Verification composed only of repository-wide gates and
   working-tree cleanliness checks, because that sequence passes most easily
   when nothing was done.
4. MUST accept a Verification that contains a repository-wide gate **and** at
   least one command asserting the Task's own effect. The rule targets the
   composition, not the presence of any particular command.
5. MUST leave every active and archived Spec in the corpus checking exactly as
   it does today, asserted rather than assumed.
6. MUST keep `TestCheckCorpusBudget` passing, so the sweep stays within budget.
7. MUST NOT change the loop order statements or add the divergence rule; those
   are task_01 and task_04.

## Subtasks

- [ ] Add the rule and its finding.
- [ ] Replay Spec 0060's `task_03` as a fixture and assert refusal.
- [ ] Build the false-positive table over legitimate Verifications.
- [ ] Assert corpus non-regression and the budget test.

## Acceptance Criteria

- [ ] A Verification of only `make verify` plus a clean-tree check is refused.
- [ ] Spec 0060's `task_03`, replayed as written, is refused by this rule.
- [ ] A Verification with a repository gate plus an effect-asserting command
      passes.
- [ ] A Verification of only effect-asserting commands passes.
- [ ] Every Spec in the existing corpus checks as it does today.
- [ ] `TestCheckCorpusBudget` passes.

## Context

- interface: `internal/speccheck/citations.go`
- instruction: `docs/adr/0093-spec-consistency-is-checked-by-citation-never-by-inference.md`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/speccheck -count=1 -run 'WorkIndependent|Verify' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the rule's tests ran and passed.
- `go test -count=1 -parallel=1 ./internal/speccheck -run '^TestCheckCorpusBudget$'`
  — expected: exit 0; the sweep stays within budget.
- `go test -parallel 16 ./... 2>&1 | grep -q "^FAIL" && exit 1 || exit 0`
  — expected: exit 0.

## References

- `_prd.md` → Core Feature 3; Success Metric 2.
- `_techspec.md` → Interfaces; Build Order 2; Risks & Considerations.
- ADR-0093.

## Result

### Implementation

- Added `SC-VERIFY-WORK-INDEPENDENT` to the existing Spec Consistency Check
  finding surface. The finding points to the first declared Verification
  command and tells the author to add a command that asserts the Task's own
  effect.
- The detector classifies only `spec.Task.Verification` commands. It refuses a
  non-completed Task only when every declared command is a repository-wide
  Make/Go gate or an explicit working-tree cleanliness check; any other
  declared command keeps the Verification accepted.
- Completed Tasks remain historical evidence and are not retroactively
  reported by `Check`. The active-plus-archived corpus golden now includes the
  new code with zero findings in both sets.
- Added an authoring-time replay of Spec 0060 `task_03` with its exact
  Verification commands, a table covering refusal and false-positive cases,
  and text-rendering assertions for code, file, line, and fix output.
- No loop-order statement or divergence rule changed in this slice.

### Focused checks

- Red signal: `rtk env GOCACHE=/private/tmp/roundfix-task02-gocache go test
  ./internal/speccheck -count=1 -run
  '^TestCheckReplay0060Task03RefusesWorkIndependentVerification$'` — exit 1
  before the detector existed; the replay returned no
  `SC-VERIFY-WORK-INDEPENDENT` finding.
- `rtk env GOCACHE=/private/tmp/roundfix-task02-gocache go test
  ./internal/speccheck -count=1 -run
  '^(TestWorkIndependentVerificationRefusesOnlyWorkIndependentCommands|TestCheckReplay0060Task03RefusesWorkIndependentVerification)$'`
  — exit 0 after implementation.
- `rtk env GOCACHE=/private/tmp/roundfix-task02-gocache go test
  ./internal/speccheck -count=1 -run '^TestCheckCorpusGolden$'` — exit 0; the
  active and archived finding counts remain unchanged and the new code is zero
  in both maps.
- `rtk env GOCACHE=/private/tmp/roundfix-task02-gocache go test
  ./internal/speccheck -count=1` — exit 0; the package suite passed. The budget
  test intentionally skipped because this was not its dedicated invocation.
- `rtk gofmt -l internal/speccheck/verification.go
  internal/speccheck/verification_test.go internal/speccheck/citations.go
  internal/speccheck/constraints_characterization_test.go` — exit 0 with no
  output after the final code and test edits.
- `rtk git diff --check` — exit 0 after the final Result update.
- The commands under this Task's `## Verification` were not run. The Daemon
  owns them, including the dedicated `TestCheckCorpusBudget` timing command.

### Acceptance evidence

- Only `make verify` plus `git status --porcelain` is refused: covered by the
  `repository_gate_and_clean_tree` table case.
- Spec 0060 `task_03` is refused as an authoring-time replay: covered by
  `TestCheckReplay0060Task03RefusesWorkIndependentVerification`, including the
  exact `task_03.md:44` location and rendered fix line.
- A repository gate plus an effect assertion passes: covered by the
  `repository_gate_plus_focused_effect_assertion` table case.
- Effect-only commands pass: covered by the `only_effect_assertions` table
  case. A repository-wide `go test` carrying a declared `-run` effect selector
  is also accepted.
- Every active and archived Spec checks with its prior finding counts: covered
  by `TestCheckCorpusGolden` and the zero entries for the new code.
- `TestCheckCorpusBudget` has not been run in this Agent turn because its exact
  command is Daemon-owned Verification. No timing verdict is claimed here.
