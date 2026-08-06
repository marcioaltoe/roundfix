---
task: task_03
spec: 0065-loop-order-and-verification-honesty
status: completed
type: backend
complexity: medium
---

# Task 03: Refuse a Task that contradicts itself

## Overview

Spec 0060's `task_03` also carried requirements that contradicted each other,
so the rehearsal could not be performed as written. An Agent Session spent a
turn discovering that, and the Task settled `completed` anyway.

Two rules close it: one for requirements that cannot all hold, one for a
rehearsal Task that never declares what it must exercise.

## Requirements

1. MUST add `SC-REQUIREMENT-CONTRADICTORY`, reporting one requirement
   forbidding a state another requirement needs.
2. MUST decide contradiction from declared MUST and MUST NOT clauses over the
   same named subject, and MUST report nothing when the subject cannot be
   identified — silence is the correct answer to an undecidable pair, per
   ADR-0093.
3. MUST add `SC-REHEARSAL-UNDECLARED`, refusing a Task whose stated purpose is
   proving a gate fires but which declares no cases and no observation for
   them.
4. MUST define the authored section a rehearsal Task uses to declare each case
   and how it is observed, so the rule reads a declaration rather than guessing
   intent.
5. MUST leave every active and archived Spec checking exactly as it does today.
6. MUST keep `TestCheckCorpusBudget` passing.
7. MUST reuse the requirement parsing task_02 introduces rather than adding a
   second parser.

## Subtasks

- [ ] Add both rules and their findings.
- [ ] Replay Spec 0060's `task_03` and assert both refusals.
- [ ] Assert an undecidable pair reports nothing.
- [ ] Assert corpus non-regression and the budget test.

## Acceptance Criteria

- [ ] A MUST and a MUST NOT over the same subject are refused.
- [ ] A pair whose subject cannot be identified reports nothing.
- [ ] Spec 0060's `task_03`, replayed as written, is refused for its
      contradictory requirements.
- [ ] A rehearsal Task with no declared cases is refused.
- [ ] A rehearsal Task declaring its cases and their observation passes.
- [ ] Every Spec in the existing corpus checks as it does today.
- [ ] `TestCheckCorpusBudget` passes.

## Context

- interface: `internal/speccheck/citations.go`
- instruction: `docs/adr/0093-spec-consistency-is-checked-by-citation-never-by-inference.md`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/speccheck -count=1 -run 'Contradict|Rehearsal' -v | grep -q -- "--- PASS"`
  — expected: exit 0; both rules' tests ran and passed.
- `go test -count=1 -parallel=1 ./internal/speccheck -run '^TestCheckCorpusBudget$'`
  — expected: exit 0.
- `go test -parallel 16 ./... 2>&1 | grep -q "^FAIL" && exit 1 || exit 0`
  — expected: exit 0.

## References

- `_prd.md` → Core Features 4 and 5; Success Metric 2.
- `_techspec.md` → Interfaces; Build Order 3.
- ADR-0093.

## Result

### Implementation

- Extended the canonical Task document parser with numbered Requirements and
  `## Rehearsal Cases` bullets. Both new rules consume these declarations from
  `spec.Task`; `speccheck` does not add a second Markdown parser.
- Added `SC-REQUIREMENT-CONTRADICTORY`. It compares only declared `MUST` and
  `MUST NOT` clauses, normalizes their written words, and reports the first
  shared named state. A pair with no identifiable shared state remains silent.
- Added `SC-REHEARSAL-UNDECLARED`. A Task title that declares rehearsing or
  proving a gate must provide at least one complete
  `- Case: <case>; Observation: <observation>` entry under
  `## Rehearsal Cases`; ordinary Tasks remain unaffected.
- Integrated both findings into the existing authoring check for non-completed
  Tasks, with the declaring requirement lines or Task title as locations.
- Added unit coverage, the Spec 0060 replay assertion, and zero-count corpus
  characterization entries for both rule identifiers.

### Focused checks

- Red signal after the Verification feedback: the Daemon-selected pattern had
  no matching tests, and the focused test-first command failed to compile
  because the declaration fields and both detectors did not exist.
- `rtk env GOCACHE=/private/tmp/roundfix-task03-gocache go test
  ./internal/spec ./internal/speccheck -count=1 -run
  'TaskDocumentDeclarations|ContradictoryRequirements|UndeclaredRehearsal|Replay0060Task03RefusesContradictory'`
  — exit 0 after implementation.
- `rtk env GOCACHE=/private/tmp/roundfix-task03-gocache go test
  ./internal/speccheck -count=1 -run
  '^(TestCheckCorpusGolden|TestCheckActiveCorpusHasNoErrors)$'` — exit 0; all
  active and archived corpus counts remain unchanged and both new codes have
  zero findings.
- `rtk env GOCACHE=/private/tmp/roundfix-task03-gocache go test
  ./internal/spec ./internal/speccheck -count=1` — exit 0; both affected
  package suites passed after the final implementation edit.
- No command from this Task's `## Verification` section was rerun. The Daemon
  owns the configured retry, including the dedicated budget invocation.

### Acceptance evidence

- Same-subject `MUST`/`MUST NOT` refusal is covered by
  `TestContradictoryRequirementsRefusesSameNamedSubject/same_subject_is_required_and_forbidden`.
- Undecidable-pair silence is covered by
  `TestContradictoryRequirementsRefusesSameNamedSubject/subject_cannot_be_identified`.
- The Spec 0060 replay produces both `SC-REQUIREMENT-CONTRADICTORY` at its
  declared commit clauses and `SC-REHEARSAL-UNDECLARED` at its title, covered
  by
  `TestCheckReplay0060Task03RefusesContradictoryRequirementsAndUndeclaredRehearsal`.
- A rehearsal without cases is refused by
  `TestUndeclaredRehearsalRequiresCasesAndObservations/rehearsal_has_no_cases`.
- Complete case and observation declarations pass under
  `TestUndeclaredRehearsalRequiresCasesAndObservations/rehearsal_declares_cases_and_observations`;
  an incomplete entry has a negative companion case.
- `TestCheckCorpusGolden` and `TestCheckActiveCorpusHasNoErrors` both passed in
  the focused corpus command.
- `TestCheckCorpusBudget` was not invoked locally because it is a declared
  Daemon Verification command; the Daemon retry owns that acceptance evidence.
