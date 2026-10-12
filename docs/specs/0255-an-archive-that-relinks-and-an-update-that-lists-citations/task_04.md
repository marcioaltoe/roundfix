---
task: task_04
spec: 0255-an-archive-that-relinks-and-an-update-that-lists-citations
status: completed
type: backend
complexity: high
---

# Task 04: Baseline update lists History Citations, promotes on request and names Stale History Citations

## Overview

Gives the history section of `roundfix baseline update` the citations it has
lacked (ADR-0260). The plan lists every History Citation the selected units'
removal breaks, marks the Live Document ones and offers `--promote`, and binds
them into the Plan Digest. `--promote <path>` copies a cited file of a selected
folder to `docs/references/` on apply, and nothing else changes in the citing
files. With the History Full Tag present, the update also names each Stale
History Citation and the command that restores its target. The Task is
verifiable on its own through command tests on synthetic adopted repositories.

## Requirements

1. MUST compute the History Citations of `_techspec.md` → API Contract 7 in
   the history section: `spec.HistoryCitations` over the files of each selected
   folder and each selected `remove` kind, then over each selected folder path,
   dropping a folder citation whose file and line already hold one; `reduce`
   kinds are left out. Each citation is marked live with `spec.IsLiveDocument`
   from task_02, and a live citation of a non-core file of a selected folder
   carries `offer` equal to its target unless it is promoted, which sets
   `promotedTo`. A detector failure blocks the history section with a one-line
   message and never the Baseline part.
2. MUST add the always-present `citations` and `staleCitations` arrays of
   `_techspec.md` → Data Models to the `history` object, never null, with
   exactly the stated keys.
3. MUST extend the history digest payload with `citations` (path, line,
   target, in report order) and, per folder unit, `promoted`, each omitted when
   empty, so a plan without citations and promotions keeps today's digest
   (`_techspec.md` → Invariant 7). A changed citing line changes the digest.
4. MUST add the repeatable `--promote <path>` flag of `_techspec.md` → API
   Contract 8. It is refused with `--no-history` as invalid input. After
   planning and before any write, each value is checked: a clean
   repository-relative file inside a selected folder, outside every Refused
   Unit, not a core artifact, with a basename no other promotion uses and no
   existing `docs/references/` file. A failed check exits 2 with category
   `invalid`, names the value and writes nothing. Valid values re-plan their
   folder through `spec.PlanLegacyConversion` with `Promote`, as
   `roundfix history sanitize` does. The update MUST NOT write to any citing
   file (Invariant 8), and the history write-failure next action MUST also
   name `docs/references`. The hazard to name: the copy's source is a
   committed history file at the planned revision, and its destination is
   read by every later Spec, so a promotion never overwrites an existing file.
5. MUST create `spec.StaleHistoryCitations` in the new
   `internal/spec/stale_citations.go` as `_techspec.md` → API Contract 10
   states, and call it from the history section with the tag tree the section
   already reads, only when `history-full` exists and the status is neither
   `blocked` nor `skipped`. The scan MUST look each candidate up in a set and
   MUST never count a directory. The report never enters the digest, state,
   category or exit code (Invariant 9).
6. MUST print the text lines of `_techspec.md` → API Contract 9 after the
   `History refused` lines and before `History tag`, and print the history
   block also when only citations or stale citations exist. The usage text of
   `baseline update` in `internal/cli/cli.go` gains `[--promote <path> ...]` on
   a line of its own and one sentence on History Citations, the offer and
   Stale History Citations; every snippet
   `internal/docscontract/publicdocs_test.go` pins for that usage stays present.
7. MUST add these tests to `internal/spec/stale_citations_test.go`:
   - `TestStaleHistoryCitationsReportsFilesTheTagHolds`: a relative link and a
     repository path to removed files the held set names are reported with
     file and line.
   - `TestStaleHistoryCitationsIgnoresDirectoriesAndPresentFiles`: a cited
     directory, a cited file still in the working tree, and a citation from a
     file under the History Root are not reported.
