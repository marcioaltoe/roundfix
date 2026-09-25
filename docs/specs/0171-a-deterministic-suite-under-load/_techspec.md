---
spec: 0171-a-deterministic-suite-under-load
status: active
created: 2026-09-25
surfaces: [backend]
---

# A deterministic suite under load

## Executive Summary

Add `internal/testwait`, a small helper package whose waits are bound by the
test deadline and end at once when the watched work ends; move the Implement,
Daemon and worktree tests onto it and raise the product timeouts they do not
examine; answer the release lookup offline in every `internal/cli` test process
and seed a fresh version cache for built binaries; resolve the adapter
fixture's link target to an absolute path. The only production change is the
corrective serialization of Git worktree administration in `internal/worktree`
(task_07), added after the first QA gate reproduced a Git race.

## Project Constraints

- Identifier strategy: not applicable — no identifier. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: applicable — the release lookup's `gh api` call
  leaves every test process; no credential is read. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0034, ADR-0056, ADR-0080, ADR-0089,
  ADR-0091, ADR-0093, ADR-0096, ADR-0097, ADR-0104, ADR-0117, ADR-0125,
  ADR-0126, ADR-0130, ADR-0155 and ADR-0156 hold. Source:
  `docs/agents/domain.md`.
- Tooling authority: not applicable — empty intersection with `GovernedPath`.
  Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## The wait helper

`internal/testwait/testwait.go` is an ordinary package with no production
caller, like `internal/testfixture`. It exports:

- `Margin`, a named constant of a few seconds left before the deadline so a
  failing wait can still report.
- `Bound(t testing.TB) time.Duration`, the time until `t.Deadline()` minus
  `Margin`; when the test has no deadline (`-timeout 0`), a documented fallback
  equal to `go test`'s default timeout of ten minutes. It never returns a
  smaller fixed duration.
- `Until[T, R any](t testing.TB, what string, ready <-chan T, ended <-chan R) T`,
  which returns the first value from `ready`; fails the test at once when
  `ended` delivers first, printing that value; and fails at `Bound(t)` naming
  `what` and dumping every goroutine's stack, as `failTestWait` in the Daemon
  tests does today. A nil `ended` watches nothing.
- `Poll[R any](t testing.TB, what string, ended <-chan R, condition func() (bool, string))`,
  which re-evaluates `condition` at a short interval until it holds; fails at
  once when `ended` delivers first, printing that value; and at `Bound(t)` fails
  with `what` and the condition's last observation.

When the awaited value and the end of the work are both ready, the awaited
value wins, and `Poll` evaluates its condition once more before it reports an
ended work: a command that produces its last effect and then exits must not fail
a wait that raced the two.

A failing wait calls `t.Fatalf`, so it must run on the test goroutine. It
abandons nothing it does not own: cancelling the watched work stays with the
test's context and cleanups. The package spawns no process, so ADR-0126 does
not require it to install the suite guard.

## The Implement tests

In `internal/cli/implement_test.go`, `implementWaitBudget`,
`detachStartupBudget` and the literal 90-second waits go, and every wait routes
through `testwait`: `waitImplementCommandResult`, `waitImplementAgentStarts`,
`waitForImplementJournal`, `waitForFile`, `waitForFileContains`,
`waitForRunState`, `waitForCleanOutcomeEvent`, `readLineWithTimeout` and
`waitProcessForTest`. A wait that runs while `runImplementCommandAsync` holds a
command in flight takes that command's result channel as `ended`, so a command
that ends first fails the wait with its exit code, stdout and stderr. The three
tests that configure `bootstrap_timeout: 1s` configure a timeout derived from
`testwait.Bound(t)` instead. `attachDetachBudget` stays: it bounds attach's own
follow loop, and Spec 0124 owns it.

## The Daemon and worktree tests

In `internal/daemon/task_engine_test.go`, `testWaitBound`, `testWaitFallback`,
`testWaitDeadlineMargin` and `failTestWait` are replaced by `testwait`. The
waits that run while a `TaskCycle` goroutine is in flight — `waitSchedulerStarts`,
`taskCapacityVerifier.waitStart`, `waitPublishedVerificationPhase` and
`waitPublishedStopEvent` — take that goroutine's result channel as `ended`
wherever the calling test has one.

In `internal/worktree/worktree_test.go`, every `BootstrapSpec.Timeout` that is
not the subject of its test is derived from `testwait.Bound(t)`.
`TestRunBootstrapReturnsBootstrapErrorOnTimeout` keeps its 10 ms bound against
`sleep 1`, because load can only lengthen the sleep. Each round of
`TestBootstrapSerializesAcrossSiblings` waits for its sibling starts and results
through `testwait`, and releasing a sibling cannot block past the deadline when
that sibling has already exited.

## The release lookup

A new file, `internal/cli/version_freshness_isolation_test.go`, assigns in
`init` the package default `versionFreshnessDeps.latestRelease` to an offline
lookup that returns the tag `v0.0.0-offline-test-lookup` and starts no process.
`init` runs before `TestMain`, so the default covers in-process tests, the
`ROUNDFIX_CLI_TEST_HELPER` subprocess and the detach children, which are all the
same binary. The tag sorts below every release, so a cache holding it never
prints an upgrade warning. Nothing reassigns it; per-test fakes keep using
`withVersionFreshnessFakeDeps`.

