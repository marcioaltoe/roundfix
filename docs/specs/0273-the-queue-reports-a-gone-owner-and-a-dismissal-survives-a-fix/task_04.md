---
task: task_04
spec: 0273-the-queue-reports-a-gone-owner-and-a-dismissal-survives-a-fix
status: pending
type: docs
complexity: low
---

# Task 04: The Roundfix skill and the glossary describe the owner report, the refusal and the dismissal

## Overview

task_01 to task_03 change what `deliver status`, the queue commands in an item
worktree and `review dispose --dismiss` do, and each documented its change in
the user guide. The Roundfix skill's `deliver` and `review` references, which
cover those surfaces in the Skill Coverage Map, and the glossary entries
**Queue Owner** and **Review Finding Disposition** still describe the old
behavior. This Task brings them up to date. It is verifiable on its own by its
phrases, the mirror check, the skill version record and the Behavior Surface
Record.

## Requirements

1. MUST answer both Backlog Entries of 2026-10-14,
   [deliver status does not say the queue owner is gone](references/2026-10-14-deliver-status-does-not-say-the-queue-owner-is-gone.md)
   and
   [a dismissal after a committed fix needs a reset](references/2026-10-14-a-dismissal-after-a-committed-fix-needs-a-reset.md),
   for the skill and glossary text their fixes need.
2. MUST add to `.agents/skills/roundfix/references/deliver.md`, where it
   describes `deliver status`, the `Owner:` line with the exact text
   `Owner: Delivery Queue owner PID <pid> is not running; run roundfix deliver resume`
   and its different-identity state, and, beside the retry refusal, the item
   worktree refusal with the phrase `this checkout is the item worktree of`.
3. MUST replace in `.agents/skills/roundfix/references/review.md` the sentence
   that a dismissal requires an unchanged reviewed `HEAD` and the words
   `moved-head dismissals` with the rule of API Contract 3, with the phrase
   `anchored lines are unchanged` and the `currentHead` field. No text inside
   any `### QA settlement` section may change.
4. MUST run `make skills-sync`, raise and record the Roundfix skill's version
   with `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   and re-record the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
   No mirror, version record or Behavior Surface Record is hand-edited.
5. MUST, through `domain-modeling`, revise **Queue Owner** in `CONTEXT.md` to
   add that `roundfix deliver status` reports a recorded Queue Owner proven
   not running, by its PID and start-time identity, while an item can still
   advance, and never reclaims it; and revise **Review Finding Disposition**
   to add that an evidence-backed dismissal may be recorded at a later HEAD
   that descends from the reviewed head when the finding's anchored lines are
   unchanged, staying tied to the reviewed head. Both cite ADR-0279, and
   every capitalized phrase the revisions add is a term `CONTEXT.md` already
   defines.

## Subtasks

- [ ] Update the `deliver` reference of the Roundfix skill.
- [ ] Update the `review` reference of the Roundfix skill.
- [ ] Sync the mirrors, record the version and the Behavior Surface Record.
- [ ] Revise the two glossary entries.

## Acceptance Criteria

- [ ] The `deliver` reference names the `Owner:` line and the item worktree refusal.
- [ ] The `review` reference states the dismissal after a fix and no longer requires an unchanged reviewed `HEAD`.
- [ ] Mirrors, the skill version and the Behavior Surface Record are current.
- [ ] `CONTEXT.md` revises **Queue Owner** and **Review Finding Disposition**.

## Context

- instruction: `docs/adr/0279-status-reports-a-queue-owner-proven-gone-and-a-dismissal-survives-a-fix-that-leaves-its-anchor.md`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `.agents/skills/roundfix/references/review.md`
- interface: `skills/roundfix/references/review.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Consumers

None. This Task changes documentation, the skill text and its derived records, which no code consumes; the skills and docscontract checks it runs read them.

## Verification

- `for f in .agents/skills/roundfix/references/deliver.md skills/roundfix/references/deliver.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'Owner: Delivery Queue owner PID <pid> is not running; run roundfix deliver resume' || { printf 'missing Owner line in %s\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'this checkout is the item worktree of' || { printf 'missing item worktree refusal in %s\n' "$f" >&2; exit 1; }; done; for f in .agents/skills/roundfix/references/review.md skills/roundfix/references/review.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'anchored lines are unchanged' || { printf 'missing dismissal rule in %s\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'currentHead' || { printf 'missing currentHead in %s\n' "$f" >&2; exit 1; }; if tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'an unchanged reviewed'; then printf '%s still requires an unchanged reviewed HEAD\n' "$f" >&2; exit 1; fi; done; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Queue Owner**:' || { printf 'CONTEXT.md does not define **Queue Owner**\n' >&2; exit 1; }; grep -A3 '^[*][*]Queue Owner[*][*]:' CONTEXT.md | grep -qF -- 'not running' || { printf 'CONTEXT.md does not revise **Queue Owner**\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Review Finding Disposition**:' || { printf 'CONTEXT.md does not define **Review Finding Disposition**\n' >&2; exit 1; }; grep -A3 '^[*][*]Review Finding Disposition[*][*]:' CONTEXT.md | grep -qF -- 'anchored lines' || { printf 'CONTEXT.md does not revise **Review Finding Disposition**\n' >&2; exit 1; }; make skills-sync-check || exit 1; go test -count=1 -run '^TestEveryOwnedSkillVersionIsRecorded$' ./skills || exit 1; go test -count=1 -tags docscontract -run '^TestTheSkillCoverageMapIsCurrent$' ./internal/docscontract || exit 1` — expected: exit 0. Before this Task neither skill reference names the new behaviors and the glossary entries lack them, so the command fails.

## References

- `_prd.md` → Core Feature 4; Glossary
- `_techspec.md` → API Contracts 1-3; Surface Transcripts 1-5; Build Order 4; Glossary
- ADR-0279; ADR-0256; ADR-0244