8. MUST add these tests to `internal/cli/baseline_update_citations_test.go`,
   each on `newBaselineUpdateRepository` with the history of
   `_techspec.md` → Surface Transcripts committed under `docs/history/`, a
   temporary home and `--no-skills`:
   - `TestBaselineUpdateListsHistoryCitationsTranscript`: asserts Surface
     Transcript 2 whole, ignoring placeholders, and the JSON `citations` with
     `live` and `offer`; the preview writes nothing.
   - `TestBaselineUpdatePromotesOnlyOnRequestTranscript`: asserts Surface
     Transcript 3 whole; the promoted copy equals the removed blob, the record
     lists it, and both citing files are byte-identical. A confirmed update
     without `--promote` copies nothing.
   - `TestBaselineUpdateReportsStaleHistoryCitationsTranscript`: asserts
     Surface Transcript 4 whole and the JSON `staleCitations`; state `current`
     and exit 0.
   - `TestBaselineUpdateHistoryDigestBindsCitations`: moving one citing line
     changes `planDigest`, and `--confirm-plan` with the old digest exits 3
     with category `approval` and writes nothing.
   - `TestBaselineUpdateRefusesAnInvalidPromotion`: a core artifact, a path
     outside the selected units, an existing destination and `--promote` with
     `--no-history` each exit 2 with category `invalid` and write nothing.
   - `TestBaselineUpdateWithoutCitationsKeepsDigestAndText`: Pending History
     with no citing file yields the digest and text block of today's code, and
     `citations` and `staleCitations` are empty arrays.
9. MUST NOT change any existing assertion; every existing `baseline update`
   and History Sanitize Command test passes unchanged, and
   `roundfix history sanitize` output does not change.

## Subtasks

- [ ] List and classify the History Citations of the selected units.
- [ ] Bind them and the promotions into the digest.
- [ ] Add and validate `--promote`, and re-plan promoted folders.
- [ ] Report Stale History Citations from the tag.
- [ ] Print the lines and the usage, and write the tests.

## Acceptance Criteria

- [ ] The preview lists every History Citation, marks Live Document ones with the offer, and its digest binds them.
- [ ] Only `--promote` copies a file upstream, and no citing file changes.
- [ ] With the tag present, the update names each Stale History Citation and its restore command without changing state or exit code.
- [ ] A repository without citations keeps today's digest and text.

## Context

