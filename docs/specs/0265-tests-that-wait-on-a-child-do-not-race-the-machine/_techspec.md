---
spec: 0265-tests-that-wait-on-a-child-do-not-race-the-machine
prd: _prd.md
created: 2026-10-13
---

# Tests that wait on a child do not race the machine — Technical Spec

## Executive Summary

Every test the authoring sweep found waits on real time for work it does not
control: a re-executed test binary, a `git` or `sh` child, a goroutine. This
design makes each wait end on the event its assertion reads or on the test
deadline (`internal/testwait`), makes every asserted budget race work that
blocks, and adds a `go/ast` repository rule that holds each test file to a
ceiling of short constant budgets. Two production files gain an unexported or
defaulted dependency so tests can reach a budget they cannot reach today; no
production value changes.

The trade-off accepted is cost against proof in one test:
`TestRunDetachedCommandMonitorReportAllowsRunCreationPastLivenessDeadline` must
show that Run creation may arrive after the liveness deadline, and the liveness
deadline necessarily bounds the child's startup. The test therefore uses the
production liveness budget, ten seconds, and the child writes Run creation ten
seconds after its liveness byte. The test takes about ten seconds instead of
one, in parallel with the rest of `internal/cli`. The second trade-off is that
the rule counts budgets per file rather than marking each one, so a new budget
can hide behind a removed one in the same file; the ceilings only fall.

## Project Constraints

- Identifier strategy: not applicable — no identifier, schema or record format
  changes. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request, credential or network
  call; tests use temporary homes, repositories and fake runners. Source:
  `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0271 governs every decision below:
  "An asserted budget races only work that cannot finish". ADR-0259: "a
  ceiling only falls". ADR-0213: "must end when the test binary that started it
  ends". ADR-0266: "A Proven Signal is the only signal Roundfix code or tests
  send". ADR-0125: "A spawned fixture is a compiled binary, never a written
  script". ADR-0126: "The guard therefore runs per package". ADR-0089: "Code
  under test takes its environment explicitly". ADR-0244: "A Spec declares the
  domain terms it introduces". ADR-0257: "Inside a Run the Daemon is the only
  full gate". ADR-0258: "A Spec is authored against the QA rerun classes".
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AF grant of 2026-10-12 ("Sim, como no AE"); bounded files:
  `docs/agents/specific-repository.md`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0265-tests-that-wait-on-a-child-do-not-race-the-machine/_authorization.md`.

## Current behavior

Measured on 2026-10-13 at `94fbf85b` on the maintainer's 10-core arm64 Mac.

**Reproduction.** Pure CPU burners do not reproduce the detach failure: the
detach tests passed 30 times under `-race -parallel 16` with 40 `yes` processes
and an exec storm at load average 30. A full suite under 20 burners at load
average 24 to 139 failed only `TestBudgetExceededRunRecordsItsReason`, the case
0259 already fixed. The failure needs the full suite's mix of fresh test
binaries, forks and Git, so the inventory below comes from a static sweep of all
636 test files, read test by test, and the QA gate measures the result under
load. The load run was stopped early because it slowed other Runs on the same
machine; the QA gate's load windows are therefore niced, PID-tracked and
trap-bounded.

**Inventory.** Each row gives a short real-time budget the test does not
assert, or an asserted one whose competitor can finish. All rows are fixed by
the Task named.

