---
task: task_07
spec: 0082-the-manifest-already-answered-that
status: pending
type: docs
complexity: medium
---

# Task 07: Teach the update path to the docs and the owned skills

## Overview

The update command is only useful if the operating contract, the skill that
dispatches Baseline work, and the guidance generated into adopting repositories
all name it. This is an authorized protected-tooling task: it edits
Roundfix-owned Skills and Baseline module assets, and it may change nothing
else. Its boundary is the exact file list recorded in the authorization.

## Requirements

1. MUST update the durable operating contract in the user guide to document the
   update command, its flags, its exit categories, and the managed-refresh
   guarantee that non-managed bytes are preserved.
2. MUST update the `setup-context-driven` skill so its recipes route an already
   adopted repository to the update command and stop implying that every refresh
   requires the full interview.
3. MUST update the `roundfix` skill so the shipped CLI surface it describes
   matches the binary, satisfying the repository's CLI skill-sync rule.
4. MUST update any Baseline module asset that names the Baseline command family
   so generated guidance teaches the update path, and MUST change module assets
   for no other reason.
5. MUST regenerate derived digest pins through the sanctioned regeneration
   command rather than transcribing any value by hand.
6. MUST change only these repository-relative paths plus this Task file:
   `skills/setup-context-driven/**` and `.agents/skills/setup-context-driven/**`;
   `skills/roundfix/**` and `.agents/skills/roundfix/**`;
   `internal/baseline/assets/modules/*.json`;
   `docs/user-guide/context-driven-development.md`; and the pins the sanctioned
   regeneration command rewrites. Any other changed path fails this Task.
7. MUST NOT weaken any module's Normative Clauses, decisions, capabilities, or
   template selection.

## Subtasks

- [ ] Document the update command in the durable operating contract.
- [ ] Route adopted repositories to the update path in the setup skill.
- [ ] Sync the roundfix skill with the shipped CLI surface.
- [ ] Name the update path in the module assets that name the command family.
- [ ] Run the sanctioned digest regeneration and keep its output unedited.
- [ ] Confirm the changed-file set matches the authorized boundary exactly.

## Acceptance Criteria

- [ ] The user guide documents the command, every flag, every exit category, and
      the preservation guarantee.
- [ ] The setup skill no longer presents the full interview as the only refresh
      route for an adopted repository.
- [ ] The roundfix skill's described CLI surface matches the binary's usage
      output for the Baseline command family.
- [ ] Generated guidance produced from the module assets names the update path.
- [ ] Every derived pin equals the value the sanctioned regeneration command
      produces; no pin was hand-edited.
- [ ] The changed-file set is a subset of the authorized boundary.

## Context

- instruction: `docs/workflow/authorizations/2026-08-07-baseline-update-command.md`
- instruction: `docs/agents/specific-repository.md`
- instruction: `docs/agents/skill-dispatch.md`
- interface: `docs/user-guide/context-driven-development.md`
- interface: `skills/setup-context-driven/SKILL.md`
- interface: `skills/roundfix/SKILL.md`

## Verification

- `go build -buildvcs=false ./...` — expected: exits 0.
- `go run -buildvcs=false ./cmd/roundfix skills check` — expected: exits 0.
- `grep -q 'baseline update' docs/user-guide/context-driven-development.md` — expected: exits 0.
- `grep -q 'baseline update' skills/setup-context-driven/SKILL.md` — expected: exits 0.
- `git diff --name-only HEAD | grep -v -E '^(skills/(setup-context-driven|roundfix)/|\.agents/skills/(setup-context-driven|roundfix)/|internal/baseline/assets/|internal/baseline/testdata/|docs/user-guide/context-driven-development\.md|docs/specs/0082-the-manifest-already-answered-that/task_07\.md)' | grep . ; test $? -eq 1` — expected: exits 0, proving no path outside the authorized boundary changed.
- `go test ./internal/baseline/ ./internal/cli/ ./skills/ -count=1` — expected: exits 0.

## References

- `_techspec.md` → Build Order 9; Project Constraints: Tooling authority.
- `_prd.md` → Project Constraints: Tooling authority.
- ADR-0081.
