---
task: task_03
spec: 0287-reports-that-count-what-they-say
status: pending
type: backend
complexity: medium
---

# Task 03: An Archive Record summary carries what its colon introduces

## Overview

The Archive Record of Spec 0282 ends with "Two environment gaps failed it:",
because its summary is the PRD's title paragraph and the summary reader stops
at the blank line before the list that paragraph introduces (`_prd.md` →
Sources 3). This Task makes `archiveParagraph` carry the one list or paragraph
that follows a paragraph ending with a colon, for every record writer, while
every record already written reads back and renders to the same bytes. It is
verifiable on its own on temporary Spec folders.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-20,
   [`references/2026-10-20-an-archive-record-summary-ends-at-the-colon-that-introduces-its-list.md`](references/2026-10-20-an-archive-record-summary-ends-at-the-colon-that-introduces-its-list.md)
   (`_prd.md` → Sources 3), by implementing API Contract 4 of `_techspec.md`
   in `archiveParagraph` in `internal/spec/archive_record.go`, keeping its
   signature.
2. MUST carry exactly one following block, as API Contract 4 states: a list,
   whose items lose their markers, are whitespace-normalized and gain `;`, or
   `.` when last, if they lack closing punctuation; or a paragraph; and
   nothing before a heading, a fence, a table row or the end of the text. The
   carried text follows the paragraph on the same line.
3. MUST keep every caller on the one rule: the PRD title paragraph, a Legacy
   Archive Folder's section fallback after its block quote markers are
   removed, the QA Report's Outcome, the Cancel Command
   (`internal/spec/cancel.go`) and `ParseArchiveRecord`. The outcome
   precedence, Legacy Link Rebase and the size cut at
   `ArchiveRecordTargetBytes` do not change.
4. MUST state the hazard in the code's behavior: the data comes from the
   Spec's PRD, QA Report or written record, and its destination is the Archive
   Record's single summary line. No Archive Record already written is
   rewritten, and every one, including one whose summary ends with a colon,
   parses and renders to the same bytes, so the History Sanitize Command's and
   the pull request record update's round-trip checks keep holding.
5. MUST add these tests, each top-level test calling `t.Parallel()` first,
   compiling for Windows and using temporary folders; every existing test
   keeps its name and passes:
   - in `internal/spec/archive_record_test.go`,
     `TestArchiveSummaryCarriesWhatItsColonIntroduces`, through
     `BuildArchiveRecord`, with subtests: a title paragraph ending "Two gaps
     failed it:" before a tight two-item list whose items end with `.` gives
     "Two gaps failed it: `A` failed with a wrapped line. `B` timed out.";
     a loose list; an ordered list; items without closing punctuation give
     `;` between them and `.` after the last; a following paragraph; a
     following fence, a following heading and the end of the text each leave
     the paragraph unchanged; and a QA Report Outcome ending with a colon
     before an ordered list carries it;
   - in the same file, `TestAnArchiveRecordEndingWithAColonRoundTrips`: a
     rendered record whose summary is "Two environment gaps failed it:"
     parses to that summary and renders to the same bytes;
   - in `internal/spec/archive_record_cancelled_test.go`,
     `TestACancelledRecordSummaryCarriesItsList`: a cancelled Spec whose PRD
     title paragraph ends with a colon before a list leaves a record whose
     summary holds every item;
   - in `internal/spec/legacy_record_summary_test.go`,
     `TestLegacySummaryCarriesTheListItsColonIntroduces`: a Legacy Archive
     Folder whose block-quoted title paragraph ends with a colon before a
     block-quoted list.
6. MUST document the change in `docs/user-guide/commands/archive.md`,
   `docs/user-guide/commands/history.md` and
   `.agents/skills/roundfix/references/archive.md`, each with the sentence "A
   summary paragraph that ends with a colon carries the list or paragraph that
   follows it, on the same line." beside its description of the record or of
   where the summary is read. MUST search the guides, the agent guides and the
   skill copies for wording that says the summary is one paragraph only, and
   change each. No text inside any `### QA settlement` section may change.
7. MUST revise **Archive Record** in `CONTEXT.md` through domain-modeling:
   its definition gains that its summary is the QA Report's Outcome paragraph
   or the PRD's title paragraph, and that a summary paragraph that ends with a
   colon carries the list or paragraph that follows it.
