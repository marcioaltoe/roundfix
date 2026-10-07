---
task: task_04
spec: 0246-a-sanitize-that-reads-older-folders-and-names-its-refusals
status: pending
type: backend
complexity: high
---

# Task 04: The History Sanitize Command lists every Refused Unit and fills the batch past it

## Overview

Stops the History Sanitize Command from aborting at the first unit it cannot
convert. The plan names each Refused Unit with its reason and the tolerances
of each folder it converts. `--apply --batch <n>` converts the next `n` units
it can convert, leaves each Refused Unit untouched and lists it. This Task
answers the "Report every refusal" direction of the Backlog Entry "History
sanitize refuses older archived Specs and aborts the whole plan" of
2026-10-07, under the maintainer's decision "Pular e seguir". It is verifiable
on its own through the CLI tests below.

## Requirements

1. MUST split `PlanHistoryKinds` into
   `PlanHistoryKind(repositoryRoot, revision string, kind ArchiveKind) (HistoryKindPlan, error)`.
   `PlanHistoryKinds` keeps its behavior by looping over it.
2. MUST plan units as `_techspec.md` → Invariant 9 states. These errors each
   make that unit a Refused Unit whose reason is the error text, and planning
   continues:
   - an inventory error of a folder or kind;
   - a `PlanLegacyConversion` error;
   - a `PlanHistoryKind` error.

   With `--batch <n>`, planning stops after the `n`th convertible unit, and
   Refused Units do not count toward `n`.
3. MUST keep the whole-command refusals of `_techspec.md` → Invariant 10. The
   tag's existence, annotation and ancestry are checked before planning, and
   its coverage of every path of the selected convertible units after
   planning. A promotion is matched to its inventoried folder. A promotion in
   a Refused Unit refuses with
   `promotion "<p>" is in refused unit <unit>: <reason>`, and one outside the
   selected convertible units refuses with today's `is outside the batch`
   message.
4. MUST print the plan per `_techspec.md` → Invariant 11 and Surface
   Transcript 1: the
   `; <r> unit(s) refused` suffix, the `tolerates <folder>: <tolerance>` lines
   from `LegacyConversion.Tolerated` and the `refused <unit>: <reason>` lines.
   With no Refused Unit and no tolerance, the output is byte-identical to
   today's.
5. MUST print apply per `_techspec.md` → Invariants 12 and 13 and Surface
   Transcripts 2 and 3. Refused lines
   come first, then the confirmation with `; <r> unit(s) refused` before the
   remaining count. Apply exits 2 with the Preflight reason
   `history sanitize --apply found no convertible unit; <r> unit(s) refused`
   when every examined unit is refused. A Refused Unit's bytes never change.
6. MUST add to the usage text that Refused Units are listed with their reason,
   left in place and not counted toward `--batch`.
7. MUST add `internal/cli/history_refusal_test.go`, using the existing
   fixture helpers and synthetic folders only:
   - `TestHistorySanitizePlanListsEveryRefusedUnit`: two malformed folders
     among clean ones both appear as `refused` lines, the first line carries
     `; 2 unit(s) refused`, the exit is 0 and nothing is written.
   - `TestHistorySanitizeApplySkipsARefusedUnitAndFillsTheBatch`: the first
     of three folders is refused, and `--batch 2` writes the other two records
     and leaves it byte-identical.
   - `TestHistorySanitizeApplyRefusesWhenEveryExaminedUnitIsRefused`.
   - `TestHistorySanitizeRefusesAPromotionInARefusedUnit`.
   - `TestHistorySanitizeRefusesAKindUnitAndPlansTheRest`: a retired Finding
     without front matter refuses `findings`, and the folders still plan.
   - `TestHistorySanitizeConvertsTheFourLegacyShapes`. It uses synthetic
     folders shaped like the adopter's: graph nodes commented out but still in
     the table, several rows left out of the graph, a `refactor` type row, and
     three `verdict: fail` reports without an override. All four convert in one
     apply, and the plan prints their `tolerates` lines and `failed-qa`.
8. MUST change the expectation of `TestHistorySanitizePreflightsWholeBatch` in
   `internal/cli/history_test.go`, and only that test. With a malformed `bbb`,
   `--apply --batch 2` now exits 0, writes `aaa.md` and `ccc.md`, prints
   `refused docs/history/specs/bbb: ` and leaves `bbb` byte-identical. The
   `--batch 1` assertion is unchanged.
9. MUST NOT change the record builder, the parser, the lenient reading or any
   other existing test expectation.

## Subtasks

- [ ] Split kind planning per kind.
- [ ] Keep a Refused Unit instead of returning, and fill the batch past it.
- [ ] Move tag coverage after planning, and refuse promotions into Refused
      Units.
- [ ] Print the `tolerates` and `refused` lines and the counts.
- [ ] Update the usage text and write the CLI tests.

## Acceptance Criteria

- [ ] The plan lists every Refused Unit with its reason and exits 0.
- [ ] A batch converts the next units it can convert and leaves each Refused
      Unit untouched.
- [ ] The four adopter shapes convert in one apply.
- [ ] Every other existing history test passes unchanged.

## Context

- interface: `internal/cli/history.go`
- interface: `internal/spec/history_entries.go`
- interface: `internal/cli/history_test.go`
- creates: `internal/cli/history_refusal_test.go`
- instruction: `internal/spec/history_sanitize.go`
- instruction: `internal/spec/archive_record.go`
- instruction: `docs/user-guide/commands/history.md`
- instruction: `docs/adr/0251-a-legacy-archive-folder-is-read-leniently-and-a-failed-qa-keeps-its-verdict.md`
- instruction: `docs/adr/0248-existing-history-is-sanitized-in-batches-after-a-history-full-tag.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestHistorySanitizePlanListsEveryRefusedUnit|TestHistorySanitizeApplySkipsARefusedUnitAndFillsTheBatch|TestHistorySanitizeApplyRefusesWhenEveryExaminedUnitIsRefused|TestHistorySanitizeRefusesAPromotionInARefusedUnit|TestHistorySanitizeRefusesAKindUnitAndPlansTheRest|TestHistorySanitizeConvertsTheFourLegacyShapes|TestHistorySanitizePreflightsWholeBatch)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for t in TestHistorySanitizePlanListsEveryRefusedUnit TestHistorySanitizeApplySkipsARefusedUnitAndFillsTheBatch TestHistorySanitizeApplyRefusesWhenEveryExaminedUnitIsRefused TestHistorySanitizeRefusesAPromotionInARefusedUnit TestHistorySanitizeRefusesAKindUnitAndPlansTheRest TestHistorySanitizeConvertsTheFourLegacyShapes TestHistorySanitizePreflightsWholeBatch; do printf '%s\n' "$out" | grep -q -- "--- PASS: $t " || { printf 'missing passing test %s\n' "$t" >&2; exit 1; }; done` — expected: exit 0. Before this Task the six new tests do not exist. After it, the command lists every Refused Unit, a batch fills itself past one, and the changed whole-batch test passes under its new expectation.

## References

- `_prd.md` → Core Feature 3; Goals; Success Metric 1; Success Metric 3; Success Metric 5
- `_techspec.md` → API Contract 3; Invariant 9; Invariant 10; Invariant 11; Invariant 12; Invariant 13; Surface Transcript 1; Surface Transcript 2; Surface Transcript 3; Build Order 4
- ADR-0251
- ADR-0248
