---
task: task_02
spec: 0243-a-history-that-holds-only-records
status: completed
type: backend
complexity: medium
---

# Task 02: Retired Findings and Backlog Entries reduce, and retired reviews and handoffs leave

## Overview

Besides Spec folders, `docs/history` holds four other retired kinds. The
Spec check still reads retired Findings and Backlog Entries: an
`absorbed_by` license, a Rollup member and a `closure_evidence` path each
name one by file. Nothing reads retired Review Artifacts or handoffs. This
Task plans and applies each kind by what still reads it. A Finding or
Backlog Entry becomes a Reduced History Entry: the front matter unchanged,
the title, the first paragraph and one line naming the revision and path of
its full text. Review Artifacts and handoffs are removed. Retired ADRs are
never touched. The Task also scans Markdown outside the History Root for
citations of removed paths, so the plan can show them. It answers the
Backlog Entry "History keeps only what the Secondbrain needs" of
2026-10-06, Expected 1 ("the same rule applies to the other history kinds").

## Requirements

1. MUST add `ReduceHistoryEntry`, `IsReducedHistoryEntry`,
   `PlanHistoryKinds`, `ApplyHistoryKind` and `HistoryCitations` in
   `internal/spec/history_entries.go` with the signatures of `_techspec.md`
   → Interfaces and the behavior of Invariants 2, 7, 8 and 9.
   - Kinds come in the order `findings`, `backlog`, `reviews`, `handoffs`,
     from `ArchiveDir`. A kind with nothing pending is omitted.
   - `adr/` is never planned, read for rewriting or removed.
   - `ReduceHistoryEntry` keeps the front matter byte for byte and refuses a
     file without front matter, without a `# ` title, or already reduced.
   - `ApplyHistoryKind` writes only the planned files and removes
     directories it leaves empty, and nothing else.
   - `HistoryCitations` reads tracked `*.md` files outside the History Root
     through `git ls-files` in the given repository and reports
     `path:line` and the removed path it names. It never refuses.
2. MUST name every unexported helper this Task adds with a `historyEntry`
   prefix, because task_01 adds files to the same package in parallel.
3. MUST add the `internal/spec` tests named in Verification to
   `internal/spec/history_entries_test.go`. Each builds its own temporary
   tree or repository and never reads this repository's `docs/history`:
   - a Finding and a Backlog Entry reduce to their exact expected bytes;
   - a reduced entry is not pending again, and each refusal of
     Requirement 1 holds;
   - the kind plan's order, file lists and byte counts before and after,
     with an ADR left byte-identical;
   - applying `reviews` and `handoffs` removes their files and empty
     directories and nothing else;
   - the citation scan reports a link and an inline path to a removed
     folder, and ignores a file under the History Root.
4. MUST add `internal/speccheck/reduced_history_entry_test.go` with the
   `internal/speccheck` tests named in Verification. They run the Spec check
   over a fixture whose retired Findings were reduced:
   - an `absorbed_by` naming an Archive Record's slug, an active Rollup
     member, and a `closure_evidence` naming a reduced Backlog Entry all
     still pass;
   - a reduced Finding whose `absorbed_by` names nothing still fails with
     the archive-license code, so the reduction cannot hide a broken
     license.
5. MUST NOT change any file under this repository's `docs/history`, any
   Spec check rule, or `internal/spec/archive.go`.

## Subtasks

- [ ] Add the reduction and its detector.
- [ ] Add the kind plan, apply and the citation scan.
- [ ] Add the spec and speccheck tests.

## Acceptance Criteria

- [ ] A reduced Finding or Backlog Entry keeps every front-matter field and
      still satisfies the Spec check.
- [ ] Reviews and handoffs leave, and ADRs stay byte-identical.
- [ ] Citations of removed paths are reported, never refused.

## Context

- creates: `internal/spec/history_entries.go`
- creates: `internal/spec/history_entries_test.go`
- creates: `internal/speccheck/reduced_history_entry_test.go`
- instruction: `docs/adr/0248-existing-history-is-sanitized-in-batches-after-a-history-full-tag.md`
- instruction: `internal/spec/archive.go`
- instruction: `internal/speccheck/citations.go`
- instruction: `docs/agents/docs-layout.md`

