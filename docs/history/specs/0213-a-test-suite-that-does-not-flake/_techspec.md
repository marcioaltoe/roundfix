---
spec: 0213-a-test-suite-that-does-not-flake
prd: _prd.md
created: 2026-10-01
---

# A test suite that does not flake — Technical Spec

## Executive Summary

The six flakes have five causes, and none of them is in production code. Each
fix makes the test observe the event its assertion depends on, never elapsed
time. The linger subtest waits for the commit. The Run Budget tests decide
expiry on the injected clock and order the renewing settlement after the
stalled start. `internal/cli` script fixtures become re-executions of the
compiled test binary, as ADR-0125 requires. Fixture processes watch their test
binary (ADR-0213). The Assets Sync template moves from disk into memory. The
primary trade-off is scope: a repository guard freezes the written-executable
class at a named residual inventory of ten written executables, rather than
converting them: eight in four other packages, one in `upgrade_test.go` that is
never executed, and one in the governed `cli_test.go`, which this Spec has no
grant to edit.

## Project Constraints

- Identifier strategy: not applicable — no identifier is created or changed.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local fixtures only, no
  credential and no network call. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0125 and ADR-0213 bind the fixture
  changes; ADR-0126, ADR-0028, ADR-0098, ADR-0158 and ADR-0164 hold
  unchanged; ADR-0093, ADR-0117, ADR-0168, ADR-0176 and ADR-0183 check
  consistency; the gate is bound by ADR-0080, ADR-0088, ADR-0091, ADR-0104,
  ADR-0155, ADR-0156 and ADR-0167; ADR-0030, ADR-0096, ADR-0097, ADR-0170, ADR-0194, ADR-0195, ADR-0210,
  ADR-0182 and ADR-0184 do not apply, as the PRD records. Source: `docs/agents/domain.md`.
- Tooling authority: not applicable — the intersection of this Spec's changed
  paths with `GovernedPath` is empty; `internal/cli/cli_test.go`, the
  Makefile, CI workflows and `go.mod` stay untouched. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## System Architecture

No production file changes and no package is added. Every change is a test or
a test helper:

| Component | Where | Change |
| --- | --- | --- |
| Linger wait | `internal/store/journal_batch_test.go` | Waits for the committed event through `testwait.Poll` |
| Budget clock | `internal/daemon/task_budget_renewal_test.go` | Both short-budget tests decide expiry on `budgetTestClock` |
| Script fixture | new `internal/cli/script_fixture_test.go`; `TestMain` in `internal/cli/implement_test.go` | A symlink to the test binary plus a `/bin/sh` sidecar |
| Fixture call sites | `internal/cli/adapter_floor_test.go`, `implement_test.go`, `carryforward_hooks_test.go`, `settle_test.go` | Nine written executables become script fixtures |
| Written-executable guard | new `internal/testfixture/written_executable_test.go` | Scans `internal/**/*_test.go` against a residual inventory |
| Owner watch | `implement_test.go` fake ACPX, `internal/cli/detach_test.go` children | Waits exit when the test binary is gone |
| Detach reaping | `implement_test.go`; new `internal/cli/implement_detach_teardown_test.go` | Reaps by the Run's Owner PID; proves death of the test binary ends the child |
| Assets Sync template | `internal/baseline/assets_sync_test.go` | Built without `testing.T`, held in memory, on-disk build removed |

## Implementation Design

### Interfaces

```go
// internal/cli/script_fixture_test.go
const scriptFixtureSuffix = ".fixture.sh"

// writeScriptFixture makes path run body through /bin/sh: path becomes a
// symlink to the absolute compiled test binary, and body is written to
// path+scriptFixtureSuffix with mode 0o600. No test writes a file it executes.
func writeScriptFixture(t testing.TB, path string, body string)

// runScriptFixture replaces the process with /bin/sh <sidecar> <args...> when
// the invoked name has a sidecar; otherwise it returns. TestMain calls it first,
// before any environment-selected mode.
func runScriptFixture()
```