- interface: `internal/cli/baseline_update_history.go`
- interface: `internal/cli/baseline_update.go`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/history.go`
- creates: `internal/spec/stale_citations.go`
- creates: `internal/spec/stale_citations_test.go`
- creates: `internal/cli/baseline_update_citations_test.go`
- instruction: `internal/spec/history_entries.go`
- instruction: `internal/cli/baseline_update_history_test.go`
- instruction: `internal/cli/baseline_update_test.go`
- instruction: `internal/docscontract/publicdocs_test.go`
- instruction: `docs/adr/0260-an-archive-relinks-live-citations-and-a-baseline-update-lists-the-citations-it-breaks.md`
- instruction: `docs/adr/0254-baseline-update-sanitizes-pending-history-in-the-change-it-plans.md`
- instruction: `docs/adr/0173-a-baseline-plan-reports-the-citations-its-history-relocations-break.md`
- instruction: `docs/user-guide/commands/baseline.md`
- instruction: `.agents/skills/roundfix/references/baseline.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestStaleHistoryCitationsReportsFilesTheTagHolds|TestStaleHistoryCitationsIgnoresDirectoriesAndPresentFiles)$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for t in TestStaleHistoryCitationsReportsFilesTheTagHolds TestStaleHistoryCitationsIgnoresDirectoriesAndPresentFiles; do printf '%s\n' "$out" | grep -q -- "--- PASS: $t " || { printf 'missing passing test %s\n%s\n' "$t" "$out" >&2; exit 1; }; done` — expected: exit 0. Before this Task neither test exists, so it fails.
- `out="$(go test -count=1 -v -run '^(TestBaselineUpdateListsHistoryCitationsTranscript|TestBaselineUpdatePromotesOnlyOnRequestTranscript|TestBaselineUpdateReportsStaleHistoryCitationsTranscript|TestBaselineUpdateHistoryDigestBindsCitations|TestBaselineUpdateRefusesAnInvalidPromotion|TestBaselineUpdateWithoutCitationsKeepsDigestAndText)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for t in TestBaselineUpdateListsHistoryCitationsTranscript TestBaselineUpdatePromotesOnlyOnRequestTranscript TestBaselineUpdateReportsStaleHistoryCitationsTranscript TestBaselineUpdateHistoryDigestBindsCitations TestBaselineUpdateRefusesAnInvalidPromotion TestBaselineUpdateWithoutCitationsKeepsDigestAndText; do printf '%s\n' "$out" | grep -q -- "--- PASS: $t " || { printf 'missing passing test %s\n%s\n' "$t" "$out" >&2; exit 1; }; done; help="$(go run -buildvcs=false ./cmd/roundfix baseline update --help 2>&1)" || { printf '%s\n' "$help"; exit 1; }; printf '%s\n' "$help" | grep -q -- '--promote <path>' || { printf 'baseline update help does not name --promote\n%s\n' "$help" >&2; exit 1; }; go test -count=1 -run '^(TestBaselineUpdate|TestHistorySanitize|TestHistorySnapshot)' ./internal/cli` — expected: exit 0. Before this Task the six tests do not exist and the usage has no `--promote`, so it fails; after it they pass, Surface Transcripts 2 to 4 are asserted whole, the built help names the flag, and every existing `baseline update` and History Sanitize Command test passes unchanged.

## References

- `_prd.md` → Core Feature 4; Core Feature 5; Core Feature 6; Goals; User Story 3; User Story 4; User Story 5; Success Metric 5; Success Metric 6
- `_techspec.md` → Data Models; API Contract 7; API Contract 8; API Contract 9; API Contract 10; Invariants 7-9; Surface Transcript 2; Surface Transcript 3; Surface Transcript 4; Build Order 4
- ADR-0260; ADR-0254; ADR-0173; ADR-0248

## Result

Implemented the task_04 slice for Daemon Verification. Task status remains
Daemon-owned; this record makes no terminal Task verdict.

### Implementation

- The history plan scans selected folders and removal kinds with
  `spec.HistoryCitations`, excludes reduction kinds, and suppresses duplicate
  folder citations on an already cited file/line. It classifies Live Documents
  and reports promotion offers. A detector failure blocks only history, clears
  selected writes and a proposed tag, and retains a one-line diagnosis.
- `history.citations` and `history.staleCitations` are always arrays. Citation
  objects have the specified lowercase fields; `offer` and `promotedTo` are
  omitted when empty. Text lists citations and stale restore commands between
  refusals and the tag, including an otherwise-current history section.
- The digest includes ordered path/line/target citations and each folder's
  sorted promotions. Both additions are omitted when empty; stale advice stays
  outside the digest and the update's state, category and exit code.
- Repeatable `--promote` rejects `--no-history`, unclean/absolute paths, core
  artifacts, non-selected files, Refused Units, existing destinations and
  duplicate basenames before any write. Valid requests re-plan through
  `PlanLegacyConversion`. A promotion reads committed history at the planned
  revision, copies its bytes into references read by later Specs, and uses the
  conversion writer's exclusive creation so an existing reference is never
  overwritten. Citing files are not rewritten. History recovery now names
  `docs/references` as well as the History Root.
