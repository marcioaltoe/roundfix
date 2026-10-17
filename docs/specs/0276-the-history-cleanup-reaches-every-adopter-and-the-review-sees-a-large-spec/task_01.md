---
task: task_01
spec: 0276-the-history-cleanup-reaches-every-adopter-and-the-review-sees-a-large-spec
status: pending
type: backend
complexity: high
---

# Task 01: The history cleanup writes a successor's lineage line, refuses only an unsucceeded ADR and repairs live links

## Overview

Today one superseded ADR whose successor does not name it refuses the whole
`adr` unit, and no unit repairs a Live Document's link to a file an earlier
layout removed (`_techspec.md` → Current behavior). This Task makes the history
plan that `roundfix history sanitize` and `roundfix baseline update` share
append the missing lineage line to the successor, refuse only an ADR without
an accepted successor by its path, and add the `live-links` unit. It is
verifiable on its own through planner and command tests on temporary
repositories.

## Requirements

1. MUST answer the adopted Backlog Entry of 2026-10-17,
   [one unnamed superseded ADR keeps an adopter's whole history](references/2026-10-17-one-unnamed-superseded-adr-keeps-an-adopters-whole-history.md),
   by implementing API Contracts 1 and 2 of `_techspec.md` in
   `internal/spec/history_entries.go`: `HistoryKindPlan` gains `Lineage` and
   `Refused`, `PlanHistoryKind` fills them for the `adr` kind instead of
   returning `historyADRSuccessor`'s error, and `ApplyHistoryKind` appends each
   lineage line after the removals, never twice (ADR-0282).
2. MUST write the lineage line with exactly the bytes of API Contract 1:
   `Supersedes ADR-<number> (<date>); Git keeps its text at <locator>.`, the
   date being the committer date of `HEAD` in UTC and the locator
   `ResolveGitLocator(<removed path>, "", "").Text(false)`; MUST NOT edit any
   other byte of the successor; MUST list the line as a History Citation in
   both commands, so the `baseline update` Plan Digest binds it, and include
   the successor in the unit's uncommitted-change check.
3. MUST report each superseded ADR without an accepted successor as a Refused
   Unit named by its repository-relative path with the reason of API Contract
   2, in text and JSON of both commands, while the `adr` unit removes every
   other retired ADR. Refusals never change either command's exit code
   (ADR-0254).
4. MUST implement API Contract 3 in a new `internal/spec/live_links.go`
   (`PlanLiveLinks`), reusing `MissingRelativeLinks`, `ResolveGitLocator` and
   the link-span rewrite of `record_links.go`, and wire the `live-links` unit
   after `record-links` into `historyInventory`, the selection loops, the
   frozen rewrite plan and both printers of `internal/cli/history.go` and
   `internal/cli/baseline_update_history.go`, with the text lines and JSON
   shape of API Contract 3. The selection of pages MUST equal the pages
   `SC-HISTORY-RETAINED`'s link scan reads other than Archive Records; the
   gate itself MUST NOT change.
5. MUST add the tests named in Verification, each calling `t.Parallel()`
   first: in the new `internal/spec/history_lineage_test.go`,
   `TestHistoryADRPlanWritesTheLineageASuccessorLacks` (exact line bytes,
   date, locator, single application) and
   `TestHistoryADRPlanRefusesOnlyAnADRWithoutAnAcceptedSuccessor` (no
   `superseded_by`, and a `proposed` successor, each refused by path while a
   rejected ADR in the same unit is planned); in the new
   `internal/spec/live_links_test.go`,
   `TestLiveLinksPlanRewritesMissingLinksInLiveDocumentsAndGuides` and
   `TestLiveLinksPlanLeavesRecordsCodeAndRemovedFiles`; in the new
   `internal/cli/history_adopter_test.go`,
   `TestHistorySanitizeAppliesTheAdopterShapes` (plan and `--apply` on a
   repository holding the four shapes of Success Metric 1, then
   `speccheck.Check` reports `SC-HISTORY-RETAINED` only for the kept ADR, and
   a second plan lists only that refusal) and
   `TestBaselineUpdatePlansLineageAndLiveLinks` (JSON lists the `live-links`
   unit with its paths, the lineage citation and the refusal, and the Plan
   Digest differs when the successor's planned line would differ).
6. MUST change `TestHistorySanitizeRefusesAKindUnitAndPlansTheRest` in
   `internal/cli/history_refusal_test.go`, keeping its name: its successor
   becomes `status: proposed` and it asserts `refused docs/adr/0001-old.md:
   superseded ADR has no accepted successor` while the other units plan. This
   is the declared break of its old expected text. MUST NOT rename or remove
   any other test; when a fixture of `internal/cli/baseline_update_paths_test.go`
   gains a `live-links` unit, its expected units gain it.
7. MUST document the lineage line, the refusal by path and the `live-links`
   unit, with its exact text line, in `docs/user-guide/commands/history.md`,
   `docs/user-guide/commands/baseline.md`,
   `.agents/skills/roundfix/references/archive.md`,
   `.agents/skills/roundfix/references/baseline.md` and the history-unit
   sentence of `.agents/skills/roundfix/references/setup.md`, and replace the
   sentence that says a superseded ADR whose successor does not name it refuses
   the unit, searching the guides and skill copies for that old wording. The
   `history sanitize` usage text names the `live-links` repairs. No text inside
   any `### QA settlement` section may change.
8. MUST run `make skills-sync`, raise and record the Roundfix skill version
   with `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   and re-record the Behavior Surface Record with
   `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`,
   because this Task edits user guide pages and command sources. No record or
   mirror is hand-edited.

## Subtasks

- [ ] Plan and apply the lineage line and the per-file refusals.
- [ ] Add `PlanLiveLinks` and wire the `live-links` unit into both commands.
- [ ] Add the six tests and move the refusal test to the new text.
- [ ] Document the units in the guides and the Roundfix references.
- [ ] Sync the mirrors, record the version and the Behavior Surface Record.

## Acceptance Criteria

- [ ] A superseded ADR whose accepted successor does not name it is removed, and the successor ends with the exact lineage line once.
- [ ] A superseded ADR without an accepted successor is refused by its path while every other retired ADR leaves.
- [ ] Missing links in Live Documents and user guide pages become the Git Locators the history gate prints, and a second plan has no `live-links` unit.
- [ ] Both commands list the lineage citation, the refusal and the unit in text and JSON, and the Plan Digest binds the line.

## Consumers

- The history plan, its text and its `baseline update` JSON — consumers: `./internal/spec`, `./internal/cli`, `./internal/docscontract`

## Context

- instruction: `docs/adr/0282-the-history-cleanup-names-a-superseded-adr-in-its-successor-and-a-large-candidate-is-reviewed-in-segments.md`
- instruction: `docs/adr/0278-a-retirement-deletes-and-the-history-root-holds-only-archive-records.md`
- instruction: `docs/adr/0254-baseline-update-sanitizes-pending-history-in-the-change-it-plans.md`
- instruction: `internal/speccheck/history_retained.go`
- instruction: `internal/spec/record_links.go`
- interface: `internal/spec/history_entries.go`
- creates: `internal/spec/live_links.go`
- creates: `internal/spec/live_links_test.go`
- creates: `internal/spec/history_lineage_test.go`
- interface: `internal/cli/history.go`
- interface: `internal/cli/baseline_update_history.go`
- interface: `internal/cli/history_refusal_test.go`
- interface: `internal/cli/baseline_update_paths_test.go`
- creates: `internal/cli/history_adopter_test.go`
- interface: `docs/user-guide/commands/history.md`
- interface: `docs/user-guide/commands/baseline.md`
- interface: `.agents/skills/roundfix/references/archive.md`
- interface: `skills/roundfix/references/archive.md`
- interface: `.agents/skills/roundfix/references/baseline.md`
- interface: `skills/roundfix/references/baseline.md`
- interface: `.agents/skills/roundfix/references/setup.md`
- interface: `skills/roundfix/references/setup.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`

## Verification

- `out="$(go test -count=1 -v -run '^(TestHistoryADRPlanWritesTheLineageASuccessorLacks|TestHistoryADRPlanRefusesOnlyAnADRWithoutAnAcceptedSuccessor|TestLiveLinksPlanRewritesMissingLinksInLiveDocumentsAndGuides|TestLiveLinksPlanLeavesRecordsCodeAndRemovedFiles)$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestHistoryADRPlanWritesTheLineageASuccessorLacks TestHistoryADRPlanRefusesOnlyAnADRWithoutAnAcceptedSuccessor TestLiveLinksPlanRewritesMissingLinksInLiveDocumentsAndGuides TestLiveLinksPlanLeavesRecordsCodeAndRemovedFiles; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; out="$(go test -count=1 -v -run '^(TestHistorySanitizeAppliesTheAdopterShapes|TestBaselineUpdatePlansLineageAndLiveLinks|TestHistorySanitizeRefusesAKindUnitAndPlansTheRest)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestHistorySanitizeAppliesTheAdopterShapes TestBaselineUpdatePlansLineageAndLiveLinks TestHistorySanitizeRefusesAKindUnitAndPlansTheRest; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; for f in docs/user-guide/commands/history.md docs/user-guide/commands/baseline.md .agents/skills/roundfix/references/archive.md .agents/skills/roundfix/references/baseline.md .agents/skills/roundfix/references/setup.md internal/cli/history.go; do grep -qF -- 'live-links' "$f" || { printf 'missing live-links in %s\n' "$f" >&2; exit 1; }; done; for f in docs/user-guide/commands/history.md docs/user-guide/commands/baseline.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'Git keeps its text at' || { printf 'missing the lineage line in %s\n' "$f" >&2; exit 1; }; done; make skills-sync-check || exit 1; go test -count=1 ./skills || exit 1; go test -count=1 ./internal/spec || exit 1; go test -count=1 -run 'History|BaselineUpdate' ./internal/cli || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the four planner tests do not exist, so the command fails at its first pass check; afterwards the planner, command, documentation, skill version and Behavior Surface Record checks pass.

## References

- `_prd.md` → Goals; User Story 1; User Story 2; Core Feature 1; Core Feature 2; Core Feature 9; Success Metric 1
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; API Contract 10; Invariants 1, 2 and 3; Current behavior; Testing Approach; Build Order 1
- ADR-0282; ADR-0278; ADR-0254; ADR-0260
