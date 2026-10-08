---
spec: 0254-a-faster-make-test
prd: _prd.md
created: 2026-10-08
---

# A faster make test — Technical Spec

## Executive Summary

The suite is slow because it waits, not because it computes. Go runs every
top-level test that never calls `t.Parallel()` alone, before the package's
parallel tests start, and `internal/cli` holds 592 of them. This Spec adds a
repository test that makes a Parallel Test Package's tests parallel by
default, converts `internal/cli` and six more packages, removes the slowest
sequential tests that only needed a per-test working directory, narrows two
pinned-history fixtures, and runs both contract tags in one `go test`. An
authoring prototype in a scratch clone converted the seven packages and
measured the suite at 59 % of the unchanged tree's wall time on the same
machine and load.

The trade-off accepted is that a converted test can now race another one on
state Go does not guard. Go refuses `t.Setenv` and `t.Chdir` in a parallel
test, but the prototype still met two tests that swap `app.BuildCommit`, and
a search found one more that swaps `app.Version`. Each converting Task therefore runs its packages
once under the race detector and once with `-shuffle=on`. The race detector
stays out of `make test`. The other trade-off is scope: after this Spec the
suite is bound by the work it does, at least 83,660 Git processes per run, so
it does not reach 60 s.

## Project Constraints

- Identifier strategy: not applicable. No identifier, schema or file format
  changes. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable. No request, credential or network
  call is added; tests keep their temporary repositories, homes and fake
  runners. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0259 (this Spec) governs every
  decision below. ADR-0089: "it keeps the global and states in one line why it
  stays sequential". ADR-0125: "Fixtures are therefore compiled once", so no
  conversion makes a test execute a file it wrote. ADR-0126: "The guard
  therefore runs per package", so every converted package keeps
  `suiteguard.Main`. ADR-0213: "Every fixture that waits therefore watches its
  owning test binary's process and exits when that process is gone". ADR-0090:
  "Caching facts across mutation boundaries was rejected"; pinned history is a
  fixed commit. ADR-0252: "Each Repository Contract Test declares when the
  selective gate runs it", and ADR-0253: "CI runs the Full Contract Run on
  every push to main"; both stay, and only the invocation's shape changes.
  ADR-0257: "the Agent runs focused tests of the packages it changed".
  ADR-0244: "A Spec declares the domain terms it introduces". ADR-0170 and
  ADR-0258 do not apply: this Spec changes neither Delivery Retry nor how
  history is written. Source: `docs/agents/domain.md`.
- Tooling authority: applicable. No Makefile, CI workflow, lint, formatter or
  test-runner configuration changes. Express maintainer authorization: the
  2026-10-08 grant of the governed paths this Spec declares. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.
  Spec-contained authorization record:
  `docs/specs/0254-a-faster-make-test/_authorization.md`. Bounded files:
  `docs/agents/specific-repository.md`,
  `internal/baseline/derived_ownership_test.go`,
  `internal/cli/cli_test.go`, `internal/spec/coverage_test.go`,
  `internal/speccheck/mechanical_test.go`.

## Current behavior

Measured on 2026-10-08 on the maintainer's 10-core machine at `2ce5abe8`, with
the worktree's own `GOCACHE`, by `go test -count=1 -parallel 16 -json ./...`.
A Delivery Queue Run of Spec 0253 was active the whole time, so the load
average stayed between 8 and 25; every comparison below is a back-to-back pair
under that load.

| Run | Wall | user + sys CPU | `internal/cli` | its sequential phase |
| --- | --- | --- | --- | --- |
| full suite, 16:54 | 374.8 s | 426 s + 688 s | 373.7 s | 325.6 s, 592 tests |
| `internal/cli` alone, 17:02 | 342.1 s | 196 s + 284 s | 341.4 s | 279.4 s |
| full suite, 17:31 (pair A) | 384.2 s | 400 s + 665 s | 383.1 s | 323.9 s |

The sequential phase is the summed duration of the top-level tests that never
paused, read from the `pause` events of the JSON stream. In the 17:31 run the
other slow packages were `internal/daemon` 183.6 s (sequential 144.8 s),
`internal/speccheck` 147.6 s (138.4 s), `internal/spec` 134.1 s (119.0 s),
`internal/baseline` 125.6 s (21.9 s) and `internal/worktree` 118.0 s (46.2 s).

