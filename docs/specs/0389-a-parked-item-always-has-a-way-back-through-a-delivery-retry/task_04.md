---
task: task_04
spec: 0389-a-parked-item-always-has-a-way-back-through-a-delivery-retry
status: pending
type: docs
complexity: low
---

# Task 04: The Roundfix deliver reference describes the three changes

## Overview

The Roundfix skill's `deliver` reference, which agents read to operate the
Delivery Queue, carries the four texts the `deliver` guide gained in tasks 01
to 03: the requeue of a branchless item, the authorization read without a
Spec Root directory, the return of an unpushed candidate to the gate, and the
`review-stale` next action. Adopters receive it through `baseline update` and
the Roundfix skills.

## Requirements

1. MUST change the canonical `.agents/skills/roundfix/references/deliver.md`
   with Texts A, B, C and D of `_techspec.md` → API Contract 6, word for word,
   at the anchors `_techspec.md` → API Contract 7 names. MUST NOT add or change
   text inside any `### QA settlement` section, and MUST add no capitalized
   phrase that `CONTEXT.md` does not define.
2. MUST run `make skills-sync` to copy the reference to
   `skills/roundfix/references/deliver.md`, then raise the Roundfix skill
   version to the next version with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
   on the tree the Task runs on, after Spec 0388 raised it; never hand-edit a
   mirror, a version field or the record. Of the Governed Paths the Task may
   change only `.agents/skills/roundfix/SKILL.md`,
   `.agents/skills/roundfix/references/deliver.md` and
   `skills/roundfix/SKILL.md`.
3. MUST, as its last step, re-record the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.

## Subtasks

- [ ] Add the four texts to the canonical reference.
- [ ] Sync the mirror and record the next skill version.
- [ ] Re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] Both copies of the `deliver` reference carry the four texts and hold the same bytes.
- [ ] The Roundfix skill version is raised and recorded; the skills and documentation contracts pass.

## Context

- instruction: `docs/adr/0398-a-parked-item-always-has-a-way-back-through-a-delivery-retry.md`
- instruction: `docs/user-guide/commands/deliver.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The `deliver` reference and the skill version record — consumers: `./skills`
- The Behavior Surface Record — consumers: `./internal/docscontract`

## Verification

- `for f in .agents/skills/roundfix/references/deliver.md skills/roundfix/references/deliver.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'An item parked before it has an Item Branch re-enters at' || { echo "deliver reference: Text A missing: $f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'The grant is read from Git history, so it needs no Spec Root in the working tree' || { echo "deliver reference: Text B missing: $f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'whose repository gate runs again before the push' || { echo "deliver reference: Text C missing: $f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'merge into the item branch any commit its Pull Request' || { echo "deliver reference: Text D missing: $f" >&2; exit 1; }; done; cmp -s .agents/skills/roundfix/references/deliver.md skills/roundfix/references/deliver.md || { echo 'deliver reference mirror differs' >&2; exit 1; }; go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task neither copy carries the texts; afterwards both do with the same bytes, and the skills suite with its recorded version and the documentation contracts with the Behavior Surface Record and the glossary phrase contract pass.

## References

- `_prd.md` → Core Feature 5; Success Metric 5
- `_techspec.md` → API Contracts 6 and 7; Build Order 4
- ADR-0398
