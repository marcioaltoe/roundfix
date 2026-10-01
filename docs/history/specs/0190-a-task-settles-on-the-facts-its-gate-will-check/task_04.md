---
task: task_04
spec: 0190-a-task-settles-on-the-facts-its-gate-will-check
status: completed
type: docs
complexity: medium
---

# Task 04: The skills, the commands guide and the glossary describe Settlement Checks

## Overview

Tasks 01 to 03 changed what the Implement Command does when a Task settles.
The Roundfix skill and `docs/user-guide/commands.md` still say a Task settles
on its declared Verification alone, and the `write-tasks` skill lets an author
split a change from the files its contract tests compare. This Task describes
the stage where operators and authors read it, and adds the glossary term
**Settlement Check**. It changes no code.

## Requirements

1. MUST add a `### Settlement Checks` section to
   `.agents/skills/roundfix/SKILL.md`, under `## Assigned Task Batches`. It
   states:
   - the stage applies to every non-QA Task of a Task Graph that has a QA gate
     Task;
   - the three checks, in order, with the labels
     `settlement check: spec consistency` and
     `settlement check: authorization`;
   - a failed check returns as Verification Feedback under the single repair,
     and a final failure settles the Task `failed`;
   - `verification.repository_at_settlement` turns off only the repository
     Verification at settlement;
   - the recorded limits of a parallel Wave and of a repository that is red on
     entry.
2. MUST add nothing inside the `### QA settlement` section of any skill. That
   section stays byte-identical in the three skills that carry it.
3. MUST add to `.agents/skills/write-tasks/SKILL.md`, in its decomposition
   rules, the rule that every Task of a graph that includes the gate
   `leaves the repository Verification green`, because the Daemon runs that
   command when each Task settles. The rule names
   `verification.repository_at_settlement` as the operator's switch, not the
   author's. Every phrase `TestWriteTasksSkillStatesTheDeclaredPathRules` and
   `TestTaskAuthoringGuidanceNamesDeclarations` require MUST stay.
4. MUST describe the stage in `docs/user-guide/commands.md`, where the guide
   describes how the Implement Command verifies and settles a Task: the three
   checks, the two labels, API Contract 3's failure reason and the switch.
5. MUST add the term **Settlement Check** to `CONTEXT.md`, next to
   **Verification Feedback**, in the glossary's existing form with an
   `_Avoid_` line.
6. MUST regenerate the mirrors with `make skills-sync` and run
   `make baseline-digests`. The only generated changes are
   `skills/roundfix/SKILL.md` and `skills/write-tasks/SKILL.md`.
7. MUST change only the four bounded skill files, `docs/user-guide/commands.md`,
   `CONTEXT.md` and this Task file.
8. MUST NOT edit `docs/agents/autonomous-work.md`, another skill, a Go file or
   a test.

## Subtasks

- [ ] Describe the stage in the Roundfix skill and the commands guide.
- [ ] Add the authoring rule to the `write-tasks` skill.
- [ ] Add the glossary term.
- [ ] Regenerate the mirrors and the digests.

## Acceptance Criteria

- [ ] The Roundfix skill and its mirror carry the `### Settlement Checks`
      section, both labels and the switch.
- [ ] The `write-tasks` skill and its mirror carry
      `leaves the repository Verification green`.
- [ ] The commands guide carries `Settlement Checks`, both labels, the failure
      reason and the switch.
- [ ] `CONTEXT.md` carries `**Settlement Check**:`.
- [ ] `make skills-sync-check` passes, and the skill contract tests that read
      the `### QA settlement` section and the `write-tasks` rules pass.

## Context

- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `.agents/skills/write-tasks/SKILL.md`
- interface: `skills/write-tasks/SKILL.md`
- interface: `docs/user-guide/commands.md`
- interface: `CONTEXT.md`
- instruction: `docs/adr/0182-a-task-settles-on-the-facts-its-gate-will-check.md`

## Verification

- `for pair in ".agents/skills/roundfix/SKILL.md|### Settlement Checks" ".agents/skills/roundfix/SKILL.md|settlement check: spec consistency" ".agents/skills/roundfix/SKILL.md|settlement check: authorization" ".agents/skills/roundfix/SKILL.md|verification.repository_at_settlement" "skills/roundfix/SKILL.md|### Settlement Checks" "skills/roundfix/SKILL.md|verification.repository_at_settlement" ".agents/skills/write-tasks/SKILL.md|leaves the repository Verification green" "skills/write-tasks/SKILL.md|leaves the repository Verification green" "docs/user-guide/commands.md|Settlement Checks" "docs/user-guide/commands.md|settlement check: spec consistency" "docs/user-guide/commands.md|settlement check: authorization" "docs/user-guide/commands.md|Settlement check failed:" "docs/user-guide/commands.md|verification.repository_at_settlement" "CONTEXT.md|**Settlement Check**:"; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done; make skills-sync-check || exit 1; out="$(go test -count=1 -v -run "^(TestSettlementGuidanceIsOneTable|TestWriteTasksSkillStatesTheDeclaredPathRules|TestTaskAuthoringGuidanceNamesDeclarations)$" ./skills 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestSettlementGuidanceIsOneTable TestWriteTasksSkillStatesTheDeclaredPathRules TestTaskAuthoringGuidanceNamesDeclarations; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task none of the fourteen phrases is present, so the command fails at the first one.

## References

- [_prd.md](_prd.md) — User Story 4; Core Feature 6; User Experience; Recorded limits
- [_techspec.md](_techspec.md) — API Contract 1; API Contract 2; API Contract 3; Testing Approach 5; Build Order 4
- [_authorization.md](_authorization.md)
- ADR-0182; ADR-0014; ADR-0038

## Result

Implemented the Settlement Checks documentation slice. The Roundfix skill now
describes the gated non-QA settlement stage, ordered checks, repair/failure
behavior, repository switch, parallel-Wave tree boundary, and red-on-entry
precondition limit. The write-tasks skill now requires every Task in a gated
graph to leave repository Verification green and identifies the switch as
operator-owned. The commands guide documents the labels, failure reason, and
switch, and the glossary defines Settlement Check next to Verification
Feedback. The mirror files carry the same authored guidance.

Focused implementation checks (run before Daemon Verification):

- `git diff --check` — passed.
- A repository-local phrase and mirror comparison check confirmed the required
  section, labels, switch, authoring rule, API failure reason, glossary term,
  and byte-identical source/mirror skill content.

Acceptance evidence:

- [x] Roundfix skill and mirror carry `### Settlement Checks`, both labels, and
  `verification.repository_at_settlement`.
- [x] write-tasks skill and mirror carry `leaves the repository Verification
  green` and identify the operator switch.
- [x] Commands guide carries `Settlement Checks`, both labels, the failure
  reason, and the switch.
- [x] `CONTEXT.md` carries `**Settlement Check**:` beside `Verification
  Feedback`.
- [ ] The Daemon must run the declared `make skills-sync-check`, digest,
  phrase, and skill contract Verification commands.

## Carry-forward provenance

- Source Run: `run_20260930T230750Z_197aa9aff15490a9`
- Source commit: `d37a1af926ab45ef062aa49e464e2d840b423359`