- The slowest single tests were sequential ones that run fast alone:
  `TestQAReportReaderAgreesWithTheDerivedVerificationOnTheArchive` took 59.0 s
  in the suite and 16.1 s alone, because it materializes 6,643 pinned files to
  read 759 QA Reports; `TestEvidenceRecordWritesOneLinePerInputWhateverItMatches`
  took 48.9 s in the suite and 2.6 s alone; the ten serial subtests of
  `TestPriorQAPassSkipNeverReportsPathDiffers` took 22.0 s.
- A `PATH` shim counted at least 83,660 `git` processes in one full run, at
  least 17,342 of them from test helpers (about 2,000 `git init`), and 388
  `python3` processes. Seven test sites build `./cmd/roundfix`; a warm build
  took 3.6 s to 4.2 s. Fixed sleeps are not a cost: the seven packages hold
  five `time.Sleep` calls, all in detach fixture children.
- Of 4,603 top-level tests, the JSON stream shows 1,042 sequential in the
  seven packages: 592 in `internal/cli`, 133 in `internal/speccheck`, 108 in
  `internal/baseline`, 86 in `internal/daemon`, 76 in `internal/spec`, 25 in
  `internal/store` and 22 in `internal/worktree`. Thirty-one `// Sequential:`
  comments already exist.
- The Full Contract Run's two invocations, run one after the other, took 109 s
  and 91 s; one `go test -tags docscontract,repocontract` with the union of
  names and packages took 62 s and 63 s in the same rounds.
- GitHub's runners measured the suite at 185 s, 203 s and 223 s of a 240 s
  budget (runs 37834066256, 37835405144 and 37821753645).

The authoring prototype, a scratch clone at `2ce5abe8`, inserted `t.Parallel()`
into every top-level test of the seven packages that reached no `Setenv`,
`Chdir`, nested `Parallel` or known global swap through same-package helpers,
and fixed what failed. It found the hidden cases this Spec's requirements
name: `setOwnerBuildCommit` swaps `app.BuildCommit`,
`TestSupersedeAcceptsAnArchiveRecord` calls another test that is already
parallel, and `testUnnamedTaskStaysBlockedByRedPrecondition` calls `t.Parallel()`
itself. After those fixes every package passed. A search for assignments to
another package's variables then found
`TestPreconditionRefusalReportNamesItsAuditor`, which swaps `app.Version`,
`app.BuildCommit` and `app.BuildTime`; the prototype had converted it, and it
passed only because no other test read them at that moment.

| Prototype run (pair A, same load) | Wall | `internal/cli` | its sequential phase |
| --- | --- | --- | --- |
| full suite, 17:28 | 227.4 s | 226.3 s | 123.8 s, 43 tests |
| `internal/cli` alone, 17:25 | 132.2 s | 129.3 s | 15.8 s |

Under load the 43 tests left sequential in `internal/cli` still cost 124 s.
The six repository-profile tests, which call `t.Chdir` only so that
`baseline profile init` finds the repository, cost 78 s of that, and the
command environment already carries a per-test working directory
(`setCommandWorkDirForTest`). Narrowing two pinned-history fixtures took
`internal/spec` alone from 15.1 s to 12.6 s and 13.3 s; narrowing two more
failed their tests, because they read whole archived Spec folders.

## System Architecture

- **Parallel Test Package rule** — a new repository test in
  `internal/testfixture`, beside the written-executable inventory it
  resembles. It walks the test files of each listed package with `go/ast`,
  skips files under a `docscontract` or `repocontract` build constraint, and
  checks each top-level test.
- **Conversions** — test-only edits in the seven packages: `t.Parallel()` as
  the first statement, or a `// Sequential:` reason.
- **Repository rule** — one bullet in `docs/agents/specific-repository.md`,
  outside setup markers, tells whoever writes a test in a Parallel Test
  Package to call `t.Parallel()` first or name why not.
- **Residue** — `internal/cli` tests that stay sequential only for the working
  directory take `setCommandWorkDirForTest` instead of `t.Chdir`, and two
  `internal/spec` tests pass a narrower pathspec to `gittest.PinnedHistory`.
- **Contract invocation** — `verifyselect.ContractInvocations` prints one line
  for both tags. The Makefile loop that reads `tag tests packages` is
  unchanged, because `-tags` accepts a comma-separated list.

