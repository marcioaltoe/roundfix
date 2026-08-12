---
task: task_02
spec: 0058-npm-trusted-publishing-and-release-preflight
status: completed
type: infra
complexity: high
---

# Task 02: Preflight the release set against the registry

## Overview

Publication currently starts before anything proves the target version is
available on every coordinate, so a single blocked name can leave installable
platform packages with no launcher. This task adds a read-only Publication
Preflight between the Verification gate and cross-compilation that classifies
all six coordinates and stops the run before the first irreversible byte.
Verifiable on its own: the classifier resolves every state against committed
fixture packuments, and the stage sits in the required position.

## Requirements

1. MUST derive the Release Set from the same sources the publish loop uses —
   the five platform packages in `dist/npm/platforms.json` plus the launcher —
   so the checked set can never drift from the published set.
2. MUST classify each coordinate from the unauthenticated public packument,
   requiring no credential, resolving exactly these states: `eligible`,
   `used` when the target version appears in the version map or in the release
   time map, `cooldown` when the package carries a full-unpublish marker, and
   `absent` when the registry has no such package.
3. MUST treat any transport failure, non-success status, or unparseable body as
   `undetermined` — the run stops, but the coordinate is never reported as
   ineligible.
4. MUST expose the classification program on a single line as a single-quoted
   `CLASSIFY_JQ='...'` shell assignment so one source of truth serves both the
   workflow and this task's fixture checks.
5. MUST print one row per coordinate naming the package, the target version,
   and the resolved state, and MUST emit a `registry:` prefixed error naming
   each blocked coordinate, including the cooldown expiry timestamp when that
   is the cause.
6. MUST stop the run before `Cross-compile and stage` when any coordinate is
   not `eligible`, and MUST allow the run to continue only when every
   coordinate is `eligible`.
7. MUST commit fixture packuments under this Spec's folder covering the used,
   single-version-unpublished, cooldown, eligible, and malformed cases.
8. MUST confine workflow changes to the bounded authorized path
   `.github/workflows/release.yml`; fixture data belongs under this Spec's
   folder and is not repository tooling.

## Subtasks

- [ ] Write the classification program and expose it as `CLASSIFY_JQ`.
- [ ] Commit the five fixture packuments covering every resolvable state.
- [ ] Add the preflight stage that iterates the Release Set and prints the
      eligibility table.
- [ ] Wire the blocking decision so a non-eligible or undetermined coordinate
      stops the run before cross-compilation.
- [ ] Confirm the stage order and that no other stage changed.

## Acceptance Criteria

- [ ] The classifier resolves `used` for a fixture whose version map contains
      the target version, and `used` for a fixture where the version was
      removed from the version map but remains in the release time map.
- [ ] The classifier resolves `cooldown` for a fixture carrying a full-unpublish
      marker and `eligible` for a fixture with neither condition.
- [ ] The classifier exits non-zero on a malformed body rather than reporting a
      state, so malformed input can never read as eligible.
- [ ] The preflight stage appears after the Verification gate and before
      `Cross-compile and stage` in the workflow.
- [ ] Blocked coordinates are reported with a `registry:` prefix and
      undetermined reads with an `undetermined:` prefix.
- [ ] `git status --porcelain` shows no path outside
      `.github/workflows/release.yml`, this Spec's folder, and this task file.

## Context

- interface: `.github/workflows/release.yml`
- interface: `dist/npm/platforms.json`
- interface: `dist/npm/roundfix/package.json`

## Verification

- `grep -q "CLASSIFY_JQ='" .github/workflows/release.yml` — expected: exit 0;
  the classifier is exposed as a single-quoted single-line assignment.
- `jq -r --arg tag 0.0.2 "$(sed -n "s/^[[:space:]]*CLASSIFY_JQ='\(.*\)'$/\1/p" .github/workflows/release.yml)" docs/specs/0058-npm-trusted-publishing-and-release-preflight/fixtures/version-used.json`
  — expected: prints `used`.
- `jq -r --arg tag 0.0.2 "$(sed -n "s/^[[:space:]]*CLASSIFY_JQ='\(.*\)'$/\1/p" .github/workflows/release.yml)" docs/specs/0058-npm-trusted-publishing-and-release-preflight/fixtures/version-unpublished-single.json`
  — expected: prints `used`; a single-version unpublish can never be reused.
- `jq -r --arg tag 0.0.2 "$(sed -n "s/^[[:space:]]*CLASSIFY_JQ='\(.*\)'$/\1/p" .github/workflows/release.yml)" docs/specs/0058-npm-trusted-publishing-and-release-preflight/fixtures/package-cooldown.json`
  — expected: prints `cooldown`.
