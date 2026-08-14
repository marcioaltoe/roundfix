---
task: task_04
spec: 0103-a-suite-that-leaks-nothing
status: pending # pending | in_progress | completed | failed — only implement-task changes this
type: test
complexity: medium
---

# Task 04: Install the guard where the suite spawns

## Overview

A guard nobody installs proves nothing. This slice puts it in the packages whose
tests spawn processes or apply Baselines, and adds the repository-contract test
that enumerates which packages carry it — so an unguarded package is a visible
fact rather than an assumption.

## Requirements

1. MUST install the guard from `TestMain` in every package whose tests spawn a
   process or write outside their own temporary directory.
2. MUST add a repository-contract test that enumerates the guarded packages and
   fails when a package that spawns is not among them.
3. MUST leave every guarded package's existing tests passing unchanged.
4. MUST NOT silence a violation the guard finds by excluding the path; a real
   violation is repaired in the test that causes it.

## Subtasks

- [ ] Install the guard in the spawning packages.
- [ ] Add the enumeration contract test.
- [ ] Repair whatever the guard finds.

## Acceptance Criteria

- [ ] Every package whose tests spawn a process installs the guard.
- [ ] The contract test names any spawning package that does not.
- [ ] The full suite passes with the guard installed.
- [ ] No violation was resolved by adding an exclusion.

## Verification

- `go test -count=1 -tags repocontract -run 'TestEverySpawningPackageInstallsTheSuiteGuard' ./internal/... -v > /tmp/0103-t04.log 2>&1; s=$?; grep -q '^--- PASS: TestEverySpawningPackageInstallsTheSuiteGuard' /tmp/0103-t04.log || { cat /tmp/0103-t04.log; exit 1; }; exit $s` — expected: exits 0 and the log names the passing contract test; fails today, where it does not exist.
- `! grep -qi 'no tests to run' /tmp/0103-t04.log` — expected: exits 0, refusing a vacuous run.
- `n=$(grep -rl 'suiteguard.Main' internal/ | wc -l | tr -d ' '); test "$n" -ge 4 || { echo "expected the guard installed in at least the four spawning packages, found $n"; grep -rl 'suiteguard.Main' internal/; exit 1; }` — expected: exits 0, proving the guard reached the packages rather than only the contract test. Fails today, where the count is zero.
- `for p in internal/agent internal/cli internal/baseline; do grep -rq 'suiteguard.Main' "$p" || { echo "FAIL: $p does not install the guard"; exit 1; }; done; go build -buildvcs=false ./... && go test -count=1 ./internal/agent ./internal/cli ./internal/baseline > /tmp/0103-t04-suite.log 2>&1; s=$?; grep -q '^ok' /tmp/0103-t04-suite.log || { cat /tmp/0103-t04-suite.log; exit 1; }; grep -q 'wrote inside the repository root' /tmp/0103-t04-suite.log && { echo 'the guard found an unrepaired violation:'; grep -B2 -A2 'wrote inside the repository root' /tmp/0103-t04-suite.log; exit 1; }; exit $s` — expected: exits 0, proving the three packages install the guard, that they pass with it, and that nothing was left violating. The installation check leads, because a green suite on an unguarded tree proves nothing.

## Context

- interface: `internal/cli/implement_test.go`

## References

`_techspec.md` → Build Order 4; Risks & Considerations, the cross-package
sighting. `_prd.md` → Core Feature 1; Goal 1; Open Questions. ADR-0126.
