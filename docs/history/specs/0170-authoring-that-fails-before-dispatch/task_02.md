---
task: task_02
spec: 0170-authoring-that-fails-before-dispatch
status: completed
type: backend
complexity: medium
---

# Task 02: Report a CLI change that ships without its guide

## Overview

A Task changes a command's flags, output or exit codes while the skill or guide
that describes it stays behind, and the pre-PR review then forces a corrective
Task. This Task adds the gap `SC-CLI-UNDOCUMENTED`, reported when a pending Task
names a CLI surface and neither it nor a Task it depends on names a guide.

## Requirements

1. MUST declare `CodeCLIUndocumented = "SC-CLI-UNDOCUMENTED"` in
   `internal/speccheck/coherence.go` beside the other Task Graph codes, with
   `SeverityGap`, and register it in `stagedDetectors` as a Task-stage
   detector.
2. MUST treat as a CLI surface a Context entry of kind `interface` or
   `creates` naming a `.go` file under `internal/cli/` or `cmd/roundfix/` whose
   name does not end in `_test.go`, or naming `internal/cli/cli_test.go`.
3. MUST treat as a guide a Context entry of kind `interface` or `creates`
   under `.agents/skills/`, `skills/`, `docs/user-guide/` or `docs/agents/`.
4. MUST report one `SC-CLI-UNDOCUMENTED` per pending non-QA Task that names a
   CLI surface when neither that Task nor any Task in its transitive `needs`
   names a guide, locating the Task file and the surface entry's line.
5. MUST call the detector from `Check` in `internal/speccheck/constraints.go`
   and list the code as skipped when the Spec has no `_tasks.md` or no
   `_prd.md`.
6. MUST add `SC-CLI-UNDOCUMENTED` to `corpusFindingCodes` in
   `internal/docscontract/corpus_test.go`, with count `0` in
   `internal/docscontract/testdata/corpus-golden.json` and in the pin in
   `internal/spec/archive_layout_characterization_test.go`, leaving every other
   count and the `update` text unchanged.
7. MUST document `SC-CLI-UNDOCUMENTED` in the Spec Consistency Check section
   of `.agents/skills/roundfix/SKILL.md` and as a row of the authoring-rule
   table in `.agents/skills/qa-gate/SKILL.md`, add the glossary entry
   **Undocumented CLI Surface** to `CONTEXT.md`, and regenerate the mirrors
   with `make skills-sync`.

## Subtasks

- [ ] Add the code, its stage registration and its skip.
- [ ] Add the detector in `internal/speccheck/surface.go` and wire it into `Check`.
- [ ] Update the corpus golden, its pin and the characterized code list.
- [ ] Document the code in both skills and the glossary, then run `make skills-sync`.

## Acceptance Criteria

- [ ] A CLI surface without a guide is reported, and `--strict` makes it an error.
- [ ] A guide in the Task or in a Task it depends on clears it; a guide only in
      a dependent Task does not.
- [ ] `internal/cli/cli_test.go` counts as a surface and another CLI test file
      does not.

## Context

- interface: `internal/speccheck/coherence.go`
- interface: `internal/speccheck/constraints.go`
- creates: `internal/speccheck/surface.go`
- creates: `internal/speccheck/surface_test.go`
- interface: `internal/docscontract/corpus_test.go`
- interface: `internal/docscontract/testdata/corpus-golden.json`
- interface: `internal/spec/archive_layout_characterization_test.go`
- interface: `.agents/skills/qa-gate/SKILL.md`
- interface: `skills/qa-gate/SKILL.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `CONTEXT.md`

## Verification

- `out="$(go test -count=1 -tags docscontract -v -run "^(TestCLISurfaceWithoutAGuideIsReported|TestStrictPromotesAnUndocumentedCLISurface|TestCLISurfaceNamingItsGuidePasses|TestCLISurfaceWhoseDependencyNamesTheGuidePasses|TestGuideOnlyInADependentTaskDoesNotCount|TestCLIContractTestIsASurface|TestOrdinaryCLITestIsNotASurface|TestMissingTaskGraphListsTheCLISurfaceSkip|TestArchiveLayoutCharacterizationPinsCorpusGoldenAfterSpec0095|TestCheckCorpusGolden|TestCheckActiveCorpusHasNoErrors)$" ./internal/speccheck ./internal/spec ./internal/docscontract 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestCLISurfaceWithoutAGuideIsReported TestStrictPromotesAnUndocumentedCLISurface TestCLISurfaceNamingItsGuidePasses TestCLISurfaceWhoseDependencyNamesTheGuidePasses TestGuideOnlyInADependentTaskDoesNotCount TestCLIContractTestIsASurface TestOrdinaryCLITestIsNotASurface TestMissingTaskGraphListsTheCLISurfaceSkip TestArchiveLayoutCharacterizationPinsCorpusGoldenAfterSpec0095 TestCheckCorpusGolden TestCheckActiveCorpusHasNoErrors; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done && grep -q '"SC-CLI-UNDOCUMENTED": 0' internal/docscontract/testdata/corpus-golden.json && grep -q "SC-CLI-UNDOCUMENTED" .agents/skills/roundfix/SKILL.md && grep -q "SC-CLI-UNDOCUMENTED" .agents/skills/qa-gate/SKILL.md && grep -q "Undocumented CLI Surface" CONTEXT.md && diff -r .agents/skills/roundfix skills/roundfix >/dev/null && diff -r .agents/skills/qa-gate skills/qa-gate >/dev/null` — expected: exit 0; before this Task none of the new tests exists and the golden carries no `SC-CLI-UNDOCUMENTED` entry, so the command fails.

