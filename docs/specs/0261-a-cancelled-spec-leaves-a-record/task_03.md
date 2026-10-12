---
task: task_03
spec: 0261-a-cancelled-spec-leaves-a-record
status: pending
type: backend
complexity: high
---

# Task 03: The Cancel Command cancels an active Spec and refuses while work holds it

## Overview

This Task adds `roundfix cancel <slug> --reason <text>`. It checks what only
the command can see — the Run Database, the Delivery Queue, the other active
Specs and the working tree — then calls task_01's `Cancel` and prints the
result and its warnings. It answers
[the adopted Backlog Entry](references/2026-10-12-a-decomposed-spec-cannot-be-cancelled.md)
of 2026-10-12 on the command side (ADR-0262). It is verifiable on its own
through whole-transcript tests over temporary repositories and homes.

## Requirements

1. MUST add the Cancel Command with the help of `_techspec.md` → Surface
   Transcript 6, the dispatch case, and the two root help lines that
   `_techspec.md` → Surface Transcripts names.
2. MUST check, in the order of `_techspec.md` → API Contract 1, every refusal
   it lists, each with exit 2 through the Preflight failed block and no write.
   The Active Run check MUST open the Run Database only when its file exists,
   reclaim an orphaned active-run lock as `settle` does, and refuse a live
   one. The queue check MUST refuse an item for the Spec at any stage other
   than `merged` and `parked`. The `requires` check MUST read every active
   Spec's Task Graph manifest and refuse when one names the Spec.
3. MUST print, on success, the lines of `_techspec.md` → Surface Transcript 1,
   and the warnings of API Contract 2 in the order of Surface Transcript 2,
   with exit 0. Citations MUST come from `spec.HistoryCitations` over the
   removed paths and the folder path. It MUST exit 1 only when `Cancel` fails
   after validation.
4. MUST create `internal/cli/cancel_test.go` with one test per Surface
   Transcript, each asserting every line of standard output and standard
   error in order and the exit code, copied from the TechSpec with only its
   declared variations: `TestCancelTranscriptCancelsASpecWithSources` (1),
   `TestCancelTranscriptWarnsAboutAParkedItemAndACitation` (2),
   `TestCancelTranscriptRefusesAnUnfinishedQueueItem` (3),
   `TestCancelTranscriptRefusesAnActiveRun` (4),
   `TestCancelTranscriptRefusesARequiringSpec` (5) and
   `TestCancelTranscriptPrintsHelp` (6). It MUST also add
   `TestCancelRefusesEveryUnfinishedQueueStage`, covering every
   `store.DeliveryStage` except `merged` and `parked`;
   `TestCancelRefusalsLeaveTheTreeUnchanged`, covering a missing reason, an
   unknown Spec, uncommitted changes, a supersession and a non-Markdown file
   that names the active directory, with `git status --porcelain` and a hash
   of the Spec folder before and after; and
   `TestRootHelpNamesTheCancelCommand`. Every test MUST use a temporary
   repository and a temporary home and MUST NOT read the real `~/.roundfix`.
   Every new top-level test MUST call `t.Parallel()` first or carry a
   `// Sequential: <reason>` comment.
5. MUST add the Skill Coverage Map entry `command cancel`, covered by the
   Roundfix archive reference with `docs/user-guide/commands/cancel.md` and
   `internal/cli/cancel.go` as sources, then run
   `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`.
   It MUST NOT hand-edit the Behavior Surface Record.
6. MUST print the help exactly as `_techspec.md` → Surface Transcript 6,
   the same text task_02 wrote into the cancel user guide, and MUST NOT edit
   that guide.
7. MUST NOT write a Task status, commit, push, create a Run, or change
   `internal/spec`, the Supersede Command or any delivery file.

## Subtasks

- [ ] Add the command, its help, its dispatch and the root help lines.
- [ ] Implement the refusals in the contract's order.
- [ ] Print the result, the returned sources and the warnings.
- [ ] Write the transcript and refusal tests.
- [ ] Add the coverage entry and record the Behavior Surface Record.

## Acceptance Criteria

- [ ] Each of the six Surface Transcripts is asserted whole in its named test.
- [ ] Every unfinished queue stage refuses, and a parked item only warns.
- [ ] Every refusal leaves the working tree byte-identical.
- [ ] The root help names the command, and the coverage map is current.

## Context

- instruction: `docs/adr/0262-a-cancelled-spec-leaves-an-archive-record-and-returns-its-sources.md`
- instruction: `internal/cli/archive.go`
- instruction: `internal/cli/settle.go`
- instruction: `internal/cli/reopen.go`
- instruction: `internal/cli/archive_record_test.go`
- instruction: `internal/store/delivery.go`
- instruction: `internal/spec/requires.go`
- instruction: `internal/spec/history_entries.go`
- instruction: `internal/speccheck/active_spec_paths.go`
- interface: `internal/cli/cli.go`
- interface: `docs/references/skill-coverage.json`
- interface: `docs/references/behavior-surfaces.json`
- creates: `internal/cli/cancel.go`
- creates: `internal/cli/cancel_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestCancelTranscriptCancelsASpecWithSources|TestCancelTranscriptWarnsAboutAParkedItemAndACitation|TestCancelTranscriptRefusesAnUnfinishedQueueItem|TestCancelTranscriptRefusesAnActiveRun|TestCancelTranscriptRefusesARequiringSpec|TestCancelTranscriptPrintsHelp)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestCancelTranscriptCancelsASpecWithSources TestCancelTranscriptWarnsAboutAParkedItemAndACitation TestCancelTranscriptRefusesAnUnfinishedQueueItem TestCancelTranscriptRefusesAnActiveRun TestCancelTranscriptRefusesARequiringSpec TestCancelTranscriptPrintsHelp; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task none of the transcript tests exists, so the command fails.
- `out="$(go test -count=1 -v -run '^(TestCancelRefusesEveryUnfinishedQueueStage|TestCancelRefusalsLeaveTheTreeUnchanged|TestRootHelpNamesTheCancelCommand)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestCancelRefusesEveryUnfinishedQueueStage TestCancelRefusalsLeaveTheTreeUnchanged TestRootHelpNamesTheCancelCommand; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; help="$(go run -buildvcs=false ./cmd/roundfix cancel --help 2>&1)" || { printf '%s\n' "$help"; exit 1; }; printf '%s\n' "$help" | grep -qF -- "--reason  Why the Spec is cancelled, recorded in the Archive Record" || { printf 'cancel help lacks the option\n' >&2; exit 1; }; out="$(go test -count=1 -tags docscontract -v -run '^TestTheSkillCoverageMapIsCurrent$' ./internal/docscontract 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- "--- PASS: TestTheSkillCoverageMapIsCurrent" || exit 1` — expected: exit 0; before this Task the refusal tests do not exist and the command has no help, so the command fails.

## References

- [_prd.md](_prd.md) — Goal 1; Goal 4; User Story 1; User Story 4; Core Feature 1; Core Feature 4; Core Feature 6; Success Metric 1; Success Metric 2
- [_techspec.md](_techspec.md) — API Contract 1; API Contract 2; Surface Transcript 1; Surface Transcript 2; Surface Transcript 3; Surface Transcript 4; Surface Transcript 5; Surface Transcript 6; Integration Points; Testing Approach; Build Order 3
- ADR-0262; ADR-0184; ADR-0256; ADR-0259
