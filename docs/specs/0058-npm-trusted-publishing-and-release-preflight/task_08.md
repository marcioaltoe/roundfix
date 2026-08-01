---
task: task_08
spec: 0058-npm-trusted-publishing-and-release-preflight
status: pending
type: infra
complexity: medium
---

# Task 08: Confine the retained token without interpolating it

## Overview

Two correct constraints currently conflict. Security review requires the
secret to reach the script through the step environment rather than being
interpolated into the generated shell text. Task 04 requires the retained token
to stay confined to the fallback branch. Because a step's environment is
inherited by every command in that step, satisfying the first widened the
second: the six OIDC publish attempts now run with the token in their
environment, and `npm publish` executes package lifecycle scripts that would
see it. This task satisfies both — the secret never appears in script text, and
no command except the fallback publish can read it.

## Requirements

1. MUST NOT interpolate any `${{ secrets.* }}` expression into a `run:` script
   body; the secret reaches the script only through the step environment.
2. MUST prevent every command in the publish step other than the fallback
   publish from inheriting the token, including the six OIDC publish attempts
   and any package lifecycle script they invoke.
3. MUST make the token available to the fallback publish command itself, so the
   bounded fallback of ADR-0084 still works.
4. MUST keep the token out of logs, the fallback record, the job summary, and
   any file on disk.
5. MUST preserve every behavior tasks 04 and 06 established: OIDC-first
   publication, evidence-based failure attribution, the closed-window
   `identity:` error, the `publish:` path with no retry, the fallback record,
   and platform-before-launcher ordering.
6. MUST confine every change to the bounded authorized path
   `.github/workflows/release.yml`.

## Subtasks

- [ ] Move the secret into a step environment variable that the script consumes
      once.
- [ ] Transfer it into a non-exported shell variable and remove it from the
      exported environment before any publish runs.
- [ ] Supply it to the fallback publish command only.
- [ ] Confirm the OIDC attempts observe no token in their environment.

## Acceptance Criteria

- [ ] No `run:` script body in the workflow contains a `${{ secrets.` /
      expression; the secret is referenced only in an `env:` mapping.
- [ ] A probe of the publish step shows the token absent from the environment
      observed by each of the six OIDC publish attempts.
- [ ] The fallback publish still receives the token and still records the
      coordinate that used it.
- [ ] With the window closed, an authentication failure still produces the
      `identity:` error and no retry.
- [ ] A non-authentication failure still produces `publish:` and no retry.
- [ ] The token appears in no log line, no fallback record, no job summary, and
      no file.
- [ ] The five platform packages still publish before the launcher.
- [ ] `git status --porcelain` shows no path outside
      `.github/workflows/release.yml` and this task file.

## Context

- instruction: `docs/agents/agent-instructions.md`
- interface: `.github/workflows/release.yml`
- interface: `docs/specs/0058-npm-trusted-publishing-and-release-preflight/qa/evidence/2026-08-01-qa-03/secret_boundary_probe.rb`

## Verification

- `! grep -nE 'NODE_AUTH_TOKEN="\$\{\{|\$\{\{ *secrets\.' .github/workflows/release.yml | grep -v '^ *[A-Z_]*:' `
  — expected: exit 0 for the negation; no secret expression sits inside a
  script body.
- `grep -q 'secrets.NPM_TOKEN' .github/workflows/release.yml` — expected: exit
  0; the secret is still mapped through `env:`.
- `grep -q 'unset ' .github/workflows/release.yml` — expected: exit 0; the
  exported copy is removed before publication runs.
- `rtk ruby docs/specs/0058-npm-trusted-publishing-and-release-preflight/qa/evidence/2026-08-01-qa-03/secret_boundary_probe.rb`
  — expected: the probe that proved the step-wide exposure now reports the
  token confined to the fallback command.
- `grep -qE 'ENEEDAUTH|E401|Unable to authenticate' .github/workflows/release.yml`
  — expected: exit 0; evidence-based attribution survives.
- `grep -q 'identity:' .github/workflows/release.yml && grep -q 'publish:' .github/workflows/release.yml`
  — expected: exit 0; both failure paths survive.
- `! grep -q 'echo .*NPM_TOKEN\|echo .*NODE_AUTH_TOKEN' .github/workflows/release.yml`
  — expected: exit 0 for the negation; the token is never echoed.
- `grep -n 'dist/npm/packages/cli-\|dist/npm/roundfix' .github/workflows/release.yml`
  — expected: the platform loop still precedes the launcher publish.

## References

- `_prd.md` → Core Features 4; Project Constraints: Authentication and HTTP.
- `_techspec.md` → Project Constraints: Authentication and HTTP.
- `qa/qa-report-2026-08-01-03.md` → QA-005.
- `docs/specs/_reviews/pr-61/` → the security review that introduced the step
  environment mapping.
- ADR-0084.
