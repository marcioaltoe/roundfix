---
task: task_03
spec: 0261-a-cancelled-spec-leaves-a-record
status: completed
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

## Result

Implemented the Cancel Command within this Task's CLI slice. Dispatch and
root help name `cancel`, and command help matches Surface Transcript 6.
Preflight checks arguments, the Git working tree, the active folder,
non-Markdown directory pins, committed Spec bytes, the Active Run, the
Delivery Queue and active Specs' `requires`, before calling `spec.Cancel`.
The Run Database is opened only if it exists; orphan reclaim uses the same
helper as `settle`. Every unfinished queue stage refuses; `merged` permits
cancellation and `parked` permits it with the prescribed warning.

Success prints the cancellation summary and returned sources, followed on
stderr by the parked-item warning and citation advice. Advice originates in
`spec.HistoryCitations` over the planned removed files and folder, collected
before the folder cut because the helper reads tracked working-tree files.
The shared Archive Relink plan filters citations that will remain, including
code examples on a line whose other link is relinked. The returned
`RemovedPaths` limits the warnings to paths actually removed. Validation
errors use exit 2; existing Cancel mutation-operation labels distinguish
write/removal failures for exit 1 without changing `internal/spec`.

Added the `command cancel` coverage entry and generated the Behavior Surface
Record with the required record command. The cancel user guide, Supersede
Command, delivery implementation and Spec library were not changed.

### Focused checks

- Before command implementation,
  `go test ./internal/cli -run '^TestCancelTranscriptPrintsHelp$' -count=1`
  exited 1 with `unknown command "cancel"`, establishing the missing surface.
- After the final code and test edits,
  `go test ./internal/cli -run '^(TestCancel|TestRootHelpNamesTheCancelCommand|TestRunHelp|TestRunCommandHelp|TestHelpTokenAsAFlagValueIsNotHelp|TestLeadingHelpTokenPrintsUsage)' -count=1`
  exited 0 (`ok roundfix/internal/cli`). This selects every Cancel test and
  the directly related existing help checks.
- `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`
  exited 0 after the final command edit (`ok roundfix/internal/docscontract`)
  and generated only the new Cancel surface entry. This is Requirement 5's
  implementation generator, not a run of the declared Verification commands.
- `git diff --check` exited 0. Commands were executed with `rtk proxy` for
  exact output; the configured sandbox root contains a symlink, so the
  executor required approved execution for workspace access.

### Acceptance evidence

| Acceptance criterion | Implementation and focused evidence |
| --- | --- |
| Six whole Surface Transcripts | All six prescribed `TestCancelTranscript...` tests compare complete stdout, stderr and exit code to independent transcript expectations, substituting only fixture names, file counts, byte counts, revision and Run ID. All passed in the focused selection. |
| Every unfinished stage refuses; parked only warns | `TestCancelRefusesEveryUnfinishedQueueStage` covers queued, running, reviewing, archiving, gating, publishing, checking and merging. Transcript 2 asserts parked warning order; `TestCancelAcceptsAMergedQueueItem` asserts merged produces no warning or commit. All passed. |
| Refusals preserve the working tree | `TestCancelRefusalsLeaveTheTreeUnchanged` covers missing reason, unknown Spec, uncommitted bytes, supersession and a non-Markdown pin. Its shared assertion compares the whole file snapshot, `git status --porcelain --untracked-files=all` and a SHA-256 digest of Spec paths and bytes. The queue, Active Run, requiring Spec, malformed requires, occupied Source Return destination and argument refusals use the same assertion. All passed. |
| Root help and coverage | `TestRootHelpNamesTheCancelCommand` asserts both new lines immediately after the archive lines; the whole command-help transcript passes. The sanctioned coverage record command passed and wrote the new Behavior Surface entry. |

Additional passing checks cover an absent Run Database, orphaned Run reclaim,
non-Git refusal, whitespace-normalized reasons, cancellation without a Task
Graph, exit 1 with unchanged bytes after a destination-creation failure,
History/untracked citation exclusions, one warning per line, and a surviving
citation on a partly relinked line. The last case initially reproduced a
missing warning and passes with the shared relink-plan filter.

The pre-existing Task status remains `in_progress`. Only this `## Result`
was appended to the Task file. No Task Graph or other Task file was changed.
No commit, push, Pull Request, selected repository Verification or declared
Task Verification was run. The Daemon owns Verification and settlement;
reader-command changes remain task_04's slice.
