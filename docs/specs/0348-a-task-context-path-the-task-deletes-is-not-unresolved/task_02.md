---
task: task_02
spec: 0348-a-task-context-path-the-task-deletes-is-not-unresolved
status: completed
type: docs
complexity: low
---

# Task 02: The guides list a removed file only under deletes and state the base rule

## Overview

`write-tasks` says a removed path is declared under `deletes:`, "never
`instruction:`", and says nothing of `interface:`, which is the label the
Fiscus Task used for the file it removed. This Task makes the `write-tasks`
skill and its Task template say that a file the Task removes is listed only
under `deletes:`, and makes the Roundfix `spec` reference, the spec command
guide and the **Delivery Base** glossary entry state when a missing Task
Context path resolves and what the finding's fix says. It answers the Fiscus
report of 2026-11-04 (ai-memory message `01a1a5c5`, adopted as
[references/2026-11-04-a-task-context-path-the-task-deletes-fails-its-own-settlement.md](references/2026-11-04-a-task-context-path-the-task-deletes-fails-its-own-settlement.md)).

## Requirements

1. MUST change the canonical `.agents/skills/write-tasks/SKILL.md`, bullet
   "Declared edits and governed paths are explicit", so it reads "every path
   it removes is declared only under `deletes:`, never `interface:` and never
   `instruction:`. A Task that reads a file before removing it still lists it
   once, under `deletes:`." in place of "every path it removes is declared
   under `deletes:`, never `instruction:`." (`_techspec.md` → API
   Contract 3), keeping "never `instruction:`" on one line, because
   `TestWriteTasksSkillStatesTheDeclaredPathRules`
   (`skills/baseline_skill_contract_test.go`) requires that exact text in
   the canonical file and stays unchanged. Never touch any `### QA
   settlement` section.
2. MUST add to the Context comment of the canonical
   `.agents/skills/write-tasks/references/task-template.md` the line
   "A file the Task removes is listed only as deletes:, never as interface: or instruction:."
3. MUST add to the paragraph on a missing Task Context entry in the
   canonical `.agents/skills/roundfix/references/spec.md` and in
   `docs/user-guide/commands/spec.md` the same text: an `interface:` or
   `instruction:` path missing from the checked tree also resolves when
   `HEAD` or the Delivery Base tracks it and the checked root is the top
   level of its Git work tree, so a Task that removes a file it lists
   settles and the QA gate accepts the committed removal; any other missing
   path reports `SC-REF-UNRESOLVED`, whose fix says to remove the Task
   Context entry or declare the path under `deletes:`.
4. MUST change **Delivery Base** in `CONTEXT.md` through domain-modeling by
   adding, before its ADR link, the sentence "The Spec Consistency Check also
   resolves a path a Task lists under `interface:` or `instruction:` that the
   checked tree no longer holds, when HEAD or the Delivery Base tracks it."
   It keeps this repository's glossary contract (no undefined capitalized
   phrase, the existing `_Avoid_` line).
5. MUST add `skills/task_context_deletion_test.go` with
   `TestTheWriteTasksGuidesListARemovedFileOnlyUnderDeletes`
   (`_techspec.md` → Testing Approach): both copies of the `write-tasks`
   `SKILL.md` and of `references/task-template.md` carry the sentences of
   requirements 1 and 2 on whitespace-normalized text, and neither copy
   carries "declared under `deletes:`, never `instruction:`". It calls
   `t.Parallel()` first.
6. MUST write the mirrors only with `make skills-sync`, then raise and record
   the `write-tasks` and Roundfix owned-skill versions with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   never by hand.
7. MUST re-record the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
   as the Task's last step, after the last edit of the spec command guide
   and the skills, and never hand-edit it.
8. MUST NOT change any Go source outside `skills/task_context_deletion_test.go`,
   any other skill, any guide under `docs/agents/`, the Makefile, CI
   workflows, `go.mod` or `go.sum`.

## Subtasks

- [ ] Edit the canonical `write-tasks` skill and Task template.
- [ ] Edit the Roundfix `spec` reference, the spec command guide and `CONTEXT.md`.
- [ ] Add the skills test, sync the mirrors and record the owned-skill versions.
- [ ] Re-record the Behavior Surface Record last.

