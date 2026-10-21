---
task: task_01
spec: 0295-the-archive-rewrites-the-satisfied-by-links-it-records
status: pending
type: backend
complexity: medium
---

# Task 01: The archive rebases each unproven entry it records

## Overview

`BuildArchiveRecord` rebases the record's outcome with the Legacy Link Rebase
but copies the `unproven` entries, the `satisfied-by` text of the Spec's
Unreachable Acceptance declarations, unchanged, so a relative link in an entry
breaks once it lands in `docs/history/specs/<slug>.md` and the strict check
fails right after the archive. This Task applies the same rebase to each
entry, proves it through `spec.Archive`, the Legacy Archive Folder conversion
and the commands, and states it in the `archive` guide and the glossary.

## Requirements

1. MUST answer the Inbox Entry of 2026-10-20,
   [`references/2026-10-20-o-archive-copia-links-relativos-das-declaracoes-para-o-archive-record.md`](references/2026-10-20-o-archive-copia-links-relativos-das-declaracoes-para-o-archive-record.md)
   (`_prd.md` → Sources 1), by implementing API Contract 1 of `_techspec.md`
   in `internal/spec/archive_record.go` `BuildArchiveRecord`: when `in.Rebase`
   is not nil, after the outcome, pass each `unproven` entry to
   `in.Rebase.Rebase` with `<in.Rebase.Folder>/_prd.md` as its source and
   `in.Rebase.RecordPath` as its destination, write the results into a new
   slice (never into the caller's `in.Unproven` or the PRD's slice), keep the
   order and count, and return `rebase archive unproven entry <n>: <error>`
   on failure. A nil `Rebase` keeps today's behavior. MUST NOT add a second
   rewrite: the entries and the outcome use the same `LegacyLinkRebase` value
   (Invariant 3).
2. MUST state the hazard in the code's behavior: the entries come from the
   Spec's `_prd.md`, written by the Spec's author against the Spec folder;
   the record that receives them is read from `docs/history/specs/` by
   `roundfix spec check --strict`, the History Sanitize Command and every
   later reader after the folder is gone. A destination inside the folder
   MUST end as the promoted copy or a Git Locator, never as a link that
   resolves under the archive root.
3. MUST keep everything else (API Contracts 2-3): `internal/spec/archive.go`,
   `internal/cli/archive.go` and the Delivery Queue are not changed; the
   archive's output, exit codes, refusals, record fields and outcome bytes are
   unchanged; an entry without a relative link keeps its bytes; every
   existing test name and assertion is unchanged.
4. MUST add, in the new file `internal/spec/archive_unproven_links_test.go`,
   the two tests `_techspec.md` → Testing Approach names, each top-level test
   calling `t.Parallel()` first, compiling for Windows, and using temporary
   repositories and Git configuration:
   - `TestArchiveRebasesTheUnprovenLinksItRecords`, through `spec.Archive`
     with a 40-hex `SourceRevision` at `HEAD` and the repository root, on a
     Spec whose newest QA Report is a qualifying `partial` with its declared
     rows, with subtests that each assert the record's whole `unproven` list:
     archived on the default branch, a removed folder file becomes
     `` `references/removed.md` (`<folder>/references/removed.md@<SourceRevision>`) ``;
     archived with a pull request from a branch whose commit no default-branch
     commit or tag reaches, a folder file only that branch holds becomes the
     label followed by the folder path and `(pull request #<n>)`; a file promoted with `Promote` becomes
     `` [`references/promoted.md`](../../references/promoted.md) ``; an entry
     without a link keeps its bytes; a link to a repository file outside the
     folder keeps naming that file; and under an external Spec Root, whose
     record is written to `<root>/_archived/<slug>.md`, the promoted link is
     the relative path from that record to `docs/references/promoted.md`.
     Each subtest also asserts that `MissingRelativeLinks` finds no link in
     the written record.
   - `TestLegacyConversionRebasesTheUnprovenLinksItRecords`, through the
     Legacy Archive Folder conversion, on a legacy `_prd.md` whose
     `unproven` holds a link to `references/<file>` and a plain line: the
     link becomes the folder file's Git Locator and the plain line keeps its
     bytes.
5. MUST add, in the new file `internal/cli/archive_unproven_links_test.go`,
   `TestArchiveOfAQualifyingPartialPassesTheStrictHistoryCheck`, calling
   `t.Parallel()` first, which builds the Characterization Repository of
   `_techspec.md` → Current behavior (with every file on `main`), runs
   `archive 0001-demo --promote references/promoted.md` and then, after
   committing the archive, `spec check --strict 0002-other` through the
   command entry point, and asserts that the archive exits 0, the record's
   two entries are the ones Success Metric 1 names, and the check's output
   holds no `SC-HISTORY-RETAINED` line.
6. MUST document the change (`_prd.md` → Core Feature 4). In
   `docs/user-guide/commands/archive.md`, change "In the outcome and each
   promoted copy, a citation of a removed folder file" to "In the outcome,
   each `unproven` entry and each promoted copy, a citation of a removed
   folder file", and right after the sentence ending "so it names what was
   never verified." add: "Each entry is read as text of the Spec's
   `_prd.md`, so its relative links are rewritten like the outcome's, and an
   entry without one keeps its bytes (ADR-0308)." MUST re-record the
   Behavior Surface Record, whose `archive` guide fingerprint changes, by
   running
   `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`.
   The Agent runs it itself before handing back; the record is never
   hand-edited, and `docs/references/skill-coverage.json` does not change.
