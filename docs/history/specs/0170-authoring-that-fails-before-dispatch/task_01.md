---
task: task_01
spec: 0170-authoring-that-fails-before-dispatch
status: completed
type: backend
complexity: high
---

# Task 01: Refuse an undeclared Governed Path at authoring

## Overview

`roundfix spec check` passes a Task that declares a Governed Path its Spec never
authorized, and the Daemon's changed-path audit refuses the commit only after
the Agent has done the work. This Task adds `SC-TOOLING-UNDECLARED`, which reads
what each pending Task declares and refuses before dispatch, and documents the
code where the other `spec check` codes are documented.

## Requirements

1. MUST declare `CodeToolingUndeclared = "SC-TOOLING-UNDECLARED"` in
   `internal/speccheck/constraints.go` beside the other tooling codes, with
   `SeverityError`.
2. MUST extend the Tooling authority row parsing in `constraints.go` so a row
   carries the backticked paths of its `bounded files:` list, read up to
   `Sanctioned regeneration:` or `Source:`, while `recordsBoundedFiles` keeps its
   current result for every existing row.
3. MUST, when the Spec has a Task Graph, collect for every Task that is neither
   `completed` nor `type: qa`: each Context entry of kind `interface` or
   `creates`; each Verification operand that resolves to an existing repository
   file; and each Verification operand the Task declares under `creates:`.
   Entries of kind `instruction` MUST NOT be collected.
4. MUST read Verification operands through one function exported from
   `internal/spec/collision.go`, and `declaredTaskTouches` MUST use that same
   function, so `spec.Collisions` returns exactly what it returns today.
5. MUST report `SC-TOOLING-UNDECLARED` for each collected path for which
   `GovernedPath` holds and that is absent from the selected authorization
   record's `paths:` and from its sanctioned regeneration outputs, resolved with
   `mechanicalRegenerationOutputs`. A Spec whose Tooling authority row cites no
   granted record MUST report every collected Governed Path.
6. MUST report `SC-TOOLING-UNDECLARED` for each collected Governed Path in the
   record's `paths:` that a present PRD or TechSpec Tooling authority row omits
   from its `bounded files:` list, and for each row path for which
   `GovernedPath` holds and that the record does not grant.
7. MUST give each finding the Task file and the line of the Context entry or
   Verification command as its first location, the record or row that omits
   the path as its second, and a Fix naming the path and every declaration to
   add it to.
8. MUST register `SC-TOOLING-UNDECLARED` in `stagedDetectors` in
   `internal/speccheck/coherence.go` as a Task-stage detector, and list it as
   skipped when the Spec has no `_tasks.md` or no `_prd.md`.
9. MUST add `SC-TOOLING-UNDECLARED` to `corpusFindingCodes` in
   `internal/docscontract/corpus_test.go`, with count `0` in
   `internal/docscontract/testdata/corpus-golden.json` and in the pin in
   `internal/spec/archive_layout_characterization_test.go`, leaving every other
   count and the `update` text unchanged.
10. MUST document `SC-TOOLING-UNDECLARED` in the Spec Consistency Check section
    of `.agents/skills/roundfix/SKILL.md` and as a row of the authoring-rule
    table in `.agents/skills/qa-gate/SKILL.md`, add the glossary entry
    **Undeclared Governed Path** to `CONTEXT.md`, and regenerate the mirrors
    with `make skills-sync`.
11. MUST NOT change `GovernedPath`, the governed set, or the changed-path audit
    in `internal/speccheck/mechanical.go`.

## Subtasks

- [ ] Carry the bounded paths on the parsed Tooling authority row.
- [ ] Export the Verification operand reader and reuse it in `spec.Collisions`.
- [ ] Add the detector, its stage registration and its skip.
- [ ] Update the corpus golden, its pin and the characterized code list.
- [ ] Document the code in both skills and the glossary, then run `make skills-sync`.

## Acceptance Criteria

- [ ] An undeclared `interface:` path, `creates:` path and Verification path are
      each refused, and the fully declared set passes.
- [ ] A path missing from one row, a row path the record does not grant, and a
      Spec without a granted record are each refused.
- [ ] `instruction:` entries, completed Tasks, ordinary paths and sanctioned
      regeneration outputs are not refused.
- [ ] The corpus golden and its pin record `0` for the new code and the active
      corpus has no errors.

## Context

