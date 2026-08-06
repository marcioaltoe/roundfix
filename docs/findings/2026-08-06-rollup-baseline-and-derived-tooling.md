---
status: pending
created_at: 2026-08-06
updated_at: 2026-08-06
kind: rollup
members:
  - 2026-07-23-setup-context-driven-adoption-process-improvements.md
  - 2026-07-24-greenfield-agent-guidance-acceptance-target.md
  - 2026-07-26-baseline-profile-refresh-retention-gap.md
  - 2026-07-26-vortex-baseline-capability-remediation.md
  - 2026-07-27-bare-go-build-writes-an-untracked-root-binary.md
  - 2026-07-27-derived-skill-digest-pins-have-no-regeneration-path.md
  - 2026-07-27-sandboxed-agents-cannot-reach-the-default-go-cache.md
  - 2026-07-28-tooling-tasks-need-a-green-repo-and-an-undocumented-commit-choreography.md
  - 2026-07-29-doctor-requires-roundfix-own-development-skills.md
  - 2026-07-30-baseline-digest-regeneration-cannot-bootstrap.md
  - 2026-08-01-characterization-corpus-is-outside-the-regeneration-command.md
  - 2026-08-05-what-this-repository-should-change-after-a-full-queue-night.md
---

# Baseline and derived tooling — one owner must reach every consequence (2026-08-06)

The Baseline and tooling findings describe one contract boundary: a declared
source change must preserve semantic guidance, identify every derived artifact,
and provide one sanctioned path that can regenerate the complete consequence
set. File hashes, scattered pins, and repository-specific recovery knowledge
do not provide that guarantee.

## Consolidated learning

- Baseline adoption must account for retained Normative Clauses and
  capabilities, not only file identity or a successful apply.
- A regeneration command must own the whole derivation chain, including the
  pins and characterization corpus that would otherwise block its own run.
- Tooling Tasks need explicit authority, a green precondition, deterministic
  local caches, and commit choreography that distinguishes prerequisite work
  from consequences of the authorized change.
- Readiness diagnostics must derive the Repository Skill Set and name the
  failed probe and a next action that can actually reach green.

## Live edge

Several member defects shipped through Specs 0054, 0057, 0061, 0062, and 0067.
The rollup remains `pending` because the accumulated evidence still asks for
one mechanically owned derivation path and semantic retention proof across the
whole Baseline lifecycle.
