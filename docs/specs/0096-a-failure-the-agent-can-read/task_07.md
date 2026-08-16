---
task: task_07
spec: 0096-a-failure-the-agent-can-read
status: completed # pending | in_progress | completed | failed — only implement-task changes this
type: docs
complexity: low
---

# Task 07: Write the ceiling's exits into the authoring contract

## Overview

When a gate returns more corrective work than the contract's ceiling allows, the
contract says the decomposition was wrong and does not say what to do. The loop
stops and asks a human for a policy decision that has two sanctioned answers.

## Requirements

1. MUST state, in the authoring contract, what to do when corrective work exceeds
   the ceiling: amend the technical spec and recut the graph from it, or promote
   the excess to its own Spec with the gate failing the discovered story
   explicitly.
2. MUST state that reaching the ceiling is a decision inside the loop's authority,
   so the loop continues rather than stopping for a human.
3. MUST NOT change the ceiling's value or how many repair turns a Task gets.
4. MUST NOT change any repository path outside the bounded scope below plus this
   Task file; stop and fail the Task if a changed-file check finds another path.

## Subtasks

- [ ] Write the two sanctioned exits.
- [ ] State that the choice is inside the loop's authority.

## Acceptance Criteria

- [ ] The authoring contract names both exits.
- [ ] It states that the loop chooses without stopping for a human.
- [ ] Neither the ceiling's value nor the repair-turn count changed.
- [ ] The changed-file set is the bounded scope plus this Task file.

## Bounded scope

This Task may create or modify only:

- `.agents/skills/write-tasks/SKILL.md`
- `docs/specs/0096-a-failure-the-agent-can-read/task_07.md`

Express maintainer authorization:
`docs/workflow/authorizations/2026-08-12-the-authoring-and-baseline-corrections.md`,
whose per-Spec section for this Spec names exactly this file and is binding. The
skill mirror rewritten by `make skills-sync` is sanctioned fallout under ADR-0081
and is declared in that record's `## Sanctioned regeneration` block.

## Verification

- `f=.agents/skills/write-tasks/SKILL.md; grep -qi 'ceiling' "$f" || { echo "FAIL: $f does not mention the corrective ceiling"; exit 1; }; grep -qiE 'amend the (technical spec|TechSpec)|recut' "$f" || { echo "FAIL: $f does not name the amend-and-recut exit"; exit 1; }; grep -qiE 'promote .*(own Spec|new Spec)' "$f" || { echo "FAIL: $f does not name the promote-to-its-own-Spec exit"; exit 1; }` — expected: exits 0, proving both exits are written. Fails today.
- `make skills-sync-check > /tmp/0096-t07.log 2>&1 && grep -qi 'ceiling' skills/write-tasks/SKILL.md || { cat /tmp/0096-t07.log; exit 1; }` — expected: exits 0, proving the mirror carries the clause and matches its source. The sync check alone passes before any work, so it is anchored to the clause it guards.
- `git diff --name-only HEAD > /tmp/0096-t07-all.txt; test -s /tmp/0096-t07-all.txt || { echo 'no file changed'; exit 1; }; grep -v -e '^\.agents/skills/write-tasks/SKILL\.md$' -e '^skills/write-tasks/SKILL\.md$' -e '^docs/specs/0096-a-failure-the-agent-can-read/task_07\.md$' /tmp/0096-t07-all.txt > /tmp/0096-t07-scope.txt; test ! -s /tmp/0096-t07-scope.txt || { echo 'out of bounds:'; cat /tmp/0096-t07-scope.txt; exit 1; }` — expected: exits 0, proving work happened and every changed path is in bounds.

## Context

- instruction: `docs/workflow/authorizations/2026-08-12-the-authoring-and-baseline-corrections.md`

## References

`_techspec.md` → Build Order 7; Risks & Considerations, the tooling Task.
`_prd.md` → Core Feature 3; Goal 3; User Story 3. ADR-0081.

## Result

The authoring contract now keeps the two-corrective-Task ceiling actionable. If
QA findings would require a third corrective Task, the loop either amends the
TechSpec and recuts the Task Graph from it, or promotes the excess corrective
work to its own Spec while leaving the gate failing the discovered story
explicitly. The clause makes that choice part of the loop's authority and tells
the loop to continue without stopping for a human. `make skills-sync` copied the
authorized source change to the sanctioned shipped-skill mirror.

Focused checks and acceptance evidence:

- Before implementation,
  `rtk rg -n "promote .*own Spec|amend the (technical spec|TechSpec)|recut|corrective-Task ceiling" .agents/skills/write-tasks/SKILL.md`
  found no match, establishing that neither exit nor the ceiling clause existed
  in the authoring contract.
- Fixed-string `rtk rg` checks found the amend-and-recut exit and the
  promote-to-its-own-Spec exit in both `.agents/skills/write-tasks/SKILL.md` and
  `skills/write-tasks/SKILL.md`. A focused authority search found “inside the
  loop's authority” and “continues without stopping for a human” in both copies.
- `rtk rg -n -F 'more than two corrective Tasks'` found the existing threshold
  in `docs/agents/autonomous-work.md` and the same threshold in the new clause;
  `rtk rg -n -F 'allows one Verification repair'
  docs/adr/0038-daemon-allows-one-verification-repair.md` found the unchanged
  one-repair policy. Neither governing file appears in the changed-file set.
- `rtk cmp -s .agents/skills/write-tasks/SKILL.md
  skills/write-tasks/SKILL.md` and `rtk git diff --check` exited 0.
- `rtk git status --short` and `rtk git diff --name-status` showed only the
  authorized canonical skill, this Task file, and the declared generated mirror;
  no other path changed.
- `rtk make verify-incremental` first reached two `internal/cli` integration
  failures because the sandbox denied process-table inspection. The approved
  rerun with process-table access exited 0; the Go suite, skill checks, and build
  passed.

The Daemon-owned commands under `## Verification` were not run during this Agent
turn. The change introduces no new or changed glossary term.
