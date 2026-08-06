---
task: task_02
spec: 0073-skill-versions-decoupled-from-the-binary
status: completed
type: backend
complexity: high
---

# Task 02: Compare a version instead of matching bytes

## Overview

A byte-exact pin answers "is this the same content?". The question Roundfix
needs answered is "is this content new enough to work with me?". Only the second
survives the skill evolving on its own schedule.

This slice adds the comparison and the minimum it compares against, and it
establishes the three states that every later surface reports.

## Requirements

1. MUST add a readiness comparison: a declared version at or above the declared
   minimum satisfies.
2. MUST carry Roundfix's minimum per owned skill in the setup snapshot, in
   place of the content pin, so a profile carries a declaration rather than a
   digest.
3. MUST keep three states distinct and never collapse any two: satisfies,
   below the minimum, and unversioned-or-unresolvable.
4. MUST report an unreachable source as unresolvable, never as a missing skill.
   They are different facts and the operator acts differently on each.
5. MUST make a version above the minimum satisfy with no change in Roundfix.
6. MUST NOT derive the minimum from any skill's declared version, and MUST NOT
   invent a version for a skill that declares none.
7. MUST NOT yet change what gates on content — task_03 owns removing the
   digest from the compatibility path.
8. MUST NOT assert a recorded version in any test. Assert the comparison's
   outcome, which holds at any version.

## Subtasks

- [ ] Add the readiness comparison and its three states.
- [ ] Carry the minimum per owned skill in the setup snapshot.
- [ ] Table-test each state, including one minor above the minimum.

## Acceptance Criteria

- [ ] A declared version equal to the minimum satisfies.
- [ ] A declared version one minor above the minimum satisfies, with no
      Roundfix change.
- [ ] A declared version below the minimum reports `below`.
- [ ] A skill declaring no version reports `unversioned`.
- [ ] An unreachable source reports unresolvable, distinctly from missing.
- [ ] The minimum is never derived from a declared version, asserted.
- [ ] No test asserts a recorded version literal.

## Context

- interface: `internal/baseline/catalog_validate.go`
- interface: `internal/baseline/assets_sync.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `output="$(go test ./internal/baseline -count=1 -run 'Readiness|Version|Minimum' -v 2>&1)"; status=$?; printf '%s\n' "$output"; [ "$status" -eq 0 ] || exit "$status"; printf '%s\n' "$output" | grep -q -- "--- PASS"`
  — expected: exit 0; the comparison tests ran and passed.
- `go test ./internal/baseline ./skills -count=1` — expected: exit 0.
- `go test -parallel 16 ./...` — expected: exit 0.

## References

- `_prd.md` → Core Features 2, 3 and 5; Success Metrics 3 and 6.
- `_techspec.md` → Interfaces; Version identity; Build Order 2.
- ADR-0085.

## Result

### Implementation

- `Readiness` compares strict three-part versions and returns the distinct
  `satisfies`, `below`, or `unversioned` state. An unreachable or malformed
  source returns `unversioned` with the matchable
  `ErrSkillVersionUnresolvable`; absence remains outside the comparison and is
  never reported as missing.
- Every Roundfix-owned setup-snapshot entry now declares
  `minimumVersion`. Catalog validation requires a valid minimum only for those
  repository-owned entries, while external skills remain outside the version
  contract.
- Both setup-snapshot producers preserve the Roundfix-declared minimum: Baseline
  Assets Sync carries it forward from the existing snapshot, and the
  authorial-skill digest regenerator round-trips it across skill content edits.
  Existing content digests remain present and enforced for task_03 to remove
  from the compatibility path.
- The sanctioned ADR-0081 regeneration updated catalog, parity,
  plan-characterization, and diagnostic artifacts after the snapshot contract
  changed.

### Focused checks

- Red signal: the first focused readiness test failed to compile with the new
  symbols undefined, proving the comparison did not exist before this slice.
- `GOCACHE=<worktree>/.gocache rtk go test ./internal/baseline -run
  '^(TestReadinessComparesDeclaredVersionToMinimum|TestCatalogRejectsMissingOwnedSkillMinimum|TestBaselineAssetsSyncRefreshProducesCanonicalTreeAndIsIdempotent|TestAssetsSyncCompatibilityMatchesMaintainedPythonContract)$'
  -count=1`: exit 0, 10 tests passed.
- `GOCACHE=<worktree>/.gocache rtk go test ./skills -run
  '^TestAuthorialSkillSyncUpdateModeRoundTrip$' -count=1`: exit 0. Its red
  predecessor dropped the minimum on regeneration; the repaired producer now
  preserves an independently supplied fixture minimum across two skill-content
  digests.
- `rtk make baseline-digests`: final exit 0 with `ok: true` and
  `changed: true`. The first attempt exposed the missing producer field and
  stopped at strict catalog validation; the rerun passed every regeneration
  step after that root cause was repaired.
- Snapshot audit: `go-cli` and `rust-cli` each contain 13 owned entries, and
  `typescript-bun` contains 14; every owned entry has a valid three-part
  minimum and the existing content pin, while zero external entries declare a
  minimum.
- `rtk git -c core.fsmonitor=false diff --check`: exit 0. A new-Go-line audit
  found no assertion containing the repository's recorded skill version.

### Acceptance evidence

- Equal and one-minor-above declarations both produce `satisfies` in the
  readiness table without reading a repository-recorded version.
- The below-minimum row produces `below`; an empty declaration with a reachable
  source produces `unversioned`.
- The unreachable-source row produces `unversioned`, matches
  `ErrSkillVersionUnresolvable`, and explicitly rejects a `missing` report.
- Catalog validation rejects an owned snapshot entry with no declared minimum.
  The two regeneration tests preserve the snapshot's independently supplied
  minimum instead of reading or deriving it from skill content.
- No Go test added by this slice contains the recorded skill-version literal;
  tests assert readiness outcomes from local comparison inputs.
- The Task's declared `## Verification` commands were not run; the Daemon owns
  those commands and terminal settlement.
