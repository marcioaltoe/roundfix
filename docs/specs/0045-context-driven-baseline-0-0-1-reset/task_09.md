---
task: task_09
spec: 0045-context-driven-baseline-0-0-1-reset
status: completed
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

- [x] Reset Go and Rust profile generation metadata and universal capability
      bindings.
- [x] Retire the old TypeScript profile identity from current-state loaders.
- [x] Regenerate deterministic snapshots for every maintained profile.
- [x] Extend Repository Skill Set membership and owned-version validation.
- [x] Align trusted restoration metadata and exact-byte verification.
- [x] Add upstream-metadata and external-schema preservation assertions.
- [x] Synchronize and run both setup skill copies in place.

## Acceptance Criteria

- [x] Go and Rust generated content differs from its prior governed content
      only where the 0.0.1 generation and universal capabilities require it.
- [x] Every maintained profile reports the same exact owned skill membership
      as its snapshot and fails on missing, extra, duplicate, or mismatched
      members.
- [x] The former TypeScript profile ID is rejected as current state and routed
      through Baseline Readoption when encountered as source evidence.
- [x] Skill restoration verifies the immutable source commit and exact restored
      bytes before reporting success.
- [x] Upstream skill metadata and external `skills-lock.json` schema fixtures
      remain byte-identical.
- [x] Canonical and embedded test suites are runnable from their own trees and
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

## Result

Completed the maintained-profile and Repository Skill Set reset to generation
`0.0.1`.

- Go and Rust retain their existing profile modules, rules, formatter choice,
  and entry decisions; their owned changes are the `0.0.1` schema/version and
  marker metadata plus the universal capability and activation-bundle bindings.
  The regenerated formatter corpus confirms that substantive TypeScript
  guidance changed only for the universal Exa, Firecrawl, and Roundfix dispatch
  ownership required by this Task.
- The current profile catalog now contains exactly `go-cli-tui`, `rust-cli`,
  and `standard-typescript-monorepo`. The former
  `typescript-bun-monorepo` asset was removed, and the alignment test proves
  that legacy occurrences are rejected as current state and surfaced through
  Baseline Readoption as source evidence.
- Every setup snapshot now uses `setup-context-driven/setup-snapshot/0.0.1`,
  has deterministic activation bundles, includes the universal Context7/Exa
  requirements and Firecrawl recommendation, and carries exactly the 14
  Roundfix-owned skills at version `0.0.1`. Mutation tests reject missing,
  unexpected, duplicate, bundle-mismatched, and version-disagreeing members.
- Restoration remains bound to immutable Git commit and tree evidence. The
  distributed restoration suite verifies the commit identity and exact restored
  bytes; the Go installer test additionally compares every installed owned file
  with its trusted embedded bytes across all supported targets.
- The profile-alignment test pins the external `skills-lock.json` and isolated
  compatibility-fixture bytes and validates its external schema. No
  upstream-managed skill metadata was added to Roundfix ownership or edited by
  this Task.
- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_profile_alignment.py'` — PASS, 10 tests.
- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s skills/setup-context-driven/tests -p 'test_restore_skills.py'` — PASS, 14 tests.
- `rtk go test ./skills` — PASS, 25 tests.
- `rtk make skills-sync-check` — PASS; canonical and distributed trees are byte-identical.
- `rtk make verify` — PASS; both setup trees passed 239 tests each, Go passed
  1,699 tests across 20 packages, the owned-skill check passed, and the CLI
  build completed.
