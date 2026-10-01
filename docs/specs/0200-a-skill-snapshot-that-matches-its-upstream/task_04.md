---
task: task_04
spec: 0200-a-skill-snapshot-that-matches-its-upstream
status: completed
type: chore
complexity: medium
---

# Task 04: This repository takes its own update

## Overview

After task_02 and task_03 this repository's catalog requires `context7-cli`, `exa-web-search` and three Go skills that its lock does not track, and it still holds the old `context7` directory and lock entry. This Task runs the adopter's path on this repository: the public update restores the required skills, the reconcile removes the obsolete lock entry, and the old directory is deleted. It documents the same path for adopters and corrects the Doctor line the guides print.

This is an authorized tooling Task. It may change only the files in its Context and this Task file. The `.agents/skills/` files it may change are exactly the twelve the grant lists.

## Requirements

1. MUST run step 1 of the TechSpec's repository procedure, `roundfix baseline update --repo . --yes --skills-source-dir <dir>`, with the directory chosen as the TechSpec's refresh procedure chooses it. MUST NOT copy skill trees by hand and MUST NOT write `~/dev/skills`.
2. MUST run step 2, `roundfix baseline skills reconcile` at `a4e18e4fa223196b51d0fd8224e5a33b84f97717`, preview then confirm, so the `context7` lock entry is removed.
3. MUST delete `.agents/skills/context7/`, which no command removes.
4. MUST make `skills/recommended.txt` the sorted keys of `skills-lock.json`, and set `upstreamManagedSkillTreeDigest` in `skills/baseline_skill_contract_test.go` to the value `TestAuthorialSkillSync` reports. No other line of that file changes.
5. MUST stop without committing when step 1 or 2 changes a `.agents/skills/` file outside the twelve the grant lists, or any Roundfix-owned skill file.
6. MUST add `internal/cli/this_repository_skill_set_test.go` with the two tests of the TechSpec's Testing Approach 4. Both call one helper that resolves the external set with `resolveExternalSkillRequirement`, checks readiness with `skills.CheckRepositoryWithExternal`, and reports each missing or outdated skill and each lock entry or directory named `context7`, `feature-systems-pattern` or `rust`.
7. MUST replace `skills: ok (39 required: 14 Roundfix-owned, 25 external)` with `skills: ok (42 required: 14 Roundfix-owned, 28 external)` in `README.md`, `docs/user-guide/commands.md` and `docs/user-guide/usage.md`, after confirming with `go run -buildvcs=false ./cmd/roundfix doctor` that this repository prints that line.
8. MUST add to the "Repository Skill Set restoration" section of `docs/user-guide/context-driven-development.md` a subsection titled "When an upstream skill is renamed", stating: the next `roundfix baseline update --yes` installs the new name; `roundfix baseline skills reconcile --profile <id> --source <owner/repo> --revision <commit>` then removes the old lock entry; the old directory under `.agents/skills/` is the repository's to delete, because Roundfix never deletes an installed skill tree; and a second update reports `current`.
9. A second `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --format text` MUST report `current`.

## Subtasks

- [ ] Run the public update with the upstream source directory.
- [ ] Reconcile the lock and delete the old directory.
- [ ] Align the recommended list and the upstream digest pin with the lock.
- [ ] Add the repository skill-set tests.
- [ ] Document the rename path and correct the Doctor line.

## Acceptance Criteria

- [ ] Every required external skill of this repository is installed with its lock hash.
- [ ] No lock entry or directory names `context7`, `feature-systems-pattern` or `rust`, and a repository missing a required skill is reported by the same helper.
- [ ] The recommended list equals the lock, and the upstream digest pin matches.
- [ ] The guides print the Doctor line this repository prints, and the user guide explains the rename path.
- [ ] A second update reports `current`.

## Context

