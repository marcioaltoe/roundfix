---
task: task_02
spec: 0191-claims-with-receipts-and-contracts-as-they-ship
status: pending
type: backend
complexity: high
---

# Task 02: A TechSpec's Surface Transcripts are declared, well-formed and traced

## Overview

A TechSpec describes a command's new output in prose, and the Task author, the QA gate and the reviewer each rebuild the exact text from it. This Task adds the Surface Transcript of ADR-0184 to the Spec Consistency Check: a numbered declaration, one `transcript` block per item with a command line, the two output streams and the exit code, a reference from a Task that is not the gate, and a Requirement in the QA Task. The check reads shape and references. It never runs the command.

## Requirements

1. MUST add `internal/speccheck/transcripts.go` with `Transcript` and `SurfaceTranscripts`, as `_techspec.md` → Interfaces and Surface Transcripts in a TechSpec state: the section heading at level two or three, numbered `Surface Transcript:` items or `None.` with a reason, numbered lines inside a fenced block never read as items, exactly one block with the info string `transcript` per item, and the four line rules of the block.
2. MUST declare `SC-TRANSCRIPT-UNDECLARED`, `SC-TRANSCRIPT-MALFORMED` and `SC-TRANSCRIPT-UNGATED` in the code block of `internal/speccheck/citations.go`, add them to `citationCoverageDetectorCodes`, and register them in `internal/speccheck/coherence.go`: the first two at the `techspec` stage and the third at the `tasks` stage. `internal/speccheck/constraints.go` MUST NOT change.
3. MUST report `SC-TRANSCRIPT-UNDECLARED` as a gap only for a held Spec whose TechSpec is present, using the contract horizon of task_01, and list it as skipped with the horizon's reason for a Spec that is not held.
4. MUST report `SC-TRANSCRIPT-MALFORMED` as an error for every declared transcript that breaks the block rules, whether or not the Spec is held, naming one of the seven reasons in `_techspec.md` → Surface Transcripts in a TechSpec.
5. MUST make each declared transcript a coverage unit named `Surface Transcript <n>`, parsed from Task References like the other coverage references. `SC-COVERAGE-UNTASKED` MUST report a transcript that no Task other than the QA gate references, with a summary that says so, and MUST keep its behavior and text for every other kind of unit. A Surface Transcript MUST need no Coverage Map entry.
6. MUST report `SC-TRANSCRIPT-UNGATED` as an error when the Task Graph includes the gate, the QA Task is not completed and none of its Requirements names the transcript. A declined gate MUST raise nothing.
7. MUST give the findings the summaries, locations and fix texts of `_techspec.md` → Surface Transcript 3, and the skip line of Surface Transcript 4.
8. MUST keep every existing test green, rename or remove no top-level test, and change no exported function signature.
9. MUST add the three codes to `corpusFindingCodes`, to the corpus golden at `0` with an `update` sentence that names this change, to the golden pinned in `internal/spec/archive_layout_characterization_test.go`, and to the identifier list of the Roundfix skill's Spec Consistency Check section with one line each, then run `make skills-sync`. Every earlier count MUST stay unchanged.
10. MUST put the new tests in `internal/speccheck/transcripts_test.go`, over plain temporary directories and `gittest` repositories. Each code MUST have one test where the defect is planted and the finding is reported, and one where it is absent and nothing is reported.

## Subtasks

- [ ] Parse the Surface Transcripts section and each block.
- [ ] Report the declaration gap through the contract horizon.
- [ ] Report malformed blocks with their reason.
- [ ] Trace each transcript to a non-QA Task and to the QA Task's Requirements.
- [ ] Add the codes to the corpus golden, its pin and the Roundfix skill.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] A held Spec whose TechSpec has no Surface Transcripts declaration reports `SC-TRANSCRIPT-UNDECLARED` as a gap, `None.` with a reason is accepted, and a Spec that is not held lists the code as skipped.
- [ ] A well-formed block reports nothing, and each of the seven malformed reasons is reported by name, also in a Spec that is not held.
- [ ] A numbered line inside a transcript block is not read as a declared item.
- [ ] A transcript named only in the QA Task's References reports `SC-COVERAGE-UNTASKED`, and one named by a non-QA Task does not.
- [ ] A transcript the pending QA Task names in no Requirement reports `SC-TRANSCRIPT-UNGATED`, and a declined gate reports nothing.
- [ ] The rendered text holds the finding lines of Surface Transcript 3 and the skip lines of Surface Transcript 4.
- [ ] The corpus golden holds the three codes at `0`, every earlier count is unchanged, this Spec's own transcripts report nothing, and the Roundfix skill and its mirror name the three codes.