- `spec.StaleHistoryCitations` reads tracked Markdown outside the History Root,
  resolves relative destinations and bounded repository paths, and looks each
  candidate up in the existing tag-file set. Directories are not candidates
  from that set. Present files and history sources are excluded. Stale scan
  failures are advisory and preserve the selected plan and update outcome.
- Command usage includes its own `[--promote <path> ...]` line, the citation
  offer, stale advice, and the promotion's provenance and collision hazard.
  Existing usage snippets and History Sanitize code/assertions are unchanged.

### Acceptance evidence

| Acceptance criterion | Implementation and focused evidence |
| --- | --- |
| Preview lists every citation, Live Document offers, and binds the digest | `TestBaselineUpdateListsHistoryCitationsTranscript` asserts Surface Transcript 2 through actual `go run`, including stderr `exit status 3` and process exit 1, exact JSON keys, and unchanged tree/index/refs. `TestBaselineUpdateHistoryDigestBindsCitations` moves a citing line, observes a new digest, and rejects old approval without writes. `TestBaselineUpdateHistoryCitationBoundaries` covers folder deduplication, core/live text, removals and omitted reductions. |
| Only requested promotion copies upstream; citing files stay byte-identical | `TestBaselineUpdatePromotesOnlyOnRequestTranscript` asserts Surface Transcript 3 through actual `go run`, the record's promotion, equality with the removed committed blob, unchanged citing files, a read-only promotion preview, a changed promotion digest, and no copy without the flag. `TestBaselineUpdateRefusesAnInvalidPromotion` covers eleven refusals with pending guidance writes, preserving tree/index/refs. Repeated promotions and shared basename refusal have an additional focused test. |
| Existing tag names stale targets and restore commands without changing outcome | `TestBaselineUpdateReportsStaleHistoryCitationsTranscript` asserts Surface Transcript 4 through actual `go run`, exact stale JSON, state `current`, empty category and exit 0. Changing stale advice leaves the digest and outcome identical. The two authored `TestStaleHistoryCitations...` tests cover relative/path citations, directories, present files, untracked sources and excluded history sources. Separate command tests cover local citation-detector blocking and advisory stale-detector failure. |
| No citations preserves today's digest and text | `TestBaselineUpdateWithoutCitationsKeepsDigestAndText` compares with the pre-ADR-0260 digest payload and field order, asserts the unchanged text block, and checks both new JSON fields are empty arrays. Selected existing update and History Sanitize tests pass without assertion changes. |

### Focused checks

All successful commands below exited 0. These are implementation checks,
not execution of the Task's declared Verification commands.

- `go test -count=1 -run '^TestStaleHistoryCitationsReportsFilesTheTagHolds$' ./internal/spec`
  initially failed compilation because `StaleHistoryCitations` did not exist,
  establishing the missing API before implementation.
- `go test -count=1 -v -run '^TestStaleHistoryCitations' ./internal/spec`:
  both authored tests passed after implementation.
- `go test -count=1 -run 'TestBaselineUpdate.*(Citations|OnRequestTranscript|InvalidPromotion)' ./internal/cli`:
  the six authored command tests passed.
- `go test -count=1 -v -run '^TestBaselineUpdate.*Transcript$' ./internal/cli`:
  all three executable transcript checks passed with private adopter homes,
  reused Go build/module caches and dependency downloads disabled.
- `go test -count=1 -v -run 'TestBaselineUpdate.*(Citations|OnRequestTranscript|InvalidPromotion|CitationBoundaries|CitationDetectorFailureIsLocal|StaleDetectorFailureKeepsOutcome|PromotionsRepeatAndRefuseSharedBasenames)' ./internal/cli`:
  the six authored tests, four added boundary tests, and the existing relocation
  citation digest test passed after the production refinements.
- `go test -count=1 -v -run '^TestBaselineUpdateRefusesAnInvalidPromotion$' ./internal/cli`:
  all eleven cases passed after strengthening the fixture to include pending
  guidance writes.
