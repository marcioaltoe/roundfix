---
task: task_04
spec: 0058-npm-trusted-publishing-and-release-preflight
status: pending
type: infra
complexity: high
---

# Task 04: Publish through OIDC with a bounded fallback

## Overview

Publication authenticates with a long-lived repository secret exposed to six
sequential publishes. This task moves authentication to GitHub Actions OIDC and
keeps the release set whole while trusted-publisher configuration is still
unproven: a coordinate that fails OIDC authentication retries once with the
existing token during a bounded window, and every coordinate that needed the
retry is recorded so the window has an evidence-based exit condition.

## Requirements

1. MUST publish through OIDC Trusted Publishing by default, with no
   unconditional token exposed to the publish stage.
2. MUST retry a coordinate whose OIDC publish fails with the existing npm token
   only while the fallback window is switched on, and MUST fail the run with an
   `identity:` prefixed error naming the coordinate when the window is off.
3. MUST record every coordinate that required the fallback and surface that
   list in the run's summary, so an empty list is the proof required to close
   the window.
4. MUST NOT read, print, echo, or persist the token value; the token may be
   referenced only as a secret expression inside the fallback branch, and the
   recorded list carries coordinate names only.
5. MUST preserve the existing publication contract exactly: all five platform
   packages publish before the launcher, package names and versions are
   unchanged, and the GitHub Release stage still runs last.
6. MUST make the fallback window switchable without editing the workflow, so
   closing it is a configuration change rather than a code change.
7. MUST confine every change to the bounded authorized path
   `.github/workflows/release.yml`.

## Subtasks

- [ ] Remove the unconditional token from the publish stage environment.
- [ ] Add the per-coordinate publish helper with its OIDC-first attempt.
- [ ] Add the window-gated token retry and the `identity:` failure path.
- [ ] Record coordinates that used the fallback and write them to the summary.
- [ ] Confirm platform-before-launcher ordering is untouched.

## Acceptance Criteria

- [ ] The publish stage does not set the npm auth token unconditionally for all
      publishes; the secret appears only inside the fallback branch.
- [ ] A coordinate failing OIDC while the window is on is retried with the
      token and recorded; the run continues.
- [ ] A coordinate failing OIDC while the window is off produces an `identity:`
      prefixed error naming that coordinate and fails the run.
- [ ] The fallback window is controlled by configuration outside the workflow
      file, so closing it requires no workflow edit.
- [ ] The run summary lists coordinates that used the fallback, and lists none
      when every coordinate published through OIDC.
- [ ] No command in the workflow prints the token value.
- [ ] The five platform packages still publish before the launcher, and the
      GitHub Release stage still runs after publication.
- [ ] `git status --porcelain` shows no path outside
      `.github/workflows/release.yml` and this task file.

## Context

- instruction: `docs/agents/agent-instructions.md`
- interface: `.github/workflows/release.yml`

## Verification

- `grep -q 'identity:' .github/workflows/release.yml` — expected: exit 0; the
  identity failure vocabulary exists.
- `grep -q 'id-token: write' .github/workflows/release.yml` — expected: exit 0;
  OIDC remains granted.
- `grep -c 'NODE_AUTH_TOKEN' .github/workflows/release.yml | grep -qx '1'` —
  expected: exit 0; the token is referenced exactly once, in the fallback
  branch, rather than as stage-wide environment.
- `grep -n 'dist/npm/packages/cli-\|dist/npm/roundfix' .github/workflows/release.yml`
  — expected: the platform package loop appears before the launcher publish,
  preserving publication order.
- `grep -q 'gh release create' .github/workflows/release.yml` — expected: exit
  0; the GitHub Release stage survived.
- `! grep -q 'echo .*NPM_TOKEN\|echo .*NODE_AUTH_TOKEN' .github/workflows/release.yml`
  — expected: exit 0 for the negation; the token is never echoed.
- `grep -q 'vars\.' .github/workflows/release.yml` — expected: exit 0; the
  fallback window is read from configuration, not hardcoded.
- `grep -q 'GITHUB_STEP_SUMMARY' .github/workflows/release.yml` — expected: exit
  0; fallback usage reaches the run summary.

## References

- `_prd.md` → User Stories 1, 3, 4; Core Features 1, 3, 4, 5.
- `_techspec.md` → Implementation Design: Interfaces; Project Constraints:
  Authentication and HTTP; Build Order 5.
- ADR-0031, ADR-0082, ADR-0084.