## Context

- instruction: `docs/adr/0184-a-techspec-states-a-command-surface-as-a-transcript.md`
- instruction: `docs/adr/0156-a-promise-a-spec-declares-names-a-consuming-task.md`
- interface: `internal/speccheck/citations.go`
- interface: `internal/speccheck/coherence.go`
- creates: `internal/speccheck/transcripts.go`
- creates: `internal/speccheck/transcripts_test.go`
- interface: `internal/docscontract/corpus_test.go`
- interface: `internal/docscontract/testdata/corpus-golden.json`
- interface: `internal/spec/archive_layout_characterization_test.go`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestAHeldTechSpecWithoutSurfaceTranscriptsIsAGap|TestSurfaceTranscriptsNoneWithAReasonIsAccepted|TestAWellFormedSurfaceTranscriptReportsNothing|TestEachMalformedSurfaceTranscriptNamesItsReason|TestNumberedOutputInsideATranscriptIsNotAnItem|TestASurfaceTranscriptNamedOnlyByTheQATaskIsUntasked|TestASurfaceTranscriptTheQATaskDoesNotNameIsUngated|TestADeclinedGateRaisesNoUngatedTranscript|TestASpecThatIsNotHeldSkipsTheTranscriptDeclarationGap|TestAMalformedTranscriptIsReportedInASpecThatIsNotHeld|TestTranscriptFindingsRenderSurfaceTranscriptsThreeAndFour|TestPromiseCoverageReachesItsTask)$" ./internal/speccheck 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestAHeldTechSpecWithoutSurfaceTranscriptsIsAGap TestSurfaceTranscriptsNoneWithAReasonIsAccepted TestAWellFormedSurfaceTranscriptReportsNothing TestEachMalformedSurfaceTranscriptNamesItsReason TestNumberedOutputInsideATranscriptIsNotAnItem TestASurfaceTranscriptNamedOnlyByTheQATaskIsUntasked TestASurfaceTranscriptTheQATaskDoesNotNameIsUngated TestADeclinedGateRaisesNoUngatedTranscript TestASpecThatIsNotHeldSkipsTheTranscriptDeclarationGap TestAMalformedTranscriptIsReportedInASpecThatIsNotHeld TestTranscriptFindingsRenderSurfaceTranscriptsThreeAndFour TestPromiseCoverageReachesItsTask; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the eleven new named tests do not exist, so the command fails.
- `out="$(go test -count=1 -tags docscontract -v -run "^(TestCheckCorpusGolden|TestCheckActiveCorpusHasNoErrors)$" ./internal/docscontract 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; pin="$(go test -count=1 -v -run "^TestArchiveLayoutCharacterizationPinsCorpusGoldenAfterSpec0095$" ./internal/spec 2>&1)" || { printf "%s\\n" "$pin"; exit 1; }; for code in SC-TRANSCRIPT-UNDECLARED SC-TRANSCRIPT-MALFORMED SC-TRANSCRIPT-UNGATED; do grep -qF -- "\"$code\": 0" internal/docscontract/testdata/corpus-golden.json || { printf 'missing golden code: %s\n' "$code" >&2; exit 1; }; for file in .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$code" || { printf 'missing %s in %s\n' "$code" "$file" >&2; exit 1; }; done; done && make skills-sync-check` — expected: exit 0; before this Task the corpus golden and the Roundfix skill do not name any of the three codes, so the command fails.

## References

- [_prd.md](_prd.md) — Goals 2 and 3; Core Features 3, 4 and 5; Success Metrics 2 and 3
- [_techspec.md](_techspec.md) — Interfaces; Surface Transcripts in a TechSpec; API Contracts 1-3; Surface Transcripts 3-4; Testing Approach 4, 5 and 6; Build Order 2
- ADR-0184; ADR-0156; ADR-0094

## Result