## Verification

- `out="$(go test -count=1 -v ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestReduceHistoryEntryKeepsFrontMatterTitleAndFirstParagraph TestReduceHistoryEntryRefusesWhatItCannotReduce TestReducedHistoryEntryIsNotPendingAgain TestPlanHistoryKindsOrdersAndMeasures TestApplyHistoryKindRemovesReviewsAndHandoffsOnly TestHistoryCitationsNameMarkdownThatCitesARemovedPath; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name (" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task none of the six tests exists, so the command fails; after it the whole `internal/spec` package passes with them.
- `out="$(go test -count=1 -v -run '^(TestReducedFindingsKeepTheirLicensesAndMembers|TestReducedFindingWithABrokenLicenseStillFails)$' ./internal/speccheck 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestReducedFindingsKeepTheirLicensesAndMembers TestReducedFindingWithABrokenLicenseStillFails; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name (" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task neither test exists, so the command fails.

## References

- `_prd.md` → Goal 4; User Story 1; Core Feature 4
- `_techspec.md` → Measured inventory; Interfaces; Invariants 2, 7, 8 and 9; Data Models; Build Order 2
- ADR-0248; ADR-0120; ADR-0215

## Result

Implemented this Task's slice for Daemon Verification. Status remains
Daemon-owned. The initial worktree had only the Daemon's `pending` to
`in_progress` status change in this file; `internal/spec/history_entries.go`
was absent.

- Added reduction and detection with byte-identical front matter, the first
  title, the first non-heading paragraph with collapsed whitespace, and the
  exact Git provenance line. Missing front matter, missing title and already
  reduced entries are refused. The detector accepts 40- and 64-hex revisions
  only on the last non-empty line.
- Added ordered kind planning through `ArchiveDir`, sorted pending paths,
  before/after byte counts and planned replacement bytes. Applying a kind
  uses only those paths, refuses ADR and out-of-kind plans, and removes only
  parent directories left empty by its deletions. Already reduced entries
  are omitted from later plans.
- Added tracked Markdown citation reporting through `git ls-files`, outside
  the History Root, with source path, line and removed target. Folder-prefix
  citations, inline paths and relative Markdown links are reported as advice.
- Added the six requested `internal/spec` tests and two requested
  `internal/speccheck` tests. Every fixture uses a temporary tree or repository.

Focused check after the final code edits:

```sh
GOCACHE=/private/tmp/roundfix-task02-gocache rtk proxy go test ./internal/spec ./internal/speccheck -run 'Test(ReduceHistoryEntry|ReducedHistoryEntry|PlanHistoryKinds|ApplyHistoryKind|HistoryCitations|ReducedFindings|ReducedFindingWith)' -count=1 -v
```

Exit 0; all eight requested test functions and their subtests passed.

Acceptance evidence:

- Metadata and Spec check: exact-byte Finding and Backlog reduction fixtures
  passed, including CRLF front matter and comments. The Spec check fixture
  accepted a reduced Finding licensed by an Archive Record, a reduced member
  licensed by its active Rollup, and closure evidence naming a reduced
  Backlog Entry, with no findings. The broken-license fixture still reported
  `CodeArchiveLicense` at the preserved `absorbed_by` line.
- Removal and ADR preservation: kind ordering and byte counts matched exact
  expectations; applying plans left the ADR byte-identical. Removal deleted
  nested Review Artifact files, handoffs and their empty parent directories,
  preserving other families, reference files and a file added after planning.
- Advisory citations: the scan reported the link, inline path and relative
  link to the removed folder, ignored Markdown under the History Root and
  untracked Markdown, and avoided a similarly prefixed sibling folder.

No repository `docs/history` files, Spec check rules or `internal/spec/archive.go`
were changed. No follow-up work was identified. The declared Verification
commands and repository-wide gates were not run; the Daemon owns the declared
Verification and settlement. No commit, push or Pull Request was made.

## Carry-forward provenance

- Source Run: `run_20261007T040343Z_bb67f44effbee11e`
- Source commit: `6845a2ba5c0f5a06019f99e36063ab99d19c86be`
