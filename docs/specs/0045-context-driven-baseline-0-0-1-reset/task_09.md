---
task: task_09
spec: 0045-context-driven-baseline-0-0-1-reset
status: pending
type: backend
complexity: high
---

# Task 09: Align maintained profiles with the Repository Skill Set

## Overview

Move every maintained setup profile onto the new owned generation while
preserving its substantive language. Align the shipped Repository Skill Set,
profile snapshots, and restoration checks with the exact 0.0.1 activation and
capability contract.

## Requirements

1. MUST migrate maintained Go and Rust profiles to 0.0.1 without changing
   their content except for owned generation metadata and universal Repository
   Capability guidance.
2. MUST replace the superseded TypeScript profile identity with
   `standard-typescript-monorepo` and reject the former identity as current
   state.
3. MUST align every profile snapshot with exact Skill Activation membership,
   universal Context7/Exa requirements, and Firecrawl/`rtk`/`rg` recommendations.
4. MUST validate the complete Roundfix-owned Repository Skill Set and reject
   missing, unexpected, duplicate, or version-disagreeing owned skills.
5. MUST make skill restoration reproduce exact trusted bytes and commit
   evidence for every repo-owned skill required by the profiles.
6. MUST preserve upstream-managed skill metadata and the external
   `skills-lock.json` schema unchanged.
7. MUST keep canonical and embedded setup-context-driven trees byte-identical
   and runnable in place.

## Subtasks

- [ ] Reset Go and Rust profile generation metadata and universal capability
      bindings.
- [ ] Retire the old TypeScript profile identity from current-state loaders.
- [ ] Regenerate deterministic snapshots for every maintained profile.
- [ ] Extend Repository Skill Set membership and owned-version validation.
- [ ] Align trusted restoration metadata and exact-byte verification.
- [ ] Add upstream-metadata and external-schema preservation assertions.
- [ ] Synchronize and run both setup skill copies in place.

## Acceptance Criteria

- [ ] Go and Rust generated content differs from its prior governed content
      only where the 0.0.1 generation and universal capabilities require it.
- [ ] Every maintained profile reports the same exact owned skill membership
      as its snapshot and fails on missing, extra, duplicate, or mismatched
      members.
- [ ] The former TypeScript profile ID is rejected as current state and routed
      through Baseline Readoption when encountered as source evidence.
- [ ] Skill restoration verifies the immutable source commit and exact restored
      bytes before reporting success.
- [ ] Upstream skill metadata and external `skills-lock.json` schema fixtures
      remain byte-identical.
- [ ] Canonical and embedded test suites are runnable from their own trees and
      produce the same result.

## Context

- instruction: `docs/agents/skill-governance.md`
- interface: `.agents/skills/setup-context-driven/assets/profiles`
- interface: `.agents/skills/setup-context-driven/scripts/context_setup.py`
- interface: `.agents/skills/setup-context-driven/tests/test_macro_profiles.py`
- interface: `.agents/skills/setup-context-driven/tests/test_restore_skills.py`
- interface: `skills/skills.go`
- interface: `skills/skills_test.go`

## Verification

- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_profile_alignment.py'` — expected: maintained profiles, snapshots, universal capabilities, and exact owned skill bundles agree.
- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s skills/setup-context-driven/tests -p 'test_restore_skills.py'` — expected: the distributed copy runs in place and restoration proves trusted exact bytes.
- `rtk go test ./skills` — expected: Repository Skill Set membership and owned-version validation pass while upstream metadata remains outside ownership.
- `rtk make skills-sync-check` — expected: canonical and distributed setup skill trees are byte-identical.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → Goals 1 and 4; Core Features 1, 10, and 14; User Stories 2
  and 4.
- `_techspec.md` → Data Models; Integration Points; Build Order 6.
- ADR-0061 → maintained profile policy.
- ADR-0062 → owned-version scope and protected upstream contracts.