| # | Site | Test | Budget | What the test asserts | Task |
| - | --- | --- | --- | --- | --- |
| 1 | `internal/cli/detach_test.go:61` | `TestRunDetachedCommandMonitorReportAllowsRunCreationPastLivenessDeadline` | liveness 500 ms; Run creation 2 s against a 750 ms child delay | exit 0 and the exact detach report | task_01 |
| 2 | `internal/cli/detach_test.go:88` | `TestRunDetachedCommandTimesOutWaitingForRunCreation` | liveness 500 ms; the asserted 50 ms races a child that sleeps 2 s and exits | the Run-creation timeout text | task_01 |
| 3 | `internal/cli/detach_test.go:41` | `TestRunDetachedCommandTimesOutWaitingForLiveness` | the asserted 50 ms races a child that sleeps 2 s and exits | the liveness timeout text | task_01 |
| 4 | `internal/cli/detach_test.go:108` | `TestRunDetachedCommandReportsInvalidLivenessHandshake` | liveness 1 s over child startup | the invalid-marker failure | task_01 |
| 5 | `internal/cli/detach_test.go:128` | `TestRunDetachedCommandReportsMalformedRunCreationHandshake` | liveness and Run creation 1 s | the malformed-handshake failure | task_01 |
| 6 | `internal/cli/detach_test.go:148` | `TestRunDetachedCommandReportsChildExitBeforeHandshakeWithOutput` | liveness 1 s | exit 2 and the relayed console | task_01 |
| 7 | `internal/cli/detach_test.go:168` | `TestRunDetachedReviewRefusalLeavesArtifactDirectoryAbsent` | liveness and Run creation 1 s | exit 2, no Artifact Directory | task_01 |
| 8 | `internal/cli/detach_test.go:212` | `TestRunDetachedCommandReportsSilentChildExitBeforeHandshake` | liveness 1 s | the silent-exit text | task_01 |
| 9 | `internal/cli/implement_test.go:1501`, `:1935` | `TestRunImplementDetachSurvivesCallerProcessGroupKill` | `attachDetachBudget` 2 s over opening the store and following | attach detaches with `Detached; Run … keeps going.` | task_01 |
| 10 | `internal/cli/implement_test.go:6422` | `TestImplementRunBudgetBoundsSetupAndIntegration/inside budget` | real 10 s Run Budget over worktree, commit, verification and integration | exit 0, state Clean | task_01 |
| 11 | `internal/cli/orphan_unix_test.go:62`, `:120` | `TestRunForceStopOwnerProcessIntegrationProvesExitBeforeStoreCompletion` | 5 s monitor over a stop whose production window is 2 s grace plus 10 s; 2 s for the owner to exit | Stopped, owner dead | task_01 |
| 12 | `internal/cli/orphan_unix_test.go:256` | `TestRunForceStopLegacyRunWithoutOwnerIdentityStillStopsOwner` | 2 s for the owner to exit | owner dead | task_01 |
| 13 | `internal/cli/readiness_forge_test.go:369`, `:380` | `TestForgeProbesAreBoundedAndNeverPromptOrPrintAToken` | 3 s elapsed ceiling; 2 s probe that must cover child startup; 5 s ceiling | deadline ends the probe; the child started and died | task_01 |
| 14 | `internal/cli/review_dispose_lock_test.go:31`, `:193` | `TestConcurrentDisposeOfOneFindingAppendsOnce`, `TestConcurrentDisposeOfTwoFindingsBothSucceed` | 10 s start barrier for 32 re-executed children, on both sides | one append; both succeed | task_01 |
| 15 | `internal/agent/acpx_runner_test.go:2981` | `TestACPXRunCancelsPromptCooperatively` | harness stop grace 2 s, a real timeout on the `cancel` child | cooperative stop | task_02 |
| 16 | `internal/agent/acpx_runner_test.go:3033` | `TestACPXRunClosesSessionAfterCancelGracePeriod` | same 2 s on `cancel` and `sessions close` | killed, close ran, order | task_02 |
| 17 | `internal/agent/acpx_runner_test.go:3083` | `TestACPXRunCancellationCommandFailuresWarnAndContinue` | same | warnings, close ran | task_02 |
| 18 | `internal/agent/acpx_runner.go:56` via `fallback.go:95`, `selection_assignment.go:162`, `sealed.go:107` | `TestProveExactSelectionCancelCleanup` (`:1459`), `TestProveExactSelectionTimeoutCleanup` (`:1485`), `TestACPXProbeCancellationStillClosesDisposableSession` (`:1697`), `TestACPXProbeFallbackCancellationStillClosesCandidateSession` (`:1842`), `TestSealedACPXPromptCancellationCancelsAndClosesSession` (`sealed_test.go:124`), and about forty passing selection, probe and sealed-prompt tests | production 5 s cleanup timeout around the `sessions close` child | `sessions close` ran; no cleanup error | task_02 |
| 19 | `internal/agent/acpx_integration_test.go:224`, `:366` (opt-in) | `TestRealACPXCommandOverrideRoundTripCancelCrashResume` | stop grace 2 s; 2 s around a `ps` child | cooperative cancel; no orphan | task_02 |
| 20 | `internal/agent/acpx_integration_test.go:51` | `TestSpawnedFixtureDeathFailsFast` | asserted 9 s measured from `cmd.Start`, so it includes the fixture's startup | death noticed within 9 s | task_02 |
| 21 | `internal/agent/agent_test.go:1070` | `TestConsoleDisplaySinkSerializesConcurrentPublish` | 1 s for a goroutine's first write | serialized output | task_02 |
| 22 | `internal/spec/source_revision_test.go:183` | `TestSourceRevisionRecoveryCancelsTheFetchProcess` | 2 s over `git cat-file` and the fetch fixture's startup | fetch cancelled, marker written, process dead | task_03 |
| 23 | `internal/worktree/worktree_test.go:710` | `TestRunBootstrapReturnsBootstrapErrorOnTimeout` | asserted 1 s; the test also requires the shell's start marker before the timeout | `BootstrapError` timed out after 1s | task_03 |
| 24 | `internal/daemon/task_engine_test.go:7675` | `TestTaskCycleEndsRunAtBudgetDeadline` | asserted 3 s real Run Budget, raised from 100 ms by Spec 0213, raced by real Git setup | task_01 cancelled once, task_02 pending | task_03 |
| 25 | `internal/store/journal_test.go:499` | `TestRetentionScanOutsideWriteTransaction` | 500 ms scan while the write lock is held | the scan does not need the lock | task_03 |
| 26 | `internal/store/prune_reclaimable_test.go:64` | `TestPruneTerminalRunsNeedsNoWriteLockWhenNoCandidateHasEvents` | 250 ms prune | the prune does not need the lock | task_03 |

