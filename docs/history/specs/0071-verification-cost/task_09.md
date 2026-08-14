---
task: task_09
spec: 0071-verification-cost
status: completed
type: infra
complexity: high
---

# Task 09: Split the gate into a fast local tier and a full CI tier

## Overview

`make verify` is the only gate that runs tests, and no workflow runs them on a
Pull Request — only the release workflow does, on a tag, which is far too late
to catch anything. That is the hole this Task closes.

**Measurement moved where the split falls.** The Task was written to strip the
heavy end-to-end journeys out of the local gate. After task_08, a second
consecutive `make verify` costs **5.2s** — Go's test result cache re-runs only
the packages whose compiled output actually changed, and the prefactors in
tasks 02 and 04 removed the `t.Setenv` and `t.Chdir` calls that made caching
untrustworthy. An implementation loop editing `internal/cli` pays 47.9s; one
editing any of the twenty light packages pays under 5s; one editing nothing
pays 5.2s.

So the local gate needs no exclusions. Stripping tests from a gate that already
costs seconds would remove coverage from the loop and buy nothing measurable —
and this Spec must evolve the gate, not weaken it. The tier boundary is the
cache, not a test list: local trusts it, CI refuses it.

## Requirements

1. MUST add a Pull Request workflow that runs the complete verification gate,
   so the suite has a tier that judges the whole tree before merge.
2. MUST make the CI tier force a fresh run, so its verdict never rests on a
   cache entry from another commit.
3. MUST leave the local gate running every test it runs today, and MUST let it
   use the test result cache so it re-runs only what changed.
4. MUST state in the workflow why the two tiers differ, so the next reader does
   not "fix" the local gate by disabling its cache.
5. MUST NOT delete, skip unconditionally, or weaken any test; both tiers run
   everything the suite runs today.
6. MUST measure and record the local gate's wall clock against the recorded
   baseline, for both a changed tree and an unchanged one.
7. MUST change only `Makefile` and `.github/workflows/` among protected
   tooling, per this Spec's Tooling authority row and the 2026-08-03 addendum
   to the queued-Spec tooling authorization.

## Subtasks

- [x] Add the Pull Request workflow running the full gate.
- [x] Force a fresh run in the CI tier and record why in the workflow.
- [x] Confirm the local gate keeps its full test set and uses the cache.
- [x] Measure the local gate on an unchanged tree and on a changed package.
- [x] Record the corrected premise in the Spec's Task manifest.

## Acceptance Criteria

- [x] A Pull Request runs a workflow that executes the complete suite.
      `.github/workflows/ci-verify.yml` runs `make verify` and `make test-race`
      on every Pull Request to `main` and on push to `main`.
- [x] The CI tier runs with the test result cache disabled, via
      `GOFLAGS: -count=1`. Confirmed that `build` and `vet` ignore the flag, so
      only `go test` is affected.
- [x] A second consecutive local gate completes in **5.1s**, against **56.5s**
      on a changed tree and the recorded **142.2s** baseline.
- [x] The local gate still runs every test in the suite — no tier exclusion,
      no build tag, no skip.
- [x] The coverage record from task 01 is unchanged.
- [x] Three consecutive fresh full runs pass at 89.3s / 87.9s / 88.9s.

## Defect surfaced

Forcing a fresh run exposed a flake the cached gate had been hiding:
`TestRunResolveDetachedChildReportsProfileProofFailure` failed once under
`go test ./...` with `TempDir RemoveAll cleanup: directory not empty`. It is
the only test that drives the detached-child branch, and that branch mutates
process-global state — the test hands production code a raw file descriptor
number, and `newDetachChildFromEnv` calls `os.Unsetenv` twice. Go refuses
`t.Parallel()` next to `t.Setenv` for exactly this reason but cannot see the
mutation when it lives in production code, so the hazard ran unguarded.

The test is now sequential with that reason stated inline, per this Spec's rule
that a test legitimately holding process-global state keeps it and says why.
Cost: about one second of a 47-second package.

## Context

- instruction: `docs/workflow/authorizations/2026-08-02-queued-spec-tooling.md`
- interface: `Makefile`
- interface: `.github/workflows/ci-conventions.yml`

## Verification

- `make verify` — expected: exit 0.
- `make verify` — expected: exit 0 on a second consecutive run, and visibly
  faster than the first because unchanged packages come from the cache.
- `go test ./... -count=1` — expected: exit 0; a fresh run still passes
  everything.
- `go test ./internal/spec -run '^TestCoverageEquivalence$' -count=1 -v | grep -q -- "--- PASS: TestCoverageEquivalence"`
  — expected: exit 0; no test disappeared.
- `grep -rqE "go test|make verify" .github/workflows/` — expected: exit 0; a
  workflow runs the suite.
- `grep -rq "pull_request" .github/workflows/ci-verify.yml` — expected: exit 0;
  it runs on Pull Requests.
- `grep -rq -- "-count=1" .github/workflows/ci-verify.yml` — expected: exit 0;
  the CI tier forces a fresh run.
- `grep -A1 '^test:' Makefile | grep -q -- '-count=1'` — expected: exit 1; the
  local test target does not disable the cache.

## References

- `_prd.md` → Goals (verification proportional to what changed); Non-Goals (no
  test deleted, skipped, or weakened).
- `baseline/2026-08-03-after.md` → the measurement behind the tiering.
