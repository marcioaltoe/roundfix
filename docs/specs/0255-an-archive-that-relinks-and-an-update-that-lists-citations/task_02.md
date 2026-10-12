---
task: task_02
spec: 0255-an-archive-that-relinks-and-an-update-that-lists-citations
status: completed
type: backend
complexity: high
---

# Task 02: The owner's archive promotes cited files and relinks Live Documents

## Overview

Gives `roundfix archive` its Archive Relink (ADR-0260). Before it writes, the
archive finds every citation of the Spec folder from tracked Markdown outside
it, promotes each file outside the core artifacts that a Live Document cites,
points citations of the folder and its core artifacts at the Archive Record,
and rewrites each Live Document citation in place. It refuses a colliding
promotion or a dirty Live Document, warns about every citation it leaves, and
records the rewritten files under `relinked`. The shared classifier and rewrite
live in a new `internal/spec` file that task_03 and task_04 reuse. The Task is
verifiable on its own through unit tests and command tests on temporary
repositories.

## Requirements

1. MUST create the functions of `_techspec.md` → Interfaces in the new
   `internal/spec/live_citations.go`: `IsLiveDocument`,
   `ArchiveCitationMoves`, `RelinkContent`, `ArchiveRelinkPlan` and
   `PlanArchiveRelink`, behaving as `_techspec.md` → API Contracts 1, 2 and 3
   state. `RelinkContent` MUST read link destinations with the archive's
   existing Markdown scanner and match repository paths with the boundary rule
   of `spec.HistoryCitations`, and MUST change no byte other than a matched
   destination or path.
2. MUST call `PlanArchiveRelink` from `spec.Archive` after the eligibility and
   outward-link checks and before any write, only when `RepositoryRoot` is
   inside a Git work tree. Outside one, `spec.Archive` MUST behave exactly as
   today, so every existing `internal/spec` archive test passes unchanged
   (`_techspec.md` → Invariant 2).
3. MUST refuse as `_techspec.md` → API Contract 4 states, before any write:
   an existing promotion destination, two promotions with one basename, and a
   relinked path with changes not committed at `HEAD`, with the stated message
   for the last. A Git failure inside a work tree refuses. The command exits 2
   and changes no file (Invariant 4).
4. MUST add `Relinked []string` to `ArchiveRecord`, rendered as `relinked:`
   after `promoted` only when non-empty and read by `ParseArchiveRecord` when
   present, and MUST list automatic and manual promotions together, sorted,
   under `promoted`. A record without relinked files renders byte-identical to
   today (Invariant 1).
5. MUST write the record, the promoted copies and the relinked bytes before
   removing the folder, and on a write failure restore every relinked file and
   remove every file it created (`_techspec.md` → API Contract 5). Every
   promoted copy MUST equal its removed blob (Invariant 5).
6. MUST implement the command output of `_techspec.md` → API Contract 5 in
   `internal/cli/archive.go`: the `; relinked <m> file(s)` suffix, one
   `roundfix: warning: archive leaves <path>:<line> naming removed <target> unrelinked`
   line on stderr per Unrelinked citation, and the three `--plan` lines before
   `promote with:`. An archive without inbound citations prints exactly what it
   prints today. The hazard to name: a relinked file is a Live Document other
   people edit, so its bytes come only from `HEAD` (the dirty refusal) and it is
   never under the History Root (Invariant 10).
7. MUST add these tests to `internal/spec/live_citations_test.go`:
   - `TestIsLiveDocumentClassifiesLiveFamilies`: true for an active Spec file,
     an open Backlog Entry, an unresolved Finding and a `docs/references/` file;
     false for a handoff, the README, a guide, a file under the History Root and
     one under a non-built-in archive root.
   - `TestRelinkContentRewritesOnlyDestinationsAndPaths`: an inline link, an
     angle-bracket link, a reference definition, a link with a fragment, a
     folder link with a trailing slash, a repository path in prose, in a code
     span and in a fenced block are rewritten as API Contract 2 states; a path
     that shares the prefix of a move but is longer, an absolute URL and an
     unrelated link keep their bytes; the count equals the replacements.
   - `TestPlanArchiveRelinkPromotesLiveCitedFilesAndTargetsTheRecord`: in a
     temporary Git repository, a reference file cited by an active Spec is
     promoted, a citation of the owner's `_prd.md` targets the record, and a
     handoff's citation is Unrelinked.
   - `TestPlanArchiveRelinkIsEmptyOutsideAGitWorkTree`: the plan is empty and
     no error is returned.
   - `TestArchiveRecordRendersRelinkedOnlyWhenPresent`: a record with
     relinked files round-trips through render and parse, and one without
     renders byte-identical to a record built before the field existed.
