---
task: task_04
spec: 0208-a-baseline-that-follows-the-reshaped-skills-catalog
status: pending
type: chore
complexity: low
---

# Task 04: This repository drops its `review` lock entry and tree

## Overview

After task_01 no module requires `review`, and upstream `b3c45a4` no longer holds it, but this repository still tracks it in `skills-lock.json` and keeps `.agents/skills/review/`. This Task runs the adopter's path: `roundfix baseline skills reconcile` removes the lock entry, and the repository deletes its own tree, because Roundfix never deletes an installed skill tree. `triage` was never installed here; the test proves it stays absent.

This is an authorized tooling Task. It may change only the files in its Context and this Task file. The only `.agents/skills/` path it may change is `.agents/skills/review/SKILL.md`, which it deletes.

## Requirements

1. MUST run steps 2 and 3 of the TechSpec's repository procedure: the reconcile preview, then the confirmed run with its Plan Digest, using a clone of the local `~/dev/skills` at `b3c45a45f1bccd3b33aaecaaa22947d942f2fc02`, then `git rm -r .agents/skills/review`. The confirmed plan MUST contain only `remove-lock-entry skills-lock.json [review]`; any other edit MUST stop the Task.
2. MUST make `skills/recommended.txt` the sorted keys of `skills-lock.json`, and set `upstreamManagedSkillTreeDigest` to the value `TestAuthorialSkillSync` reports. No other line of that file changes.
3. MUST add `review` and `triage` to the names whose lock entry or directory `checkThisRepositorySkillSet` reports, and change no other line of `internal/cli/this_repository_skill_set_test.go`.
4. A second reconcile MUST report no changes, and `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --format text` MUST report `current`.
5. MUST NOT change a Baseline asset, a guide under `docs/agents/`, or any other `.agents/skills/` file.

## Subtasks

- [ ] Preview and confirm the reconcile.
- [ ] Delete the `review` tree.
- [ ] Align the recommended list and the digest pin with the lock.
- [ ] Extend the repository skill-set check.

## Acceptance Criteria

- [ ] No lock entry or directory names `review` or `triage`, and the same helper reports one that does.
- [ ] Every required external skill is still installed with its lock hash.
- [ ] The recommended list equals the lock, and the digest pin matches.
- [ ] A second reconcile reports no changes, and the update reports `current`.

## Context

- instruction: `docs/adr/0191-a-setup-snapshot-follows-its-upstream-by-name.md`
- instruction: `docs/adr/0206-the-baseline-takes-upstream-setup-names-and-drops-a-removed-skill.md`
- interface: `skills-lock.json`
- interface: `skills/recommended.txt`
- interface: `skills/baseline_skill_contract_test.go`
- interface: `.agents/skills/review/SKILL.md`
- interface: `internal/cli/this_repository_skill_set_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestThisRepositoryHoldsEveryRequiredExternalSkill|TestARepositoryMissingARequiredExternalSkillIsReported)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestThisRepositoryHoldsEveryRequiredExternalSkill TestARepositoryMissingARequiredExternalSkillIsReported; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && out="$(go test -count=1 -v -run '^(TestRecommendedSkillsMatchLock|TestAuthorialSkillSync|TestAuthoringConstraintOwnership|TestUpstreamADRFormatUnchanged)$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestRecommendedSkillsMatchLock TestAuthorialSkillSync TestAuthoringConstraintOwnership TestUpstreamADRFormatUnchanged; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && test ! -e .agents/skills/review && test ! -e .agents/skills/triage && ! grep -q '"review"' skills-lock.json && ! grep -qx 'review' skills/recommended.txt && go run -buildvcs=false ./cmd/roundfix baseline update --repo . --format json >/dev/null` — expected: exit 0; before this Task the lock and the tree still hold `review`, so the extended test and the absence checks fail.

## References

- `_prd.md` → Goal 5; Core Feature 5; Success Metric 5; Decisions (dogfood the adopter path)
- `_techspec.md` → The repository procedure; API Contract 5; Testing Approach 4; Build Order 4
- ADR-0191, ADR-0206
