# Test suite audit, 2026-10-24

The maintainer asked on 2026-10-23, at about 23:55 BRT, for faster and more
dependable tests. This audit measured the whole Go suite at `56dd1a76`
(v0.116.0) on the maintainer's 10-core Mac mini (Go 1.26.5, macOS 27), in
scratch clones and worktrees, while Delivery Queues of other repositories ran
on the same machine. Load averages are given with each number, because the
suite is bound by process starts and the machine was shared. Times are BRT.

## Flakes that cost Runs: macOS killed the fake acpx child

### What failed

- 2026-10-23 18:14 BRT, Run `run_20261023T205240Z_cb3296cbba44e77f` of Spec
  0304, QA precondition `make verify-changed`:
  `TestACPXExitCodeMapping/agent_protocol_failure` (`internal/agent`) failed
  reading `<t.TempDir()>/001/args.json`, "no such file or directory".
- 2026-10-23 23:10 BRT, Run `run_20261024T020951Z_493078906e450701` of Spec
  0306, same precondition: `TestRunPromptIgnoresAMalformedUsagePayload`
  failed the same way.

Both tests start the test binary again as a fake acpx (`ROUNDFIX_FAKE_ACPX=1`)
through `os.Args[0]`; the fake writes `args.json` first. Neither failure
reproduced in twenty sequential runs on main.

### What the evidence shows

1. `runFakeACPXPrompt` reads `args.json` before it returns the `RunPrompt`
   error, so the logs never showed why the child wrote nothing. The
   operator's note that `RunPrompt` returned no error was an inference from
   that order, not an observation.
2. Twelve concurrent runs of the compiled `internal/agent` test binary,
   `-test.count=3 -test.parallel 16`, failed in four of eight rounds (load
   averages 30 to 60). The failures moved between tests that start the fake
   acpx (sealed prompt tests, profile proof, warm session, exit
   classification), always with an acpx exit code of -1.
3. With a probe that reports the child's exit state, the child had died of
   `signal: killed` (`SIGKILL`). A marker written as the first statement of
   `TestMain` was absent: the child died before any test code ran.
4. A Go standard library overlay that logged every `SIGKILL` sent through
   `os.Process` in all twelve processes showed none aimed at the killed PID.
5. The macOS unified log (`/usr/bin/log show`; `log` alone is a zsh builtin)
   named the sender for every lost child, for example:
   - 2026-10-24 00:17:52.047 `syspolicyd`: "Could not create URL for pid:
     10102 process path:  library path: (null)" and "Terminating process due
     to Malware rejection: 10102,", PID 10102 being the killed child;
   - 2026-10-24 00:12:16.120: "Unable (errno: 2) to read file at <private>
     for pid: 80567" and "Terminating process due to Malware rejection:
     80567", the probe's killed child of that round;
   - 2026-10-23 18:14:41.768: "Terminating process due to Malware
     rejection: 90369", 79 s before the Spec 0304 Run wrote its failing log;
   - 2026-10-23 23:10:10.974: "Unable (errno: 2) to read file at <private>
     for pid: 14242" and "Terminating process due to Malware rejection:
     14242", one second before the Spec 0306 Run's failure line.
6. `syspolicyd` checks each new process of an ad-hoc signed binary, which
   every Go test binary is. It resolves the executable's path from the
   process's file; `errno: 2` is ENOENT. The `internal/agent` `TestMain`
   hard-links the test binary to three paths under a temporary
   `node_modules` tree and removes them when the process ends, and the
   re-executed children of `TestSpawnedFixtureDeathFailsFast` (`wait` and
   `fixture`) do the same while their parent starts other children. A path
   resolved to a link that has just been removed no longer exists.
7. A standalone program in a scratch directory started itself as a child,
   eight at a time, for 60 s. With four goroutines creating and removing hard
   links to its executable it lost 48 of 101,438 children to `SIGKILL`; the
   log held exactly 48 `syspolicyd` rejections in that window. Without the
   links it lost none of 115,022.
