---
task: task_04
spec: 0385-a-queue-item-stops-only-for-its-own-work
status: completed
type: docs
complexity: medium
---

# Task 04: The Roundfix deliver and review references describe the three changes

## Overview

The Roundfix skill's `deliver` and `review` references carry, word for word,
the four guide texts tasks 01 to 03 wrote: the owner's record, the queue log,
the anchor that may end a sentence and the validator's working directory. The
skill's version rises once, through the record command.

## Requirements

1. MUST change the canonical `.agents/skills/roundfix/references/deliver.md`
   with the two texts of `_techspec.md` → API Contract 6, word for word, at the
   anchors that contract names, which this reference holds with the same
   words as the guide.
2. MUST change the canonical `.agents/skills/roundfix/references/review.md`
   with the two texts of API Contract 7, word for word, at the anchors that
   contract names. MUST NOT add or change text inside any `### QA settlement`
   section, and MUST add no capitalized phrase that `CONTEXT.md` does not
   define.
3. MUST run `make skills-sync` to copy both references to
   `skills/roundfix/references/deliver.md` and
   `skills/roundfix/references/review.md`, then raise the Roundfix skill version
   with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`;
   never hand-edit a mirror, a version field or the record. Of the Governed
   Paths the Task may change only `.agents/skills/roundfix/SKILL.md`,
   `.agents/skills/roundfix/references/deliver.md`,
   `.agents/skills/roundfix/references/review.md` and
   `skills/roundfix/SKILL.md`.
4. MUST, as its last step, re-record the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.

## Subtasks

- [ ] Add the two `deliver` and the two `review` texts.
- [ ] Sync, record the skill version, then re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] Both copies of both references carry the four texts.
- [ ] The Roundfix skill version is recorded, and the skills, documentation and glossary contracts pass.

## Context

- instruction: `docs/adr/0395-a-queue-item-stops-only-for-its-own-work.md`
- instruction: `docs/user-guide/commands/deliver.md`
- instruction: `docs/user-guide/commands/review.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `.agents/skills/roundfix/references/review.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/review.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The `deliver` and `review` references and the skill version record — consumers: `./skills`
- The Behavior Surface Record — consumers: `./internal/docscontract`

## Verification

- `for f in .agents/skills/roundfix/references/deliver.md skills/roundfix/references/deliver.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'because its engine returned an error or the process panicked, keeps its owner record' || { echo "deliver reference: owner record text missing: $f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qE -- 'The queue log is the Queue Owner.s console log' || { echo "deliver reference: queue log text missing: $f" >&2; exit 1; }; done; for f in .agents/skills/roundfix/references/review.md skills/roundfix/references/review.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'An anchor followed by a period or a semicolon reads as one followed by a space' || { echo "review reference: anchor text missing: $f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'Roundfix copies it into that directory, which holds nothing else' || { echo "review reference: validator text missing: $f" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task neither reference carries the texts; afterwards both copies do, and the skills suite with its recorded version and the documentation contracts with the Behavior Surface Record and the glossary phrase contract pass.

## References

- `_prd.md` → User Story 2; Core Feature 5; Success Metric 5
- `_techspec.md` → API Contracts 6, 7 and 8; Build Order 4
- ADR-0395

## Result

Skills loaded: implement-task, tech-writer, ai-memory-handoff.

Starting-tree characterization: the four API Contract 6 and 7 texts were
present at their named anchors in the two user guides, but absent from both
Roundfix reference copies.

Implementation evidence:

- Added the owner-record and queue-log texts to the canonical `deliver`
  reference at the corresponding anchors.
- Added the punctuation-anchor and validator-working-directory texts to the
  canonical `review` reference at the corresponding anchors; no `### QA
  settlement` text changed.
- `make skills-sync` completed successfully; `cmp -s` confirmed both synced
  reference pairs are identical.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'
  -record-skill-versions` passed and recorded Roundfix skill version `0.1.224`
  in both skill copies and the owned-skill version record.

Focused checks:

- `git diff --check` passed.
- Direct normalized-text checks found all four required texts in both copies
  of the `deliver` and `review` references.
- `go test -count=1 -tags docscontract ./internal/docscontract -run
  '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage` passed and
  re-recorded the Behavior Surface Record.

The Task status remains `in_progress` and Daemon-owned. The declared
`## Verification` commands were not run; Daemon Verification remains pending.
