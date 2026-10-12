---
task: task_04
spec: 0255-an-archive-that-relinks-and-an-update-that-lists-citations
status: pending
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