`runRoundfixBinaryMacro` in `internal/cli/implement_test.go` and the built
binary in `TestRunDispositionCharacterizationPreflightRefusesOnAnUnintegratedBranch`
(`internal/daemon/run_disposition_characterization_test.go`) seed
`<HOME>/.roundfix/version-check.json` with a current `checked_at` and no
`latest_version` before they run the binary, which then takes
`readVersionFreshnessCache`'s fresh path and performs no lookup.

## The adapter fixture

`provisionFakeAdapter` in `internal/agent/acpx_runner_test.go` links an absolute
path to the compiled test binary, resolved from `os.Args[0]` before any link is
made. The fixture stays the compiled test binary that ADR-0125 requires.

## API Contracts

1. `testwait.Bound` returns the time until the test deadline minus `Margin`, and
   the ten-minute fallback only when the test has no deadline.
2. `testwait.Until` and `testwait.Poll` fail at once when the watched work ends
   first, naming its result, and otherwise never before `Bound`.
3. The default release lookup of every `internal/cli` test process returns
   `v0.0.0-offline-test-lookup` and starts no process.

## Coverage Map

- Goal 1 → The wait helper, The Implement tests, The Daemon and worktree tests.
- Goal 2 → The wait helper; API Contract 2.
- Goal 3 → The release lookup; API Contract 3.
- Goal 4 → The adapter fixture.
- Core Feature 1 → The wait helper.
- Core Feature 2 → The Implement tests.
- Core Feature 3 → The Daemon and worktree tests.
- Core Feature 4 → The release lookup.
- Core Feature 5 → The adapter fixture.
- Success Metric 1 → Testing Approach 2, 3.
- Success Metric 2 → Testing Approach 2, 3.
- Success Metric 3 → Testing Approach 1.
- Success Metric 4 → Testing Approach 4.
- Success Metric 5 → Testing Approach 5.
- API Contracts 1-2 → The wait helper.
- API Contract 3 → The release lookup.

## Integration Points

- **Spec 0124.** Owns `attachDetachBudget` and the attach deadline split.
- **Spec 0165 and Spec 0167.** Made the Implement budget test deterministic;
  this Spec leaves it alone.
- **`docs/references/coverage-record.json`.** Every recorded top-level test
  keeps its name, so the governed record needs no re-record.

## Testing Approach

1. **Wait helper.** Unit tests with a `testing.TB` wrapper that overrides
   `Deadline`: `Bound` follows a deadline and falls back without one; `Until`
   returns a ready value, prefers it when `ended` is ready too, fails within one
   second when `ended` fires first and prints its value, and fails at the
   deadline and not before; `Poll` returns when the condition holds, fails
   within one second when `ended` fires first, and fails at the deadline with
   the last observation.
2. **Implement tests.** A sweep proves `implement_test.go` holds no
   `time.After`, `time.NewTimer`, `implementWaitBudget` or `bootstrap_timeout: 1s`;
   a `-count=5 -cpu 1,4` stress of the four bootstrap and cancellation tests
   passes.
3. **Daemon and worktree tests.** A sweep proves the Daemon file holds no
   `time.After` or `testWaitBound` and the worktree file no one- or five-second
   bootstrap timeout; a `-count=10 -cpu 1,4` stress of the queued-cancellation
   and bootstrap tests passes.
4. **Release lookup.** Guard tests run `fetch`, `resolve`, `watch` and
   `implement` in process and `implement` in the helper subprocess, each in a
   fresh test HOME, and read the offline tag from the cache they wrote; a built
   binary leaves a seeded cache byte-identical; a per-test fake still overrides
   the default.
5. **Adapter fixture.** The compiled test binary re-executed through a relative
   `argv[0]` provisions and runs a fake adapter; the absolute case still runs.
6. **Repository gate.** The terminal QA Task records the Daemon's repository
   Verification result as a fact.

## Worktree administration lock

The first QA gate reproduced the concurrent-bootstrap failure: one `git worktree
add` read a sibling's administrative files while Git was still writing them
(`failed to read .git/worktrees/<task>/commondir: Result too large`).
`internal/worktree` now serializes `git worktree add`, `remove`, `prune` and
`move` per Git common directory, inside one process and across processes, for
the duration of the Git command only, bounded by the caller's context. No
output, exit code or path layout changes.

## Build Order

1. The wait helper (depends on: none).
2. The Implement tests (depends on: 1).
3. The Daemon and worktree tests (depends on: 1).
4. The release lookup (depends on: 2).
5. The adapter fixture (depends on: none).
6. Terminal QA (depends on: 1, 2, 3, 4, 5, 7, 8).
7. The worktree administration lock (depends on: 3).
8. The bootstrap timeout test (depends on: 7).

The release lookup waits for the Implement tests only because both edit
`internal/cli/implement_test.go`.

## Risks & Considerations

- **A rare flake that is not one of the three shapes.** No reproduction was
  obtained, so a remaining cause cannot be excluded; the waits now fail with the
  work's result or a full goroutine dump, which names it on its next occurrence.
- **A test without a test-owned HOME.** An in-process test that runs an
  operational command against the real HOME writes the offline tag into the
  maintainer's cache, which silences the upgrade notice for up to 24 hours;
  today such a test writes GitHub's answer there.
- **A slower failure.** A genuinely stuck wait now fails near the package
  deadline instead of after 90 s; the goroutine dump makes that one failure
  worth its time.
