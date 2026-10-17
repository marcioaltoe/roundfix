---
task: task_03
spec: 0277-a-run-never-reaches-the-judge-and-an-amended-plan-keeps-its-settled-tasks
status: completed
type: backend
complexity: high
---

# Task 03: Carry-forward keeps a settled Task whose Binding Text an amendment left unchanged

## Overview

Task Carry-Forward treats the PRD, the TechSpec and the Task Graph manifest as
inputs of every Task, so any amendment of them refuses every settled Task of
the Spec, which happened three times on 2026-10-17 (`_techspec.md` → Current
behavior). This Task adds the Binding Text of a Task and makes carry-forward
compare only that text for the three planning files, so an amendment that
leaves a settled Task's own contract unchanged carries it, while every other
moved input and every changed Binding Text still refuses the whole set. It is
verifiable on its own through planner and command tests on temporary
repositories.

## Requirements

1. MUST answer the operator observations of 2026-10-17 cited in `_prd.md` →
   Sources 4 (Spec 0272 after pull request #551, Spec 0274 after pull request
   #559) by implementing API Contracts 6 and 7 of `_techspec.md` (ADR-0284).
2. MUST add `spec.PlanningBinding` in the new
   `internal/spec/carry_forward_binding.go` with exactly the rules of API
   Contract 6: for `_tasks.md`, the Task's graph node (`id`, `file`, `needs` in
   order) and its projection table row; for `_prd.md` and `_techspec.md`, the
   sections that the units of the Task's References bullet for that file name,
   with headings of level 2 to 4 outside fenced blocks, case-insensitive
   singular or plural matching and numbers, ranges, commas and `and` dropped;
   `whole` when the node or row is missing, the Task has no bullet for the
   file, or any unit names no heading.
3. MUST make `inspectCarryForwardCandidate` tolerate a moved `_prd.md`,
   `_techspec.md` or `_tasks.md` of the Task's own Spec only when the bindings
   at the settlement parent and on the staged target, both read with the Task
   file at the parent, are equal and neither is `whole`; a tolerated path
   leaves `MovedInputs` and joins the new `spec.CarryForward.AmendedInputs`
   (`amendedInputs` in JSON, `[]` when empty), and the reconcile text prints
   `  amended-inputs: [<paths>]` after `  moved-inputs:`. The Task file,
   `AGENTS.md`, `CONTEXT.md`, the `implement-task` skill and every Context path
   MUST still refuse on any change, and every refusal text and recovery MUST
   stay as today. The rule serves the Reconcile Command, the Implement
   Command's Preflight and a Delivery Retry, and the whole remaining set still
   refuses when one Task refuses (ADR-0170).
4. MUST add the tests named in Verification, each calling `t.Parallel()`
   first: in the new `internal/spec/carry_forward_binding_test.go`,
   `TestPlanningBindingReadsTheTaskNodeAndRow` and
   `TestPlanningBindingReadsTheSectionsTheReferencesName` (plural forms,
   numbered units, a `:` unit, a fenced heading, an unresolved unit and a
   missing bullet each give the stated binding or `whole`); in the new
   `internal/cli/carryforward_amendment_test.go`,
   `TestCarryForwardCarriesATaskOverAnAmendmentOutsideItsBindingText` (a Run
   with one settled Task whose References name neither Project Constraints nor
   another Task, then an amendment that adds a Tooling authority row to both
   planning files and a corrective Task node and row to the manifest: the Task
   is carried, `amendedInputs` lists the three files and the provenance is
   written as today) and
   `TestCarryForwardRefusesATaskWhoseBindingTextChanged` (an amendment inside a
   section the Task names refuses the whole set with
   `Task task_01 declared input(s) moved: <path>`). Every new test compiles
   for Windows: no Unix-only system call outside a `unix`-constrained file.
   `TestRetryRefusalNamesTheAmendingCommitAndTheRecovery` and every other
   existing test keep their names and expected text.
5. MUST document the Binding Text rule, the amended inputs and their text and
   JSON lines in `docs/user-guide/commands/reconcile.md` and
   `.agents/skills/roundfix/references/reconcile.md`, replacing the sentence
   that says a Task whose Spec moved is never replayed, searching the guides
   and skill copies for that old wording. No text inside any
   `### QA settlement` section may change.
6. MUST add **Binding Text** to `CONTEXT.md` and revise **Task Carry-Forward**
   through domain-modeling, with ADR-0284's link.