8. MUST add these tests to `internal/cli/archive_relink_test.go`, built on the
   existing archive fixtures (`newImplementWorkspace`, `writeArchiveQAReport`,
   `commitArchiveFixture`) with a temporary home and a judge transport that
   rejects every request:
   - `TestArchiveRelinksLiveCitationsTranscript`: the repository of
     `_techspec.md` → Surface Transcripts; asserts Surface Transcript 1 whole
     (stdout line, stderr line, exit code), the promoted copy equal to the
     removed blob, each relinked file equal to its committed bytes with only
     the destination or path changed, the handoff unchanged, and the record's
     `promoted` and `relinked` lists.
   - `TestArchiveRefusesACollidingPromotion`: with
     `docs/references/2026-01-01-source.md` already present, exit 2 and every
     file and the Spec folder unchanged.
   - `TestArchiveRefusesADirtyLiveDocument`: with an uncommitted edit in the
     secondary Spec's `_prd.md`, exit 2 with the stated message and no file
     changed.
   - `TestArchiveWithoutInboundCitationsIsUnchanged`: no other file cites the
     Spec; stdout, stderr and the record equal what today's archive produces.
   - `TestArchivePlanListsRelinks`: `--plan` prints the promotes, relinks and
     leaves lines and writes nothing.
9. MUST NOT change any existing assertion in `internal/spec` or `internal/cli`
   tests. Every existing archive and archive record test passes unchanged.

## Subtasks

- [ ] Create the classifier, the moves and the rewrite with their unit tests.
- [ ] Plan the Archive Relink and wire it into the archive with its refusals.
- [ ] Add `relinked` to the Archive Record.
- [ ] Print the suffix, the warnings and the plan lines.
- [ ] Write the command tests, Surface Transcript 1 included.

## Acceptance Criteria

- [ ] Archiving the owner leaves every Live Document citation resolving, promotes the cited source with its bytes, and warns about the handoff.
- [ ] A colliding promotion or a dirty Live Document refuses with no file changed.
- [ ] An archive without inbound citations, and any archive outside a Git work tree, behaves exactly as today.

## Context

- creates: `internal/spec/live_citations.go`
- creates: `internal/spec/live_citations_test.go`
- interface: `internal/spec/archive.go`
- interface: `internal/spec/archive_record.go`
- interface: `internal/cli/archive.go`
- creates: `internal/cli/archive_relink_test.go`
- instruction: `internal/spec/history_entries.go`
- instruction: `internal/spec/archive_test.go`
- instruction: `internal/spec/archive_record_test.go`
- instruction: `internal/cli/archive_output_test.go`
- instruction: `internal/cli/archive_plan_test.go`
- instruction: `internal/cli/archive_links_test.go`
- instruction: `docs/adr/0260-an-archive-relinks-live-citations-and-a-baseline-update-lists-the-citations-it-breaks.md`
- instruction: `docs/adr/0247-an-archive-leaves-an-archive-record-and-the-spec-folder-stays-in-git.md`
- instruction: `docs/user-guide/commands/archive.md`
- instruction: `.agents/skills/roundfix/references/archive.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestIsLiveDocumentClassifiesLiveFamilies|TestRelinkContentRewritesOnlyDestinationsAndPaths|TestPlanArchiveRelinkPromotesLiveCitedFilesAndTargetsTheRecord|TestPlanArchiveRelinkIsEmptyOutsideAGitWorkTree|TestArchiveRecordRendersRelinkedOnlyWhenPresent)$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for t in TestIsLiveDocumentClassifiesLiveFamilies TestRelinkContentRewritesOnlyDestinationsAndPaths TestPlanArchiveRelinkPromotesLiveCitedFilesAndTargetsTheRecord TestPlanArchiveRelinkIsEmptyOutsideAGitWorkTree TestArchiveRecordRendersRelinkedOnlyWhenPresent; do printf '%s\n' "$out" | grep -q -- "--- PASS: $t " || { printf 'missing passing test %s\n%s\n' "$t" "$out" >&2; exit 1; }; done; go test -count=1 -run '^(TestArchive|TestParseArchiveRecord|TestBuildArchiveRecord|TestRenderArchiveRecord)' ./internal/spec` — expected: exit 0. Before this Task the five tests do not exist, so the first check fails; after it they pass and every existing archive test in the package passes unchanged.
- `out="$(go test -count=1 -v -run '^(TestArchiveRelinksLiveCitationsTranscript|TestArchiveRefusesACollidingPromotion|TestArchiveRefusesADirtyLiveDocument|TestArchiveWithoutInboundCitationsIsUnchanged|TestArchivePlanListsRelinks)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for t in TestArchiveRelinksLiveCitationsTranscript TestArchiveRefusesACollidingPromotion TestArchiveRefusesADirtyLiveDocument TestArchiveWithoutInboundCitationsIsUnchanged TestArchivePlanListsRelinks; do printf '%s\n' "$out" | grep -q -- "--- PASS: $t " || { printf 'missing passing test %s\n%s\n' "$t" "$out" >&2; exit 1; }; done; go test -count=1 -run '^TestArchive' ./internal/cli` — expected: exit 0. Before this Task the five command tests do not exist, so it fails; after it they pass, Surface Transcript 1 is asserted whole, and every existing archive command test passes unchanged.

