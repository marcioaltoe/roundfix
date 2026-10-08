---
task: task_02
spec: 0249-a-baseline-update-that-sanitizes-pending-history
status: pending
type: backend
complexity: medium
---

# Task 02: A legacy unproven list of maps converts and every refusal prints on one line

## Overview

Answers the Backlog Entry "A legacy `unproven` list of maps refuses its
folder" of 2026-10-08
([the adopted entry](references/2026-10-08-legacy-unproven-maps-refuse-a-folder.md)).
Under Lenient Legacy Reading, `BuildArchiveRecord` accepts the list-of-maps
`unproven` that older Roundfix versions wrote and turns each map into one
stable text line in the Archive Record. The History Sanitize Command prints
every Refused Unit reason on one line. Active Specs keep the string-only list.
The Task is verifiable on its own through spec and CLI tests on synthetic
folders.

## Requirements

1. MUST implement `_techspec.md` → API Contract 8 in `BuildArchiveRecord`.
   Only when `ArchiveRecordInput.Legacy` is true, read the PRD front matter
   `unproven` as a sequence of YAML nodes, keeping each scalar's literal text
   (a row `03` stays `03`). A scalar item keeps its text. A mapping item
   becomes:
   - `row <row>: ` when `row` exists;
   - then the value of `claim`;
   - then every other key in lexical order as `<key>: <value>` joined by `; `.
     Those keys are wrapped in ` (` and `)` after a claim, and stand alone
     without one.

   Values are whitespace-normalized, and a sequence of scalars is joined by
   `, `. An empty mapping, or a value that is a mapping or holds one, returns
   the error `legacy unproven item <i> cannot be read as text` with the 1-based
   item index. The resulting record MUST pass `PlanLegacyConversion`'s
   render-parse-render round trip unchanged.
2. MUST keep the non-legacy path byte-for-byte in behavior: without `Legacy`,
   `unproven` is decoded as `[]string`, and a list of maps fails with today's
   `parse archive PRD: yaml: unmarshal errors:` error. The Archive Command,
   `spec.Load` and the Archive Record schema do not change.
3. MUST make the folder branch of `ReadArchivedSpec` read the same lines.
   It already calls `BuildArchiveRecord` with `Legacy: true`, so this holds
   when requirement 1 lives in that builder and not in the sanitize caller.
4. MUST add `historyRefusalLine(err error) string` in the History Sanitize
   Command code. It joins the error text's whitespace-separated fields with
   single spaces. `printHistoryRefused` MUST print every
   `refused <unit>: <reason>` line through it. No other output of
   `roundfix history sanitize` changes.
5. MUST add these tests to the new
   `internal/spec/history_sanitize_unproven_test.go`, on folders in
   `t.TempDir()` built like `history_sanitize_legacy_test.go`'s fixtures:
   - `TestLegacyUnprovenMapsBecomeOneLineEach` uses the two key sets of the
     Fluxus report: `row`, `goal`, `claim` and `satisfied-by` (with a sequence
     value), and `row`, `claim` and `reason`. It adds a string item and a
     multi-line claim. It asserts the exact `Record.Unproven` lines, for
     example `row 03: Catalog verification runs alone (goal: G2; satisfied-by: task_04, task_05)`.
     It also asserts that `ApplyLegacyConversion` writes a record
     `ReadArchivedSpec` reads back with the same lines.
   - `TestLegacyUnprovenRefusesANestedMap` covers an empty mapping and a
     nested mapping value. Each refuses with the exact message and leaves the
     folder byte-identical.
   - `TestActiveSpecStillRefusesUnprovenMaps` gives an active Spec the same
     front matter. `BuildArchiveRecord` without `Legacy` fails with an error
     containing `cannot unmarshal !!map into string`, exactly as before.
6. MUST add `TestHistorySanitizeConvertsLegacyUnprovenMaps` and
   `TestHistorySanitizePrintsEachRefusalOnOneLine` to
   `internal/cli/history_refusal_test.go`, using its existing fixtures:
   - The first test plans and applies, with a tag, a folder whose `unproven`
     is a list of maps. It asserts the record's lines.
   - The second test plans a folder whose refusal error spans lines, such as
     a PRD with a malformed `unproven` that is not a sequence. It asserts that
     the plan prints exactly one `refused ` line for it, with no newline inside
     the reason, and that the reason still names the cause.
7. MUST NOT change any existing assertion in `internal/spec` or
   `internal/cli` tests.

## Subtasks

- [ ] Read legacy `unproven` as nodes and render each map as one line.
- [ ] Keep the active path string-only.
- [ ] Print Refused Unit reasons through `historyRefusalLine`.
- [ ] Write the spec and CLI tests.

## Acceptance Criteria

- [ ] A Legacy Archive Folder with a list-of-maps `unproven` converts, and its record holds one stable line per map.
- [ ] A nested or empty map refuses with a named one-line reason, and an active Spec still fails as before.
- [ ] Every `refused` line of the History Sanitize Command is one line.

## Context

- interface: `internal/spec/archive_record.go`
- creates: `internal/spec/history_sanitize_unproven_test.go`
- interface: `internal/cli/history.go`
- interface: `internal/cli/history_refusal_test.go`
- instruction: `internal/spec/history_sanitize.go`
- instruction: `internal/spec/history_sanitize_legacy_test.go`
- instruction: `docs/adr/0254-baseline-update-sanitizes-pending-history-in-the-change-it-plans.md`
- instruction: `docs/adr/0251-a-legacy-archive-folder-is-read-leniently-and-a-failed-qa-keeps-its-verdict.md`
- instruction: `docs/user-guide/commands/history.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestLegacyUnprovenMapsBecomeOneLineEach|TestLegacyUnprovenRefusesANestedMap|TestActiveSpecStillRefusesUnprovenMaps|TestHistorySanitizeConvertsLegacyUnprovenMaps|TestHistorySanitizePrintsEachRefusalOnOneLine)$' ./internal/spec ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for t in TestLegacyUnprovenMapsBecomeOneLineEach TestLegacyUnprovenRefusesANestedMap TestActiveSpecStillRefusesUnprovenMaps TestHistorySanitizeConvertsLegacyUnprovenMaps TestHistorySanitizePrintsEachRefusalOnOneLine; do printf '%s\n' "$out" | grep -q -- "--- PASS: $t " || { printf 'missing passing test %s\n%s\n' "$t" "$out" >&2; exit 1; }; done; go test -count=1 -run '^(TestLegacy|TestActiveSpecStillRefusesWhatLegacyReadingTolerates|TestReadArchivedSpec|TestArchiveRecord)' ./internal/spec && go test -count=1 -run '^TestHistorySanitize' ./internal/cli` — expected: exit 0. Before this Task the five new tests do not exist, so the first check finds no passing test and fails. After it, the new tests pass, and so do the existing legacy, archive record and History Sanitize Command tests.

## References

- `_prd.md` → Core Feature 6; Goals; User Story 5; Success Metric 6
- `_techspec.md` → API Contract 8; Invariant 11; Invariant 12; Build Order 2
- ADR-0254; ADR-0251; ADR-0248
- [Adopted Backlog Entry](references/2026-10-08-legacy-unproven-maps-refuse-a-folder.md)