## References

- `_prd.md` → Goal 2; Core Feature 2; Success Metrics 3-4.
- `_techspec.md` → CLI changes without their guide; API Contract 2; Testing
  Approach 2; Vocabulary Contract; ADR-0093; ADR-0094; ADR-0117.

## Result

Implemented the Task-stage `SC-CLI-UNDOCUMENTED` gap detector. It checks only
pending non-QA Tasks, recognizes the declared CLI source and contract-test
surfaces, walks transitive `needs` ancestors for a declared guide, and reports
the Task file at the CLI Context entry's line. Missing PRD and Task Graph inputs
now record the detector as skipped.

The stable code is registered in the staged detector set and the corpus code
list. The active-corpus golden and archive-layout pin each carry count `0`
without changing any other count or the golden's `update` text. The Roundfix
and QA-gate source skills document the rule, `CONTEXT.md` defines
**Undocumented CLI Surface**, and `make skills-sync` regenerated both mirrors.
The repository-required `make baseline-digests` reported no derived changes.

Focused checks:

- Red signal: the focused new-test command initially failed to build because
  `speccheck.CodeCLIUndocumented` did not exist.
- `GOCACHE=/tmp/roundfix-task02-gocache go test -count=1 ./internal/speccheck`
  passed after the final detector and test edits.
- `GOCACHE=/tmp/roundfix-task02-gocache go test -count=1 -run
  '^TestArchiveLayoutCharacterizationPinsCorpusGoldenAfterSpec0095$'
  ./internal/spec` passed.
- `GOCACHE=/tmp/roundfix-task02-gocache go test -count=1 -tags docscontract
  -run '^(TestCheckCorpusGolden|TestCheckActiveCorpusHasNoErrors)$'
  ./internal/docscontract` passed.
- `diff -r .agents/skills/roundfix skills/roundfix` and `diff -r
  .agents/skills/qa-gate skills/qa-gate` both exited `0`.
- `GOCACHE=/tmp/roundfix-task02-gocache make baseline-digests` exited `0`
  and reported `changed:false`.
- `GOCACHE=/tmp/roundfix-task02-gocache make verify-incremental` reached the
  repository-wide test phase but could not finish because the sandbox blocked
  an existing integration path's access to `api.github.com`; the escalation
  request was rejected. This is recorded as blocked, not passing.
- The daemon-owned command under `## Verification` was not run.

Acceptance evidence:

- `TestCLISurfaceWithoutAGuideIsReported` locates the surface declaration and
  reports a gap; `TestStrictPromotesAnUndocumentedCLISurface` proves strict
  promotion to an error.
- `TestCLISurfaceNamingItsGuidePasses` and
  `TestCLISurfaceWhoseDependencyNamesTheGuidePasses` cover local and
  two-edge transitive ancestry, while `TestGuideOnlyInADependentTaskDoesNotCount`
  proves dependencies are followed in only the permitted direction.
- `TestCLIContractTestIsASurface` proves `internal/cli/cli_test.go` is a
  surface, and `TestOrdinaryCLITestIsNotASurface` proves another CLI test file
  is excluded.