```go
// internal/testfixture/written_executable_test.go
type writtenExecutable struct{ file string; line int }

// writtenExecutables lists every os.WriteFile and os.OpenFile call in
// root/internal/**/*_test.go whose literal mode sets any 0o111 bit.
func writtenExecutables(root string) ([]writtenExecutable, error)

// writtenExecutableResidue is the frozen inventory, file -> exact count.
var writtenExecutableResidue = map[string]int{ /* see The residual inventory */ }
```

```go
// internal/baseline/assets_sync_test.go
type assetsSyncTemplateEntry struct { path string; mode fs.FileMode; data []byte }

// buildAssetsSyncTemplate builds target and source repositories in a
// temporary directory, reads them into memory and removes the directory.
// It takes no testing.T and returns its error.
func buildAssetsSyncTemplate() (target, source []assetsSyncTemplateEntry, revision string, buildDir string, err error)
```

### The linger wait

The linger subtest publishes one event to a writer with batch size 100 and
linger 50 ms. It then calls `testwait.Poll` until `RunEventsAfter` returns one
event, and after that asserts `pendingCount() == 0`. Nothing else can commit
that event: the batch size is not reached, the event is not immediate, and no
flush is called. A commit therefore still proves the linger closed the batch.
The 2-second `time.Now` deadline and the `time.Sleep` loop go.

### The Run Budget tests

Both tests use a two-node graph with independent `task_01` and `task_02`,
`Concurrency: 2`, and `engine.deps.Now = clock.Now` with
`clock := &budgetTestClock{now: startedAt}`. `RunStartedAt` is `startedAt`,
and `MaxRunDuration` is one hour, so the first allowance cannot expire during
the test. `settleThenStallRunner` stalls `task_02` and closes `started`. Its
`task_01` branch waits for `started`, sets the clock to
`time.Now().Add(-maximum)`, and then settles. The renewal therefore computes a
deadline equal to the real present. The real watchdog, `taskCycleBudget.watch`
in `internal/daemon/task_engine.go`, fires at once and cancels the stalled
Task with `context.DeadlineExceeded`. Its renewed deadline came from the
settlement, so the reason ends ` since Task task_01 settled.` Real elapsed
time no longer decides any outcome, and the watchdog's real timer stays under
test. The other tests in the file already use `budgetTestClock` and stay
unchanged.

### The script fixture and its call sites

`TestMain` in `implement_test.go` records the absolute test binary path, by
resolving `os.Args[0]` against the start directory, and then calls
`runScriptFixture`. The sidecar lookup tries `os.Args[0]` with the suffix. When
`os.Args[0]` has no separator, it uses the `exec.LookPath` result, because a
PATH lookup leaves the bare name in `argv[0]`. `syscall.Exec` runs `/bin/sh`
with the sidecar and the original arguments, and the environment unchanged.
The check comes before `cliTestHelperEnv` and `detachTestChildModeEnv`,
because a fixture inherits a helper's environment.

These nine write sites, in five helpers, become `writeScriptFixture`:

1. `doctorWithVersionFixture` in `adapter_floor_test.go`, split into
   `writeAdapterVersionFixture(t, pkg, version) string` so the test can read
   the path.
2. `fakeACPXCommand`'s `acpx`, `codex-acp` and `claude-agent-acp`.
3. `newMacroFakeACPX`'s `acpx` and its `codex-acp`, `claude-agent-acp`,
   `opencode` and `npx` adapters, written in one loop, and the non-Darwin
   branch of `buildMacroCodexExecutable`.
4. The four hooks of `writeCarryForwardHookFixtures`.
5. The pre-commit hook of `settle_test.go`.

Git runs a hook through `access(X_OK)` and `execve` on its path. Both follow
the symlink, and Git ignores the sidecar because its name is not a hook name.
`upgrade_test.go` writes an executable that is never executed, and its mode is
what the upgrade replaces, so it stays as it is and in the inventory.

### The residual inventory

After task_02, the guard expects exactly these counts:

