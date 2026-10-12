---
task: task_02
spec: 0259-proven-signals-and-a-shared-build-cache
status: pending
type: test
complexity: high
---

# Task 02: No test helper outlives its test binary

## Overview

The helpers that wait for a signal watch nothing, one of them ignores
`SIGTERM`, and the detach children watch their owner by PID alone, so six
`store.test` helpers lived for 45 hours (`_techspec.md` → Current behavior).
This Task gives every waiting helper a Helper Lifeline, makes the suite guard
report and kill Helper Residue after each package, and returns the test that
saw F6 to running in parallel. It is verifiable on its own: killing a test
binary with `SIGKILL` leaves none of its helpers alive, and a package that
leaves a marked process behind fails.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-09,
   [a test's child process is killed by a signal nobody identified](references/2026-10-09-a-test-child-killed-by-an-unknown-signal.md),
   whose direction asks that helper processes be joined or killed so none
   outlives its test.
2. MUST give the `graceful`, `ignore`, `tree-root` and `tree-grandchild` modes
   of `TestOwnerProcessHelper`, and `TestCLIForceStopOwnerProcessHelper`, the
   standard-input lifeline of `_techspec.md` → API Contract 5: the test keeps
   the write end of the pipe until its cleanup, and `tree-middle` passes its
   standard input to the grandchild. The modes' existing signal behavior MUST
   NOT change; `ignore` still survives `SIGTERM` while its lifeline is open.
3. MUST make the detach test children compare their owner's PID and start
   identity, carried in `ROUNDFIX_DETACH_TEST_OWNER_IDENTITY`, instead of the
   PID-only liveness probe.
4. MUST add `TestOwnerProcessHelpersEndWithTheirTestBinary` and
   `TestCLITestHelpersEndWithTheirTestBinary`: each starts an inner test binary
   that starts every helper family of its package and records each PID and
   start identity, kills the inner binary with `SIGKILL`, and passes only when
   every recorded identity is gone within the `testwait` bound. They MUST NOT
   signal a helper themselves, except a Proven Signal in cleanup after the
   assertion failed.
5. MUST add the Helper Residue check to `suiteguard.Main` following API
   Contract 6: mark, scan after the tests with a 10 s grace, one
   `suiteguard: helper residue: pid <pid>: <command>` line per process,
   `killResidue` that rereads the mark before `SIGKILL`, exit 1, and the
   `unavailable` line, never a failure, when the process table cannot be read
   or the platform is unsupported. It MUST NOT import `internal/store`.
6. MUST add `TestGuardReportsAndKillsHelperResidue`,
   `TestGuardPassesWhenEveryHelperEnded` and `TestGuardMarksItsTestBinary`
   with the suite guard's existing fixture-binary pattern; the first proves
   the residue line, exit 1 and that the left process is dead afterwards.
7. MUST add `killResidue` to the allowlist of
   `TestEverySignalGoesToAProvenProcess`, and nothing else.
8. MUST make `TestDeliveryStepFallsBackWhenTheItemBinaryWouldMigrate` call
   `t.Parallel()` as its first statement, removing its `Sequential:` comment,
   and lower the `internal/cli` ceiling of the Parallel Test Package rule from
   11 to 10. The package MUST pass once with `-race` for the delivery step
   tests (ADR-0259).
9. MUST fix, not suppress, any Helper Residue the new check finds in a
   package, and record each one in the Result. The grace and the mark MUST
   NOT be widened to make a package pass.
10. MUST add to `docs/agents/specific-repository.md`, outside every setup
    marker, one rule: a test signals a process only through the `*os.Process`
    of an unreaped child it started or a Proven Signal, and a helper that
    waits has a Helper Lifeline; the Signal rule and the suite guard enforce
    both.
11. MUST add **Helper Lifeline** and **Helper Residue** to `CONTEXT.md`
    through `domain-modeling`, linked to ADR-0266, and keep Helper Residue
    distinct from **Process Residue**.
12. MUST make `gittest.Harden` and `gittest.InitRepo` write
    `core.fsmonitor = false` and `core.untrackedCache = false` into the
    repository's own config, as `_techspec.md` → API Contract 6 states, and add
    `TestHardenedRepositoryNeverStartsAFilesystemMonitor` to
    `internal/gittest`: with a global config file that sets
    `core.fsmonitor=true` (through `GIT_CONFIG_GLOBAL` pointing at a temporary
    file), a hardened repository reports `false` for both keys, and a
    `git add --all` followed by a commit of a file rewritten with the same size
    records the change. This answers the operator's 2026-10-12 observation that
    fixture seed commits failed with "nothing to commit" inside the Daemon.

## Subtasks

- [ ] Give the store and CLI helpers their standard-input lifeline.
- [ ] Make the detach children watch their owner's identity.
- [ ] Add the two death tests.
- [ ] Add the Helper Residue check and its tests to the suite guard.
- [ ] Return the F6 test to parallel and lower the ceiling.
- [ ] Add the repository rule and the glossary terms.
- [ ] Keep every hardened fixture repository free of a filesystem monitor.

## Acceptance Criteria

- [ ] A test binary killed with `SIGKILL` leaves no helper of either package alive.
- [ ] A package that leaves a marked process fails with the residue line and the process is dead.
- [ ] A package whose helpers all ended passes and prints no residue line.
- [ ] The F6 test runs in parallel and the `internal/cli` ceiling is 10.
- [ ] A hardened fixture repository never starts a filesystem monitor, even with `core.fsmonitor=true` in the global config.

## Context

- instruction: `docs/adr/0266-a-signal-goes-only-to-a-proven-process-and-worktrees-share-one-build-cache.md`
- instruction: `docs/adr/0213-a-test-fixture-process-ends-with-the-test-binary-that-started-it.md`
- instruction: `docs/adr/0126-the-suite-proves-its-own-isolation-per-package.md`
- instruction: `docs/adr/0259-a-test-runs-in-parallel-unless-it-names-why-it-cannot.md`
- instruction: `internal/worktree/staging_test.go`
- creates: `internal/suiteguard/residue.go`
- creates: `internal/suiteguard/residue_darwin.go`
- creates: `internal/suiteguard/residue_linux.go`
- creates: `internal/suiteguard/residue_other.go`
- creates: `internal/store/helper_lifeline_unix_test.go`
- creates: `internal/cli/helper_lifeline_unix_test.go`
- interface: `internal/suiteguard/suiteguard.go`
- interface: `internal/suiteguard/suiteguard_test.go`
- interface: `internal/store/process_unix_test.go`
- interface: `internal/cli/orphan_unix_test.go`
- interface: `internal/cli/detach_test.go`
- interface: `internal/cli/deliver_item_binary_test.go`
- interface: `internal/testfixture/parallel_tests_test.go`
- creates: `internal/testfixture/signal_rule_test.go`
- interface: `docs/agents/specific-repository.md`
- interface: `CONTEXT.md`
- interface: `internal/gittest/gittest.go`
- interface: `internal/gittest/gittest_test.go`

## Verification

- `for term in '**Helper Lifeline**:' '**Helper Residue**:'; do tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "$term" || { printf 'CONTEXT.md does not define %s\n' "$term" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < docs/agents/specific-repository.md | grep -qF -- 'Helper Lifeline' || { printf 'the repository rule does not name the Helper Lifeline\n' >&2; exit 1; }; out="$(go test -count=1 -v -run '^(TestGuardReportsAndKillsHelperResidue|TestGuardPassesWhenEveryHelperEnded|TestGuardMarksItsTestBinary)$' ./internal/suiteguard 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestGuardReportsAndKillsHelperResidue TestGuardPassesWhenEveryHelperEnded TestGuardMarksItsTestBinary; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing passing test %s\n' "$name" >&2; exit 1; }; done; GOOS=windows go build -buildvcs=false ./internal/suiteguard || exit 1; GOOS=linux go build -buildvcs=false ./internal/suiteguard || exit 1` — expected: exit 0. Before this Task the terms, the rule and the three tests are absent, so the command fails.
- `out="$(go test -count=1 -v -run '^TestOwnerProcessHelpersEndWithTheirTestBinary$' ./internal/store 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestOwnerProcessHelpersEndWithTheirTestBinary ' || { printf 'the store death test did not pass\n' >&2; exit 1; }; out="$(go test -count=1 -v -run '^TestCLITestHelpersEndWithTheirTestBinary$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestCLITestHelpersEndWithTheirTestBinary ' || { printf 'the CLI death test did not pass\n' >&2; exit 1; }` — expected: exit 0. Before this Task neither death test exists, so no PASS line is printed.
- `grep -A1 '^func TestDeliveryStepFallsBackWhenTheItemBinaryWouldMigrate(' internal/cli/deliver_item_binary_test.go | grep -q 't.Parallel()' || { printf 'the F6 test does not call t.Parallel first\n' >&2; exit 1; }; grep -Eq '"internal/cli": +10,' internal/testfixture/parallel_tests_test.go || { printf 'the internal/cli ceiling is not 10\n' >&2; exit 1; }; go test -count=1 -run '^(TestEveryTestInAParallelTestPackageRunsInParallel|TestEverySignalGoesToAProvenProcess)$' ./internal/testfixture || exit 1; go test -count=1 -race -run '^TestDeliveryStep' ./internal/cli || exit 1` — expected: exit 0. Before this Task the F6 test is sequential and the ceiling is 11, so the first check fails.
- `out="$(go test -count=1 -v -run '^TestHardenedRepositoryNeverStartsAFilesystemMonitor$' ./internal/gittest 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestHardenedRepositoryNeverStartsAFilesystemMonitor ' || { printf 'the filesystem monitor test did not pass\n' >&2; exit 1; }; grep -q 'fsmonitor = false' internal/gittest/gittest.go || { printf 'gittest does not disable the filesystem monitor\n' >&2; exit 1; }` — expected: exit 0. Before this Task the test is absent and `gittest` writes no `core.fsmonitor`, so the command fails.

## References

- `_prd.md` → Goals; User Story 1; User Story 3; Core Feature 3; Core Feature 4; Core Feature 5; Core Feature 8; Success Metric 3
- `_techspec.md` → Current behavior; API Contract 5; API Contract 6; Invariants 3-4; Testing Approach; Risks & Considerations; Build Order 2
- ADR-0266; ADR-0213; ADR-0126; ADR-0259
