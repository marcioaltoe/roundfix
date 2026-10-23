---
task: task_04
spec: 0302-a-subsumed-test-selector-is-flagged-and-a-qa-refusal-names-its-commit
status: completed
type: docs
complexity: low
---

# Task 04: The qa-gate refusal template, the Context-driven development guide and the glossary name the audited commit

## Overview

task_03 makes the Daemon's precondition refusal and blocking mechanical
reports record `build`, and the gate's own refusal template in the `qa-gate`
skill still has none, so an Agent that refuses at its own strict check writes
the old shape (oraculum, 2026-10-21; the adopted Backlog Entry
`references/2026-10-22-a-qa-precondition-refusal-report-records-no-audited-commit.md`).
This Task adds `build` to that template, names it in the refusal paragraph of
the Context-driven development guide, revises **Precondition Refusal** and
**QA Report** in the glossary, raises the `qa-gate` skill's version and
re-records the Behavior Surface Record (`_prd.md` → Core Feature 4).

## Requirements

1. MUST change the refusal template in section `## 1. Resolve scope and
   preconditions` of the canonical `.agents/skills/qa-gate/SKILL.md` by adding
   the line `build: "<audited-commit>"` right after the template's `verdict:
   fail` line, with the template's indentation, and MUST add, right after the
   paragraph that ends "no rows for the checks that never executed.", with the
   same two-space indentation, this paragraph:
   "A refusal records `build` as a full report does: the commit the gate
   audited, quoted, so a reader can tell this refusal from one an earlier Run
   left on the branch. A refusal the Daemon wrote already carries it; keep its
   value."
   No text inside any `### QA settlement` section of any skill may change, and
   section `## 6. Close the report` keeps its bytes.
2. MUST change, in `docs/user-guide/context-driven-development.md`, the
   refusal paragraph's list "verdict `fail`, `rows_blocked_precondition`" to
   "verdict `fail`, `build` naming the commit the gate audited, as a full
   report does, `rows_blocked_precondition`", and add, right after that
   paragraph's last sentence "A gate that reached its matrix writes none of
   those three keys.", the sentence "The Daemon also records `build` in the
   report it closes for a blocking mechanical result, and the pending report
   it seeds for the gate carries none, because the gate's full report writes
   its own."
3. MUST revise **Precondition Refusal** in `CONTEXT.md` through
   `domain-modeling` by appending to its definition, after "and one terminal
   row `0 | blocked | precondition` in the Results table.", the sentence "It
   also records `build`, the commit the gate audited, as a full QA Report
   does, so a reader can tell a refusal an earlier Run left on the branch from
   the current Run's result; a refusal written before this rule keeps its
   bytes." Its `_Avoid_` line keeps its bytes.
4. MUST revise **QA Report** in `CONTEXT.md` through `domain-modeling` by
   appending to its definition, after "a gate that reached its matrix writes
   none of those three.", the sentence "A report the Daemon closes without an
   Agent turn, a Precondition Refusal or a blocking mechanical result, records
   the audited commit as `build`, as the gate's full report does; the pending
   report the Daemon seeds for the gate carries none." Its `_Avoid_` line
   keeps its bytes.
5. MUST, because this Task edits an owned skill and a guide, raise both
   version fields (`metadata.version` and the top-level `version`) of the
   canonical `SKILL.md` of the `qa-gate` skill by one patch over the value the
   tree holds when this Task starts, and then run these commands itself, in
   this order, before handing back; no mirror or record is hand-edited, and
   the Behavior Surface Record is re-recorded after the last guide edit:
   - `make skills-sync`;
   - `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`;
   - `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`,
     which re-records `docs/references/behavior-surfaces.json` because the
     Context-driven development guide's fingerprint changes.
6. MUST NOT change any Go source, ADR, `docs/references/skill-coverage.json`,
   any other guide page, or the Roundfix and `write-tasks` skills.

## Subtasks