| File | Count | Why it stays |
| --- | --- | --- |
| `internal/cli/cli_test.go` | 1 | Governed; no grant (`fakeACPXVersionCommand`) |
| `internal/cli/upgrade_test.go` | 1 | Never executed |
| `internal/baseline/plan_characterization_test.go` | 1 | Outside this Spec's package scope |
| `internal/baseline/profile_alignment_test.go` | 1 | Outside scope |
| `internal/baseline/skills_lock_read_test.go` | 1 | Outside scope |
| `internal/daemon/commit_hook_test.go` | 1 | Outside scope |
| `internal/daemon/task_engine_test.go` | 2 | Outside scope |
| `internal/preflight/preflight_test.go` | 1 | Outside scope |
| `internal/store/process_unix_test.go` | 1 | Outside scope |

A site outside the inventory fails with `file:line`. A count above or below
the recorded one fails too, so a later Spec that converts a residual updates
the table in the same commit. The guard reads literal modes only. A mode held
in a constant, or an `os.Chmod` that sets an execute bit, is outside its
Boundary and is named in its suite header.

### The owner watch and reaping

- `cliHelperEnv` adds `ROUNDFIX_FAKE_ACPX_OWNER_PID=<os.Getpid()>`. The fake
  ACPX release loop runs `kill -0 "$ROUNDFIX_FAKE_ACPX_OWNER_PID"` on every
  pass, and exits 1 once the process is gone.
- `detach_test.go` passes the test binary's PID to its children. Both
  `waitForDetachTestSentinel` and the `detachTestChildIgnoreSentinel` loop
  return `exitRunFailed` once `store.ProcessAlive(owner)` is false. The second
  child still ignores the sentinel; it only stops ignoring its owner's death.
- `TestRunImplementDetachSurvivesCallerProcessGroupKill` reads the Run's
  `OwnerPID` once the caller prints the Run ID. It registers, after every
  `t.TempDir` it uses, a cleanup that sends `SIGKILL` to the process group
  `-OwnerPID`, treats `ESRCH` as gone, and polls `store.ProcessAlive` to false
  within the test deadline. When `ROUNDFIX_DETACH_TEST_OWNER_RECORD` names a
  file, the test writes the Owner PID there for the death test.

### The Assets Sync template

`assetsSyncTemplateRoot` keeps its `sync.Once`. The once now calls
`buildAssetsSyncTemplate`, which takes no `testing.T`, runs Git through an
error-returning runner with the same `-c` flags, and reads both repositories,
`.git` included, into entries with their modes. It then removes its build
directory and records the path. `newAssetsSyncTarget` and
`newAssetsSyncSource` write the entries into `t.TempDir()`. A build error is
stored once and reported by every caller as
`build Assets Sync template: <error>`. `removeAssetsSyncTemplate` and its
`TestMain` call go.

### Data Models

None. No schema, Run record or file format changes.

### API Contracts

None. No command, flag, output or exit code changes; every change is test
code.

### Surface Transcripts

None. No command surface changes.

## Coverage Map

- Goal 1 → The linger wait; The Run Budget tests; The script fixture and its
  call sites; The owner watch and reaping; The Assets Sync template.
- Goal 2 → The script fixture and its call sites; The residual inventory.
- Goal 3 → The owner watch and reaping.
- Goal 4 → The Assets Sync template.
- Core Feature 1 → The linger wait.
- Core Feature 2 → The Run Budget tests.
- Core Feature 3 → The script fixture and its call sites; The residual
  inventory.
- Core Feature 4 → The owner watch and reaping.
- Core Feature 5 → The Assets Sync template.
- Success Metric 1 → Testing Approach 1.
- Success Metric 2 → Testing Approach 1.
- Success Metric 3 → Testing Approach 2.
- Success Metric 4 → Testing Approach 3.
- Success Metric 5 → Testing Approach 4.

## Integration Points

- **The Run Budget watchdog** (`internal/daemon/task_engine.go`) is exercised
  unchanged, through its real timer and the injected `Now`.
- **Git** executes hook fixtures through the symlink; nothing else changes in
  how tests drive Git.
- **The suite guard** (ADR-0126) keeps fingerprinting each package. The death
  tests start their inner test binary with a `TMPDIR` inside the outer test's
  temporary directory, so nothing reaches the repository.

## Testing Approach

1. **Wall-clock waits.** The edited linger subtest and both Run Budget tests,
   repeated with `-count` and `-cpu 1,4`. The QA gate repeats the authoring
   stress shape: 8 workers × 10 rounds × `-count=50` for the linger subtest,
   and 8 × 6 × `-count=20` beside ten CPU burners for the budget pair. The
   file holds no `time.Sleep` or 2-second deadline, and no Run Budget below
   one hour.
