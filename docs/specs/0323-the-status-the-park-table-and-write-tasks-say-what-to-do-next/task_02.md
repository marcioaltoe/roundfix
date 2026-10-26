---
task: task_02
spec: 0323-the-status-the-park-table-and-write-tasks-say-what-to-do-next
status: pending
type: docs
complexity: medium
---

# Task 02: The Park Class table lists blockers and states the secret-found action

## Overview

The `finding` row of the Park Class table ends with the `secret-found`
rewrite instruction, which reads as the action for every `finding` park. This
Task makes the row list its blockers only, states the `secret-found` next
action in a sentence of its own in the `deliver` guide and both Roundfix
reference copies, adds a documentation contract that compares each
blocker-specific sentence with the queue's `ClassifyPark`, and records the
skill version and the Behavior Surface fingerprint.

## Requirements

1. MUST answer the Backlog Entry
   [`references/2026-10-26-the-park-class-table-joins-the-secret-rewrite-to-every-finding-park.md`](references/2026-10-26-the-park-class-table-joins-the-secret-rewrite-to-every-finding-park.md)
   of 2026-10-26 (`_prd.md` → Sources 2) with `_techspec.md` → Interfaces
   and API Contract 3.
2. MUST end the `finding` row at `secret-found` in
   `docs/user-guide/commands/deliver.md` and the canonical
   `.agents/skills/roundfix/references/deliver.md`, and add the
   `secret-found` sentence of Interfaces directly after the
   `archived-item-unproven` sentence in each, changing no other row, no
   blocker and no other sentence of the section.
3. MUST add `TestTheParkActionsAreDocumentedAsTheQueueClassifiesThem`,
   `TestTheParkActionContractReportsAnActionInTheClassTable` and their helper
   `checkParkActionDocumentation` to the existing `//verify:always` file
   `internal/docscontract/command_documentation_test.go`, as
   `_techspec.md` → Testing Approach describes, reading the queue's text from
   `delivery.ClassifyPark` and the blocker names from the `delivery`
   constants, never as copied literals.
4. MUST run `make skills-sync`, then
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   then `make skills-sync` again so the mirror matches.
5. MUST, as the Task's last step, after every guide and skill edit,
   re-record the Behavior Surface fingerprint ("re-record after the last
   guide edit") with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
6. MUST NOT change `internal/delivery`, any Park Class, blocker or
   queue-printed next action, any other skill, `docs/agents/`, the Makefile
   or the CI workflows.

## Subtasks

- [ ] Split the row and add the `secret-found` sentence in the guide and the canonical reference.
- [ ] Add the two contract tests.
- [ ] Sync and record the skill, then re-record the Behavior Surface fingerprint last.

## Acceptance Criteria

- [ ] The `finding` row lists only blocker names in the guide, the canonical reference and its mirror.
- [ ] The `secret-found` sentence equals the queue's next action with `<worktree>` and `<slug>`, and the contract refuses the joined row.
- [ ] The skill version and the Behavior Surface record are current.

## Consumers

- `checkParkActionDocumentation` — consumers: `internal/docscontract`

## Context

- instruction: `internal/delivery/park_class.go`
- instruction: `skills/skills_test.go`
- instruction: `internal/docscontract/skill_coverage_test.go`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `internal/docscontract/command_documentation_test.go`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`

## Verification

- `grep -q 'func TestTheParkActionsAreDocumentedAsTheQueueClassifiesThem' internal/docscontract/command_documentation_test.go && go test -timeout 15m -count=1 -tags docscontract -run '^(TestTheParkAction.*|TestTheSkillCoverageMapIsCurrent)$' ./internal/docscontract && go test -timeout 10m -count=1 ./skills` — expected: exit 0. Before this Task the contract test does not exist, so the first check fails; afterwards the guide and both reference copies list blockers only in the `finding` row and state both blocker-specific next actions as the queue prints them, the contract refuses the joined row, the owned-skill version is recorded and the Behavior Surface record is current.

## References

- `_prd.md` → Sources 2; Goal 3; Core Feature 2; Success Metric 3
- `_techspec.md` → Current behavior 4; Interfaces; API Contract 3; Invariants 3-4; Testing Approach; Build Order 2
