---
task: task_06
spec: 0272-the-history-holds-only-archive-records
status: pending
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
