---
type: fix
status: promoted
created: 2026-09-30
spec: 0195-owned-skills-and-a-release-step-that-follow-the-bundle
reason: null
---

# An owned skill older than the bundle passes readiness

## Symptom

The binary holds one minimum version, `0.0.2`, for all 14 Roundfix-owned skills. The bundle it carries is already ahead of that: `qa-gate` and `write-tasks` are `0.0.3`, and `implement-spec` is `0.1.0`. A repository that still has `qa-gate` `0.0.2` installed passes Doctor's `skills:` line, and `roundfix baseline update` without `--yes` answers `current` for it. The repository runs an older gate skill than the binary ships, and nothing says so.

Three related gaps:

- Nothing stops an owned skill's content from changing while its version stays the same. Two copies at one version can differ, so a version comparison proves nothing about content.
- The minimum is a second list. Raising a skill's version does not raise its minimum, and three existing tests in `./skills` fail when the Roundfix skill's version moves past the literal minimum, because they look for that literal in the skill file.
- The `go-cli` and `rust-cli` setups do not list the Roundfix skill, and no module gives it a dispatch trigger. The `typescript-bun` setup lists it, also without a trigger.

## Where

`skills/skills.go` (`ownedSkillMinimumVersions`), `internal/cli/baseline_update.go` (the preview returns before it reads any skill), `internal/baseline/assets/setups/go-cli.json` and `rust-cli.json`, and every module's `skillDispatch`.

## Expected

The minimum for an owned skill is the version the binary carries. Doctor and the `baseline update` preview name an installed owned skill that is older. An owned skill's content cannot change unless its version changes, and a test refuses it. Every setup lists the Roundfix skill, and every profile tells an Agent when to load it.

## Evidence

Measured on `9e439dbb` on 2026-09-30. Raising both version fields of the Roundfix skill to `0.0.3` in a scratch clone made `TestOwnedSkillBundleReadinessKeepsStatesDistinct`, `TestOwnedSkillContractRejectsSetAndVersionDisagreement` and `TestCheckRepositoryClassifiesMissingAndOutdatedSkills` fail. The upstream skill lists `setups/go-cli.txt` and `setups/rust-cli.txt`, read through the Secondbrain mirror `projects/skills/mirror/`, omit `skills/06-review-repair/roundfix`; `setups/typescript-bun.txt` lists it on line 107.