- `jq -r --arg tag 0.0.2 "$(sed -n "s/^[[:space:]]*CLASSIFY_JQ='\(.*\)'$/\1/p" .github/workflows/release.yml)" docs/specs/0058-npm-trusted-publishing-and-release-preflight/fixtures/eligible.json`
  — expected: prints `eligible`.
- `! jq -r --arg tag 0.0.2 "$(sed -n "s/^[[:space:]]*CLASSIFY_JQ='\(.*\)'$/\1/p" .github/workflows/release.yml)" docs/specs/0058-npm-trusted-publishing-and-release-preflight/fixtures/malformed.json`
  — expected: exit 0 for the negation; malformed input fails instead of
  resolving a state.
- `grep -q 'registry:' .github/workflows/release.yml` — expected: exit 0.
- `grep -q 'undetermined:' .github/workflows/release.yml` — expected: exit 0.
- `grep -n 'make verify\|Publication preflight\|Cross-compile and stage' .github/workflows/release.yml`
  — expected: the three markers appear in that order, proving the preflight
  sits between the Verification gate and cross-compilation.
- `grep -q 'platforms.json' .github/workflows/release.yml` — expected: exit 0;
  the Release Set is derived, not hardcoded.

## References

- `_prd.md` → User Story 2; Core Features 2; Success Metrics (simulated
  ineligible coordinate).
- `_techspec.md` → Implementation Design: Interfaces; Testing Approach;
  Build Order 2–3.
- ADR-0082.

## Result

Implemented the read-only Publication Preflight in the authorized release
workflow. It materializes one validated six-coordinate Release Set from
`dist/npm/platforms.json` and the launcher manifest, classifies unauthenticated
public packuments, prints every coordinate and state to stdout and the job
summary, reports every blocked coordinate, and exits before cross-compilation
unless all six states are `eligible`. The platform-first publish loop consumes
the same materialized Release Set so publication cannot drift from preflight.

Added five Spec-owned packument fixtures for a used version, a
single-version unpublish retained in the release time map, a full-package
unpublish cooldown, an eligible package, and malformed JSON.

Focused checks:

- `rtk ruby -ropen3 -e '<classifier fixture matrix>'` extracted the workflow's
  `CLASSIFY_JQ` assignment, resolved the four valid fixtures as `used`, `used`,
  `cooldown`, and `eligible`, and observed a non-zero classifier exit for
  `malformed.json`.
- `rtk ruby -ropen3 -rtmpdir -ryaml -e '<stubbed preflight response matrix>'`
  executed the actual YAML-embedded preflight with `curl` stubbed at the
  registry boundary. Six eligible reads exited 0; used, cooldown, malformed,
  absent, HTTP 503, and transport-failure scenarios each exited 1 with their
  expected `registry:` or `undetermined:` diagnostic. The cooldown diagnostic
  included `2026-07-26T17:14:13Z`.
- `rtk ruby -ropen3 -rtmpdir -ryaml -e '<publish-loop Release Set probe>'`
  parsed and executed the actual publish block with `npm` stubbed, observing
  five platform package directories before the launcher.
- `rtk ruby -ropen3 -ryaml -e '<workflow structure comparison>'` compared
  named stages with `HEAD`: `Validate tag`, `Verify gate`,
  `Cross-compile and stage`, and `GitHub Release` remain unchanged, and the
  order is `Verify gate` → `Publication preflight` →
  `Cross-compile and stage`.
- `rtk ruby -e '<CLASSIFY_JQ assignment assertion>' .github/workflows/release.yml`
  observed exactly one single-line, single-quoted classifier assignment.
- `rtk git -c core.fsmonitor=false diff --check` exited 0.

Acceptance evidence:

1. The classifier fixture matrix resolves both a target in `.versions` and a
   target retained only in `.time` as `used`.
2. The same matrix resolves the full-unpublish marker as `cooldown` and the
   unblocked packument as `eligible`.
3. The malformed fixture makes the classifier exit non-zero; the stubbed
   workflow probe maps that failure to `undetermined` and blocks the step.
4. The structural comparison proves the preflight sits after `Verify gate`
   and before `Cross-compile and stage`; the unrelated named stages remain
   byte-equivalent to `HEAD`.
5. The stubbed response matrix observes `registry:` diagnostics for used,
   cooldown, and absent coordinates, and `undetermined:` diagnostics for
   malformed, HTTP 503, and transport-failure reads.
6. `rtk git -c core.fsmonitor=false status --short` lists only
   `.github/workflows/release.yml`, this Spec's `fixtures/` directory, and this
   task file. The task file's pre-existing `status: in_progress` change remains
   Daemon-owned.

The commands under `## Verification` were not run; the Daemon owns that gate
and Task settlement.
