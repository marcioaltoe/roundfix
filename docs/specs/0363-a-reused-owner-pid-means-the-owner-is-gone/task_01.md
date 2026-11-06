---
task: task_01
spec: 0363-a-reused-owner-pid-means-the-owner-is-gone
status: completed
type: docs
complexity: low
---

# Task 01: The doctor and deliver guides and the Roundfix references say a reused owner PID is gone

## Overview

Tasks 02 and 03 make a proven start-identity mismatch at a recorded owner
PID mean the owner is gone for Doctor's residue check, the Active Run lock
reclaim, the Run slots and `deliver start`. This Task writes that rule first
into the doctor and deliver guides and the Roundfix skill's setup and
deliver references, so each code Task names its guide through this
dependency, records the skill's next version and re-records the Behavior
Surface Record last.

## Requirements

1. MUST add, right after the bullet text "a partial result when Doctor
   cannot inspect every lineage." in `docs/user-guide/commands/doctor.md`,
   on the same bullet, the sentence: "A terminal Run whose recorded owner PID is alive with a different start identity adds nothing: its owner is gone and the PID was reused, so `partial` is left for a process table Doctor cannot read."
2. MUST add the same sentence, as a new paragraph, to the canonical
   Roundfix reference `.agents/skills/roundfix/references/setup.md`, right
   after the paragraph that begins "Use `roundfix doctor` when you only need
   a read-only readiness report.".
3. MUST, in `docs/user-guide/commands/deliver.md` and the canonical Roundfix
   reference `.agents/skills/roundfix/references/deliver.md`:
   - replace "Only a holder whose recorded owner PID is no longer alive is
     excluded;" with "Only a holder whose owner is gone is excluded: its
     recorded owner PID is no longer alive, or it is alive with a different
     start identity because the PID was reused;";
   - add, as a new paragraph right after the paragraph that ends "for
     `deliver resume` to recover.", the paragraph: "`deliver start` releases the recorded owner of an existing queue whose owner is gone, because its PID is no longer alive or is alive with a different start identity, before it starts the new queue; a recorded owner proven live is left as it is."
4. MUST raise both version fields of the Roundfix skill by one patch step
   over the values the tree holds when the Task runs, run `make
   skills-sync`, and re-record the version record with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`.
   No mirror or record is hand-edited.
5. MUST re-record the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
   as the Task's LAST step, after the last edit of the guides and the
   skill.
6. MUST NOT change any other paragraph of these files, `CONTEXT.md`, any Go
   file or any other skill.

## Subtasks

- [ ] Add the residue sentence to the doctor guide and the setup reference.
- [ ] Change the slot sentence and add the start paragraph in the deliver guide and reference.
- [ ] Raise the skill version, sync the mirrors and re-record the version record.
- [ ] Re-record the Behavior Surface Record last.

## Acceptance Criteria

- [ ] The doctor guide and the setup reference say a reused owner PID adds no residue.
- [ ] The deliver guide and reference say the Run slots and `deliver start` treat a reused owner PID as gone.
- [ ] The skill records its next version, and the skill and documentation contracts pass.

## Context

- interface: `docs/user-guide/commands/doctor.md`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `.agents/skills/roundfix/references/setup.md`
- interface: `skills/roundfix/references/setup.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The Roundfix skill version and the Behavior Surface Record — consumers: `./skills`, `./internal/docscontract`

## Verification

- `tr -s '[:space:]' ' ' < docs/user-guide/commands/doctor.md | grep -qF -- 'its owner is gone and the PID was reused' || exit 1; tr -s '[:space:]' ' ' < skills/roundfix/references/setup.md | grep -qF -- 'its owner is gone and the PID was reused' || exit 1; tr -s '[:space:]' ' ' < docs/user-guide/commands/deliver.md | grep -qF -- 'Only a holder whose owner is gone is excluded' || exit 1; tr -s '[:space:]' ' ' < skills/roundfix/references/deliver.md | grep -qF -- 'releases the recorded owner of an existing queue whose owner is gone' || exit 1; go test -timeout 20m -count=1 ./skills && go test -timeout 20m -count=1 -tags docscontract ./internal/docscontract -run '^(TestTheSkillCoverageMapIsCurrent|TestGlossaryPhrasesAreDefinedTerms)$'` — expected: exit 0. Before this Task none of the guides or references carries the sentences, so the first check fails; afterwards the mirrors match, the version record is current and the Behavior Surface Record matches the guides.

## References

- `_prd.md` → Goals; Core Feature 4
- `_techspec.md` → API Contracts 2-4; Build Order 1
- ADR-0373

## Result

Skills loaded: implement-task, tech-writer

Implemented the documentation slice: Doctor now says a terminal Run with a
reused owner PID adds no residue; both deliver guides treat a reused owner PID
as gone for Run slots and `deliver start`; and `deliver start`'s takeover rule
leaves a proven-live owner in place. The canonical Roundfix skill was raised
from `0.1.198` to `0.1.199` in both version fields, synchronized to its mirror,
and its owned-skill version record was regenerated.

Focused evidence:

- `make skills-sync` — passed; canonical and mirror `SKILL.md`, setup, and
  deliver references are byte-identical.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'
  -record-skill-versions` — passed (`1 passed in 1 packages`).
- `go test -count=1 -tags docscontract ./internal/docscontract -run
  '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage` — passed
  (`1 passed in 1 packages`); the Behavior Surface Record was regenerated after
  the final guide and skill edits.
- Focused phrase checks and diff inspection found the prescribed Doctor,
  setup-reference, Run-slot, and `deliver start` wording in both canonical
  surfaces.

Acceptance evidence:

- Doctor guide and setup reference: reused owner PID is described as gone and
  contributes no residue; the full declared Verification remains for the
  Daemon.
- Deliver guide and reference: Run-slot holders and `deliver start` use the
  owner-gone wording, including the reused-PID case; the full declared
  Verification remains for the Daemon.
- Skill/version and generated records: the sanctioned version recorder and
  Behavior Surface recorder passed; the declared skill and documentation
  contract Verification remains for the Daemon.
