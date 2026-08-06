---
task: task_02
spec: 0073-skill-versions-decoupled-from-the-binary
status: pending
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