- instruction: `docs/adr/0191-a-setup-snapshot-follows-its-upstream-by-name.md`
- instruction: `internal/cli/doctor.go`
- interface: `skills-lock.json`
- interface: `skills/recommended.txt`
- interface: `skills/baseline_skill_contract_test.go`
- creates: `.agents/skills/context7/SKILL.md`
- creates: `.agents/skills/context7-cli/SKILL.md`
- creates: `.agents/skills/context7-cli/references/docs.md`
- creates: `.agents/skills/context7-cli/references/setup.md`
- creates: `.agents/skills/context7-cli/references/skills.md`
- interface: `.agents/skills/golang-dependency-management/SKILL.md`
- interface: `.agents/skills/golang-safety/SKILL.md`
- interface: `.agents/skills/golang-safety/references/nil-safety.md`
- interface: `.agents/skills/golang-safety/references/slice-map-safety.md`
- interface: `.agents/skills/golang-structs-interfaces/SKILL.md`
- creates: `.agents/skills/golang-structs-interfaces/references/struct-fields.md`
- creates: `.agents/skills/golang-structs-interfaces/references/type-assertions.md`
- interface: `README.md`
- interface: `docs/user-guide/commands.md`
- interface: `docs/user-guide/usage.md`
- interface: `docs/user-guide/context-driven-development.md`
- creates: `internal/cli/this_repository_skill_set_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestThisRepositoryHoldsEveryRequiredExternalSkill|TestARepositoryMissingARequiredExternalSkillIsReported)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestThisRepositoryHoldsEveryRequiredExternalSkill TestARepositoryMissingARequiredExternalSkillIsReported; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && out="$(go test -count=1 -v -run '^(TestRecommendedSkillsMatchLock|TestAuthorialSkillSync|TestAuthoringConstraintOwnership|TestUpstreamADRFormatUnchanged)$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestRecommendedSkillsMatchLock TestAuthorialSkillSync TestAuthoringConstraintOwnership TestUpstreamADRFormatUnchanged; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && test ! -e .agents/skills/context7 && test -f .agents/skills/context7-cli/SKILL.md && for pair in "README.md|skills: ok (42 required: 14 Roundfix-owned, 28 external)" "docs/user-guide/commands.md|skills: ok (42 required: 14 Roundfix-owned, 28 external)" "docs/user-guide/usage.md|skills: ok (42 required: 14 Roundfix-owned, 28 external)" "docs/user-guide/context-driven-development.md|When an upstream skill is renamed" "docs/user-guide/context-driven-development.md|never deletes an installed skill tree"; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done && go run -buildvcs=false ./cmd/roundfix baseline update --repo . --format json >/dev/null` — expected: exit 0; before this Task the two new tests do not exist, `context7-cli` is not installed and `.agents/skills/context7` exists, so the command fails.

## References

- `_prd.md` → Goal 4; User Story 4; Core Feature 5; Success Metric 5; Non-Goals (skills that match their lock)
- `_techspec.md` → Measured facts; The repository procedure; API Contract 5; Testing Approach 4; Build Order 4
- ADR-0103, ADR-0191

## Result

The public update at upstream revision `a4e18e4fa223196b51d0fd8224e5a33b84f97717` restored `context7-cli` and the three required Go skills through the Roundfix CLI. The reconcile preview produced plan digest `6f81eb25b256b41dcf90993cf2bb78e4df4b94274cd1f756831497818dad4d48`; the confirmed plan removed only the obsolete `context7` lock entry, and the old `context7` directory was then deleted.

The repository skill-set tests now share a helper that resolves the external requirement from the Setup Manifest, checks `skills.CheckRepositoryWithExternal`, logs missing and outdated skills, and reports legacy lock entries or directories. The lock-derived recommended list and upstream managed tree digest are aligned; the focused digest test reported `3b19955019e42be5f583157c12eadc58652031361b6e19af8c11a90cf6b9389f`. The guides document the upstream rename path and print `skills: ok (42 required: 14 Roundfix-owned, 28 external)`.

Focused checks passed: `go test -count=1 -run 'TestThisRepositoryHoldsEveryRequiredExternalSkill|TestARepositoryMissingARequiredExternalSkillIsReported' ./internal/cli`; `go test -count=1 -run 'TestRecommendedSkillsMatchLock|TestAuthorialSkillSync|TestAuthoringConstraintOwnership|TestUpstreamADRFormatUnchanged' ./skills`; and the required second public update reported `Baseline update: current`. The authored Verification block was not rerun by this Agent.

## Carry-forward provenance

- Source Run: `run_20261001T013631Z_c59b100cbb9f7af1`
- Source commit: `664fdb32945c79945de1b4e982950ee0ea47dea4`
