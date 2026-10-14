---
task: task_02
spec: 0270-acceptance-rules-that-match-the-authorization-and-the-queue
status: completed
type: backend
complexity: medium
---

# Task 02: Spec check reports a QA requirement that pins a byte invariant to a commit

## Overview

The Spec Consistency Check gains the Tasks-stage detector of `_techspec.md` →
Data Models and API Contract 2: `SC-QA-BYTES-PINNED`, a gap for each
requirement of a pending QA Task that names a byte-identity phrase and a
commit-shaped hexadecimal token together. The corpus pins gain the code at 0.
It is verifiable on its own through the detector tests and the consuming
suites.

## Requirements

1. MUST answer the adopted Backlog Entry of 2026-10-13,
   [a byte invariant pinned to a past commit fails once other Specs merge first](references/2026-10-13-a-byte-invariant-pinned-to-a-past-commit-fails-after-other-merges.md),
   by implementing API Contract 2 of `_techspec.md` with the grammar of
   `_techspec.md` → Data Models, in a new file beside the existing per-Task
   detectors, called from the Task loop that already skips completed Tasks
   (ADR-0276).
2. MUST read only the requirement items of the Task that the manifest's `qa:`
   names or whose type is `qa`, and MUST NOT read non-QA Tasks, completed
   Tasks, Verification commands or other sections.
3. MUST register the code at the Tasks stage in `stagedDetectors`, so the PRD
   and TechSpec stages list it as skipped, and MUST list it as skipped when the
   Task Graph is absent (ADR-0094).
4. MUST make each finding a gap that holds every active Spec with no horizon
   (ADR-0265), located at the requirement's first line, with the summary and
   fix text of API Contract 2.
5. MUST add the code at 0 to `corpusFindingCodes`, to the corpus golden and to
   the characterization that pins it, and append the identical sentence of API
   Contract 3 to the update text of both.
6. MUST add the four tests that `_techspec.md` → Testing Approach names. The
   positive test reproduces the shape of the source: a QA requirement that
   holds a file "byte-identical to its bytes at" a synthetic token, and
   another that confirms "from the diff against" a token that two files are
   byte-identical. The negative test covers a requirement that names the
   Delivery Base without a token, a replay requirement that names a token
   without a byte phrase, a byte phrase with a path such as `docs/abc1234`, a
   date, an ADR number and a 64-character digest, a non-QA Task and a
   completed QA Task. Every new top-level test MUST call `t.Parallel()` first.
7. MUST NOT execute any authored command, and MUST NOT read Git history.

## Subtasks

- [ ] Write the detector and its four tests.
- [ ] Wire it into the Task loop and the stage registry, with the skip.
- [ ] Add the code to the corpus pins.

## Acceptance Criteria

- [ ] A pending QA requirement with a byte phrase and a commit-shaped token
      yields one gap naming the file, the line and the token; `--strict` fails.
- [ ] The negative cases, a non-QA Task and a completed QA Task yield none.
- [ ] The PRD and TechSpec stages and a Spec without a Task Graph list the
      code as skipped.
- [ ] The fix names the Delivery Base and the write-tasks rule.
- [ ] The existing speccheck, spec check command and settlement tests and the
      corpus tests pass, with the new code at 0.

## Consumers

- The Spec Consistency Check finding set — consumers: `./internal/speccheck`, `./internal/cli`, `./internal/daemon`, `./internal/docscontract`

## Context

- instruction: `docs/adr/0276-a-row-the-authorization-keeps-offline-qualifies-and-a-byte-invariant-compares-with-the-delivery-base.md`
- instruction: `docs/adr/0265-a-spec-proves-what-it-claims-and-every-active-spec-is-held-to-it.md`
- creates: `internal/speccheck/pinned_bytes.go`
- creates: `internal/speccheck/pinned_bytes_test.go`
- interface: `internal/speccheck/citations.go`
- interface: `internal/speccheck/coherence.go`
- interface: `internal/docscontract/corpus_test.go`
- interface: `internal/docscontract/testdata/corpus-golden.json`
- interface: `internal/spec/archive_layout_characterization_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestPinnedByteInvariantGapNamesTheRequirementAndTheCommit|TestPinnedByteInvariantSkipsUnpinnedRequirementsAndOtherTasks|TestPinnedByteInvariantRunsAtTheTasksStageAndFailsStrict|TestPinnedByteInvariantFixNamesTheDeliveryBaseAndTheSkillRule)$' ./internal/speccheck 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestPinnedByteInvariantGapNamesTheRequirementAndTheCommit TestPinnedByteInvariantSkipsUnpinnedRequirementsAndOtherTasks TestPinnedByteInvariantRunsAtTheTasksStageAndFailsStrict TestPinnedByteInvariantFixNamesTheDeliveryBaseAndTheSkillRule; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; grep -qF -- '"SC-QA-BYTES-PINNED": 0' internal/docscontract/testdata/corpus-golden.json || { printf 'the corpus golden lacks SC-QA-BYTES-PINNED\n' >&2; exit 1; }; go test -count=1 ./internal/speccheck || exit 1; go test -count=1 -run '^TestArchiveLayoutCharacterizationPinsCorpusGoldenAfterSpec0095$' ./internal/spec || exit 1; go test -count=1 -run '^(TestRunSpecCheckCleanText|TestSpecCheckRunVerification|TestRunSpecCheckGapStrictPromotion)$' ./internal/cli || exit 1; go test -count=1 -run 'SpecConsistency' ./internal/daemon || exit 1; go test -count=1 -tags docscontract -run '^(TestCheckCorpusGolden|TestCheckActiveCorpusHasNoErrors)$' ./internal/docscontract` — expected: exit 0; before this Task the four detector tests do not exist, so the command fails at its first pass check.