**Already fixed on 0259's item branch, left alone:**
`TestBudgetExceededRunRecordsItsReason` (`internal/daemon/task_engine_test.go:7721`),
`TestRemoveRegisteredWorktreeWaitsForTheAdminLock` and
`TestAddDetachedWorktreeWaitsForTheAdminLock` (`internal/worktree/disposal_test.go:16`, `:36`).

**Kept, with the reason:**

- `detachTestSurvivor.terminate` in `detach_test.go` waits `detachWaitTimeout`
  (2 s) to reap a child after `SIGKILL`. 0259 rewrites the function, and reaping
  a killed process does not wait for a child to start.
- `detachWaitTimeout` in production `detach.go` waits 2 s to reap a child after
  its handshake ended with end of file, for rows 6 to 8. The child is already
  exiting; the constant is production behavior and 0259 changes `detach.go`.
- `TestOutcomeNotificationCarriesTerminalContextWithBoundedDeadline`
  (`internal/cli/cli_test.go:6111`) and `TestAskSendsTheMeasuredRequest`
  (`internal/judge/client_test.go:33`) allow 1 s between two in-process
  statements; no child, goroutine hand-off or Git sits between them, and
  `cli_test.go` is a Governed Path.
- Four `internal/testwait` tests bound elapsed time by 1 s; promptness is
  their subject and only goroutine scheduling competes.
- Budgets of ten seconds or more: the production detach liveness budget, Force
  Stop's 10 s window, the readiness probe's 10 s production timeout and the QA
  formatter's 10 s ready guard.
- `TestTaskBudgetWatchdogStopsWhenTheCycleReturns` counts process-wide
  watchdog goroutines while other parallel tests start theirs. That is
  interference between parallel tests, not a real-time budget; it is outside
  this class.

**Rule baseline.** A prototype of the rule found 72 constant budgets below ten
seconds in 24 test files. Most are data for fake clocks (`internal/watch`,
`internal/delivery`, `internal/notify`, the GC retention clock) or asserted
budgets against blocking work; the rest are the rows above.

**0259 overlap.** Four declared files are also changed on 0259's item branch.
The functions this Spec changes in each lie outside 0259's hunks, so the merge
0259 performs afterwards needs no hand resolution:

