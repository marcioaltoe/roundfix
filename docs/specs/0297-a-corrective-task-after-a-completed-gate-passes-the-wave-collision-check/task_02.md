---
task: task_02
spec: 0297-a-corrective-task-after-a-completed-gate-passes-the-wave-collision-check
status: pending
type: docs
complexity: low
---

# Task 02: The spec guide and the Roundfix reference describe the Wave collision finding

## Overview

Neither the `spec` user guide nor the Roundfix skill's `spec` reference
mentions `SC-WAVE-COLLISION`. This Task describes the finding there: what
makes two Tasks collide, the edge its summary names, and that the QA gate,
the settlement checks and the pre-dispatch refusal read the same working-tree
graph and settlement history. It raises the skill's version and re-records the
Behavior Surface Record (`_prd.md` → Core Feature 4).

## Requirements

1. MUST add to `docs/user-guide/commands/spec.md`, as the last row of the
   table whose header is `| Code | Reported when |`, right after the
   `SC-QA-BYTES-PINNED` row, this row:
   "| `SC-WAVE-COLLISION` | Two Tasks the current Task Graph lets run in one
   Wave are known to touch a common repository file: a literal operand of a
   Verification command, a `creates:`, `interface:` or `deletes:` path the
   tree holds, or a path the Task's newest settlement commit on any branch
   changed (`prior Run settlement commit`). The summary names each shared
   path with its source and ends with the one edge that orders the pair,
   ``add a `needs` edge from Task <later> to Task <earlier>``, where the
   later Task comes later in graph order; the fix says the same. An error
   from the Tasks stage. The QA gate's precondition, the Daemon's settlement
   checks and its pre-dispatch refusal at a Task Capacity above 1 read the
   same working-tree Task Graph and settlement history, so a corrective Task
   given that edge passes all of them. |", on one line.
2. MUST add to the canonical Roundfix reference
   `.agents/skills/roundfix/references/spec.md`, right after the
   `SC-CONSUMERS-UNPROVEN` bullet, this bullet:
   "- `SC-WAVE-COLLISION` — two Tasks the current Task Graph lets run in one
   Wave touch a common repository file: a Verification operand, a declared
   `creates:`, `interface:` or `deletes:` path the tree holds, or a path the
   Task's newest settlement commit on any branch changed. The summary ends
   with the one edge that orders the pair, ``add a `needs` edge from Task
   <later> to Task <earlier>``. The QA gate's precondition, the settlement
   checks and the pre-dispatch refusal read the same working-tree Task Graph
   and settlement history, so a corrective Task given that edge passes all of
   them." No text inside any `### QA settlement` section of any skill may
   change.
3. MUST run these commands, because this Task edits an owned skill and a
   guide. Before them, raise both version fields (`metadata.version` and the
   top-level `version`) of `.agents/skills/roundfix/SKILL.md` by one patch
   over the value the tree holds when this Task starts, because the record
   command refuses changed content under an unchanged version. The Agent runs
   them itself before handing back, because the Daemon never regenerates a
   record; no mirror or record is hand-edited.
   - `make skills-sync`;
   - `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`;
   - `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`,
     which re-records `docs/references/behavior-surfaces.json` because the
     `spec` guide's fingerprint changes.
4. MUST search the guides, the agent guides and the skill copies for wording
   about Wave collisions or the `needs` remedy that contradicts the finding
   as Task 01 writes it, and report what it found.
5. MUST NOT change any Go source or ADR, `CONTEXT.md`,
   `docs/references/skill-coverage.json` or any other guide page.

## Subtasks

- [ ] Add the `SC-WAVE-COLLISION` row to the `spec` guide and the bullet to the Roundfix `spec` reference.
- [ ] Raise the skill version, sync the mirror, record the version and re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] The guide and both copies of the reference describe `SC-WAVE-COLLISION`, its sources and its edge.
- [ ] The skill version is raised and recorded, the mirror is in sync, and the Behavior Surface Record is current.

## Consumers

None. This Task changes documentation, the skill mirror, the skill version record and the Behavior Surface Record only; no code consumes them.

## Context

- instruction: `internal/spec/collision.go`
- instruction: `internal/speccheck/coherence.go`
- instruction: `internal/docscontract/skill_coverage_test.go`
- instruction: `docs/references/skill-coverage.json`
- interface: `docs/user-guide/commands/spec.md`
- interface: `.agents/skills/roundfix/references/spec.md`
- interface: `skills/roundfix/references/spec.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`

## Verification

- `for f in docs/user-guide/commands/spec.md .agents/skills/roundfix/references/spec.md skills/roundfix/references/spec.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'SC-WAVE-COLLISION' || { printf 'missing code in %s\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'the one edge that orders the pair' || { printf 'missing edge in %s\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'so a corrective Task given that edge passes all of them' || { printf 'missing corrective Task in %s\n' "$f" >&2; exit 1; }; done; make skills-sync-check || exit 1; go test -count=1 ./skills || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task no document mentions `SC-WAVE-COLLISION`, so the command fails at its first phrase check. Afterwards the guide and both reference copies describe the finding and its edge, the mirror is in sync, the skill version is recorded, and the documentation contracts, the Skill Coverage Map among them, pass.

## References

- `_prd.md` → Core Feature 4; Goals
- `_techspec.md` → API Contract 6; Build Order 2
- ADR-0117
