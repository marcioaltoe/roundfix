---
task: task_06
spec: 0272-the-history-holds-only-archive-records
status: completed
type: backend
complexity: low
---

# Task 06: The baseline update plan names every file it will remove

## Overview

User Story 3 asks that `roundfix baseline update` show every file it will
delete before the adopter approves. The history section prints and reports
each removal unit with only its file count and bytes: the QA gate of
2026-10-17 built a synthetic adopter with six retired files, and five of the
six paths appear nowhere in the `--format json` plan, although the Plan
Digest already binds every path. This corrective Task adds each unit's paths
to the report and to the text plan, without changing the Plan Digest, the
apply path or `roundfix history sanitize`. It is verifiable on its own through
the command tests on a fixture adopter.

## Requirements

1. MUST add `paths` to each history unit of the `roundfix baseline update`
   report in `internal/cli/baseline_update_history.go`: the sorted repository
   paths the unit removes for a `remove` or `convert` unit, and the Archive
   Records it rewrites for the `record-links` unit, the same set the Plan
   Digest's unit already binds. The JSON field is `paths`, present for every
   pending unit.
2. MUST print, in the text plan, one line `  - <path>` per path directly
   under its unit's line (`- <unit>: removes …`, `- folder …`,
   `- record-links: …`), in the same order as the JSON.
3. MUST NOT change the Plan Digest, the applied report, `--yes`, the
   `roundfix history sanitize` output (Surface Transcripts 1 and 2), or any
   other line of either command.
4. MUST add `TestBaselineUpdateNamesEveryFileItWillRemove` in
   `internal/cli/baseline_update_paths_test.go`, calling `t.Parallel()`
   first: a fixture adopter with an annotated `history-full` tag, a Reduced
   History Entry, a retired ADR in `docs/adr/`, a terminal Backlog Entry, an
   `_archived/findings/` file, an uncited `superseded` record and a record with
   a broken link; the dry run exits 3, the JSON lists each removed path once
   under its unit and the rewritten record under `record-links`, the text plan
   prints each of them as a `  - ` line, and `--yes` then removes exactly the
   listed paths.
5. MUST describe the listing in `docs/user-guide/commands/baseline.md` and in
   `.agents/skills/roundfix/references/baseline.md`, each with the phrase
   `names every file it will remove`, mirror it with
   `make skills-sync`, raise the `version` of `.agents/skills/roundfix/SKILL.md`
   and its mirror `skills/roundfix/SKILL.md`, record it with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   and re-record the Behavior Surface Record with
   `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`.
6. MUST leave every existing test passing unchanged. When this Task adds a
   test, regenerate the test inventory with
   `go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1`
   after its last test change.

## Subtasks

- [ ] Add each unit's paths to the report and the text plan.
- [ ] Add the command test on a fixture adopter.
- [ ] Update the baseline guide and skill reference, the skill version and the records.

## Acceptance Criteria

- [ ] The JSON and text plans name every file a unit removes and every record it rewrites.
- [ ] The Plan Digest and the history sanitize output are unchanged.

## Consumers

- `roundfix baseline update` history plan output — consumers: `./internal/cli`

## Context

- instruction: `docs/adr/0278-a-retirement-deletes-and-the-history-root-holds-only-archive-records.md`
- instruction: `docs/agents/specific-repository.md`
- interface: `internal/cli/baseline_update_history.go`
- interface: `internal/cli/baseline_update_history_test.go`
- creates: `internal/cli/baseline_update_paths_test.go`
- interface: `docs/user-guide/commands/baseline.md`
- interface: `.agents/skills/roundfix/references/baseline.md`
- interface: `skills/roundfix/references/baseline.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `docs/references/coverage-record.json`

## Verification

- `out="$(go test -count=1 -v -run '^TestBaselineUpdateNamesEveryFileItWillRemove$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestBaselineUpdateNamesEveryFileItWillRemove ' || { printf 'missing passing test TestBaselineUpdateNamesEveryFileItWillRemove\n' >&2; exit 1; }; go test -count=1 -run 'History|BaselineUpdate' ./internal/cli || exit 1; for f in docs/user-guide/commands/baseline.md .agents/skills/roundfix/references/baseline.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'names every file it will remove' || { printf 'missing the listing in %s\n' "$f" >&2; exit 1; }; done; make skills-sync-check || exit 1; go test -count=1 -run '^TestEveryOwnedSkillVersionIsRecorded$' ./skills || exit 1; go test -count=1 -tags docscontract -run '^TestTheSkillCoverageMapIsCurrent$' ./internal/docscontract || exit 1` — expected: exit 0. Before this Task the test does not exist and neither page states the listing, so the command fails.

## References

- `_prd.md` → User Story 3; Core Feature 4; Success Metric 5
- `_techspec.md` → API Contract 12
- ADR-0278

## Result

### Implementation and acceptance evidence

- The pending JSON plan now includes sorted `paths` for every selected unit:
  removed files for `remove` and `convert`, and rewritten Archive Records for
  `record-links`. The text plan prints each path directly under its unit as
  `  - <path>` in JSON order. The digest implementation and apply implementation
  are unchanged.