## Acceptance Criteria

- [ ] Both copies of `write-tasks` and its Task template say a removed file is listed only under `deletes:`.
- [ ] Both copies of the Roundfix `spec` reference, the spec command guide and `CONTEXT.md` state the base rule.
- [ ] The owned-skill versions and the Behavior Surface Record are recorded by their commands.

## Context

- instruction: `.agents/skills/domain-modeling/SKILL.md`
- interface: `.agents/skills/write-tasks/SKILL.md`
- interface: `skills/write-tasks/SKILL.md`
- interface: `.agents/skills/write-tasks/references/task-template.md`
- interface: `skills/write-tasks/references/task-template.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/spec.md`
- interface: `skills/roundfix/references/spec.md`
- interface: `skills/testdata/owned-skill-versions.json`
- creates: `skills/task_context_deletion_test.go`
- interface: `docs/user-guide/commands/spec.md`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Consumers

- The `write-tasks` and Roundfix owned-skill texts and their version record — consumers: `./skills`
- The spec command guide, the Behavior Surface Record and the glossary contract over `CONTEXT.md` — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^TestTheWriteTasksGuidesListARemovedFileOnlyUnderDeletes$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestTheWriteTasksGuidesListARemovedFileOnlyUnderDeletes ' || { echo 'missing pass: TestTheWriteTasksGuidesListARemovedFileOnlyUnderDeletes' >&2; exit 1; }; go test -count=1 -timeout 20m ./skills` — expected: exit 0. Before this Task the test does not exist, so its pass line is missing; afterwards it passes with the whole skills suite, the owned-skill version record included.
- `grep -F -- '**Delivery Base**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Delivery Base\*\*(:| —( |$))' || { echo 'missing term in CONTEXT.md' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'no longer holds, when HEAD or the Delivery Base tracks it' || { echo 'missing base rule in CONTEXT.md' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/spec.md | grep -qF -- 'also resolves when' || { echo 'missing base rule in the spec guide' >&2; exit 1; }; go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task neither file carries the base rule, so the second check fails; afterwards both carry it and the documentation contracts, the glossary contract and the Behavior Surface fingerprint among them, pass.

## References

- `_prd.md` → Goals; Core Feature 4; Success Metric 5; Glossary
- `_techspec.md` → Current behavior 5; API Contract 3; Testing Approach; Build Order 2

## Result

Skills loaded: implement-task, domain-modeling, coding-guidelines, golang-testing, testing-boss, tech-writer, knowledge-workspace, no-workarounds, evidence-gate

Implementation:

- Updated the canonical `write-tasks` skill and Task template, then ran
  `make skills-sync` to regenerate their mirrors.
- Updated the canonical and mirrored Roundfix `spec` reference, the spec
  command guide and the **Delivery Base** glossary entry with the base rule.
- Added `skills/task_context_deletion_test.go` and recorded the
  `roundfix` and `write-tasks` owned-skill versions through the sanctioned
  recorder. It recorded `roundfix` 0.1.180 and `write-tasks` 0.0.30.

Focused checks:

- `go test -count=1 ./skills -run '^TestTheWriteTasksGuidesListARemovedFileOnlyUnderDeletes$'`: passed.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`: passed.
- `go test -count=1 -tags docscontract ./internal/docscontract -run
  '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`: passed and
  re-recorded `docs/references/behavior-surfaces.json` after the Daemon
  reported its fingerprint as stale.
- `git diff --check`: passed.
- Read-only assertions found the removed-file rule in both copies of the
  `write-tasks` skill and Task template, and the base rule in both
  Roundfix reference copies, the spec command guide and `CONTEXT.md`.

Acceptance evidence:

- Removed-file declarations: implemented and covered by the focused contract
  test in both canonical and mirrored skill trees.
- Missing Task Context base rule: implemented in both guides and the
  **Delivery Base** glossary entry.
- Owned-skill versions: recorded by the required command.
- Behavior Surface Record: the stale fingerprint reported by the Daemon was
  re-recorded through the required sanctioned command; the recorder passed.