No production Go file changes except `internal/verifyselect/contracts.go`,
which is repository tooling.

## Implementation Design

### Interfaces

```go
// internal/testfixture/parallel_tests_test.go
var parallelTestPackages = map[string]int{ // package directory → Sequential Test ceiling
	"internal/cli": 48,
}

type sequentialTest struct{ file, name, reason string; line int }

func parallelTestViolations(root string, ceilings map[string]int) ([]string, []sequentialTest, error)
func TestEveryTestInAParallelTestPackageRunsInParallel(t *testing.T)

// internal/verifyselect/contracts.go
func ContractInvocations(selected []ContractTest) []string // at most one line
```

### Data Models

None. The ceilings are a map literal in the repository test.

### API Contracts

1. API Contract: the Parallel Test Package rule. For each listed package
   directory, every `*_test.go` file without a contract build constraint is
   parsed. A top-level function named `Test…` with one `*testing.T` parameter,
   other than `TestMain`, is **parallel** when its first statement is
   `<param>.Parallel()`. It is a **Sequential Test** when it is not parallel
   and its doc comment, or a comment between its opening brace and its first
   statement, has a line beginning `Sequential:` followed by a reason of at
   least two words. A violation, reported as `<file>:<line>: <Test>` with the
   reason, is a test that is neither, a test that is both, or a package whose
   Sequential Tests exceed its ceiling. The test fails listing every violation
   and logs each package's count against its ceiling.
2. API Contract: ceilings. `internal/cli` 48 after task_01 and 40 after
   task_03; `internal/daemon` 16, `internal/baseline` 12, `internal/store` 3,
   `internal/spec` 2, `internal/speccheck` 2 and `internal/worktree` 2 after
   task_02. A ceiling may only fall. Counted by this rule, where only a first
   statement makes a test parallel, the prototype left 45 in `internal/cli`
   (39 after the six repository-profile tests), 14 in `internal/daemon`, 10 in
   `internal/baseline`, 2 in `internal/store` and none elsewhere, plus the one
   `internal/spec` test that swaps `app.Version`. Eight of the daemon's 14 call
   a shared helper that calls `t.Parallel()` itself; moving that call into
   each test makes them parallel under the rule. The margins leave room for
   tests Spec 0253 adds.
3. API Contract: `ContractInvocations`. For an empty selection it returns no
   line. Otherwise it returns one line: the sorted, comma-joined tags, the
   anchored alternation of the sorted unique quoted names, and the sorted
   unique `./<package>` arguments. One tag gives the line it gives today.
   `make verify-changed-contracts` and `make verify-contracts` run that line
   unchanged.
4. API Contract: name uniqueness. A repository test in `internal/verifyselect`
   fails when a Repository Contract Test's name is also the name of a
   top-level test in another build-tag class: the other contract tag, or no
   contract tag. Today no name crosses a class; the per-package suite-guard
   contract shares one name within the `repocontract` tag only.

### Surface Transcripts

None. `cmd/verify-select` is repository tooling the Makefile reads, not a
Roundfix command; API Contract 3 states its output and
`TestRunPrintsEveryContractWithAll` asserts it.

### Invariants

1. A conversion never changes what a test asserts. A converted test differs
   from its starting bytes only by the inserted `t.Parallel()`, a
   `// Sequential:` comment, a `t.Parallel()` call moved from a shared helper
   into each test that calls it, or, in task_03, the per-test working
   directory or pathspec that replaces process-wide state.
2. Every converted package passes once with `-race -short` and once with
   `-shuffle=on`, and keeps `suiteguard.Main`.
3. A test that sets an environment variable, the working directory, a
   package-level variable or a signal handler for the whole process stays
   sequential with that reason, unless task_03 replaces the process-wide
   value with a per-test one.
4. No test executes a file it wrote (ADR-0125), and no fixture process
   outlives its test binary (ADR-0213).
5. The merged contract invocation selects exactly the contracts the per-tag
   invocations selected.

## Coverage Map

- Goal 1, Story 3, Core Feature 1, Success Metric 3 → Parallel Test Package
  rule (task_01), extended by task_02.
- Goal 2, Story 4, Success Metric 4 → race and shuffle runs in task_01,
  task_02 and task_03.
