---
task: task_02
spec: 0286-the-history-conversion-rewrites-every-live-citation
status: pending
type: backend
complexity: medium
---

# Task 02: A lineage line names a retired ADR by its identifier and dates its supersession

## Overview

The adopter's successor of a date-named ADR received a `Supersedes` line naming
`ADR-` and the year `2026`,
dated with the day of the update rather than the day of the supersession
(`_prd.md` → Sources 1, defect 2). This Task names a retired ADR `ADR-NNNN`
only when its file name is a number followed by a slug, names any other ADR by
its file name without `.md`, and dates the line with the retired ADR's
`updated_at`. It is verifiable on its own through `PlanHistoryKind` on
temporary repositories.

## Requirements

1. MUST answer defect 2 of the Inbox Entry of 2026-10-19,
   [`references/2026-10-19-a-conversao-do-historico-deixa-citacoes-vivas-quebradas.txt`](references/2026-10-19-a-conversao-do-historico-deixa-citacoes-vivas-quebradas.txt)
   (`_prd.md` → Sources 1), by implementing API Contract 3 of `_techspec.md`
   (ADR-0299, which amends ADR-0282) in `internal/spec/history_entries.go`:
   `historyADRIdentifier` and `historyADRSupersessionDate` as
   `_techspec.md` → Interfaces sketches them, used by `PlanHistoryKind` and
   `historyADRNamesPredecessor`.
2. MUST keep the line's shape exactly `Supersedes <identifier> (<date>); Git
   keeps its text at <locator>.`, and MUST make `historyADRReferenceNames`
   accept `ADR-NNNN` and `NNNN` only for a numbered file name, while the stem,
   the file name and a path ending in the file name keep matching any name.
3. MUST take the date from the retired ADR's own front matter or Git history
   only. The fallback order is: `updated_at` converted to UTC; the committer
   date in UTC of the newest commit on `HEAD` that changed the file; `HEAD`'s
   committer date. An unparsable `updated_at` falls through to the next
   source and never refuses the unit.
4. MUST keep `TestHistoryADRPlanWritesTheLineageASuccessorLacks` unchanged and
   passing: its retired ADRs have no `updated_at`, and the newest commit that
   changed them is the dated fixture commit.
5. MUST add these tests in `internal/spec/history_adr_identity_test.go`, each
   top-level test calling `t.Parallel()` first, compiling for Windows and
   ending every child it starts in `t.Cleanup`:
   - `TestHistoryADRLineageNamesADateNamedADRByItsFileName`: the adopter's
     pair `docs/adr/2026-09-18-old-cluster.md` (`superseded_by:
     2026-10-01-new-cluster`, `updated_at: 2026-10-01T00:00:00Z`) and its
     accepted successor, committed on 2026-10-19, give exactly `Supersedes
     2026-09-18-old-cluster (2026-10-01); Git keeps its text at
     <locator>.`, and a successor that already names
     `2026-09-18-old-cluster` gets no line;
   - `TestHistoryADRLineageDatesTheSupersession`: a numbered ADR with
     `updated_at: 2026-10-05T22:30:00-03:00` is dated `2026-10-06`; one
     without `updated_at`, committed on 2026-10-03 under a later `HEAD` of
     2026-10-19, is dated `2026-10-03`; one with an unparsable `updated_at`
     uses its commit date;
   - `TestHistoryADRNumericReferencesSelectOnlyNumberedADRs`: `superseded_by:
     2026`, and the same number with the `ADR-` prefix, select no date-named
     ADR, so the ADR is refused with the existing reason, while a
     `superseded_by` naming a numbered ADR by its prefixed number still
     selects that `NNNN-<slug>.md`.

## Subtasks

- [ ] Derive the identifier from the file name.
- [ ] Derive the date from `updated_at`, then the commit dates.
- [ ] Restrict the numeric reference forms to numbered ADRs.
- [ ] Add the three tests.

## Acceptance Criteria

- [ ] A date-named retired ADR is named by its file name without `.md`, and a numbered one keeps `ADR-NNNN`.
- [ ] The date is the retired ADR's supersession day, whatever day the plan runs.
- [ ] The year alone, with or without the `ADR-` prefix, never selects a date-named ADR.

## Context

- instruction: `docs/adr/0299-the-history-conversion-rewrites-every-citation-it-lists-and-the-check-reports-a-removed-history-file-cited-without-a-git-locator.md`
- instruction: `docs/adr/0282-the-history-cleanup-names-a-superseded-adr-in-its-successor-and-a-large-candidate-is-reviewed-in-segments.md`
- instruction: `internal/spec/history_lineage_test.go`
- instruction: `internal/spec/history_successor_forms_test.go`
- interface: `internal/spec/history_entries.go`
- creates: `internal/spec/history_adr_identity_test.go`

## Consumers

- the lineage line and the History Citation it produces — consumers: `./internal/cli`

## Verification

- `grep -q 'func historyADRIdentifier' internal/spec/history_entries.go || { printf 'missing historyADRIdentifier\n' >&2; exit 1; }; grep -q 'func historyADRSupersessionDate' internal/spec/history_entries.go || { printf 'missing historyADRSupersessionDate\n' >&2; exit 1; }; out="$(go test -count=1 -v -run '^(TestHistoryADRLineageNamesADateNamedADRByItsFileName|TestHistoryADRLineageDatesTheSupersession|TestHistoryADRNumericReferencesSelectOnlyNumberedADRs|TestHistoryADRPlanWritesTheLineageASuccessorLacks|TestHistoryADRPlanRefusesOnlyAnADRWithoutAnAcceptedSuccessor)$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestHistoryADRLineageNamesADateNamedADRByItsFileName TestHistoryADRLineageDatesTheSupersession TestHistoryADRNumericReferencesSelectOnlyNumberedADRs TestHistoryADRPlanWritesTheLineageASuccessorLacks TestHistoryADRPlanRefusesOnlyAnADRWithoutAnAcceptedSuccessor; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; GOOS=windows go build -buildvcs=false ./internal/spec || exit 1; go test -count=1 -run 'History|Lineage|Successor|ADR' ./internal/spec ./internal/cli` — expected: exit 0. Before this Task neither function exists, so the command fails at its first check. Afterwards the three new tests and the two existing lineage tests pass, the package builds for Windows, and the history, lineage, successor and ADR tests of `internal/spec` and `internal/cli` pass.

## References

- `_prd.md` → Sources; Goals; Core Feature 2; Success Metric 2; Acceptance evidence
- `_techspec.md` → API Contract 3; Invariant 3; Current behavior; Testing Approach; Build Order 2
- ADR-0299; ADR-0282; ADR-0259; ADR-0285