7. MUST revise two glossary terms in `CONTEXT.md` through `domain-modeling`
   (`_prd.md` → Glossary), changing no other entry:
   - **Legacy Link Rebase**: its definition starts "The rewrite of relative
     destinations in an Archive Record outcome, its unproven entries and each
     promoted copy:", keeps the rest of today's sentence, adds "An unproven
     entry is read as text of the folder's `_prd.md`." and cites ADR-0308
     beside ADR-0278 as a link relative to `CONTEXT.md`.
   - **Lenient Legacy Reading**: its sentence "Legacy Link Rebase is part of
     this reading for the outcome and promoted copies." becomes "Legacy Link
     Rebase is part of this reading for the outcome, the unproven entries and
     promoted copies.", citing ADR-0308 beside its ADRs.

## Subtasks

- [ ] Rebase each `unproven` entry in `BuildArchiveRecord`.
- [ ] Add the package tests and the command test.
- [ ] Update the `archive` guide and re-record the Behavior Surface Record.
- [ ] Revise the two glossary terms.

## Acceptance Criteria

- [ ] A qualifying `partial` archive writes `unproven` entries with no relative link into the removed folder, and the strict check passes right after it.
- [ ] An entry without a relative link, and every existing archive, legacy conversion and eligibility test, are unchanged.
- [ ] The `archive` guide and both glossary terms state the rewrite, and the Behavior Surface Record is current.

## Consumers

- `BuildArchiveRecord` and the Archive Record it renders — consumers: `./internal/spec`, `./internal/cli`

## Context

- instruction: `docs/adr/0308-an-archive-rebases-the-unproven-links-it-records-and-the-record-link-repair-points-a-promoted-target-at-its-copy.md`
- instruction: `docs/adr/0277-a-legacy-archive-folder-conversion-keeps-the-summary-links-and-dates-its-folder-held.md`
- instruction: `docs/adr/0247-an-archive-leaves-an-archive-record-and-the-spec-folder-stays-in-git.md`
- instruction: `internal/spec/archive.go`
- instruction: `internal/spec/legacy_link_rebase.go`
- instruction: `internal/spec/history_sanitize.go`
- instruction: `internal/spec/history_sanitize_unproven_test.go`
- instruction: `internal/spec/record_links.go`
- instruction: `internal/cli/archive_test.go`
- instruction: `internal/cli/implement_test.go`
- instruction: `internal/docscontract/skill_coverage_test.go`
- creates: `internal/spec/archive_unproven_links_test.go`
- creates: `internal/cli/archive_unproven_links_test.go`
- interface: `internal/spec/archive_record.go`
- interface: `docs/user-guide/commands/archive.md`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestArchiveRebasesTheUnprovenLinksItRecords|TestLegacyConversionRebasesTheUnprovenLinksItRecords|TestLegacyUnprovenMapsBecomeOneLineEach|TestLegacyUnprovenRefusesANestedMap|TestActiveSpecStillRefusesUnprovenMaps|TestSpec0058ReplayArchivesDeclaredUnreachableRelease|TestArchiveQAOverrideRecordsARelativeOutcome|TestArchiveAppliesTheOneEligibilityDecision|TestLegacyConversionWritesARecordGitStillHolds)$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestArchiveRebasesTheUnprovenLinksItRecords TestLegacyConversionRebasesTheUnprovenLinksItRecords TestLegacyUnprovenMapsBecomeOneLineEach TestSpec0058ReplayArchivesDeclaredUnreachableRelease; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; out="$(go test -count=1 -v -run '^(TestArchiveOfAQualifyingPartialPassesTheStrictHistoryCheck|TestRunArchiveDeclaredUnreachableContract|TestRunArchiveUsesConfiguredExternalSpecRoot)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestArchiveOfAQualifyingPartialPassesTheStrictHistoryCheck TestRunArchiveDeclaredUnreachableContract; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < docs/user-guide/commands/archive.md | grep -q -- 'In the outcome, each .unproven. entry and each promoted copy, a citation of a removed folder file' || { echo 'missing archive guide sentence' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Legacy Link Rebase**: The rewrite of relative destinations in an Archive Record outcome, its unproven entries and each promoted copy:' || { echo 'missing Legacy Link Rebase revision' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Lenient Legacy Reading**:' || exit 1; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'Legacy Link Rebase is part of this reading for the outcome, the unproven entries and promoted copies.' || { echo 'missing Lenient Legacy Reading revision' >&2; exit 1; }; go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' || exit 1; GOOS=windows go vet ./internal/spec ./internal/cli` — expected: exit 0. Before this Task the new tests do not exist, so the command fails at its first pass check. Afterwards the record's entries are rebased through the package and the command, the existing legacy, replay, override and eligibility tests pass, the guide and both glossary terms carry their new text, the Behavior Surface Record is current, and both packages' tests compile for Windows.

## References

- `_prd.md` → Sources; Goals; Core Features 1 and 4; Behavior changes; Success Metric 1; Success Metric 3; Glossary
- `_techspec.md` → Current behavior; Interfaces; API Contracts 1-3; Invariants 1-3; Testing Approach; Build Order 1
- ADR-0308; ADR-0247; ADR-0277; ADR-0278; ADR-0230; ADR-0260; ADR-0261; ADR-0251; ADR-0240; ADR-0259; ADR-0271; ADR-0285; ADR-0280