## References

- `_prd.md` → Core Feature 1; Core Feature 2; Core Feature 3; Goals; User Story 1; User Story 2; Success Metric 1; Success Metric 2; Success Metric 3
- `_techspec.md` → Interfaces; API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; Invariants 1-5; Invariant 10; Surface Transcript 1; Build Order 2
- ADR-0260; ADR-0247; ADR-0083; ADR-0230

## Result

Implemented the Task 02 Archive Relink slice. The shared classifier, move map,
byte-preserving rewrite and planner live in `internal/spec/live_citations.go`.
The planner combines automatic and manual promotions, uses the history
citation detector and the archive Markdown scanner, and excludes the owner
folder, History Root and configured archive root. It reads both tracked and
HEAD paths so staged deletion or removal of a citation cannot hide a dirty
Live Document. New staged documents without relevant citations remain untouched.

`spec.Archive` plans after eligibility and outward-link checks, then refuses
promotion collisions and dirty relinked paths before writing. Relinked bytes
come from HEAD, with both Git status and a byte comparison guarding documents
other people edit. Read-only Git preflight disables optional index refreshes
and filesystem monitoring. Writes create the record and promoted copies,
rewrite Live Documents, and finally remove the owner folder. A write failure
restores attempted rewrites and removes every created file. Existing manual
promotion refusal messages retain the caller's original path spelling.

The Archive Record emits sorted promotions and relinked paths; `relinked` is
optional and follows `promoted`. Parsing accepts it when present. The command
prints the relink suffix, one warning per remaining citation, and deterministic
promotion, relink and unrelinked-citation lines in `--plan`.

### Acceptance evidence

- Live citations resolve, promoted bytes survive, and the handoff is warned:
  `TestArchiveRelinksLiveCitationsTranscript` asserts the entire stdout/stderr/
  exit transcript, the promoted source's exact bytes, all three Live Document
  rewrites against their committed bytes, the unchanged handoff and both record
  lists. The classifier, scanner forms, URL escaping, path boundaries and
  replacement counts are covered by the new shared-API tests. Configured archive
  roots and History Root exclusions have additional focused coverage.
- Collisions and dirty Live Documents refuse without changing files:
  `TestArchiveRefusesACollidingPromotion` covers occupied destinations and
  duplicate basenames. `TestArchiveRefusesADirtyLiveDocument` covers unstaged
  and staged edits, deletions with and without staging, removal of a citation,
  and a newly staged citing document, asserting the specified refusal and
  unchanged repository files. `TestArchiveRelinkRestoresFilesAfterALaterWriteFailure`
  proves restoration after an earlier rewrite succeeds and a later write fails,
  including removal of the record and promoted copy. A broken Git worktree
  pointer refuses instead of falling back to non-Git behavior.
- Archives without inbound citations and outside Git preserve existing behavior:
  `TestArchiveWithoutInboundCitationsIsUnchanged` compares stdout, stderr and
  record bytes with the existing archive construction. The record test pins
  the old rendering and checks optional-field round-trip parsing.
  `TestPlanArchiveRelinkIsEmptyOutsideAGitWorkTree` and selected existing
  library archive tests cover the non-Git path. Existing assertions were not
  changed. `TestArchivePlanListsRelinks` also asserts the added plan lines and
  proves the plan writes no repository files or provider requests.

### Focused checks

The first focused compile failed with the missing shared API and record field,
confirming the pre-change signal. After the final implementation edits:

- `rtk proxy go test ./internal/spec -run 'Test(IsLiveDocument|RelinkContent|PlanArchiveRelink|ArchiveRecordRendersRelinked|ArchiveRelinkRestores|ArchiveWritesTheArchiveRecord|ArchiveRecordRoundTrips|ArchiveRollsBackTheRecord|ArchiveKeepsLinksInside)' -count=1`
  — exit 0, `ok roundfix/internal/spec`.
- `rtk proxy go test ./internal/cli -run 'TestArchive(RelinksLive|RefusesAColliding|RefusesADirty|WithoutInbound|PlanListsRelinks|RelinkCombines|RelinkAllows|RelinkKeeps|CommandReportsRewritten|ConfirmationNames|PlanWithoutAKeyPrints|IgnoresMarkdown)' -count=1`
  — exit 0, `ok roundfix/internal/cli`.
- `rtk proxy git diff --check` — exit 0.

Command tests use disposable repositories and homes with a judge transport
that rejects every request. No provider request was made.

The declared Verification commands and repository Verification were not run;
the Daemon owns them and Task settlement. Status, `_tasks.md`, other Task
files and existing test assertions are unchanged by this implementation.
No commit, push or pull request was made. Queue staging and exact-retirement
proof belong to task_03; baseline update citation handling belongs to task_04.

## Carry-forward provenance

- Source Run: `run_20261012T110224Z_380dd005763aa861`
- Source commit: `205832921032d8636a2dec4d016031f155017da3`
