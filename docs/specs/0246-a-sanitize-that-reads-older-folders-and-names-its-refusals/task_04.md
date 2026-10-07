---
task: task_04
spec: 0246-a-sanitize-that-reads-older-folders-and-names-its-refusals
status: completed
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

## Result

Implemented the Task 04 slice for Daemon Verification; status remains
Daemon-owned.

- Split `PlanHistoryKind` from `PlanHistoryKinds`, retaining the latter's
  fail-fast behavior and existing error text. The CLI now plans each selected
  kind independently.
- Folder inventory and conversion errors and kind inventory/planning errors
  become ordered Refused Units. Selection fills `--batch` with convertible
  units and stops conversion planning at its limit. Folder inventories beyond
  that limit are not walked. Remaining counts include refused units.
- Apply checks tag existence, annotation and ancestry before conversion
  planning, and checks path coverage only for selected convertible units.
  Promotions into examined Refused Units refuse the entire command with the
  unit's reason; promotions outside selection, duplicate destinations and
  promotion planning errors still refuse before writing. Promotion replanning
  distinguishes a folder from a kind with the same name.
- Plan prints tolerance lines between each folder and its candidates, followed
  by ordered refusal lines before citations. Plan and apply counts include the
  refusal suffix only when needed. Apply prints refusals before its confirmation
  and exits 2 with the specified reason when all examined units are refused.
  Help explains that refused units remain in place and do not consume batch
  capacity.
- Added the six requested synthetic CLI tests and three focused edge-case
  tests. Changed only `TestHistorySanitizePreflightsWholeBatch` among existing
  tests; its batch-one assertion is preserved. No record builder, parser or
  lenient reader was changed.

Acceptance evidence from focused implementation checks:

| Acceptance criterion | Evidence |
| --- | --- |
| Plan lists every Refused Unit and exits 0 | `TestHistorySanitizePlanListsEveryRefusedUnit` checks two ordered refusal reasons, the suffix, exit 0 and unchanged repository bytes/status. `TestHistorySanitizeRefusesAKindUnitAndPlansTheRest` checks a malformed Finding without blocking folders or later kinds. |
| Batch fills past refusals and preserves their bytes | `TestHistorySanitizeApplySkipsARefusedUnitAndFillsTheBatch` and the updated whole-batch test check two written records and identical refused-folder bytes. The all-refused and refused-promotion tests check exit 2 and unchanged repository bytes/status. Extra tests cover inventory symlinks, the batch boundary, coverage excluding a refused folder absent from the tag, and folder/kind name collisions during promotion. |
| Four adopter shapes convert together | `TestHistorySanitizeConvertsTheFourLegacyShapes` constructs a commented node still in the table, several omitted graph rows, a retired `refactor` row, and three failing QA reports. It checks all named tolerances, a read-only plan, four applied parseable records, and the newest failed-QA report without an override or QA completion. |
| Other existing history tests retain their expectations | The focused CLI sanitize selection passes, including the exact-output plan test and unchanged tag, dirty-tree, promotion, advice, empty-plan and configured-root cases. Focused spec history/legacy/archive tests pass. Diff inspection confirms that only the requested existing test changed. |

Commands and outcomes:

- `GOCACHE=/tmp/roundfix-task04-gocache rtk proxy go test -count=1 -run 'TestHistorySanitize' ./internal/cli`
  — exit 0 on the final code/test changes (`ok roundfix/internal/cli`).
- `GOCACHE=/tmp/roundfix-task04-gocache rtk proxy go test -count=1 -run 'History|Legacy|Archive' ./internal/spec`
  — exit 0 (`ok roundfix/internal/spec`).
- Baseline reproduction using a Go overlay containing the two production
  files from committed `HEAD`:
  `GOCACHE=/tmp/roundfix-task04-gocache rtk proxy go test -overlay=/tmp/roundfix-task04-baseline/overlay.json -count=1 -run 'TestHistorySanitizePlanListsEveryRefusedUnit|TestHistorySanitizeApplySkipsARefusedUnitAndFillsTheBatch' ./internal/cli`
  — expected exit 1: both new regressions observe exit 2 at the first malformed
  PRD under the original implementation. The overlay changes no repository
  file.
- `rtk proxy git -c core.fsmonitor=false diff --check` — exit 0.

The initial worktree change was the Daemon's `status: in_progress` update in
this Task file. It is preserved. Declared Verification and repository delivery
gates were not run; those remain with the Daemon. No commit, push or Pull
Request was created. No follow-up outside this slice was implemented.
