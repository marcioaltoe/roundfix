---
task: task_01
spec: 0195-owned-skills-and-a-release-step-that-follow-the-bundle
status: pending
type: backend
complexity: high
---

# Task 01: An installed owned skill older than the bundle is reported

## Overview

The binary compares each installed owned skill with a minimum version, and that minimum is one literal for all 14 skills. The bundle is already ahead of it, so a repository with an older `qa-gate` passes Doctor and the `baseline update` preview calls it current. This Task builds the minimum from the embedded bundle, makes the preview read the installed owned skills, and changes the two documents and the two skill sentences that describe that behavior.

Doctor needs no code change: it already renders a skill below its minimum. The Task proves the new result with a test.

This is an authorized tooling Task. It may change only the files in its Context, the derived files the sanctioned regeneration rewrites, and this Task file.

## Requirements

1. MUST build `ownedSkillMinimumVersions` in `skills/skills.go` from the embedded bundle, as the TechSpec's Interfaces section states: each owned skill maps to the version its embedded `SKILL.md` declares, and a skill with no valid embedded version keeps the base constant. No literal per-skill list remains.
2. MUST make the existing subtest `owned version below minimum` of `TestOwnedSkillContractRejectsSetAndVersionDisagreement` read the minimum from the same map the code reads. No top-level test is renamed or removed, and no other assertion of an existing test changes.
3. MUST implement "The preview rule" of the TechSpec in `internal/cli/baseline_update.go`, including the exact message, next action, status value, field names and text lines it gives.
4. MUST leave the result byte-identical when no installed owned skill is older, when `--no-skills` is given, when the owned skills are absent, and when the read returns an error. `--yes` and `--confirm-plan` keep their behavior.
5. MUST add `skills/owned_skill_minimum_test.go` and `internal/cli/baseline_update_outdated_skills_test.go` with the seven tests the TechSpec's Testing Approach 1 names. Each builds its repository in a temporary directory and reads no Run Database.
6. MUST apply the four text changes the TechSpec's "Fixed texts" gives for the Roundfix skill, `docs/user-guide/commands.md` and `docs/user-guide/context-driven-development.md`. Edit `.agents/skills/roundfix/SKILL.md`, raise both of its version fields by one patch step from their value on this Task's base, and regenerate the mirror with `make skills-sync`.
7. MUST NOT change the `### QA settlement` section of any skill, any other owned skill, any `minimumVersion` in a setup snapshot, `internal/cli/cli_test.go`, or the result's schema version string.
8. MUST keep every phrase the existing contract tests require from the Roundfix skill and the user guide.

## Subtasks

- [ ] Build the minimum from the embedded bundle and adapt the one subtest.
- [ ] Read the installed owned skills in the preview and report the older ones.
- [ ] Add the minimum tests, the preview tests and the Doctor test.
- [ ] Change the two skill sentences and the two guide sentences, raise the skill version and regenerate the mirror.

## Acceptance Criteria

- [ ] For every owned skill the minimum equals the embedded version, and a lowered entry is reported by the comparison.
- [ ] A repository whose installed owned skill is one patch step below the bundle gets `skills: failed` from Doctor, with the skill, both versions and the install command.
- [ ] The preview for that repository exits `3` with `plan_ready`, `skills.status` `outdated` and one `skills.outdated` entry; with `--no-skills` it reports `current`.
- [ ] A repository whose owned skills match gets `current` and no `outdated` field.
- [ ] The Roundfix skill and its mirror are byte-identical, and both version fields are one patch step above the base.

## Context

- instruction: `docs/adr/0189-an-owned-skills-version-names-its-content-and-its-minimum-is-the-bundle.md`
- instruction: `internal/cli/doctor.go`
- instruction: `skills/repository.go`
- interface: `skills/skills.go`
- interface: `skills/skills_test.go`
- interface: `internal/cli/baseline_update.go`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `docs/user-guide/commands.md`
- interface: `docs/user-guide/context-driven-development.md`
- creates: `skills/owned_skill_minimum_test.go`
- creates: `internal/cli/baseline_update_outdated_skills_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestTheOwnedSkillMinimumIsTheEmbeddedVersion|TestAMinimumThatDiffersFromTheEmbeddedVersionIsReported|TestAnInstalledOwnedSkillOlderThanTheBundleIsBelow|TestOwnedSkillContractRejectsSetAndVersionDisagreement|TestOwnedSkillBundleReadinessKeepsStatesDistinct|TestCheckRepositoryClassifiesMissingAndOutdatedSkills|TestAuthorialSkillSync|TestBaselineUpdatePreviewReportsAnOlderOwnedSkill|TestBaselineUpdatePreviewStaysCurrentWhenOwnedSkillsMatch|TestBaselineUpdatePreviewWithNoSkillsSkipsTheOwnedSkillCheck|TestDoctorFailsForAnOwnedSkillOlderThanTheBundle|TestBaselineUpdateFleetSweep)$" ./skills ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestTheOwnedSkillMinimumIsTheEmbeddedVersion TestAMinimumThatDiffersFromTheEmbeddedVersionIsReported TestAnInstalledOwnedSkillOlderThanTheBundleIsBelow TestOwnedSkillContractRejectsSetAndVersionDisagreement TestOwnedSkillBundleReadinessKeepsStatesDistinct TestCheckRepositoryClassifiesMissingAndOutdatedSkills TestAuthorialSkillSync TestBaselineUpdatePreviewReportsAnOlderOwnedSkill TestBaselineUpdatePreviewStaysCurrentWhenOwnedSkillsMatch TestBaselineUpdatePreviewWithNoSkillsSkipsTheOwnedSkillCheck TestDoctorFailsForAnOwnedSkillOlderThanTheBundle TestBaselineUpdateFleetSweep; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && for pair in ".agents/skills/roundfix/SKILL.md|the version of that skill the running binary carries" ".agents/skills/roundfix/SKILL.md|older than the version the binary carries" "docs/user-guide/commands.md|the version of that skill the running binary carries" "docs/user-guide/context-driven-development.md|older than the version the binary carries"; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done && cmp -s .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md` — expected: exit 0; before this Task none of the seven new named tests exists and neither document carries the new sentences, so the command fails.

## References

- [_techspec.md](_techspec.md) — Interfaces; The preview rule; Fixed texts; Testing Approach 1
- `_prd.md` → Goal 1; Core Feature 1; Success Metrics 1 and 2
- `_techspec.md` → API Contracts 1 and 2
- [references/2026-09-30-an-owned-skill-older-than-the-bundle-passes-readiness.md](references/2026-09-30-an-owned-skill-older-than-the-bundle-passes-readiness.md)
- ADR-0062, ADR-0189

## Result