- interface: `internal/speccheck/constraints.go`
- interface: `internal/speccheck/coherence.go`
- creates: `internal/speccheck/undeclared.go`
- creates: `internal/speccheck/undeclared_test.go`
- interface: `internal/spec/collision.go`
- interface: `internal/spec/collision_test.go`
- interface: `internal/docscontract/corpus_test.go`
- interface: `internal/docscontract/testdata/corpus-golden.json`
- interface: `internal/spec/archive_layout_characterization_test.go`
- interface: `.agents/skills/qa-gate/SKILL.md`
- interface: `skills/qa-gate/SKILL.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `CONTEXT.md`

## Verification

- `out="$(go test -count=1 -tags docscontract -v -run "^(TestUndeclaredGovernedContextPathIsRefused|TestUndeclaredGovernedCreatesPathIsRefused|TestUndeclaredGovernedVerificationPathIsRefused|TestDeclaredGovernedPathsPass|TestGovernedPathMissingFromABoundedFilesRowIsRefused|TestBoundedFilesRowPathTheRecordDoesNotGrantIsRefused|TestGovernedPathWithoutAGrantIsRefused|TestInstructionContextPathIsNotAudited|TestCompletedTaskIsNotAudited|TestOrdinaryPathNeedsNoDeclaration|TestSanctionedRegenerationOutputCountsAsDeclared|TestMissingTaskGraphListsTheUndeclaredGovernedPathSkip|TestUndeclaredGovernedPathIsSkippedBeforeTheTaskStage|TestTaskVerificationFilesReadsExistingOperands|TestArchiveLayoutCharacterizationPinsCorpusGoldenAfterSpec0095|TestCheckCorpusGolden|TestCheckActiveCorpusHasNoErrors)$" ./internal/speccheck ./internal/spec ./internal/docscontract 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestUndeclaredGovernedContextPathIsRefused TestUndeclaredGovernedCreatesPathIsRefused TestUndeclaredGovernedVerificationPathIsRefused TestDeclaredGovernedPathsPass TestGovernedPathMissingFromABoundedFilesRowIsRefused TestBoundedFilesRowPathTheRecordDoesNotGrantIsRefused TestGovernedPathWithoutAGrantIsRefused TestInstructionContextPathIsNotAudited TestCompletedTaskIsNotAudited TestOrdinaryPathNeedsNoDeclaration TestSanctionedRegenerationOutputCountsAsDeclared TestMissingTaskGraphListsTheUndeclaredGovernedPathSkip TestUndeclaredGovernedPathIsSkippedBeforeTheTaskStage TestTaskVerificationFilesReadsExistingOperands TestArchiveLayoutCharacterizationPinsCorpusGoldenAfterSpec0095 TestCheckCorpusGolden TestCheckActiveCorpusHasNoErrors; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done && grep -q '"SC-TOOLING-UNDECLARED": 0' internal/docscontract/testdata/corpus-golden.json && grep -q "SC-TOOLING-UNDECLARED" .agents/skills/roundfix/SKILL.md && grep -q "SC-TOOLING-UNDECLARED" .agents/skills/qa-gate/SKILL.md && grep -q "Undeclared Governed Path" CONTEXT.md && diff -r .agents/skills/roundfix skills/roundfix >/dev/null && diff -r .agents/skills/qa-gate skills/qa-gate >/dev/null` — expected: exit 0; before this Task none of the new tests exists and the golden carries no `SC-TOOLING-UNDECLARED` entry, so the command fails.

## References

- `_prd.md` → Goal 1; Core Feature 1; Success Metrics 1, 2, 4.
- `_techspec.md` → Undeclared Governed Paths; API Contract 1; Testing
  Approach 1; Vocabulary Contract; ADR-0093; ADR-0094; ADR-0117; ADR-0130;
  ADR-0149.

## Result

Implemented `SC-TOOLING-UNDECLARED` as a Task-stage error. The Tooling
authority row now retains its backticked `bounded files:` paths without
changing `recordsBoundedFiles`; the detector compares pending non-QA Task
Context and Verification declarations with the selected grant, every present
row, and sanctioned regeneration outputs. Verification operands now flow
through exported `spec.TaskVerificationFiles`, which `spec.Collisions` reuses
without changing its existing touch set.

Focused checks:

- Red check: the two initial regression tests failed to compile because
  `CodeToolingUndeclared` and `TaskVerificationFiles` did not exist.
- `go test -count=1 -run '^(TestUndeclaredGoverned|TestDeclaredGovernedPathsPass|TestGovernedPathMissingFromABoundedFilesRowIsRefused|TestBoundedFilesRowPathTheRecordDoesNotGrantIsRefused|TestGovernedPathWithoutAGrantIsRefused|TestInstructionContextPathIsNotAudited|TestCompletedTaskIsNotAudited|TestOrdinaryPathNeedsNoDeclaration|TestSanctionedRegenerationOutputCountsAsDeclared|TestMissingTaskGraphListsTheUndeclaredGovernedPathSkip|TestTaskVerificationFilesReadsExistingOperands|TestCollisionsFindsTheMeasuredShape|TestCollisionsReturnsEveryPairAndSharedPath|TestCollisionsLearnsPathFromDeclaredContext)$' ./internal/speccheck ./internal/spec` — passed.
- `go test -count=1 ./internal/spec ./internal/speccheck` — passed.
- `go test -count=1 -tags docscontract ./internal/docscontract` — passed.
- `make skills-sync` — passed; `diff -r` reports both required mirrors equal
  to their canonical skill directories.
- `make baseline-digests` — passed with `changed:false`; no derived digest
  changed.
- `git diff --check` — passed, and `internal/speccheck/mechanical.go` has no
  diff.
- `make verify-incremental` — not completed: the existing full test sweep
  attempted `api.github.com`, sandbox policy blocked the network request, and
  the escalation request was rejected. `go vet ./...` had completed before
  the network block. The Daemon-owned Verification command was not run.

Acceptance evidence:

- Undeclared `interface:`, `creates:`, and Verification paths are refused;
  `TestDeclaredGovernedPathsPass` covers the fully declared case.
- The focused suite refuses a path missing from one row, a row path missing
  from the grant, and a Governed Path with no granted record.
- The focused suite exempts `instruction:` entries, completed and QA Tasks,
  ordinary paths, and sanctioned regeneration outputs.
- The docscontract suite passed with `SC-TOOLING-UNDECLARED: 0` in the corpus
  golden and its archive-layout pin, and its active-corpus check reported no
  errors.