- Goal 3, Core Feature 4 → task_03.
- Goal 4 → task_03 pathspecs.
- Goal 5, Story 5, Core Feature 5, Success Metric 5 → `ContractInvocations`
  and the name-uniqueness test (task_04).
- Goal 6, Stories 1–2, Success Metrics 1–2 → the QA gate's back-to-back
  measurement (task_05).
- Core Feature 2 → task_01. Core Feature 3 → task_02. Core Feature 6 →
  task_01.

## Integration Points

None outside the repository. The Makefile's contract loop and CI's
`make verify-contracts` consume `cmd/verify-select` output unchanged in form.

## Testing Approach

- `internal/testfixture/parallel_tests_test.go` (task_01, new) runs the rule
  on the repository and on seeded trees in `t.TempDir()`: a parallel test
  passes; a test without `t.Parallel()` and without a reason fails with its
  `file:line`; a reason of one word fails; a test that has both fails; a
  `Sequential:` reason in the doc comment and one inside the body both pass;
  `t.Parallel()` called later than the first statement fails; a
  `repocontract` file is skipped; a package over its ceiling fails naming its
  count; an unlisted package is ignored.
- Conversion proof (task_01, task_02): each package passes
  `go test -count=1 -race -short` and `go test -count=1 -shuffle=on`. The
  `-short` flag skips only the two `internal/store` journal harnesses, whose
  120 s seeding deadline the race detector's slowdown exceeds; the shuffled
  run still runs them. The prototype's race run of the seven packages took
  389 s under load and reported no data race.
- task_03 changes the six repository-profile tests and adds no test: the rule
  test's lower ceiling and a `go test -run` of the six tests prove it.
- task_04 changes `TestContractInvocationsSortAndDeduplicate`,
  `TestRunPrintsContractInvocations` and `TestRunPrintsEveryContractWithAll`
  to the one-line form, adds the name-uniqueness test to
  `internal/verifyselect/contracts_repository_test.go`, and leaves the
  Makefile wiring tests in `full_contract_run_test.go` and
  `contracts_repository_test.go` unchanged.
- No golden, assertion or fixture meaning changes in any Task.

## Build Order

1. Parallel Test Package rule, glossary terms and the `internal/cli`
   conversion (task_01).
2. The six packages join the rule (task_02) (depends on: 1).
3. The `internal/cli` residue and the pinned-history pathspecs (task_03)
   (depends on: 1, 2).
4. One contract invocation (task_04) (depends on: 1).
5. QA gate (task_05) (depends on: 3, 4).

task_02 and task_03 share the rule's test file with task_01, so they run in
series. task_04 touches only `internal/verifyselect`; it follows task_01 only
so that the converted `internal/cli` is the one its gate runs.

## Risks & Considerations

- A converted test can race on state Go does not guard. The race run catches
  an overlapping write; the requirement to search for assignments to another
  package's variables (`app.`, `config.` and the like) catches one that did
  not overlap in that run.
- A panic in a converted `internal/cli` test, such as `t.Setenv` in a parallel
  test, killed the test binary in the prototype and left a nested detach
  fixture (`TestRunImplementDetachSurvivesCallerProcessGroupKill`) running
  for its own `-test.timeout`. An agent that hits one reaps it by the PID it
  started, never with a broad `pkill`.
- `captureStderr` in `internal/worktree` already swaps `os.Stderr` inside a
  parallel test. This Spec does not change it; task_02 must not convert a
  test that calls it.
- Measurements on this machine depend on concurrent Runs. The QA gate
  measures the starting commit and the audited head back to back and records
  the load average beside each run.
- The 60 s target is not reachable here. On the same machine and load the
  converted suite used 448 s user and 873 s system CPU in 227 s; the system
  time alone is 87 s of ten cores. The next lever is fewer Git processes in
  production code paths (ADR-0090's batching, extended), which is a later
  Spec.

## Glossary

- adds: **Sequential Test**
- adds: **Parallel Test Package**

## Decisions

- The rule is a `go/ast` repository test with per-package ceilings, not a
  linter dependency; see ADR-0259.
- Each conversion runs `-race` and `-shuffle=on` once; `make test` keeps no
  race detector; see ADR-0259.
- Both contract tags run in one `go test`; see ADR-0259.
- The CI suite budget is left at 240 s until CI measures the converted suite;
  see `_prd.md` → Unreachable Acceptance.