2. **Written executables.** `TestScriptFixtureIsTheCompiledTestBinary` in
   `internal/cli` asserts that each fixture helper (adapter version, fake
   ACPX, macro adapters and hooks) leaves a symlink to the absolute test
   binary and a mode-0o600 sidecar, and that the fixture runs and prints its
   output, both by absolute path and through `PATH`.
   `TestNoTestWritesAnExecutableOutsideTheResidue` in `internal/testfixture`
   passes on the tree and fails on a temporary tree whose test file writes
   `0o755`, naming `file:line`. It also fails when an inventory count is off
   by one. The adapter floor tests repeat under `-count` and `-cpu 1,4`.
3. **Fixture processes.** `TestImplementDetachChildEndsWhenItsTestBinaryDies`
   re-executes the test binary for
   `TestRunImplementDetachSurvivesCallerProcessGroupKill`, with
   `ROUNDFIX_DETACH_TEST_OWNER_RECORD` and a private `TMPDIR`. It reads the
   recorded Owner PID, sends `SIGKILL` to the inner binary's process group,
   and polls until the Owner PID and its process group are gone.
   `TestDetachSurvivorEndsWhenItsTestBinaryDies` does the same for
   `TestDetachedChildIsTerminatedAtTeardown`, reading the survivor's PID file
   under the private `TMPDIR`. Both fail at the deadline on the old shape.
4. **The Assets Sync template.**
   `TestAssetsSyncTemplateLeavesNoDirectoryBehind` forces the build. It then
   asserts that the recorded build directory does not exist, that a target
   copy passes `git fsck` and `git rev-parse HEAD`, and that a source copy's
   `HEAD` equals the recorded revision. The Assets Sync tests repeat under
   `-count=3 -cpu 1,4`.

## Build Order

1. The linger wait and the Run Budget tests, task_01 (depends on: none).
2. The script fixture, its nine call sites and the written-executable guard,
   task_02 (depends on: none).
3. The owner watch, the reaping and the two death tests, task_03 (depends on:
   2, because both edit `internal/cli/implement_test.go` and the fake ACPX
   becomes a script fixture in step 2).
4. The in-memory Assets Sync template, task_04 (depends on: none).
5. Terminal QA, task_05 (depends on: 1, 2, 3, 4).

## Risks & Considerations

- **Fixture cost.** A script fixture adds one re-execution of the test binary
  per call, measured at about 17 ms against 8 ms for `/bin/sh` under load.
  The fake ACPX runs a few times per Run, so the cost is small. A large rise
  in `internal/cli` wall time is a reason to stop and report.
- **Linux is where ETXTBSY lives.** Darwin cannot reproduce the race. The
  structural guard and fixture test are the local proof, and the CI
  Verification gate is the Linux one.
- **The guard is a ratchet.** A Spec authored elsewhere that adds a written
  executable fails the guard. That is intended, and the fix is to use a
  script fixture or record a reasoned residual.
- **PID reuse.** The owner watch could see a recycled PID and keep waiting.
  The window is the interval between a PID's death and its reuse, which the
  death tests measure only as absence. The detach reaping polls absence, so
  reuse delays only the poll.
- **Production retries.** When its fake ACPX exits 1, the orphaned detached
  child walks the Fallback Chain before the Run fails. The death test bounds
  this by the test deadline, not by a constant.

## Decisions

- Wait on the event the assertion reads, never on a proxy that changes
  earlier.
- Decide budget expiry on the injected clock, and order the renewal after the
  stalled start. A production timer seam was rejected, because the watchdog
  has no race.
- Script fixtures re-execute the compiled test binary, as ADR-0125 requires.
  A `syscall.ForkLock` write was rejected, because ADR-0125 already decides
  the class.
- Fixture processes watch their test binary and are reaped by recorded PID.
  See ADR-0213.
- Hold the Assets Sync template in memory and keep its single build, rather
  than rebuilding it per test.
- Freeze the residue in a guard rather than widen the Spec into governed files
  and four more packages.