| File | This Spec changes | 0259 changes |
| --- | --- | --- |
| `internal/cli/detach_test.go` | the bodies of rows 1–8; the `sleep-before-liveness`, `slow-run-creation` and `run-created` cases of `runDetachTestChild`; one assignment inside and one after the `run-created` block of `runDetachParentForTest` | imports, constants, the survivor struct, `TestDetachedChildIsTerminatedAtTeardown`, the owner environment line and the survivor identity of `runDetachParentForTest`, the owner checks, `terminate`, two new helpers |
| `internal/cli/implement_test.go` | `attachDetachBudget` and the attach call of row 9; row 10's subtest | the caller kill of `TestRunImplementDetachSurvivesCallerProcessGroupKill` |
| `internal/cli/orphan_unix_test.go` | rows 11 and 12 | `TestCLIForceStopOwnerProcessHelper` |
| `internal/daemon/task_engine_test.go` | the body of `TestTaskCycleEndsRunAtBudgetDeadline` and `budgetDeadlineRunner` | the body of `TestBudgetExceededRunRecordsItsReason` |

## System Architecture

- **Detach tests** (`internal/cli/detach_test.go`): test-only. Budgets a test
  does not assert take `testwait.Bound(t)` or the production constants
  `detachHandshakeTimeout` and `detachRunCreationTimeout`; blocking child modes
  reuse `waitForDetachTestSentinel`, which ends when the owner is gone.
- **Other CLI tests**: test-only waits on events and on the test deadline.
- **Agent cleanup timeout** (`internal/agent`): an unexported `ACPXRunner`
  field beside the existing `cancelClock`, read by the three cleanup sites.
- **Run Budget timer** (`internal/daemon`): a `Deps` field beside `Now` that
  the budget watchdog uses to start its timer, defaulting to the real timer.
- **Real-Time Budget rule** (`internal/testfixture`): a repository test beside
  the Parallel Test Package rule.

## Implementation Design

### Interfaces

```go
// internal/agent/acpx_runner.go
type ACPXRunner struct {
	// existing fields …
	cleanupTimeout time.Duration // zero: acpxPreflightCleanupTimeout
}
func (runner *ACPXRunner) disposableCleanupTimeout() time.Duration

// internal/daemon/engine.go
type Deps struct {
	// existing fields, Now among them …
	BudgetTimer func(time.Duration) (<-chan time.Time, func() bool) // nil: time.NewTimer
}

// internal/testfixture/real_time_budget_test.go
const realTimeBudgetThreshold = 10 * time.Second
type realTimeBudgetCeiling struct{ max int; reason string }
var realTimeBudgetCeilings = map[string]realTimeBudgetCeiling{ /* file → ceiling */ }
func realTimeBudgetViolations(root string, ceilings map[string]realTimeBudgetCeiling) ([]string, error)
```

### Data Models

None. No record, schema or configuration changes.

### API Contracts

1. API Contract: detach test budgets. In `detach_test.go`, a phase budget the
   test does not assert is `testwait.Bound(t)`, except row 1, whose liveness
   budget is `detachHandshakeTimeout` and whose Run-creation budget is
   `detachRunCreationTimeout`. The asserted budgets keep their values (50 ms in
   rows 2 and 3), so every expected diagnostic stays byte-identical. Every
   assertion of rows 1 to 8 stays: exit status, stdout, stderr and the
   diagnostic substrings.
2. API Contract: detach child modes. `sleep-before-liveness` blocks before its
   liveness byte, and `slow-run-creation` blocks after it, both through
   `waitForDetachTestSentinel` on a sentinel path the test never writes, so
   each ends only when the parent kills it or its owner is gone.
   `run-created` waits, after its liveness byte, for the duration in
   `ROUNDFIX_DETACH_TEST_RUN_CREATION_DELAY`, 750 ms when unset, before it
   writes Run creation; `runDetachParentForTest` sets that variable to the
   liveness budget, so Run creation always arrives after the liveness deadline.
3. API Contract: CLI waits. Row 9 cancels attach once its standard output
   holds `Following live events.` instead of letting a 2 s context expire, and
   keeps both assertions. Row 10's Run Budget is one hour. Rows 11 and 12 wait
   for the owner's exit and the monitor's outcome through `testwait`. Row 13
   keeps `context.DeadlineExceeded` and the 250 ms budget in its first half,
   with the child's 30 s sleep as the elapsed ceiling, and in its second half
   waits for the child's PID file, then cancels the probe and asserts that
   cancellation and the child's death. Row 14's barriers are bounded by the
   parent's test deadline, which the parent passes to its children.
