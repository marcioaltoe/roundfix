---
task: task_06
spec: 0058-npm-trusted-publishing-and-release-preflight
status: pending
type: infra
complexity: medium
---

# Task 06: Attribute a publish failure to its actual cause

## Overview

Every nonzero first `npm publish` currently enters the identity branch and is
retried with the token, so a network timeout is reported as an authentication
problem and burns the fallback. This task makes attribution evidence-based:
only a proven authentication failure may be called `identity:` and may reach
the fallback, and every other failure keeps its own cause and fails the run
without a retry. Verifiable on its own: a simulated non-authentication failure
is reported under a non-identity prefix and triggers no token retry.

## Requirements

1. MUST capture the failing publish attempt's output and classify the cause
   before choosing a branch, rather than treating any nonzero result as an
   identity failure.
2. MUST enter the token fallback only when the captured evidence shows an
   authentication failure, matching the signals npm emits for a rejected or
   absent publishing identity.
3. MUST report a failure that is not an authentication failure under a
   distinct prefix that is not `identity:`, name the coordinate, surface the
   underlying error, and fail the run without any token retry.
4. MUST preserve the existing `identity:` behavior for genuine authentication
   failures, including the closed-window error and the fallback record.
5. MUST NOT print, echo, or persist the token value in any branch, and MUST
   NOT write the captured publish output anywhere it could contain and expose
   a credential.
6. MUST keep publication ordering, package set, and the GitHub Release stage
   unchanged.
7. MUST confine every change to the bounded authorized path
   `.github/workflows/release.yml`.

## Subtasks

- [ ] Capture the publish attempt's exit status and output for classification.
- [ ] Add the authentication-failure test that gates the fallback branch.
- [ ] Add the non-identity failure path with its own prefix and no retry.
- [ ] Confirm the closed-window and fallback-record behavior still holds.

## Acceptance Criteria

- [ ] A simulated authentication failure with the window open is reported as
      `identity:`, retried once with the token, and recorded in the fallback
      record.
- [ ] A simulated authentication failure with the window closed is reported as
      `identity:` and fails the run with no retry.
- [ ] A simulated non-authentication failure — such as a network or validation
      error — is reported under a prefix other than `identity:`, names the
      coordinate, and triggers no token retry.
- [ ] The captured publish output is surfaced for the non-identity case so the
      maintainer sees the real cause.
- [ ] No workflow command prints the token value.
- [ ] The five platform packages still publish before the launcher.
- [ ] `git status --porcelain` shows no path outside
      `.github/workflows/release.yml` and this task file.

## Context

- instruction: `docs/agents/agent-instructions.md`
- interface: `.github/workflows/release.yml`
- interface: `docs/specs/0058-npm-trusted-publishing-and-release-preflight/qa/evidence/2026-07-31-qa-01/failure_attribution_probe.rb`

## Verification

- `grep -q 'identity:' .github/workflows/release.yml` — expected: exit 0; the
  identity vocabulary survives for genuine authentication failures.
- `! grep -q 'npm publish --access public ) && return 0' .github/workflows/release.yml`
  — expected: exit 0 for the negation; the unconditional
  success-or-identity-branch shape is gone.
- `grep -qiE 'ENEEDAUTH|E401|Unable to authenticate' .github/workflows/release.yml`
  — expected: exit 0; the fallback is gated on an authentication signal.
- `grep -c 'NODE_AUTH_TOKEN' .github/workflows/release.yml | grep -qx '1'` —
  expected: exit 0; the token is still referenced exactly once.
- `! grep -q 'echo .*NPM_TOKEN\|echo .*NODE_AUTH_TOKEN' .github/workflows/release.yml`
  — expected: exit 0 for the negation; the token is never echoed.
- `rtk ruby docs/specs/0058-npm-trusted-publishing-and-release-preflight/qa/evidence/2026-07-31-qa-01/failure_attribution_probe.rb`
  — expected: the probe that proved the defect now reports a non-identity
  failure for the simulated non-authentication case.
- `grep -n 'dist/npm/packages/cli-\|dist/npm/roundfix' .github/workflows/release.yml`
  — expected: the platform loop still precedes the launcher publish.

## References

- `_prd.md` → User Story 3; Core Features 3.
- `_techspec.md` → Implementation Design: Interfaces; API Contracts: failure
  vocabulary.
- `qa/qa-report-2026-07-31.md` → QA-002.
- ADR-0084.