- `go test -count=1 -run '^(TestBaselineUpdatePlansPendingHistoryWithoutWrites|TestBaselineUpdateAppliesPendingHistoryAndCreatesTheTag|TestBaselineUpdateNoHistorySkipsTheSection)$' ./internal/cli`:
  three existing update boundaries passed unchanged.
- `go test -count=1 -run '^(TestBaselineUpdateHistoryPlanningBoundaries|TestBaselineUpdateHistoryHonorsAnExistingTag|TestBaselineUpdateReportsATagFailureBeforeAnyConversion|TestBaselineUpdateRefusesUnitsWithUncommittedChanges|TestHistorySanitizePlanWritesNothing|TestHistorySanitizePromotesIntoReferences|TestHistorySanitizeUsageErrors|TestHistorySanitizeFolderCitationsAndKindNamedSlugs|TestHistorySanitizeRefusesAPromotionInARefusedUnit|TestHistorySanitizePromotionSeparatesFolderFromKind)$' ./internal/cli`:
  ten selected existing update/sanitize tests passed unchanged.
- `go test -count=1 -tags docscontract -run '^TestBaselineDocumentationContract$' ./internal/docscontract`:
  existing pinned help and documentation snippets passed unchanged.
- `git diff --check`: no whitespace errors.

An added boundary check initially found an unused result variable in its new
test; asserting the selected-unit count fixed the compile error, and the
subsequent focused run passed. The configured sandbox could not start because
its writable worktree path contains a symlink; reviewed commands used the
assigned worktree for implementation and disposable local fixtures for tests.

### Verification feedback — attempt 1

The Daemon's configured verification reached the documentation contracts and
reported a stale `command baseline update` fingerprint in
`TestTheSkillCoverageMapIsCurrent`. The task's intentional help change had not
been reflected in the generated Behavior Surface Record. The diagnostic source
was `batch-004-attempt-1.log` under this Run's verification artifacts; this
Result records the finding, not the log body.

Regenerated `docs/references/behavior-surfaces.json` with its owning recorder.
Only the `command baseline update` fingerprint changed. The schema, surface
keys and every other fingerprint remain identical to the committed record.
This ordinary generated path is necessary fallout of task_04's help change;
no production code, assertion, tooling configuration or other Task was changed
in the repair.

Focused evidence:

- `go test -count=1 -tags docscontract -run '^TestTheSkillCoverageMapIsCurrent$' ./internal/docscontract -record-skill-coverage`:
  exited 0 and regenerated the record through its owning implementation.
- A read-only JSON comparison against `HEAD:docs/references/behavior-surfaces.json`
  confirmed that exactly `command baseline update` changed, with unchanged
  schema and surface keys.
- `go test -count=1 -tags docscontract -run '^(TestTheSkillCoverageMapIsCurrent|TestBaselineDocumentationContract)$' ./internal/docscontract`:
  exited 0 without record mode; the coverage and pinned help contracts passed.
- `git diff --check`: no whitespace errors after the Result update.

The full configured verification and the Task's declared Verification commands
were not rerun in this repair turn. Task status and the terminal verdict remain
Daemon-owned. The implementation is handed back for its configured rerun.

### Handoff

The pre-existing change was the Daemon's `pending` to `in_progress` status
transition. This session changed the declared task_04 code/test paths, refreshed
the generated baseline-update help fingerprint in
`docs/references/behavior-surfaces.json`, and updated this Result. No other Task, Task Graph, commit, push or pull request
was changed or created. Declared Task Verification, repository Verification
and terminal status settlement remain with the Daemon. No follow-up outside
this Task's slice was implemented.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `docs/references/behavior-surfaces.json`

## Carry-forward provenance

- Source Run: `run_20261012T110224Z_380dd005763aa861`
- Source commit: `fe63f87c5ec82c61bc3027e36ce4262b9f23f767`
