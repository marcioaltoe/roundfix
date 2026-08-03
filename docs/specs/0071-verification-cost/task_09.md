---
task: task_09
spec: 0071-verification-cost
status: pending
type: infra
complexity: high
---

# Task 09: Split the gate into a fast local tier and a full CI tier

## Overview

`make verify` is the only gate that runs tests, and no workflow runs them on a
Pull Request — only the release workflow does, on a tag, which is far too late
to catch anything. Every Task therefore pays for the whole suite, including
end-to-end journeys that spawn the real binary, and implementation time absorbs
that cost repeatedly.

This Task tiers the gate: the local gate runs what must be true at every step,
and a Pull Request workflow runs everything. Moving the heavy tier out of the
local gate is only safe because this Task also creates the place it lands.

## Requirements

1. MUST add a Pull Request workflow that runs the complete verification gate,
   including every test the local tier excludes, so nothing loses coverage by
   moving.
2. MUST make the local gate exclude the heavy end-to-end journeys while keeping
   build, format, lint, skill checks, and the fast tests.
3. MUST make the exclusion explicit and self-describing: a reader can tell
   which tier a test belongs to from the test itself, and the local gate
   reports what it skipped.
4. MUST keep every excluded test runnable locally by an obvious command, so a
   maintainer working on that surface is never blocked from exercising it.
5. MUST NOT delete, skip unconditionally, or weaken any test; the full tier
   runs everything the suite runs today.
6. MUST measure and record the local gate's new wall clock against the
   recorded baseline.
7. MUST change only `Makefile` and `.github/workflows/` among protected
   tooling, per this Spec's Tooling authority row and the maintainer's
   2026-08-03 direction to separate the tiers.

## Subtasks

- [ ] Add the Pull Request workflow running the full gate.
- [ ] Mark the heavy end-to-end journeys as belonging to the full tier.
- [ ] Make the local gate run the fast tier and report what it excluded.
- [ ] Document the command that runs an excluded test locally.
- [ ] Measure and record the local gate's new wall clock.

## Acceptance Criteria

- [ ] A Pull Request runs a workflow that executes the complete suite,
      including every heavy journey.
- [ ] The local gate completes materially faster than its recorded 142.2s and
      reports which tests it excluded.
- [ ] Every excluded test still runs, and passes, under the full tier.
- [ ] A documented single command runs an excluded test locally.
- [ ] The coverage record from task 01 is unchanged: nothing was deleted or
      unconditionally skipped.
- [ ] `git status --porcelain` shows no path outside `Makefile`,
      `.github/workflows/`, `internal/`, and this task file.

## Context

- instruction: `docs/workflow/authorizations/2026-08-02-queued-spec-tooling.md`
- interface: `Makefile`
- interface: `.github/workflows/ci-conventions.yml`

## Verification

- `make verify` — expected: exit 0, and materially faster than 142.2s.
- `make verify` — expected: exit 0 on a second consecutive run.
- `go test ./... -count=1` — expected: exit 0; the full tier still passes
  everything.
- `go test ./internal/spec -run '^TestCoverageEquivalence$' -count=1 -v | grep -q -- "--- PASS: TestCoverageEquivalence"`
  — expected: exit 0; no test disappeared.
- `grep -rqE "go test|make verify" .github/workflows/` — expected: exit 0; a
  workflow runs the suite.
- `grep -rq "pull_request" .github/workflows/` — expected: exit 0; it runs on
  Pull Requests.

## References

- `_prd.md` → Goals (verification proportional to what changed); Non-Goals (no
  test deleted, skipped, or weakened).
- `baseline/2026-08-03-after.md` → the measurement behind the tiering.
