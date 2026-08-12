# Scope and history audit

Build: `171f6a378c9e640a8a10c9382e28b501b21ff5a0`

`rtk ruby .../current_scope_probe.rb` exited 0.

- All seven Tasks are `completed` and each Result carries acceptance evidence.
- The PRD and TechSpec each account for identifier strategy,
  authentication/HTTP, active ADR obligations, and tooling authority through
  operative `docs/agents/` sources.
- Commit `397227ff` contains express authorization for exactly
  `.github/workflows/release.yml` and predates every protected-tooling change.
- `git diff-tree --no-commit-id --name-only -r` matched the exact bounded paths
  for Task commits `21bc4bf`, `8d14a67`, `b0052e9`, `47de307`, `b411b30`,
  `551433d`, `ab34e03`, and `4493add`.
- Consequent documentation fix `04d5bbb` changes only the release runbook;
  review fix `83fbf6b` changes only the authorized workflow; review-artifact
  commit `171f6a3` changes only the eight Spec-local review files.
- Ancestry is authorization `397227ff` → final remediation Task `4493add` →
  documentation fix `04d5bbb` → workflow review fix `83fbf6b` → separate
  review-artifact commit `171f6a3`.
- Round 002 settles two actionable issues as `resolved` and two acknowledgement
  comments as `invalid`.
- The current worktree delta contains only this QA report and its evidence.

The workflow probes independently confirm the unchanged package names,
registry, version authority, five-platform-before-launcher order, asset path,
Release Plan exclusion, and GitHub Release-last contract.
