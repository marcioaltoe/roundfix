---
task: task_08
spec: 0084-an-update-that-can-run
status: completed
type: docs
complexity: low
---

# Task 08: State the two obligations where a Task author reads them

## Overview

A clause in the catalog governs every repository that adopts the Baseline; it
does not reach the two skills a Supervisor actually reads while decomposing a
Spec and running its gate. This slice states the outside-evidence obligation and
the glossary check inside `write-tasks` and `qa-gate`, where the Rehearsal Cases
rule and the QA row contract already live, so the obligation is visible at the
moment it applies.

## Requirements

1. MUST state, in the task-authoring skill, that at least one acceptance row per
   Spec rests on evidence originating outside the Spec's own artifacts and records
   that evidence's origin, placed with the existing Rehearsal Cases contract.
2. MUST state, in the QA gate skill, that the report records the origin of that
   row's external evidence, and that a row whose external evidence cannot be
   obtained is recorded as blocked with its reason rather than dropped.
3. MUST state, in the task-authoring skill, the glossary check at the close of a
   Spec, without duplicating the catalog clause's wording so the two cannot drift
   into contradiction.
4. MUST NOT add a new required section, frontmatter field, or refusal condition to
   either skill beyond what the clauses oblige.
5. MUST edit the authoritative skill sources rather than the generated copies, and
   sync the generated copies with the sanctioned command.
6. MUST change only the paths the authorization bounds.

## Subtasks

- [x] State the outside-evidence obligation in the task-authoring skill.
- [x] State the external-evidence recording and blocked-row handling in the QA
      gate skill.
- [x] State the glossary check in the task-authoring skill.
- [x] Sync the generated skill copies with the sanctioned command.

## Acceptance Criteria

- [x] The task-authoring skill states the outside-evidence obligation adjacent to
      its Rehearsal Cases contract.
- [x] The QA gate skill states how the external-evidence row's origin is recorded
      and how an unobtainable one is blocked.
- [x] The task-authoring skill states the glossary check without restating the
      catalog clause verbatim.
- [x] Neither skill gains a new required section or frontmatter field.
- [x] The generated skill copies match their authoritative sources.

## Context

- instruction: `docs/workflow/authorizations/2026-08-08-evidence-from-outside-the-spec.md`
- instruction: `docs/agents/skill-dispatch.md`
- interface: `.agents/skills/write-tasks/SKILL.md`
- interface: `.agents/skills/qa-gate/SKILL.md`

## Verification

- `grep -q 'outside' .agents/skills/write-tasks/SKILL.md` — expected: exits 0, proving the obligation is stated in the authoritative source.
- `grep -q 'outside' .agents/skills/qa-gate/SKILL.md` — expected: exits 0, proving the QA gate states the recording rule.
- `diff -r .agents/skills/write-tasks skills/write-tasks > /tmp/0084-task-08-a.log 2>&1` — expected: exits 0, proving the generated copy matches its authoritative source.
- `diff -r .agents/skills/qa-gate skills/qa-gate > /tmp/0084-task-08-b.log 2>&1` — expected: exits 0, proving the generated copy matches its authoritative source.
- `go run ./cmd/roundfix skills check > /tmp/0084-task-08-c.log 2>&1` — expected: exits 0.

## References

- `_techspec.md` → Build Order 8.
- `_prd.md` → Core Features 8 and 9; User Story 7; Goal 5.
- ADR-0104.

## Result

Implemented 2026-08-08.

### What was stated

Three additions to the two authoritative skill sources, +33 lines, no deletions
beyond a re-wrapped paragraph boundary:

| Skill | Placement | States |
| --- | --- | --- |
| `write-tasks` | Decomposition rules, directly after **Gate rehearsals declare their evidence** (line 85) | one acceptance row rests on evidence the Spec did not author, and records its origin |
| `write-tasks` | close of **Author the QA gate decision** (line 122) | the graph's closing node carries the glossary check |
| `qa-gate` | matrix row list and section 2 body (lines 123, 129) plus one decision example (line 315) | the report records that row's evidence origin, and an unobtainable source is blocked, not dropped |

The outside-evidence text is placed where the Rehearsal Cases contract already
lives because the two are the same concern seen from opposite ends: a rehearsal
proves the code matches the requirement, and only an outside source can show the
requirement was right. That is the failure the authorization records for Spec
0082, whose gate passed with zero blocked rows on a rehearsal of its own premise.

The glossary paragraph names *where* the check happens — the authored `qa` Task
when the gate is included, otherwise the last Task in topological order — and
delegates *what it looks for* to `docs/agents/domain.md`, the guide that renders
`clause.domain.glossary-currency`. The skill therefore cannot drift into
contradicting the clause: it never states the clause's content, only its home in
the graph.