8. MUST run these commands, because this Task edits user guide pages and the
   Roundfix reference. No record or mirror is hand-edited.
   Before them, raise both version fields of the Roundfix skill's canonical
   `SKILL.md` (`metadata.version` and the top-level `version`) by one patch
   over the value the tree holds when this Task starts, because the record
   command refuses changed content under an unchanged version. The Agent runs
   these commands itself before handing back: the Daemon only runs the
   Verification and never regenerates a mirror or a record.
   - `make skills-sync`;
   - `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
     to record the raised Roundfix skill version;
   - `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`,
     to re-record the Behavior Surface Record.

## Subtasks

- [ ] Carry the list or paragraph a colon introduces in `archiveParagraph`.
- [ ] Add the archive, round-trip, cancel and legacy summary tests.
- [ ] Document the rule in the archive and history guides and the reference.
- [ ] Revise **Archive Record** in the glossary; sync and record.

## Acceptance Criteria

- [ ] A PRD title paragraph, a QA Outcome, a cancelled Spec's PRD and a Legacy Archive Folder that end with a colon before a list each give a summary holding every item.
- [ ] A paragraph ending with a colon before a fence, a heading or the end of the text stays as today.
- [ ] A record already ending with a colon parses and renders to the same bytes.
- [ ] The guides, the reference and the glossary describe the rule.

## Consumers

- the Archive Record summary that `archiveParagraph` reads and `ParseArchiveRecord` parses — consumers: `./internal/spec`, `./internal/cli`

## Context

- instruction: `docs/adr/0247-an-archive-leaves-an-archive-record-and-the-spec-folder-stays-in-git.md`
- instruction: `docs/adr/0277-a-legacy-archive-folder-conversion-keeps-the-summary-links-and-dates-its-folder-held.md`
- instruction: `docs/adr/0262-a-cancelled-spec-leaves-an-archive-record-and-returns-its-sources.md`
- instruction: `internal/spec/cancel.go`
- instruction: `internal/spec/history_sanitize.go`
- interface: `internal/spec/archive_record.go`
- interface: `internal/spec/archive_record_test.go`
- interface: `internal/spec/archive_record_cancelled_test.go`
- interface: `internal/spec/legacy_record_summary_test.go`
- interface: `docs/user-guide/commands/archive.md`
- interface: `docs/user-guide/commands/history.md`
- interface: `.agents/skills/roundfix/references/archive.md`
- interface: `skills/roundfix/references/archive.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestArchiveSummaryCarriesWhatItsColonIntroduces|TestAnArchiveRecordEndingWithAColonRoundTrips|TestACancelledRecordSummaryCarriesItsList|TestLegacySummaryCarriesTheListItsColonIntroduces|TestLegacySummaryRemovesBlockQuoteMarkers|TestLegacySummaryFallsBackToTheFirstSection|TestLegacySummaryKeepsTheQAReportOutcome|TestActiveArchiveRecordKeepsItsSummaryReading|TestArchiveRecordRoundTrips|TestArchiveRecordStaysWithinTheTargetSize)$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestArchiveSummaryCarriesWhatItsColonIntroduces TestAnArchiveRecordEndingWithAColonRoundTrips TestACancelledRecordSummaryCarriesItsList TestLegacySummaryCarriesTheListItsColonIntroduces TestActiveArchiveRecordKeepsItsSummaryReading; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; for f in docs/user-guide/commands/archive.md docs/user-guide/commands/history.md .agents/skills/roundfix/references/archive.md CONTEXT.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'ends with a colon carries the list or paragraph that follows it' || { printf 'missing summary rule in %s\n' "$f" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Archive Record**:' || exit 1; go test -count=1 -run 'Archive|Cancel|Legacy|History' ./internal/spec || exit 1; go test -count=1 -run 'Archive|Cancel|History' ./internal/cli || exit 1; GOOS=windows go build -buildvcs=false ./internal/spec || exit 1; make skills-sync-check || exit 1; go test -count=1 ./skills || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the four new tests do not exist, so the command fails at its first pass check. Afterwards the summary, round-trip, cancel and legacy tests pass with the existing archive, cancel and history tests of both packages, the guides, the reference and the glossary state the rule, the package builds for Windows, and the skill version and the Behavior Surface Record are current.

## References

- `_prd.md` → Sources; Goals; Core Feature 3; Core Feature 4; Success Metric 3; Glossary
- `_techspec.md` → API Contract 4; API Contract 5; Invariants 5-6; Interfaces; Current behavior; Testing Approach; Build Order 3
- ADR-0247; ADR-0277; ADR-0262; ADR-0259; ADR-0285; ADR-0189; ADR-0256; ADR-0280