- `TestBaselineUpdateNamesEveryFileItWillRemove` calls `t.Parallel()` first and
  exercises an adopted temporary repository with an annotated `history-full`
  tag, a Reduced History Entry, a retired active ADR, a terminal Backlog Entry,
  an `_archived/findings/` file, an uncited superseded record, a broken-link
  record and a two-file Legacy Archive Folder. It asserts exit 3, exact sorted
  unit membership without duplicates, adjacent text path lines, preview tree,
  index and ref preservation, then `--yes` approval of the preview digest and
  exactly the listed removals with the surviving record repaired.
- Acceptance criterion 1: the focused new test passed after the implementation.
  Before the implementation, it failed because the conversion unit's JSON
  paths were empty.
- Acceptance criterion 2: unchanged tests
  `TestBaselineUpdateRemovesRetiredHistoryAndRewritesCitations` (explicit digest
  payload), `TestHistorySanitizePlanTranscript`,
  `TestHistorySanitizeApplyTranscript` and
  `TestBaselineUpdateHistoryPlanningBoundaries` (including the exact applied
  text summary) passed in the focused selection below.
- The baseline guide and canonical skill reference state that the plan
  "names every file it will remove" and describe JSON/text ordering.
  `make skills-sync` regenerated the mirror. The owned-skill record command
  raised both Roundfix version fields in both copies from 0.1.81 to 0.1.82.
  The Behavior Surface Record and test inventory were regenerated by their
  owning commands.

### Focused commands and outcomes

- `go test ./internal/cli -run '^TestBaselineUpdateNamesEveryFile' -count=1`:
  initially exit 1 for missing JSON paths; after the implementation, exit 0.
- `go test ./internal/cli -run '^Test(BaselineUpdateNamesEveryFileItWillRemove|BaselineUpdateRemovesRetiredHistoryAndRewritesCitations|HistorySanitize(Plan|Apply)Transcript|BaselineUpdateHistoryPlanningBoundaries|BaselineUpdatePlansPendingHistoryWithoutWrites)$' -count=1 -v`:
  exit 1. The new test, explicit digest test, both sanitize transcripts and
  all planning-boundary subtests passed. Only the existing exact pending text
  transcript test failed; its old expectation excludes the required new lines.
- `make skills-sync`: exit 0.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`:
  exit 0 after synchronizing the mirror. The earlier invocation before sync
  exited 1 with the expected canonical/mirror difference diagnostic.
- `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`:
  exit 0; regenerated the baseline guide fingerprint.
- `make baseline-digests`: exit 0, `changed: false`.
- `go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1`:
  generated the one-entry inventory addition and the test body passed, but the
  process exited 1 because suiteguard reported the generated coverage record
  as a repository-boundary violation. The existing coverage generator never
  calls `suiteguard.DeclareSanctionedRegeneration`; this remains an out-of-slice
  follow-up, despite the command and output being authorized in this Spec.
- `go test ./internal/spec -run '^TestCoverageEquivalence$' -count=1`:
  exit 0; independently confirmed the generated inventory without rewriting it.

### Verification Feedback repair

Inspected the Daemon diagnostic artifact
`/Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20261017T074441Z_7010289b20c58c7a/verification/batch-001-attempt-1.log`.
It identifies three exact pending-plan transcript expectations that exclude
the new path lines: `TestBaselineUpdatePlansPendingHistoryWithoutWrites`,
`TestBaselineUpdateListsHistoryCitationsTranscript` and
`TestBaselineUpdateWithoutCitationsKeepsDigestAndText`.

Following the repair instruction, updated only those expected transcripts in
`internal/cli/baseline_update_history_test.go` and
`internal/cli/baseline_update_citations_test.go` to include the required path
lines. Expected paths use the fixture paths and constants, independently of
the new report field. All existing digest, exit-code, citation, refusal,
preview-preservation and applied-output assertions remain intact. The earlier
clarification request no longer blocks this repair.

Fresh focused evidence after the last test edit:

- `go test ./internal/cli -run '^TestBaselineUpdate(PlansPendingHistoryWithoutWrites|ListsHistoryCitationsTranscript|WithoutCitationsKeepsDigestAndText|NamesEveryFileItWillRemove|HistoryPlanningBoundaries|RemovesRetiredHistoryAndRewritesCitations)$' -count=1 -v`:
  exit 0; all six tests and all planning-boundary subtests passed. This covers
  the three reported mismatches, the exact JSON/text listing and `--yes`
  removal set (acceptance criterion 1), both digest compatibility oracles and
  the unchanged applied summary (acceptance criterion 2). The sanitize
  implementation and its transcript tests are unchanged by this repair;
  their passing focused evidence is recorded above.
- `go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1`:
  exit 0 after the last test change, with no further inventory changes.
  The earlier suiteguard diagnostic remains recorded above; this invocation
  reports no repository-boundary violation.
- `git diff --check`: exit 0.

The Daemon's configured Verification sequence was not rerun in this repair
turn, and no terminal Task verdict is claimed.

The Daemon's pre-existing `status: in_progress` change is preserved. No Task
Graph or other Task file was edited. Authored Verification and repository
Verification remain Daemon-owned and were not run. No commit, push or pull
request was made.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `internal/cli/baseline_update_citations_test.go`
