---
task: task_01
spec: 0239-a-glossary-every-spec-keeps-current
status: pending
type: backend
complexity: high
---

# Task 01: The Spec Consistency Check and the Archive Command hold a Glossary Declaration to the glossary

## Overview

The Spec Consistency Check has no detector that reads a Spec's domain terms
against the glossary, so twenty-one Specs archived after 2026-10-02 without
touching it ([the adopted Backlog Entry of 2026-10-06](references/2026-10-06-the-context-driven-loop-does-not-keep-context-md-current.md)).
This Task adds the glossary detector of `_techspec.md` → Invariants 1 to 10,
wires it into the full and staged checks, adds its three codes to the
characterized corpus at zero, makes the Archive Command refuse a Glossary Gap,
and documents both in the `spec` and `archive` command guides (ADR-0244). It is
verifiable on its own: temporary repositories prove every finding and the
refusal through the public entry points.

This is an authorized tooling Task for `internal/speccheck/constraints.go`,
`internal/speccheck/coherence.go`,
`internal/docscontract/testdata/corpus-golden.json` and
`internal/spec/archive_layout_characterization_test.go`. It may change only the
files in its Context and this Task file.

## Requirements

1. MUST create `internal/speccheck/glossary.go` with the constants,
   `GlossaryFindings` and the detector of `_techspec.md` → Interfaces,
   implementing Invariants 1 to 8 and 10 exactly: the glossary files, the
   declaration grammar, the horizon (reusing `prdAddingCommit` and
   `adrHorizonGitOutput` with the contract horizon's semantics), the bold rule,
   the binding Task, the `changes` rule, the presence rule and the stages.
2. MUST call the detector from `Check` in `internal/speccheck/constraints.go`
   with the Task Graph when present, and from `CheckStage` in
   `internal/speccheck/coherence.go` at the PRD and TechSpec stages, and MUST
   add `CodeGlossaryUndeclared` at `StagePRD` and `CodeGlossaryUnplanned` and
   `CodeGlossaryMissing` at `StageTasks` to `stagedDetectors`, so
   `CheckStage(StageAll)` still equals `Check` for every Spec.
3. MUST make `roundfix archive` call `speccheck.GlossaryFindings` after the
   active-path-pin preflight and before `spec.Archive`, and refuse per
   `_techspec.md` → Invariant 9 with the reason
   `Spec "<slug>" cannot archive with a Glossary Gap: <code>: <summary>`, with
   and without `--qa-override`, changing no file.
4. MUST create `internal/speccheck/glossary_test.go` with the eleven tests of
   `_techspec.md` → Testing Approach 1 and `internal/cli/archive_glossary_test.go`
   with the three tests of Testing Approach 2, each in a temporary Git
   repository it creates, with every fixture under `t.TempDir()`.
   `TestGlossaryBoldRuleIgnoresLabelsDigitsAndFences` MUST include a one-word
   label, a phrase with a digit, a phrase in parentheses, a phrase ending in a
   period and a phrase in a fenced block, none reported, beside one reported
   phrase. `TestGlossaryLegacySpecIsSkipped` MUST commit a PRD without the
   section before the commit that adds the guide and require no glossary
   finding and the recorded skip, and a second PRD committed after it require
   the missing-section error.
5. MUST add the three codes to `corpusFindingCodes` in
   `internal/docscontract/corpus_test.go`, to `active` at `0` in
   `internal/docscontract/testdata/corpus-golden.json` with one sentence
   appended to its `update` naming them, and to the wanted golden in
   `internal/spec/archive_layout_characterization_test.go` with the same
   sentence, per `_techspec.md` → Existing tests that change. No other existing
   assertion changes, except as that section allows.
6. MUST document in `docs/user-guide/commands/spec.md` the Glossary
   Declaration, the three codes and when each is reported, the bold rule and
   the glossary horizon, using the words glossary horizon; and in
   `docs/user-guide/commands/archive.md` that archive refuses a Spec with a
   Glossary Gap, also under a QA Archive Override.
7. MUST NOT read or write outside the repository and `t.TempDir()`, call the
   network, or call the Jev judge.
8. MUST prove each new gate can fail. The Result MUST record one sabotage of
   the bold rule (for example accepting one-word phrases), one of the binding
   rule (for example ignoring the Verification text) and one of the archive
   refusal (for example skipping it under `--qa-override`), each with the test
   that failed, and that the code was restored.

## Subtasks

- [ ] Write the detector and its tests.
- [ ] Wire it into the full and staged checks and the corpus code set.
- [ ] Refuse archive on a Glossary Gap.
- [ ] Document the codes and the refusal.
- [ ] Record each sabotage.

## Acceptance Criteria

- [ ] Every finding of Invariants 2 to 7 is reported at its artifact and line
      in a temporary repository, and none for a legacy Spec.
- [ ] The staged and full checks agree, and the active corpus reports the three
      codes at zero.
- [ ] `roundfix archive` exits `2` with the Glossary Gap reason, with and
      without `--qa-override`, and moves nothing.
- [ ] The `spec` and `archive` command guides describe the codes and the
      refusal.

## Context

- instruction: `docs/adr/0244-a-spec-declares-the-domain-terms-it-introduces-and-the-check-holds-them-to-the-glossary.md`
- instruction: `docs/adr/0094-the-consistency-check-is-artifact-presence-aware.md`
- instruction: `docs/adr/0117-a-defect-is-checked-by-the-stage-that-can-produce-it.md`
- instruction: `internal/speccheck/authoring_horizon.go`
- instruction: `internal/speccheck/adr_horizon.go`
- instruction: `internal/speccheck/vocabulary.go`
- creates: `internal/speccheck/glossary.go`
- creates: `internal/speccheck/glossary_test.go`
- interface: `internal/speccheck/constraints.go`
- interface: `internal/speccheck/coherence.go`
- interface: `internal/cli/archive.go`
- creates: `internal/cli/archive_glossary_test.go`
- interface: `internal/docscontract/corpus_test.go`
- interface: `internal/docscontract/testdata/corpus-golden.json`
- interface: `internal/spec/archive_layout_characterization_test.go`
- interface: `docs/user-guide/commands/spec.md`
- interface: `docs/user-guide/commands/archive.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestGlossaryDeclarationIsRequiredAfterTheHorizon|TestGlossaryMalformedEntryIsReported|TestGlossaryUndeclaredBoldTermIsReported|TestGlossaryBoldRuleIgnoresLabelsDigitsAndFences|TestGlossaryNotATermCoversAPhrase|TestGlossaryAddedTermNeedsATaskThatWritesIt|TestGlossaryChangedTermMustExist|TestGlossaryCompletedTaskMustLeaveTheTermDefined|TestGlossaryReadsMappedAndRenamedGlossaries|TestGlossaryStagesReportOnlyTheirCodes|TestGlossaryLegacySpecIsSkipped|TestArchiveRefusesASpecWithAGlossaryGap|TestArchiveWithQAOverrideStillRefusesAGlossaryGap|TestArchiveAcceptsASpecWhoseGlossaryIsCurrent|TestArchiveRefusesASpecAFileStillPins|TestStageScopeDefaultSweepIsUnchanged|TestCheckCorpusGolden|TestCheckActiveCorpusHasNoErrors|TestArchiveLayoutCharacterizationPinsCorpusGoldenAfterSpec0095)$" ./internal/speccheck ./internal/cli ./internal/docscontract ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestGlossaryDeclarationIsRequiredAfterTheHorizon TestGlossaryMalformedEntryIsReported TestGlossaryUndeclaredBoldTermIsReported TestGlossaryBoldRuleIgnoresLabelsDigitsAndFences TestGlossaryNotATermCoversAPhrase TestGlossaryAddedTermNeedsATaskThatWritesIt TestGlossaryChangedTermMustExist TestGlossaryCompletedTaskMustLeaveTheTermDefined TestGlossaryReadsMappedAndRenamedGlossaries TestGlossaryStagesReportOnlyTheirCodes TestGlossaryLegacySpecIsSkipped TestArchiveRefusesASpecWithAGlossaryGap TestArchiveWithQAOverrideStillRefusesAGlossaryGap TestArchiveAcceptsASpecWhoseGlossaryIsCurrent TestArchiveRefusesASpecAFileStillPins TestStageScopeDefaultSweepIsUnchanged TestCheckCorpusGolden TestCheckActiveCorpusHasNoErrors TestArchiveLayoutCharacterizationPinsCorpusGoldenAfterSpec0095; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the glossary and archive tests do not exist, so the command fails.
- `for pair in "docs/user-guide/commands/spec.md|SC-GLOSSARY-UNDECLARED" "docs/user-guide/commands/spec.md|SC-GLOSSARY-UNPLANNED" "docs/user-guide/commands/spec.md|SC-GLOSSARY-MISSING" "docs/user-guide/commands/spec.md|glossary horizon" "docs/user-guide/commands/archive.md|Glossary Gap"; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done; grep -qF -- 'Glossary Gap' internal/cli/archive.go || { printf 'internal/cli/archive.go names no Glossary Gap\n' >&2; exit 1; }` — expected: exit 0; before this Task neither guide nor the command names the glossary codes or the refusal, so the command fails.

## References

- [_prd.md](_prd.md) — Goal 1; User Stories 2 and 3; Core Features 1, 2, 3 and 4; Success Metric 1; Success Metric 2; Success Metric 3
- [_techspec.md](_techspec.md) — Interfaces; Invariants; API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; Testing Approach 1; Testing Approach 2; Testing Approach 3; Existing tests that change; Vocabulary Contract; Build Order 1
- ADR-0244; ADR-0094; ADR-0117; ADR-0168
