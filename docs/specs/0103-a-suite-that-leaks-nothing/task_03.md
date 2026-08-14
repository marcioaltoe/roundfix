---
task: task_03
spec: 0103-a-suite-that-leaks-nothing
status: pending # pending | in_progress | completed | failed — only implement-task changes this
type: test
complexity: medium
---

# Task 03: Give the suite a boundary it can assert

## Overview

A test that applies a Baseline against the live tree deleted tracked files while
another test was copying them. Nothing in the suite can currently say that
happened, so the damage surfaces as a stranger's failure somewhere else. This
slice builds the guard: fingerprint the repository root around a package's tests
and name every path that moved.

## Requirements

1. MUST provide a helper a package installs from `TestMain` that fingerprints the
   repository root before and after the package's tests.
2. MUST fail the package naming every path that was created, modified, or removed,
   and whether it was created, modified, or removed.
3. MUST ignore paths the repository's own ignore rules exclude, so ordinary build
   output does not read as a violation.
4. MUST cost little enough to run in every guarded package, and MUST record what
   it measured.
5. MUST NOT depend on the Git index, which is stale whenever a Task has moved
   files without committing them.

## Subtasks

- [ ] Build the fingerprint and the comparison.
- [ ] Report violations with their change kind.
- [ ] Cover a violating package and a clean one.
- [ ] Measure the guard's own cost.

## Acceptance Criteria

- [ ] A test that writes inside the repository root fails its package, and the
      failure names the path and the change kind.
- [ ] A test that writes only inside its temporary directory passes.
- [ ] A path excluded by the repository's ignore rules is not reported.
- [ ] The guard reads no Git index.
- [ ] The guard's measured cost is recorded in its own test output.

## Verification

- `go test -count=1 ./internal/suiteguard -v > /tmp/0103-t03.log 2>&1; s=$?; grep -q '^--- PASS: TestGuardNamesTheViolatingPath' /tmp/0103-t03.log || { cat /tmp/0103-t03.log; exit 1; }; grep -q '^--- PASS: TestGuardPassesOnAnIsolatedWrite' /tmp/0103-t03.log || { cat /tmp/0103-t03.log; exit 1; }; exit $s` — expected: exits 0 and the log names both passing cases; fails today, where the package does not exist.
- `test -s /tmp/0103-t03.log || { echo 'the guard suite produced no output'; exit 1; }; grep -qi 'no tests to run' /tmp/0103-t03.log && { echo 'the guard suite selected no cases'; exit 1; }; test -d internal/suiteguard || { echo 'the guard package does not exist'; exit 1; }` — expected: exits 0, refusing both a vacuous run and an absent package. Fails today, where the package does not exist and the log records why.
- `! grep -rn 'ls-files\|git.*index' internal/suiteguard/ > /tmp/0103-t03-git.txt 2>&1; test ! -s /tmp/0103-t03-git.txt || { echo 'the guard reads Git state it was told not to:'; cat /tmp/0103-t03-git.txt; exit 1; }; test -d internal/suiteguard || { echo 'the guard package does not exist'; exit 1; }` — expected: exits 0, proving the package exists and reads no index. Fails today on the second clause, which is what stops an absent package from passing the first.
- `grep -qE 'guard cost|fingerprint took|measured' /tmp/0103-t03.log || { echo 'the guard did not record what it measured'; cat /tmp/0103-t03.log; exit 1; }` — expected: exits 0, proving the cost is observable rather than asserted in prose.

## Context

- interface: `internal/gittest/gittest.go`

## References

`_techspec.md` → Build Order 3; System Architecture, the suite guard;
Implementation Design, Interfaces. `_prd.md` → Core Feature 1; Goal 1; User
Story 1; Open Questions. ADR-0126.