4. API Contract: agent cleanup timeout. `disposableCleanupTimeout` returns
   `cleanupTimeout` when it is above zero, else `acpxPreflightCleanupTimeout`
   (5 s). The cleanup contexts in `fallback.go`, `selection_assignment.go` and
   `sealed.go` use it. Every test that builds an `ACPXRunner` whose cleanup
   spawns a child sets the field to `testwait.Bound(t)`; the fake ACPX harness
   stop grace becomes `testwait.Bound(t)`, because the fake cancellation clock
   already drives the grace timer the tests assert.
5. API Contract: Run Budget timer. The watchdog in `task_engine.go` starts its
   timer through `Deps.BudgetTimer` when it is set and through `time.NewTimer`
   otherwise; the duration passed stays `time.Until(deadline)`. Row 24's test
   sets `BudgetTimer` to a fake whose channel delivers only after
   `budgetDeadlineRunner` reports that task_01 started, and keeps all four
   assertions.
6. API Contract: the remaining rows. Row 22 cancels the fetch after the fetch
   fixture reports its start and asserts the cancellation it caused, the
   marker fields and the dead process; row 23 accepts a timeout that ends before
   the shell writes its marker and keeps every `BootstrapError` assertion;
   rows 25 and 26 bound the operation by the test's context, not 500 or 250 ms;
   row 19 uses `testwait.Bound(t)` for both budgets; row 20 starts its 9 s
   measurement at the fixture's first output; row 21 waits through
   `testwait.Until`.
7. API Contract: the Real-Time Budget rule.
   `TestEveryRealTimeBudgetIsWithinItsFileCeiling` walks every `*_test.go` under
   the repository root except `.git`, `.gocache`, `bin`, `testdata`,
   `node_modules`, `.agents` and `skills`. It counts constant duration
   expressions below `realTimeBudgetThreshold`, built from numeric literals,
   `time` unit selectors, `time.Duration(…)`, `+`, `-`, `*` and `/`, or held in
   a `const` declared in a test file of the same package (a production constant
   is not counted), in three forms: the duration argument of
   `context.WithTimeout`, `time.After`, `time.NewTimer` or `time.AfterFunc`; a
   keyed composite-literal field or an assignment whose name matches
   `(?i)timeout|deadline|grace|budget|liveness|runcreation|maxrunduration`; and
   an argument to a function declared in the same package whose parameter name
   matches the same pattern. A file whose count exceeds its ceiling, or that
   has no ceiling and any count, fails with one line per budget,
   `<file>:<line>: Real-Time Budget <duration> in <form>`, after a line naming
   the file, its count and its ceiling. Every ceiling carries a non-empty
   reason. As in the Parallel Test Package rule, only a count above the ceiling
   fails; a ceiling above its file's count is logged, so a Spec delivered in
   between that removes a budget never breaks the rule, and the next change to
   that file lowers the ceiling.
   `TestRealTimeBudgetRuleRejectsSeededViolations` runs the same walk on trees
   it writes under `t.TempDir()`: `unlisted-timeout`, `over-ceiling`,
   `parameter-budget` and `named-constant` fail naming the seeded line, and
   `long-budget`, `sleep-only` and `listed` pass.

### Surface Transcripts

None. No command's output, flags or exit codes change.

### Invariants

1. No production budget value changes; every new dependency defaults to the
   value used today.
2. Every assertion of rows 1 to 8 holds byte-for-byte.
3. A test bounds a child process, a goroutine or Git by real time only when the
   budget is what it asserts, and that work then blocks until cancelled or its
   owner is gone.
4. No Task adds a signal, a written executable or a Sequential Test.
5. A ceiling in the rule only falls.

## Coverage Map

- Goal 1, Core Feature 1, Success Metric 1 → detach tests (task_01).
- Goal 2, Core Feature 2 → other CLI tests (task_01).
- Goal 2, Core Feature 3 → agent cleanup timeout and stop grace (task_02).
- Goal 2, Core Feature 4 → Run Budget timer and the store, spec and worktree
  tests (task_03).
