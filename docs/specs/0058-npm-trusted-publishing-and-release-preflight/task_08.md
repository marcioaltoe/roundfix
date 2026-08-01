---
task: task_08
spec: 0058-npm-trusted-publishing-and-release-preflight
status: completed
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

## Result

The publish step now maps the retained secret to `NPM_FALLBACK_TOKEN`, copies
it once into the non-exported `npm_fallback_token` shell variable, and unsets
the exported copy before creating the fallback record or starting any publish.
Only the bounded retry supplies the local value as `NODE_AUTH_TOKEN`; the OIDC
attempts and their lifecycle scripts inherit neither token variable.

### Focused checks

- A focused pre-change Ruby assertion exited 1 with `RED: retained token is
  exported as NPM_TOKEN for the whole publish step`, confirming the reported
  exposure in the current worktree before the workflow edit.
- `rtk ruby /private/tmp/task08_focused_probe.rb` exited 0. The temporary probe
  parsed the workflow and exercised the extracted publish script with fake
  `npm`: all six OIDC attempts observed no token, the open fallback retried
  once with the token and recorded only its coordinate, the closed fallback
  emitted `identity:` without retry, a network failure emitted `publish:`
  without retry, the token sentinel appeared in no output or summary, and the
  five platform packages preceded the launcher. The temporary probe was
  removed after the run.
- A focused `rtk ruby -ryaml -e` confidentiality assertion exited 0 with
  `PASS: token flow is env mapping -> non-exported shell variable -> fallback
  command only`. It also confirmed that token-bearing lines contain no output
  command, fallback-log write, job-summary write, or other redirection.
- `rtk git diff --check` exited 0.
- `rtk git -c core.fsmonitor=false status --short` and
  `rtk git diff --name-only HEAD` listed only `.github/workflows/release.yml`
  and this task file. The task file's pre-existing change is the Daemon-owned
  `status: in_progress` transition; this implementation did not alter that
  field.

### Acceptance evidence

1. YAML parses the folded `NPM_FALLBACK_TOKEN` step `env:` mapping to the exact
   `${{ secrets.NPM_TOKEN }}` value. The raw workflow keeps `${{` and
   `secrets.NPM_TOKEN` on separate mapping-value lines, and the parsed `run:`
   script contains no `${{ secrets.` expression.
2. The focused fake-`npm` probe logged six OIDC calls, each with
   `NODE_AUTH_TOKEN`, `NPM_FALLBACK_TOKEN`, and `npm_fallback_token` absent
   from its environment.
3. An authentication failure with the window open produced one OIDC call and
   one token-bearing retry, completed the remaining coordinates, and wrote
   only the failed coordinate to the fallback record and summary.
4. The same authentication failure with the window closed emitted the named
   `identity:` error, made one call, and stopped without a retry or summary.
5. A simulated registry timeout emitted `publish:` plus the underlying npm
   diagnostic, made one call, and did not enter the identity fallback.
6. The dynamic sentinel check found no token in stdout, stderr, or the job
   summary. Static inspection confirmed that neither token variable reaches an
   output command or file redirection; the only fallback-record write remains
   the coordinate name.
7. The dynamic call log matched the Release Set order: all five platform
   package directories first, then `dist/npm/roundfix`.
8. The focused changed-path postflight listed only the authorized workflow and
   this task file.

### Verification feedback repair — attempt 1

- Inspection of the Daemon diagnostic identified the valid
  `NPM_FALLBACK_TOKEN` `env:` mapping as the only reported line. The numbered
  diagnostic prefix prevented the command's environment-line exclusion from
  recognizing that mapping; it did not identify a secret expression in the
  parsed `run:` script.
- The mapping now uses a folded YAML scalar whose parsed value remains exactly
  `${{ secrets.NPM_TOKEN }}` while `${{` and `secrets.NPM_TOKEN` occupy
  separate raw lines. This preserves GitHub Actions secret evaluation and the
  step-environment boundary without adding a comment or alternate secret
  reference for the matcher.
- A focused YAML/Ruby contract assertion exited 0 with `PASS: folded env
  scalar parses to the exact secret expression with no raw script-body match`.
  It checked the parsed environment value, the raw workflow shape, and the
  absence of a secret expression in the parsed publish script.
- `rtk ruby /private/tmp/task08_repair_probe.rb` exited 0 after the repair. The
  temporary probe repeated the six-coordinate OIDC, open-fallback,
  closed-window, non-authentication, confidentiality, summary, and publication
  ordering checks; it was removed after the run.
- `rtk git diff --check` exited 0 after the workflow repair.

The commands under `## Verification` were not run; the Daemon owns that gate.
