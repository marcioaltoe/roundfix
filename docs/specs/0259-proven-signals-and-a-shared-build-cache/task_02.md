---
task: task_02
spec: 0259-proven-signals-and-a-shared-build-cache
status: completed
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


## Result

### Implementation

- Added stdin Helper Lifelines to the store's graceful, ignore, tree-root and
  tree-grandchild helpers and the CLI Force Stop helper. Only the spawning test
  binary holds each pipe writer; tree-root and tree-middle pass the reader on
  as stdin. Existing SIGTERM behavior remains, including ignore mode's SIGUSR1
  liveness acknowledgement while its lifeline is open.
- Store fixture entry points retain the spawning package's ownership mark
  rather than installing a second package guard. This keeps the deliberately
  outliving grandchild inside the original package's residue boundary.
- Detach fixtures now carry `ROUNDFIX_DETACH_TEST_OWNER_IDENTITY` and compare
  both PID and start identity on their owner watch.
- Added the two inner-binary death tests. They record four store helper
  identities and three CLI helper identities, prove those helpers are alive,
  SIGKILL only the inner binary, and observe every recorded identity ending
  within `testwait`. Helper rescue signals are limited to Proven Signals in
  cleanup after a failed assertion.
- Suiteguard now assigns the `<pid>.<16 random hex digits>` environment mark,
  scans same-user processes on Darwin and Linux, and uses the unchanged ten
  second grace. It reports every remaining marked process, rereads its exact
  environment assignment before SIGKILL, and returns 1 for residue.
  Unreadable process tables and unsupported platforms report unavailable
  without changing the package's test exit code. Production suiteguard imports
  no store package.
- Added guarded fixtures for residue, clean helper exit and mark shape, plus
  safety coverage for exact environment matching, refusal of a different mark,
  matching-mark SIGKILL, Darwin argument/environment separation and reaped-PID
  absence. The Signal rule's only added exception is `killResidue`.
- Returned the F6 delivery test to first-statement `t.Parallel()` and lowered
  the CLI Sequential Test ceiling from 11 to 10.
- Added the repository signal/lifeline rule outside setup markers and defined
  Helper Lifeline and Helper Residue in CONTEXT.md through domain-modeling.
  Both link to ADR-0266; Helper Residue remains distinct from Process Residue.
- Harden and InitRepo persist `core.fsmonitor = false` and
  `core.untrackedCache = false`, including bare repositories and Git-resolved
  config paths. The regression uses a temporary global config enabling both
  features, reads settings without ConfigArgs masking them, commits a same-size
  rewrite, and checks Git Trace2 for filesystem-monitor daemon launches.

### Focused implementation checks

All commands below used `GOCACHE=/private/tmp/roundfix-task02-go-build` and
`rtk proxy`; process-table access ran outside the sandbox because the assigned
writable worktree root contains a symlink that prevents sandbox startup.
These are implementation checks, not execution of the authored Verification.

- `go test -count=1 -shuffle=on ./internal/store ./internal/suiteguard ./internal/gittest ./internal/testfixture`
  — exit 0 on the final stable source tree; store 13.622 s, suiteguard
  12.614 s, gittest 1.979 s, testfixture 2.321 s.
- `go test -count=1 -shuffle=on ./internal/cli`
  — exit 0 on the final stable source tree; 149.450 s; no unexpected Helper
  Residue or repository-boundary violation.
- `go test -count=1 -race -shuffle=on -run 'TestDeliveryStep|TestCLITestHelpers|TestCLIForceStopOwner|TestRunDetachedCommand|TestDetachedChildIsTerminated|TestDetachSurvivor' ./internal/cli`
  — exit 0 on the final source tree; 7.127 s.
- `GOOS=linux go build -buildvcs=false ./internal/suiteguard ./internal/gittest`
  and `GOOS=windows go build -buildvcs=false ./internal/suiteguard ./internal/gittest`
  — both exit 0 on the final source tree. Linux and Windows runtime behavior
  was not executed on this Darwin host.
- `git diff --check` — exit 0. Source inspection confirmed both glossary
  definitions and ADR links, the repository rule, F6's first statement, and
  the ceiling of 10.

### Acceptance evidence

| Acceptance criterion | Implementation and focused evidence |
| --- | --- |
| SIGKILL of a test binary leaves neither package's helpers alive | Both named death tests ran in the passing store and CLI package checks; the CLI death test also ran under race detection. Four store and three CLI identities are observed without helper signals before the assertions. |
| Marked residue fails its package and is dead afterward | TestGuardReportsAndKillsHelperResidue requires exit 1, exactly one PID-specific residue line, and subsequent process absence or zombie state. TestKillResidueRereadsTheExactMark additionally proves mismatched-mark refusal and matching-mark SIGKILL. Both ran in the passing suiteguard check. |
| Ended helpers pass without a residue line | TestGuardPassesWhenEveryHelperEnded closes the writer and joins its helper, then requires exit 0 and no residue diagnostic; it ran in the passing suiteguard check. TestGuardMarksItsTestBinary checks the PID and nonce form. |
| F6 is parallel and the CLI ceiling is 10 | Source inspection confirmed both changes; the full testfixture check exercised its Parallel Test Package and Signal rules, and the delivery tests passed with race detection and shuffled order. |
| Hardened repositories never start a filesystem monitor despite a global true setting | TestHardenedRepositoryNeverStartsAFilesystemMonitor ran in the passing gittest check, covering InitRepo, direct Harden, bare repositories and linked worktrees; it asserts both local false values, committed same-size rewritten content and no monitor launch in Trace2. |