- [x] Add `build` to the refusal template and its paragraph.
- [x] Name `build` in the guide's refusal paragraph.
- [x] Revise **Precondition Refusal** and **QA Report** in the glossary.
- [x] Raise the skill version, sync the mirror, record the version and re-record the Behavior Surface Record.

## Acceptance Criteria

- [x] Both copies of the `qa-gate` skill carry `build` in the refusal template and the paragraph that explains it, and `### QA settlement` is unchanged.
- [x] The guide's refusal paragraph names `build` and the Daemon's reports.
- [x] `CONTEXT.md` revises **Precondition Refusal** and **QA Report**.
- [x] The skill version is raised and recorded, the mirror is in sync, and the Behavior Surface Record is current.

## Result

Implemented the assigned documentation slice:

- Added the quoted `build` line and refusal explanation to the canonical
  `qa-gate` skill, then synchronized the mirror. Both skill version fields
  are `0.0.19); the `### QA settlement` section and `## 6. Close the
  report` bytes were left unchanged.
- Updated the Context-driven development refusal paragraph to name the audited
  commit and the Daemon's blocking and pending report behavior.
- Updated the **Precondition Refusal** and **QA Report** glossary definitions,
  preserving their `_Avoid_` lines.
- Re-recorded the owned skill version inventory and the Behavior Surface Record.

Focused checks run:

- `make skills-sync` — passed.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions` —
  passed, 1 test in 1 package.
- `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1` —
  passed, 1 test in 1 package.
- `cmp .agents/skills/qa-gate/SKILL.md skills/qa-gate/SKILL.md` — passed;
  `git diff --check` — passed.

The Daemon-owned Task status remains `in_progress`. The Task's declared
Verification commands were not run.

## Consumers

None. This Task changes documentation, the skill mirror, the skill version record and the Behavior Surface Record only; no code consumes them.

## Context

- instruction: `docs/adr/0132-a-refused-gate-records-the-refusal-as-its-terminal-row.md`
- instruction: `internal/docscontract/skill_coverage_test.go`
- instruction: `docs/references/skill-coverage.json`
- interface: `.agents/skills/qa-gate/SKILL.md`
- interface: `skills/qa-gate/SKILL.md`
- interface: `docs/user-guide/context-driven-development.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Verification

- `for f in .agents/skills/qa-gate/SKILL.md skills/qa-gate/SKILL.md; do grep -qxF -- '  build: "<audited-commit>"' "$f" || { printf 'missing template build line in %s\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'so a reader can tell this refusal from one an earlier Run left on the branch' || { printf 'missing refusal paragraph in %s\n' "$f" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < docs/user-guide/context-driven-development.md | grep -qF -- 'naming the commit the gate audited, as a full report does' || { echo 'missing guide build key' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/context-driven-development.md | grep -qF -- 'the pending report it seeds for the gate carries none' || { echo 'missing guide seed sentence' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Precondition Refusal**:' || { echo 'missing Precondition Refusal' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'so a reader can tell a refusal an earlier Run left on the branch from the current Run' || { echo 'missing Precondition Refusal revision' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**QA Report**:' || { echo 'missing QA Report' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'records the audited commit as' || { echo 'missing QA Report revision' >&2; exit 1; }; make skills-sync-check || exit 1; go test -count=1 ./skills || exit 1; go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$'` — expected: exit 0. Before this Task the `qa-gate` refusal template has no `build` line, so the command fails at its first check. Afterwards both skill copies carry the line and its paragraph, the guide names `build` and the Daemon's reports, the glossary revises **Precondition Refusal** and **QA Report**, the mirror is in sync, the skill version is recorded with the settlement table unchanged, and the Skill Coverage Map is current.

## References

- `_prd.md` → Core Feature 4; Goals; Glossary
- `_techspec.md` → API Contract 8; Build Order 4
- ADR-0132

## Carry-forward provenance

- Source Run: `run_20261023T070317Z_62af7ae2274ed43c`
- Source commit: `f6f581e9a8670bb3e1cf1389f311815c4342e8b4`