- Goal 3, Core Feature 5, Success Metric 3 → Real-Time Budget rule, glossary and
  repository rule (task_04).
- Goal 4, Success Metric 4 → Invariant 1, proved from the diff by the QA gate
  (task_05).
- Success Metric 1, Success Metric 2 → QA gate under synthetic load (task_05).

## Integration Points

None outside the repository.

## Testing Approach

- Each fixed test runs on its own once with `-race` in its Task's
  Verification, and the Task's package runs once with `-race` (ADR-0259).
- task_01 and task_03 prove the detach and Run Budget phases without wall-clock
  luck: the seeded proof is that the assertions stay and the budget a test does
  not assert can no longer expire first.
- task_02 adds `TestDisposableCleanupTimeoutDefaultsToTheProductionValue`, which
  asserts that a runner without the field uses `acpxPreflightCleanupTimeout`.
- task_03 adds `TestTaskBudgetWatchdogUsesTheRealTimerByDefault`, which asserts
  that an engine without `BudgetTimer` still ends a Run at a real deadline.
- task_04's rule test runs on the repository and on seeded trees.
- The QA gate runs each fixed test 50 times under `-race -parallel 16` and
  three full suites, both while niced CPU burners run.

## Build Order

1. Detach tests and the other `internal/cli` tests of the class (task_01).
2. Agent cleanup timeout and stop grace (task_02).
3. Run Budget timer and the daemon, worktree, store and spec tests (task_03).
4. The Real-Time Budget rule with its ceilings, the **Real-Time Budget** term
   and the repository rule line (task_04) (depends on: 1, 2, 3).
5. QA gate (task_05) (depends on: 4).

Steps 1 to 3 share no file and run side by side. Step 4 measures the ceilings
on the tree the first three leave, so it follows all of them.

## Risks & Considerations

- **A test deadline as a budget hides a hang longer.** A wait bounded by
  `testwait.Bound(t)` that never ends fails at the test deadline with stacks,
  not after one second. That is the cost of not racing the machine, and
  `testwait` already prints every goroutine's stack.
- **0259's merge.** The overlap table shows disjoint regions, but 0259's
  Helper Lifeline test starts the `run-created` child without the new delay
  variable; the 750 ms default keeps that test unchanged.
- **The rule's limits.** It cannot see budgets carried by identifiers from
  another package, production constants or sleeps. Not counting production
  constants also keeps 0259's merge clean: its rewritten `terminate` waits
  `time.After(detachWaitTimeout)` twice, and `detachWaitTimeout` lives in
  `detach.go`. The inventory records the
  ones that exist; ADR-0271 states the limit.
- **Sibling Specs.** A Spec delivered after this one that adds a short budget
  fails the rule until it uses an event, the test deadline or a ceiling with a
  reason. That is intended.

## Glossary

- adds: **Real-Time Budget**
- not a term: **Current Behavior** — a section heading
- not a term: **Detach tests** — a group of test functions
- not a term: **Other CLI tests** — a group of test functions
- not a term: **Agent cleanup timeout** — a field of one runner
- not a term: **Run Budget timer** — a dependency field of the Daemon engine
- not a term: **Real-Time Budget rule** — the name of one repository rule test
- not a term: **Rule baseline** — a paragraph label
- not a term: **Already fixed on 0259's item branch, left alone** — a paragraph label
- not a term: **Kept, with the reason** — a paragraph label
- not a term: **A test deadline as a budget hides a hang longer** — a risk label
- not a term: **The rule's limits** — a risk label
- not a term: **Sibling Specs** — a risk label

## Decisions

- Budgets a test does not assert take the test deadline or the production
  constant; see ADR-0271.
- Asserted budgets race blocking work; detach children reuse the owner-watching
  sentinel wait; see ADR-0271.
- The agent cleanup timeout and the Run Budget timer become dependencies with
  production defaults, because no test can reach them otherwise; see ADR-0271.
- A per-file ceiling table holds the remaining constant budgets; see ADR-0271.
- The QA gate accepts on three full suites and 50 runs of each fixed test under
  niced synthetic load.
