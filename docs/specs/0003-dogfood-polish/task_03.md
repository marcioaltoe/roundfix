---
task: task_03
spec: 0003-dogfood-polish
status: pending
type: test
complexity: medium
---

# Task 03: Make temp-repo tests hermetic against user git config

## Overview

Six tests failed on a machine with global `commit.gpgsign=true` and produced
an environmental QA `fail` on spec 0002. Every test helper that creates a git
repository must isolate it from user, global, and system git configuration so
the suite's result reflects the code on any machine. Verifiable by a canary
test that enables signing in a scoped config and proves the helpers stay
green.

## Requirements

1. MUST make every repo-creating test helper (across the daemon, cli, store,
   and spec test suites — find them all) run git with user/global/system
   config isolated: `GIT_CONFIG_GLOBAL=/dev/null`, `GIT_CONFIG_SYSTEM=/dev/null`
   on the helper's git invocations, plus explicit `user.name`, `user.email`,
   and `commit.gpgsign=false` in the temp repo itself.
2. MUST keep production code untouched — the daemon already passes explicit
   `-c` config; this task is test-side only.
3. MUST add one canary test proving hermeticity: with signing forced on via a
   scoped config visible to an unisolated git call, the helper-created repo
   still commits without gpg.
4. MUST NOT weaken any assertion; the fix is isolation, not skipping.

## Subtasks

- [ ] Inventory of repo-creating helpers across test suites
- [ ] Config isolation applied to every helper
- [ ] Canary test for forced-signing environments
- [ ] Full-suite pass without environment overrides

## Acceptance Criteria

- [ ] `rg` shows no test helper creating a git repo without the isolation
      pattern (inventory recorded in the Result).
- [ ] The canary test passes; removing the isolation makes it fail (proven
      once during development, stated in the Result).
- [ ] The full suite passes with no `GIT_CONFIG_*` overrides in the
      environment.

## Verification

- `rtk go test ./...` — expected: full suite passes with no environment
  overrides.
- `rtk go test ./internal/daemon/ ./internal/cli/ ./internal/store/` —
  expected: all pass.

## References

`_prd.md` → User Story 3; Core Feature 3; Success Metrics. `_techspec.md` →
Interfaces (test helpers), Build Order 3. Dogfood finding 21; 0002 QA report
`qa-report-2026-07-05.md` (the environmental fail this unblocks).
