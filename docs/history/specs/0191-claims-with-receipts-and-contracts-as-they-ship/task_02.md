---
task: task_02
spec: 0191-claims-with-receipts-and-contracts-as-they-ship
status: completed
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

### Implementation

Added the exported `Transcript` reader and the three transcript diagnostics.
The full sweep and authoring-stage check share Task 01's single contract
horizon computation. Only the declaration gap depends on that horizon;
malformed declarations are checked in every Spec. Transcript blocks are read
as text and their commands are never executed.

Surface Transcripts join coverage References, including plural and ranged
references, and require no Coverage Map entry. Their implementation coverage
excludes the QA gate without changing other coverage units' behavior or text.
The pending gate must name each transcript in a Requirement; a declined gate
and completed QA Task raise no ungated error.

The corpus code list, golden and its characterization pin name all three
codes at zero. The golden's update sentence identifies Task 02, and all older
counts are unchanged. The canonical Roundfix skill and its mirror list each
code and explain the declaration horizon. Both skill version fields advance
from `0.0.5` to `0.0.6`.

### Focused evidence

All Go checks below used `GOCACHE=/tmp/roundfix-task02-gocache` after the
default cache was refused by the sandbox.

- Starting signal: `go test ./internal/speccheck -run '^TestAWellFormedSurfaceTranscriptReportsNothing$'`
  failed to compile because `SurfaceTranscripts` and the transcript codes did
  not exist.
- `go test ./internal/speccheck -count=1` — exit 0; existing tests and the new
  transcript tests passed before the final test-only expansions.
- `go test ./internal/speccheck -run 'Transcript|SurfaceTranscripts' -count=1`
  — exit 0 on the final implementation and tests, including both held and
  unheld cases for every malformed reason.
- `go test -tags docscontract ./internal/docscontract -run 'Corpus' -count=1`
  — exit 0; actual corpus counts match the golden, and the active corpus has
  no errors.
- `go test ./internal/spec -run 'ArchiveLayoutCharacterization' -count=1`
  — exit 0; the updated corpus golden matches its characterization pin.
- `make skills-sync` — exit 0; `cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md`
  — exit 0.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
  — exit 0; recorded Roundfix `0.0.6` with its canonical content digest.
- `make baseline-digests` — exit 0, `changed: false`; sanctioned regeneration
  found that derived artifacts already match their sources.
- `git -c core.fsmonitor=false diff --check` — exit 0.
- `rtk make verify-incremental` — initial sandbox run exited 2: two
  process-owner integration tests could not read the host process table, and
  the suite guard detected my test-file edits during the run. Repeated with
  required sandbox escalation and an unchanged worktree; exit 0. Formatting,
  vet, the entire package suite, skill sync/version checks, skill contracts
  and CLI build passed. No repository files changed while the successful
  rerun was running.


| Acceptance criterion | Evidence |
| --- | --- |
| Held declaration gap, reasoned None and unheld skip | `TestAHeldTechSpecWithoutSurfaceTranscriptsIsAGap`, `TestSurfaceTranscriptsNoneWithAReasonIsAccepted`, `TestASpecThatIsNotHeldSkipsTheTranscriptDeclarationGap`; the last uses a disposable `gittest` history with a PRD committed before the guide. `TestTranscriptDeclarationRejectsMissingReasonsAndWrongShapes` rejects missing reasons, wrong heading levels and fenced examples. |
| Well-formed blocks and all seven malformed reasons, including unheld Specs | `TestAWellFormedSurfaceTranscriptReportsNothing`, `TestEachMalformedSurfaceTranscriptNamesItsReason` (each reason under held=true and held=false), `TestAMalformedTranscriptIsReportedInASpecThatIsNotHeld`. `TestTranscriptBlockBoundariesAndExactStreams` covers empty streams, first stderr, indentation, fence length and type, info strings, ordering and exit boundaries. |
| Numbered output is not a declaration | `TestNumberedOutputInsideATranscriptIsNotAnItem` includes both a numbered transcript-like output line and a heading in the block. |
| Implementation coverage excludes QA-only References | `TestASurfaceTranscriptNamedOnlyByTheQATaskIsUntasked` plants a QA-only reference and checks the specific summary; `TestASurfaceTranscriptNamedByANonQATaskIsTasked` clears it with a ranged implementation reference. `TestSurfaceTranscriptsNeedNoCoverageMap` checks map exemption. Existing promise-coverage tests passed in the package check. |
| Pending QA Requirements trace transcripts; declined gate raises nothing | `TestASurfaceTranscriptTheQATaskDoesNotNameIsUngated`, `TestASurfaceTranscriptNamedInAQARequirementIsGated`, `TestADeclinedGateRaisesNoUngatedTranscript`. `TestCompletedQATaskIsHistoricalTranscriptEvidence` checks the completed-gate exception at the detector seam using a graph loaded from temporary artifacts. |
| Finding lines and skip lines match Surface Transcripts 3 and 4 | `TestTranscriptFindingsRenderSurfaceTranscriptsThreeAndFour` asserts full malformed and ungated finding blocks, source locations, fix text and the exact absent-guide skip line. `TestTranscriptDetectorStagesAndMissingArtifacts` checks stage ownership and absent-artifact skips. |
| Corpus zeros, prior counts, own transcripts and skill mirror | Focused corpus and archive-layout checks passed. `TestThisSpecsSurfaceTranscriptsAreWellFormed` checks all four authored blocks. `TestThisSpecsTranscriptsHaveImplementationAndGateReferences` copies the Spec bundle into a temporary directory and checks shape, implementation references and QA Requirements. Skill sync, version recording and byte comparison passed. |

### Scope and settlement

The initial worktree change was the Daemon's `status: in_progress` in this Task;
that field remains untouched. No other Task or Task Graph was edited,
`internal/speccheck/constraints.go` is unchanged, and no existing top-level test
or exported signature was renamed or removed. No commit, push or Pull Request
was made. The authored `## Verification` commands were not executed; their
execution and terminal Task status remain Daemon-owned.

Additional ordinary path declared for the repository's owned-skill version
contract: `skills/testdata/owned-skill-versions.json`. Its only addition is
Roundfix `0.0.6` and its recorded digest. Sanctioned digest regeneration changed
no derived paths. No follow-up implementation was added to this slice.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `skills/testdata/owned-skill-versions.json`

## Carry-forward provenance

- Source Run: `run_20261001T000219Z_0e911d72b0d72d2f`
- Source commit: `2cc943fc260e9ba5165997d21604dc27e5d1178f`