For `qa-gate` the unobtainable source is typed as `blocked (environment:
<cause>)`, which is the existing typed cause for a proved environmental
constraint, and counted in the existing `rows_blocked_environment` key. No new
typed cause, no new verdict rule.

### Acceptance criteria

- **The task-authoring skill states the outside-evidence obligation adjacent to
  its Rehearsal Cases contract.** `.agents/skills/write-tasks/SKILL.md:85-93`
  is the bullet **One acceptance row rests on evidence the Spec did not
  author**, the list entry immediately following **Gate rehearsals declare their
  evidence** (line 84). It names the row, the admissible outside sources, the
  obligation to record where the evidence came from, and the blocked-with-reason
  outcome that keeps the obligation from stalling the Spec. `grep -n outside`
  reports lines 86 and 90 in the authoritative source.

- **The QA gate skill states how the external-evidence row's origin is recorded
  and how an unobtainable one is blocked.**
  `.agents/skills/qa-gate/SKILL.md:123` adds the outside-evidence row to the
  matrix row list; lines 129-137 require recording where its evidence came from
  "named precisely enough for a later reader to reach the same source", and
  require `blocked (environment: <cause>)` counted in
  `rows_blocked_environment` when the source cannot be obtained — "Never drop
  the row, and never satisfy it with evidence the Spec authored." Lines 315-318
  carry the matching decision example. `grep -n outside` reports lines 123, 129,
  130, and 315 beyond the two pre-existing unrelated matches at 38 and 299.

- **The task-authoring skill states the glossary check without restating the
  catalog clause verbatim.** `.agents/skills/write-tasks/SKILL.md:122-128`.
  Overlap check against `clause.domain.glossary-currency` — every distinctive
  clause phrase is absent from the skill:

  ```
  "introduced, changed, or retired a term the glossary should carry" -> 0
  "update the domain context through"                                -> 0
  "The check is what is obliged"                                     -> 0
  ```

  The skill says "a term the Spec coined, changed, or dropped" and points at
  `docs/agents/domain.md` for what the check looks for.

- **Neither skill gains a new required section or frontmatter field.** Comparing
  every heading and frontmatter key against `HEAD` for both files:

  ```
  $ diff <(git show HEAD:$f | grep "^#\|^[a-z_]*:") <(grep "^#\|^[a-z_]*:" $f)
  == .agents/skills/write-tasks/SKILL.md
  identical headings + frontmatter keys
  == .agents/skills/qa-gate/SKILL.md
  identical headings + frontmatter keys
  ```

  Only line numbers moved. No refusal condition was added either: the
  outside-evidence bullet and the glossary paragraph are stated as obligations,
  not as new `Refuse …` triggers, and `qa-gate` reuses its existing typed
  blocked cause and verdict rules unchanged.

- **The generated skill copies match their authoritative sources.**
  `make skills-sync` regenerated `skills/write-tasks` and `skills/qa-gate` from
  `.agents/skills/`; `make skills-sync-check` (the `diff -r` over every owned
  skill plus the bundle contract tests) passed:

  ```
  $ make skills-sync && make skills-sync-check
  go test -count=1 ./skills -run 'TestNoPythonBaselineRuntime|TestThinSetupSkill|TestCheckRejectsExecutableSetupEngineArtifacts|TestRecommendedSkillsMatchLock'
  ok  	roundfix/skills	0.319s
  ```

  `.claude/skills` is `-> ../.agents/skills`, so the Claude-facing copy is the
  same file and needs no separate sync.

### Focused checks

- `go test ./skills -count=1` → 136 passed, including the authorial-skill
  contract tests and `TestCharacterizationCorporaDoNotRecordOwnedSkillDigests`.
- `make skills-sync-check` → passed (no drift between `.agents/skills/` and
  `skills/`).
- `git status --short` → only `.agents/skills/{write-tasks,qa-gate}/SKILL.md`,
  their two `skills/` copies, and this Task file. Both authoritative paths are
  named by `docs/workflow/authorizations/2026-08-08-evidence-from-outside-the-spec.md`,
  and the `skills/` copies are its sanctioned `make skills-sync` fallout.
- The frozen parity corpus is unaffected: its fixtures record a snapshot digest
  (`write-tasks/SKILL.md`, 9287 bytes) that has not tracked the live file for
  many commits, so no derived pin moves with this change.

### Follow-up note (not in this diff)

The owned-skill contract is a *minimum* version (`ownedSkillMinimumVersion =
"0.0.2"` in `skills/skills.go`), not content parity, and both skills already
declare `version: 0.0.2`. A repository that adopted the Baseline earlier
therefore stays compliant with the older text until it reinstalls; propagating
edited skill text to already-adopted repositories is a separate question from
this slice, and the authorization bounds this Task to the two SKILL.md files.
