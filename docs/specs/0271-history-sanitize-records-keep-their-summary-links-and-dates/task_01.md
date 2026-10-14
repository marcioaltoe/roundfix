---
task: task_01
spec: 0271-history-sanitize-records-keep-their-summary-links-and-dates
status: pending
type: backend
complexity: medium
---

# Task 01: A legacy record's summary drops block quote markers and falls back to a section

## Overview

A Legacy Archive Folder's summary is joined line by line with its block quote
markers kept, and a PRD that goes straight from its title to a section yields
no summary (`_techspec.md` → Current behavior). This Task changes the summary
reading under `Legacy` only, so the record of such a folder carries the
paragraph its author wrote. It is verifiable on its own through
`PlanLegacyConversion` on temporary repositories.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-13,
   [history sanitize records lose their summary, links and dates](references/2026-10-13-history-sanitize-records-lose-their-summary-links-and-dates.md),
   for its first two defects (the joined block quote and the missing summary).
2. MUST implement API Contract 1 in `BuildArchiveRecord` behind `in.Legacy`:
   every line of the summary paragraph, from the PRD title or from the newest
   QA Report's `## Outcome`, loses its block quote markers (up to three spaces,
   `>`, one optional space, repeated for nesting) before joining, and an empty
   title paragraph falls back to the first paragraph after the PRD's first
   `## ` heading. A QA Report outcome and a supersession reason keep their
   precedence.
3. MUST NOT change `archiveParagraph`'s result for a non-legacy call or for
   `ParseArchiveRecord` (API Contract 6): an active Spec's PRD holding the same
   block quote still yields today's joined text.
4. MUST add the tests named in Verification to the new
   `internal/spec/legacy_record_summary_test.go`, each calling `t.Parallel()`
   first: `TestLegacySummaryRemovesBlockQuoteMarkers` (a two-line and a nested
   quote, asserting the record body line equals the joined text without any
   `>`), `TestLegacySummaryFallsBackToTheFirstSection` (title then
   `## Problem`, asserting the record body holds that section's first
   paragraph and `ReadArchivedSpec` reads the same outcome),
   `TestLegacySummaryKeepsTheQAReportOutcome`, and
   `TestActiveArchiveRecordKeepsItsSummaryReading` (non-legacy
   `BuildArchiveRecord` on the same block quote returns exactly
   `> a > b`-style text as today). Every existing test MUST pass unchanged.

## Subtasks

- [ ] Add the legacy summary reader beside `archiveParagraph`.
- [ ] Use it for the PRD and the QA Report outcome under `Legacy`.
- [ ] Add the first-section fallback.
- [ ] Add the four tests.

## Acceptance Criteria

- [ ] A two-line or nested block quote summary becomes one line with no `>`.
- [ ] A PRD whose title is followed by a section yields that section's first paragraph.
- [ ] A QA Report outcome still wins over the PRD.
- [ ] A non-legacy record reads its summary exactly as before.

## Consumers

- Archive Record outcome of a Legacy Archive Folder conversion — consumers: `./internal/spec`, `./internal/cli`

## Context

- instruction: `docs/adr/0277-a-legacy-archive-folder-conversion-keeps-the-summary-links-and-dates-its-folder-held.md`
- instruction: `docs/adr/0251-a-legacy-archive-folder-is-read-leniently-and-a-failed-qa-keeps-its-verdict.md`
- interface: `internal/spec/archive_record.go`
- creates: `internal/spec/legacy_record_summary_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestLegacySummaryRemovesBlockQuoteMarkers|TestLegacySummaryFallsBackToTheFirstSection|TestLegacySummaryKeepsTheQAReportOutcome|TestActiveArchiveRecordKeepsItsSummaryReading)$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestLegacySummaryRemovesBlockQuoteMarkers TestLegacySummaryFallsBackToTheFirstSection TestLegacySummaryKeepsTheQAReportOutcome TestActiveArchiveRecordKeepsItsSummaryReading; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing passing test %s\n' "$name" >&2; exit 1; }; done; go test -count=1 -run 'Legacy|Archive|History' ./internal/spec ./internal/cli || exit 1` — expected: exit 0. Before this Task the four tests do not exist, so the PASS check fails; afterwards they and every legacy, archive and history test of both packages pass.

## References

- `_prd.md` → Goals; User Story 1; Core Feature 1; Success Metric 1; Success Metric 3
- `_techspec.md` → API Contract 1; API Contract 6; Invariants 1, 4 and 6; Current behavior; Testing Approach; Build Order 1
- ADR-0277; ADR-0251