### Findings and scope notes

- The intentional residue fixture was reported and killed. No unexpected Helper
  Residue was found in the final checks of store, CLI, gittest or suiteguard;
  no package leak was suppressed and neither the mark nor grace was widened.
- Earlier suiteguard checks exposed a Darwin enumeration race: argument-memory
  reads can return EIO during exit/exec, and the typed process lookup reports
  EIO even for an empty response from a reaped PID. The reader now proves
  absence from raw status bytes and retries live argument-memory transitions
  only within the original grace; persistent errors still report unavailable.
  The final suiteguard check passed after those changes.
- An earlier CLI package check ran while Agent edits changed suiteguard files.
  Its tests printed PASS but its repository-boundary guard correctly failed.
  The final CLI package check above ran with edits paused and exited 0.
- Additional ordinary source/test paths needed by this slice are
  `internal/store/main_test.go`, `internal/suiteguard/residue_test.go`,
  `internal/suiteguard/residue_unix_test.go` and
  `internal/suiteguard/residue_darwin_test.go`. Governed mutations stay within
  the approved repository rule and gittest source paths.
- The initial worktree change was the Daemon-owned `status: in_progress`
  field. It was preserved. No other Task or Task Graph was edited; no commit,
  push or pull request was made. Authored Verification, repository Verification
  and terminal Task settlement remain Daemon-owned.


### Verification Feedback — attempt 1

The Daemon's configured `make verify-changed` exited 2. The diagnostic artifact
is `/Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20261012T221735Z_f00baa2097cf9b98/verification/batch-002-attempt-1.log`.
The only test failure recorded there is
`TestBudgetExceededRunRecordsItsReason` in `internal/daemon`: its assertion
requires the first fixture Task and its commit to settle before a fixed
500 ms wall-clock Run Budget expires. In this gate run that settlement had
not happened when the budget expired. No diagnostic log body is embedded here.
The artifact records passing store, suiteguard, gittest and testfixture
package checks; these are Daemon-provided facts, not Agent reruns.

The same budget assertion was already reported in task_01's Verification
Feedback Result, before this Task's implementation. That prior Result records
an isolated pass and a baseline diagnostic demonstrating the assertion's
sensitivity to first-Task work taking longer than its budget; those prior
commands were read, not rerun in this feedback turn.

Current source inspection confirmed:

- The failing test uses the real wall clock and preserves all its existing
  completion, commit, reason and Task-state assertions.
- Repository fixture preparation and gittest hardening happen before the test
  assigns `RunStartedAt` and starts that budget. Task_02 adds two keys to the
  existing single config append for an ordinary fixture repository; it adds
  no subprocess to that path.
- Suiteguard calls `m.Run` before scanning Helper Residue. The scan and its
  ten-second grace therefore cannot consume this test's Run Budget from
  inside its own package lifecycle.
- This Task changes neither the daemon budget implementation nor its fixture.
  No Helper Residue or mark-related failure was reported in the artifact.
  The artifact does not establish the exact source of the Git/settlement
  latency during the configured gate.

Focused diagnostic check in this feedback turn:

- `GOCACHE=/private/tmp/roundfix-task02-go-build rtk proxy go test -count=1 -race -run '^TestBudgetExceededRunRecordsItsReason$' ./internal/daemon`
  — exit 0, 3.131 s, on one isolated run. This narrow pass is not a repair
  and does not establish that configured Verification will pass under load.

No in-slice implementation defect was identified, so this feedback turn
updates only task_02's Result. The configured gate failure remains unresolved.
The follow-up belongs to the daemon budget fixture: arrange expiration after
first settlement through a deterministic test seam, preserving its existing
outcome, reason, commit and Task-state assertions. Increasing the fixed
budget, weakening the assertions or changing helper isolation would not
repair that test contract and is not part of this Task's slice.

Authored Verification, `make verify-changed` and the full suite were not rerun
by the Agent. Task status, other Task files, the Task Graph and implementation
sources were preserved; no commit, push or pull request was created. The
Daemon retains the next configured Verification run and settlement ownership.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `internal/store/main_test.go`
- `internal/suiteguard/residue_darwin_test.go`
- `internal/suiteguard/residue_test.go`
- `internal/suiteguard/residue_unix_test.go`

## Carry-forward provenance

- Source Run: `run_20261012T221735Z_f00baa2097cf9b98`
- Source commit: `f4a2c2d419a9ea8bb89e0c0f7e06b195f10c5e5b`
