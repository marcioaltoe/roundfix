---
task: task_04
spec: 0103-a-suite-that-leaks-nothing
status: completed # pending | in_progress | completed | failed — only implement-task changes this
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

## Result

Installed `suiteguard.Main` in every internal package whose test source starts
an OS process: `agent`, `baseline`, `cli`, `daemon`, `gittest`, `preflight`,
`spec`, `specaudit`, `speccheck`, `store`, `suiteguard`, and `worktree`. Existing
`TestMain` helper-process branches and package cleanup remain in their original
order; only the normal package-test path is wrapped by the repository boundary.

Added the `repocontract`-tagged
`TestEverySpawningPackageInstallsTheSuiteGuard`. It parses every internal
`*_test.go` file, recognizes subprocess calls through their imported Go package,
compares the discovered package set with an explicit guarded-package inventory,
and verifies that each package calls `suiteguard.Main`. A synthetic package that
uses `os/exec` without an installation proves the contract reports
`internal/uninstalled spawns a process but is not enumerated as guarded`.

Focused checks:

- Pre-change source inspection found no
  `TestEverySpawningPackageInstallsTheSuiteGuard` and only the suiteguard
  fixture's conditional `suiteguard.Main`, establishing the missing package
  installations and repository contract.
- `GOCACHE=<worktree>/.gocache go test -count=1 -tags repocontract -run
  '^TestEverySpawningPackageInstallsTheSuiteGuard$' ./internal/suiteguard`
  passed in 1.316s, including the synthetic missing-installation subtest and the
  live repository inventory.
- One selected existing test per guarded package passed in a single focused
  sweep; all 12 package-level guards executed and returned success.
- Full package tests passed with the guard active: `internal/agent` (9.499s),
  `internal/baseline` (56.910s), `internal/cli` (48.072s), and
  `internal/worktree` (3.475s) passed individually; one additional focused
  command passed `internal/daemon`, `internal/gittest`, `internal/preflight`,
  `internal/spec`, `internal/specaudit`, `internal/speccheck`, `internal/store`,
  and `internal/suiteguard`.
- `git diff --check` passed. `git diff -- .gitignore
  internal/suiteguard/suiteguard.go` produced no output: this slice added no
  exclusion and did not weaken the guard.

Acceptance evidence:

- Every spawning package installs the guard: the live AST contract discovered
  the 12-package inventory above and found a `suiteguard.Main` call in each.
- Missing packages are named: the contract's negative fixture produced the
  exact unenumerated-package diagnostic before the repository inventory ran.
- Existing guarded tests remain unchanged and passed in the focused package
  runs. The Daemon still owns the declared full-suite Verification.
- No violation was silenced: every guarded package run completed without a
  repository-boundary diagnostic, and no ignore rule or guard exclusion changed.

Verification Feedback repair:

- Daemon attempt 1 reached and passed the named contract in
  `internal/suiteguard`, but the following non-vacuity check found Go's
  package-local `no tests to run` warning from the other 25 test-bearing
  internal packages selected by the same `-run` expression.
- The tagged contract is now package-scoped. A shared read-only analyzer under
  `internal/suiteguardcontract` checks the calling package's test sources and
  requires both enumeration and `suiteguard.Main` when it finds a process
  spawn. All 25 other test-bearing packages expose the same named contract;
  `internal/suiteguard` retains the full repository inventory and negative
  fixture.
- Focused pre-repair reproduction across `internal/app` and
  `internal/suiteguard` failed because `internal/app` emitted `no tests to run`
  while the suiteguard contract passed.
- The same two-package focused check passed after the repair with two named
  contract passes and no vacuous-run warning.
- Two disjoint focused package subsets covered all 26 test-bearing internal
  packages. Each subset reported 13 top-level named contract passes, exited 0,
  and contained no `no tests to run` warning.

The Daemon-owned Verification commands were not run in this turn.
