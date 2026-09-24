---
task: task_03
spec: 0163-baseline-decisions-and-regeneration
status: pending
type: backend
complexity: medium
---

# Task 03: Skill regeneration declares its outputs

## Overview

`OutputsFor` in `internal/baseline/derived_ownership.go` reads ownership records only below `internal/baseline`, so a grant naming `make skills-sync` covers none of the mirrors that command writes and every grant enumerates them.

## Requirements

1. MUST add `skills/_ownership.yml` declaring `owner: dedicated` and `command: make skills-sync`.
2. MUST make `OutputsFor(root, "make skills-sync")` return every regular file under `skills/<owned>/` for each member of the Makefile's `OWNED_SKILLS`, reading the Makefile without writing it.
3. MUST exclude `.agents/skills`, Go files, `testdata`, `recommended.txt` and any skill directory not in `OWNED_SKILLS` from that set.
4. MUST keep `OutputsFor` for `make baseline-digests` returning exactly what it returns today.
5. MUST make the mechanical audit accept a changed owned mirror under a grant naming `make skills-sync`, and still refuse it under a grant that names neither the path nor the command.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] The `make skills-sync` output set equals the owned mirrors' files in a fixture tree.
- [ ] Authorial, Go, `testdata`, `recommended.txt` and unowned files are excluded.
- [ ] `make baseline-digests` outputs are unchanged.
- [ ] The audit accepts a mirror change under `make skills-sync` and refuses it without.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/baseline/derived_ownership.go`
- interface: `internal/speccheck/mechanical.go`
- creates: `skills/_ownership.yml`

## Verification

- `out="$(go test -count=1 -v -run "^(TestOutputsForSkillsSyncListsOwnedSkillMirrors|TestOutputsForSkillsSyncExcludesAuthorialAndUnownedFiles|TestOutputsForBaselineDigestsUnchangedBySkillOwnership|TestMechanicalAuditAcceptsSkillMirrorUnderSkillsSync|TestMechanicalAuditRefusesSkillMirrorWithoutSkillsSync)$" ./internal/baseline ./internal/speccheck 2>&1)" || { printf "%s\n" "$out"; exit 1; }; for name in TestOutputsForSkillsSyncListsOwnedSkillMirrors TestOutputsForSkillsSyncExcludesAuthorialAndUnownedFiles TestOutputsForBaselineDigestsUnchangedBySkillOwnership TestMechanicalAuditAcceptsSkillMirrorUnderSkillsSync TestMechanicalAuditRefusesSkillMirrorWithoutSkillsSync; do printf "%s\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task none of the named cases exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — Skill regeneration ownership
- [2026-09-08-skill-regeneration-declares-its-owned-outputs.md](../0121-baseline-decisions-and-complete-regeneration/references/2026-09-08-skill-regeneration-declares-its-owned-outputs.md)