8. The `fixture` child ends through `os.Exit(23)` inside its test, so it
   never removes its links. On 2026-10-24, 300 roots named
   `roundfix-agent-package-adapters-*` sat under `$TMPDIR`; their links
   pinned 26 deleted test binaries, 184 MiB, 20 of them from Runs of
   2026-10-23.
9. With the package adapters copied instead of hard-linked and the fixture
   children making none (the prototype of Spec 0311), eight rounds of the
   same harness, 96 runs at load averages from 143 to 171, all passed, and no
   `syspolicyd` rejection fell in their windows.

Published reports agree: a build-system maintainer wrote on the Apple
Developer Forums (https://developer.apple.com/forums/thread/663456) that "in
case of 'File not found' macOS kills some other running process sharing the
inode with the already removed file" and avoided it "by using clonefile(2)
instead of hardlink"; the Flutter tool lost children after deleting its own
running binary (https://github.com/christopherfujino/catalina-crasher-demo).

A child of an `internal/cli` test killed by an unknown signal on 2026-10-09
(queue intervention 250) fits the same mechanism, but the unified log no
longer covers that date, so it stays unproven.

### Earlier flakes, from the operator's queue interventions

Each was fixed or explained by its own Spec; listed so a recurrence is
recognised: `TestDoctorAcceptsAnAdapterAtTheFloor/codex` (CI, 2026-10-01);
`TestImplementDetachChildEndsWhenItsTestBinaryDies` and
`TestRunImplementDetachSurvivesCallerProcessGroupKill` (EPERM on a process
group, 2026-10-02 and 2026-10-05); `TestTaskCycleVerificationCapacityCancellationWhileQueued`
(a capacity grant that ignored cancellation, Spec 0269);
`TestBootstrapSerializesAcrossSiblings` (release workflow, 2026-10-14);
`TestDeliveryStepNamesTheSignalThatKilledTheItemBinaryProbe` ("text file
busy" on Linux, go.dev/issue/22315); the Real-Time Budget class (ADR-0271).

## Suite cost

### Whole suite

`go test -count=1 -json -parallel 16 -timeout 20m ./...` in a clone, after
compiling every test binary:

| run | wall | load (1 min) at the end |
| --- | --- | --- |
| uncached, 2026-10-24 00:30 | 284 s | about 45 |
| uncached, 00:40 | 299 s | 16.6 |

43 packages, 5,195 top-level tests; the package times add up to 1,169 s, and
the top-level test times to about 11,900 s, because tests overlap inside a
package (`-parallel 16`) and packages overlap across the machine.

### Slowest packages (uncached, 00:30 run)

| package | wall | tests | sum of test times | Git starts |
| --- | --- | --- | --- | --- |
| internal/cli | 283.8 s | 1,714 | 4,351 s | 60,669 |
| internal/baseline | 164.3 s | 511 | 1,742 s | 10,549 |
| internal/daemon | 148.5 s | 437 | 1,664 s | 12,685 |
| internal/worktree | 130.4 s | 196 | 1,576 s | 14,725 |
| internal/store | 104.0 s | 240 | 226 s | 124 |
| internal/spec | 88.1 s | 380 | 961 s | 1,611 |
| internal/speccheck | 78.4 s | 429 | 886 s | 5,593 |
| internal/verifyselect | 62.6 s | 32 | — | 139 |
| internal/delivery | 20.0 s | 211 | — | 106 |
| internal/specaudit | 19.4 s | 18 | 120 s | 657 |

"Git starts" counts the opens of the `git` executable that Go's test log
records for one package run. The suite starts about 106,000 Git processes;
`internal/cli` alone also starts 597 shells. The earlier campaign's reading
holds: the suite is bound by process starts, not by computation.

### Slowest tests (uncached, 00:30 run)

| test | package | time |
| --- | --- | --- |
| TestReconcileItemBranchSafety | internal/cli | 137.9 s |
| TestJournalConsumerCorpusReplaysEveryConsumer | internal/store | 101.8 s |
| TestReviewValidatorFailureLeavesEveryAskedFindingStanding | internal/cli | 101.3 s |
| TestAnAssetSyncWithNoSourceChangeLeavesTheComposedSetup | internal/baseline | 100.6 s |
| TestQAReportAcceptArchiveAndSettleAgreeOnEveryPartialShape | internal/cli | 91.6 s |
| TestItemBranchSafetyBoundaries | internal/worktree | 88.1 s |
| TestAnAssetSyncRewritesTheComposedSetupWithItsComponents | internal/baseline | 83.2 s |
| TestAssetSyncStillRefusesAnInvalidCatalogWithoutDrift | internal/baseline | 80.5 s |
| TestQAReportReaderAgreesWithTheDerivedVerificationOnTheArchive | internal/spec | 78.9 s |
| TestTheDaemonSettlesEveryPartialShapeAsTheCommandsDo | internal/daemon | 77.3 s |
| TestReconcileKeepsTheItemBranchOfALiveQueueItem | internal/cli | 75.2 s |
| TestQAGateLegacyArchivedManifestsLoadUnchanged | internal/spec | 70.5 s |
| TestReconcileListsAndReleasesAnItemBranchOfAMergedSpec | internal/cli | 66.0 s |
| TestRepositorySpecCorpusStillLoads | internal/spec | 63.1 s |
| TestAuthoritativeGateReportsFailure | internal/spec | 62.6 s |

A parallel test's time runs from its release to its end, so it includes the
time it waited for Git and the CPU on a loaded machine. The longest test of a
package bounds the package from below: `internal/cli` cannot finish in less
than its 138 s reconcile matrix.

### Redundant work

- The `ci-verify` job runs `internal/docscontract` twice per pull request,
  once in "Verify changed" (its tests are `//verify:always`) and once in
  "Verify docs" (`docs-test`); `spec-budget` then runs
  `TestCheckCorpusBudget` a third time with `-parallel=1`. Backlog Entry
  `docs/backlog/2026-10-22-the-ci-verify-job-runs-the-documentation-contract-package-twice.md`.
- `make verify-changed` builds `bin/roundfix`, and `make verify-docs` in the
  same CI job builds it again; the second build is a cache hit.
- `internal/spec` holds four whole-archive corpus tests,
  `TestQAGateLegacyArchivedManifestsLoadUnchanged`,
  `TestRepositorySpecCorpusStillLoads`,
  `TestArchivedPassCorpusRemainsArchiveEligible` and
  `TestArchivedQAOverrideCorpusIncludesFailedSpec`, 70.5 s, 63.1 s, 41.8 s
  and 40.8 s, each reading the archive on its own, the shape Spec 0305 fixed
  in `internal/docscontract` by sharing one sweep.

### Test cache and hidden inputs

Go caches a package's test result when nothing it read changed: "Tests that
open files within the package's module or that consult environment variables
only match future runs in which the files and environment variables are
unchanged" (`go help test`). The setting `GODEBUG` is itself part of the key,
so two runs compared here always share the same environment.

| case | packages cached | wall |
| --- | --- | --- |
| Run-like worktree, no change, second run | 40 of 40 | 15 s |
| clone being edited (planning files written during the run) | 30 of 40 | 275 s |

`internal/cli` alone, `go test -parallel 16 ./internal/cli` in a worktree,
each pair of runs in the same environment:

| change between the two runs | all tests | without the 13 tests below |
| --- | --- | --- |
| one Task file's `status:` line | cached | cached |
| one line appended to `docs/user-guide/commands/gc.md` | ran again, 162.6 s and 237.0 s | cached |

A first measurement on the whole suite after a Task-file change also showed
`internal/cli` running again (204 s); it was not reproduced in the controlled
pairs above and is not relied on. The recorded inputs of the packages that
re-run after a change inside the module:

| package | inputs in the module | where |
| --- | --- | --- |
| internal/cli | about 650 | `.agents/skills` 279, `docs/user-guide` 36, and source and `testdata` paths of other packages |
| internal/spec | 793 | source files of other packages |
| internal/testfixture | 1,469 | every test file (its rules scan them) |
| internal/verifyselect | 911 | the package tree |
| internal/suiteguardcontract | 905 | the package tree |

With `docs/history`, `docs/adr`, `docs/specs`, `docs/user-guide` and
`.agents/skills` removed, and again with them unreadable, the same 13
`internal/cli` tests fail and no other:
`TestARepositoryMissingARequiredExternalSkillIsReported`,
`TestBaselineExamplesParse`,
`TestBaselineUpdatePreviewListsADifferenceWithARepositoryProfile`,
`TestBaselineUpdatePreviewListsASetupSnapshotDifference`,
`TestBaselineUpdateReachesCurrentWithARepositoryProfile`,
`TestBaselineUpdateSkillsStageRefreshesSkillsWithARepositoryProfile`,
`TestDoctorComparesARepositoryProfileWithItsSnapshot`,
`TestEventsHelpDocumentsAgentSelectionFilter`,
`TestEveryStatementOfTheHTTPDefaultNamesTheCatalogDefault`,
`TestProjectConstraintPlanningExamplesParse`,
`TestReleasePlanChecksOnAnAdoptedRepository`,
`TestReleasePlanChecksWriteNothing` and
`TestThisRepositoryHoldsEveryRequiredExternalSkill`; together about 150 s of
test time. Every Task that edits a guide or a skill, which the repository's
own rules make common, therefore re-runs the whole package, the one on the
critical path.

A Run's worktree also has its own path, and Go builds record absolute source
paths, so the first Verification of every Run compiles and tests from
nothing. `-trimpath` would remove the path from the build, but ten test files
locate the repository through `runtime.Caller` (for example
`internal/spec/spec_test.go` and `internal/docscontract/corpus_test.go`), so
it is not a flag-only change.

On this host `core.fsmonitor` is set to `true` in the user's Git
configuration, so any Git command in a checkout without
`-c core.fsmonitor=false` makes the monitor daemon write under `.git`. In a
plain clone that changes the `.git` directory that `internal/cli` tests stat;
in a Run worktree `.git` is a file and does not change.

## Disk

- The shared build cache `roundfix-go-build-501` under `$TMPDIR` held 31 GB
  on 2026-10-24 (43.7 GB on 2026-10-22, when it filled the disk); Go's
  default cache held another 10 GB. Each Run worktree compiles every package
  under its own path, so entries do not repeat across Runs, and Go deletes an
  entry only after five days without use. `roundfix gc` does not reclaim it.
- The leaked package adapter roots above held 184 MiB.

## Times outside this machine

From the operator, 2026-10-23 and 2026-10-24: the CI Verification gate job of
PR #654, after Spec 0305, took 463 s (Verify changed 271 s, Verify docs
173 s), against about 25 to 32 minutes before 0305; inside a Run,
`make verify-changed` takes about 15 minutes of a Task's or QA's wall clock
(the Spec 0306 Run's log of 2026-10-23 23:11 BRT records `internal/daemon`
75.8 s, `internal/spec` 72.5 s, `internal/store` 69.9 s,
`internal/verifyselect` 68.3 s and `internal/worktree` 63.4 s).

## Targets

| measure | from | target | owner |
| --- | --- | --- | --- |
| loaded `internal/agent` runs that fail (12 × 3, 8 rounds) | 4 of 8 rounds | 0 of 96 runs | Spec 0311 |
| hard links created by tests | 1 site | 0, enforced | Spec 0311 |
| CI Verification gate of a pull request | 463 s | at most 300 s | Spec 0312 and later work |
| `internal/cli` after a guide or skill edit | runs again, 160 to 240 s | cached | Spec 0313 |
| shared Go build cache | 31 to 43.7 GB | trimmed by `roundfix gc` above a declared limit | Spec 0314 |
| `make verify-changed` inside a Run | about 15 min | at most 10 min | Specs 0312 to 0314 and later work |

## Follow-up the audit records

- `internal/cli` tests that read the live documentation trees: Spec 0313
  (the maintainer's answer of 2026-10-24, "Escrever agora").
- The shared Go build cache: Spec 0314 (the maintainer's answer of
  2026-10-24, "roundfix gc recolhe").
- The four `internal/spec` archive corpus tests could share one sweep.
- The slowest Git-bound tests (reconcile, asset sync, item branch safety)
  build their repositories one commit at a time.