7. MUST run `make skills-sync`, raise and record the Roundfix skill version
   with `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   and re-record the Behavior Surface Record with
   `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`,
   because this Task edits a user guide page and a command source. No record
   or mirror is hand-edited.

## Subtasks

- [ ] Add `PlanningBinding` and its tests.
- [ ] Use it in carry-forward and report the amended inputs.
- [ ] Add the carry and refusal command tests.
- [ ] Document the rule in the guide, the reference and the glossary.
- [ ] Sync the mirrors, record the version and the Behavior Surface Record.

## Acceptance Criteria

- [ ] An amendment that changes only text outside a settled Task's Binding Text carries the Task and lists the amended files.
- [ ] An amendment of a section the Task names, of its own Task file or of a non-planning input still refuses the whole set with today's reason.
- [ ] The reconcile text and JSON show the amended inputs, and every earlier refusal text is unchanged.

## Consumers

- The carry-forward candidate and its JSON — consumers: `./internal/spec`, `./internal/cli`

## Context

- instruction: `docs/adr/0284-carry-forward-binds-a-task-to-the-planning-text-it-references.md`
- instruction: `docs/adr/0170-a-task-already-completed-on-its-target-is-nothing-to-carry.md`
- instruction: `internal/cli/deliver_retry_amendment_test.go`
- instruction: `internal/cli/carryforward_test.go`
- interface: `internal/spec/spec.go`
- creates: `internal/spec/carry_forward_binding.go`
- creates: `internal/spec/carry_forward_binding_test.go`
- interface: `internal/cli/carryforward.go`
- interface: `internal/cli/reconcile.go`
- creates: `internal/cli/carryforward_amendment_test.go`
- interface: `docs/user-guide/commands/reconcile.md`
- interface: `.agents/skills/roundfix/references/reconcile.md`
- interface: `skills/roundfix/references/reconcile.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestPlanningBindingReadsTheTaskNodeAndRow|TestPlanningBindingReadsTheSectionsTheReferencesName)$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestPlanningBindingReadsTheTaskNodeAndRow TestPlanningBindingReadsTheSectionsTheReferencesName; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; out="$(go test -count=1 -v -run '^(TestCarryForwardCarriesATaskOverAnAmendmentOutsideItsBindingText|TestCarryForwardRefusesATaskWhoseBindingTextChanged|TestRetryRefusalNamesTheAmendingCommitAndTheRecovery)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestCarryForwardCarriesATaskOverAnAmendmentOutsideItsBindingText TestCarryForwardRefusesATaskWhoseBindingTextChanged TestRetryRefusalNamesTheAmendingCommitAndTheRecovery; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; for f in docs/user-guide/commands/reconcile.md .agents/skills/roundfix/references/reconcile.md; do grep -qF -- 'amended-inputs' "$f" || { printf 'missing amended-inputs in %s\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'Binding Text' || { printf 'missing Binding Text in %s\n' "$f" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Binding Text**:' || { printf 'missing glossary term Binding Text\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Task Carry-Forward**:' || { printf 'missing glossary term Task Carry-Forward\n' >&2; exit 1; }; make skills-sync-check || exit 1; go test -count=1 ./skills || exit 1; go test -count=1 -run 'CarryForward|Retry|Reconcile' ./internal/cli || exit 1; go test -count=1 -run 'CarryForward|PlanningBinding' ./internal/spec || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the binding tests do not exist, so the command fails at its first pass check; afterwards the binding, carry-forward, retry, documentation, glossary, skill version and Behavior Surface Record checks pass.

## References

- `_prd.md` → Sources; Goals; User Story 4; Core Feature 7; Core Feature 9; Success Metric 4; Glossary
- `_techspec.md` → API Contract 6; API Contract 7; API Contract 9; Invariant 6; Current behavior; Testing Approach; Build Order 3
- ADR-0284; ADR-0170; ADR-0158; ADR-0280; ADR-0189; ADR-0256; ADR-0259


## Result

Implemented the Task 03 slice for Daemon Verification. The Daemon-owned
`status: in_progress` is preserved; the Task Graph and every other Task file
are unchanged. The initial worktree change was the Daemon's status edit in
this Task file.

### Acceptance evidence