## References

- [_prd.md](_prd.md) — Goals 3; User Stories 3; Core Features 3; Success Metric 3; Success Metric 4
- [_techspec.md](_techspec.md) — Data Models; API Contract 2; API Contract 3; Invariants 4-5; Testing Approach; Build Order 2
- ADR-0276; ADR-0093; ADR-0094; ADR-0265; ADR-0148; ADR-0259; ADR-0271

## Result

Implemented the Task 02 slice for Daemon Verification. The detector reads
only parsed requirement items of a non-completed Task identified by the
manifest's QA entry or its `qa` type. It emits one `SC-QA-BYTES-PINNED` gap
per matching requirement with the first commit-shaped token, the requirement's
first line, and the exact summary and fix of API Contract 2. It follows the
specified token grammar, performs no IO, executes no authored command, and
reads no Git history or contract horizon.

The existing Task loop calls the detector after its completed-Task skip.
The stage registry marks it Tasks-only; an absent Task Graph or PRD records
a skip. The corpus registry, golden and characterization include the code
at 0. Both pinned update texts append the identical API Contract 3 sentence,
and every earlier count is unchanged.

Focused evidence:

- Before implementation, `rtk proxy go test ./internal/speccheck -run TestPinnedByteInvariant -count=1`
  exited 1: the new tests could not compile because `CodeQABytesPinned` and
  `PinnedByteInvariants` did not exist.
- After implementation, the same focused command exited 0.
- `rtk proxy go test ./internal/speccheck -run 'Test(PinnedByteInvariant|Consumers|CheckStage)' -count=1`
  exited 0, covering the four new tests and neighboring consumer/stage tests.
- `rtk proxy go test -tags docscontract ./internal/docscontract ./internal/spec -run '^$'`
  exited 0 for both changed consumer packages; this was a compilation check
  with no tests selected.
- An inline `rtk proxy python3` consistency check parsed the corpus JSON,
  asserted its new count is 0, confirmed the new code in the registry and
  characterization, and compared both complete update texts: exit 0.
- `rtk proxy git diff --check` exited 0. Changed-path inspection showed only
  the seven declared implementation/test paths and this Task file.

| Acceptance criterion | Implementation and focused evidence |
| --- | --- |
| Pending QA requirement reports one gap with file, line and token; strict fails | `TestPinnedByteInvariantGapNamesTheRequirementAndTheCommit` checks both source shapes at their exact first lines, wrapped text, all eight phrases, 7/40-character bounds, trailing dots, ranges and one gap despite multiple tokens. `TestPinnedByteInvariantRunsAtTheTasksStageAndFailsStrict` checks promotion to error through the shared strict API, including a token labelled as Delivery Base. |
| Negative cases, non-QA and completed QA yield none | `TestPinnedByteInvariantSkipsUnpinnedRequirementsAndOtherTasks` checks every required exclusion plus token case/length/prefix boundaries, separate requirements, other sections and Verification text. An authored command's marker remains absent. |
| PRD/TechSpec stages and missing Task Graph skip the code | `TestPinnedByteInvariantRunsAtTheTasksStageAndFailsStrict` checks both earlier stages and both all/tasks modes without the graph. |
| Fix names Delivery Base and write-tasks rule | `TestPinnedByteInvariantFixNamesTheDeliveryBaseAndTheSkillRule` asserts the complete API Contract 2 fix exactly. |
| Existing consuming suites and corpus tests pass with the new code at 0 | Neighboring speccheck tests pass; changed spec/docscontract consumers compile; all three pins and the identical update sentence agree. The full speccheck, CLI, Daemon, characterization and active-corpus results await Daemon Verification; they were not executed or claimed here. |

The Task arrived with the Daemon-owned `in_progress` change. Its status and
authored content were preserved; only this Result was appended. Authored
Verification and repository Verification remain Daemon-owned. No commit,
push, Pull Request, manifest edit or other Task edit was performed.

Follow-ups: none within this slice. Skill, Baseline, glossary and user-guide
work remains with the other Tasks named by the graph.