| Acceptance criterion | Implementation and focused evidence |
| --- | --- |
| An amendment outside a settled Task's Binding Text carries it and lists the amended files. | `PlanningBinding` extracts the ordered node fields and projection row, or the referenced level 2–4 sections. The new carry command test adds Tooling authority to both planning files and a corrective node and row to the manifest, then proves all three paths appear in `amendedInputs`, the implementation is carried, the amendment survives, and the existing Run/commit provenance is written. |
| A changed named section, Task file or non-planning input refuses the whole set with today's reason. | The new refusal command test exercises a named TechSpec section, the Task file, `AGENTS.md`, `CONTEXT.md`, the implementation skill, an ordinary Context path, an explicit Context path to the Task's own planning file amended outside its named section, and another Spec's planning file. Each asserts the exact `Task task_01 declared input(s) moved: <path>` reason, whole-set refusal, unchanged target HEAD, and absence of either Task's implementation. Existing retry-refusal and staged-carry regressions also passed. |
| Reconcile text and JSON show amended inputs and earlier refusal text is preserved. | The carry command test checks JSON's three amended paths and the adjacent text lines `moved-inputs` then `amended-inputs`. Refusal cases check `"amendedInputs":[]`. Every candidate constructor initializes the list. The existing retry refusal tests retain their names and expectations, including the amending commits and recovery commands. |

The parser tests cover dependency order, another Task's node and row, manifest
prose, changed own row, missing node/row, fenced rows and headings, singular
and plural forms, numbered units, ASCII and Unicode ranges, colon units,
unresolved units, missing or fenced bullets, wrapped bullets, section order,
repeated headings, level boundaries and Git text line endings. Both bindings
are read using the settlement parent's Task declaration. Explicit Context
paths always retain whole-file comparison.

### Focused checks

- Before implementation, `go test ./internal/spec -run '^TestPlanningBindingReadsTheTaskNodeAndRow$' -count=1` failed to compile because `PlanningBinding` did not exist.
- Before planner integration, `go test ./internal/cli -run '^TestCarryForwardCarriesATaskOverAnAmendmentOutsideItsBindingText$' -count=1` reproduced the refusal naming all three amended planning files.
- `go test ./internal/spec -run '^TestPlanningBindingReads' -count=1` — exit 0 after the final parser/test edits.
- `go test ./internal/cli -run '^(TestCarryForwardCarriesATaskOverAnAmendmentOutsideItsBindingText|TestCarryForwardRefusesATaskWhoseBindingTextChanged|TestRetryRefusal.*|TestCarryForwardInputsResolveAgainstStagedCarries|TestInspectSpecCarryForwards)$' -count=1` — exit 0 after the final command/test edits.
- `GOOS=windows GOARCH=amd64 go test -c ./internal/spec -o /tmp/roundfix-0277-task03-spec.test.exe` — exit 0.
- `GOOS=windows GOARCH=amd64 go test -c ./internal/cli -o /tmp/roundfix-0277-task03-cli.test.exe` — compilation blocked by existing Unix-only tests assigned to Task 02: `deliver_item_probe_report_test.go` (`syscall.Stat_t`, `syscall.Getpgrp`), `detach_fixture_group_other_test.go` (`syscall.Kill`), `detach_test.go` (`store.SignalProvenProcess`, `syscall.Wait4`), and `implement_detach_teardown_test.go` (Unix process APIs). These files were not changed. The new Task 03 tests use portable APIs and call `t.Parallel()` first.
- `git diff --check` — exit 0. A read-only scope inspection found all 15 changed paths in this Task's declared Context or this Task file, preserved status and QA settlement bytes, matching canonical/embedded Roundfix skill files, and no old blanket Spec-moved wording in the guides or either skill copy.

### Documentation and generated records

The Reconcile guide and canonical reference describe Binding Text, conservative
whole-file fallback, strict Context inputs, `amended-inputs` text, and
`amendedInputs` JSON including the empty list. `CONTEXT.md` adds Binding Text
and revises Task Carry-Forward through domain-modeling, linked to ADR-0284.
No text inside a QA settlement section changed.

- `make skills-sync` — exit 0; mirrors generated from canonical skills.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions` — exit 0; generated Roundfix version `0.1.91` in both skill copies and its version record.
- `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1` — exit 0; regenerated the Reconcile guide fingerprint.
- `make baseline-digests` — exit 0, `changed:false`; no derived Baseline files changed.

No record or mirror was hand-edited. Declared Verification, selected repository
Verification and incremental Verification were left to the Daemon. No commit,
push or pull request was created. The CLI Windows compilation limit above is
follow-up evidence for Task 02, outside this diff.
