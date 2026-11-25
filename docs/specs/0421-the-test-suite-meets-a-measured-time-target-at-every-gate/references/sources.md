# Sources of Spec 0421

This file copies what the Idea and the PRD rest on: the maintainer's first decision, the measurements with the command and the machine of each, and the published passages they quote. The maintainer's later decisions of 2026-11-25 are quoted in the `## Decisions` sections of the Idea and of the PRD. The raw logs stay outside the repository, on the machine that measured; every number the PRD uses is copied here.

## 1. The maintainer's decision of 2026-11-25

Delivered as a `DECISÃO DE MARCIO` block to this project's ai-memory inbox by the Praetor, the session that coordinates the software projects, as ai-memory message `01a2121b-9124-7571-8f8f-de0f40876364`. It answers the decision request `01a2109e-f31b-74e3-a162-d0f700280f5c`. The block, verbatim:

```text
DECISÃO DE MARCIO
pedido: 01a2109e-f31b-74e3-a162-d0f700280f5c
pergunta: A verificação completa do main e do release falha em macOS porque o pacote internal/cli passa do tempo-limite padrão de 10 min. Opções: -timeout 30m nos alvos test e test-budget; voltar main e release ao Linux; runner macOS maior.
resposta (literal): "Essa verificação do roundfix de 10 minutos já seria um absurdo, ainda mais se rodar em todas as runs, 30 minutos é ainda mais absurdo. É primordial uma revisão dos testes e conjuntos de testes para que tenhamos performance decente."
alcance: NÃO aplicar -timeout 30m nem subir o SUITE_BUDGET_SECONDS. Abrir uma Spec de performance da suíte de testes do Roundfix (o que roda em cada Run, em PR, em main e em release; testes lentos, duplicados, que esperam timeout, que não são herméticos; paralelismo; o que pode sair do gate de cada Run). Autoria da Spec sem Run; mande ao Praetor o plano com uma meta de tempo medida (proposta sua, por gate: Run, PR, main, release) e, separado, o que fazer com o main vermelho e a v0.153.3 enquanto a Spec não sai, como PEDIDO DE DECISÃO com opções. Nada de mudança em Makefile ou workflow sem nova decisão.
quando: 25/11/2026 06:40 BRT
```

## 2. Machines

| Name in this file | Machine | Cores, memory | Notes |
| --- | --- | --- | --- |
| the mini | The maintainer's Mac mini (`Mac16,10`), macOS 27.2.1, Go 1.26.5 | 10 cores, 16 GB | Runs every Run and the Delivery Queue. Roundfix gives each Verification there `GOMAXPROCS=5`, the Verification Core Budget of a 10-core machine with `runs.max_verifying: 1` |
| the macOS runner | GitHub-hosted standard runner, label `macos-26`, arm64 | 3 cores, 7 GB | Verifies pull requests, pushes to main and releases since Spec 0420 merged on 2026-11-25 |
| the Linux runner | GitHub-hosted standard runner, `ubuntu-latest` | 4 cores, 16 GB | Verified them until 2026-11-24; kept here as the comparison |

"Lab" below means the planning worktree of this Spec at `33ece0cf`, on the mini with nothing else running, with `GOMAXPROCS=5` unless stated, a warm build cache and fresh test results (`-count=1`). "Field" means what the Run Database or the workflow runs recorded while the machine did its normal work.

## 3. The gates as they are

| Gate | What it runs | Where |
| --- | --- | --- |
| Run, each Task settlement and the QA gate | the Task's declared Verification, then `make verify-changed` | the mini |
| Delivery Queue, repository gate | `make verify-changed`, read from the same key of the Project Config as the settlement, `defaults.verification` | the mini |
| No local gate | `make verify`, the whole suite | the mini, when someone runs it by hand |
| Before a pull request opens | `make verify-docs` | the mini |
| Pull request | `make verify-changed`, then `make verify-docs` | the macOS runner, one job |
| Push to main | `make verify` with the timed suite, `make verify-docs`, `make verify-contracts`, the budget check | the macOS runner, one job |
| Release | `make verify`, `make verify-contracts`, then the publication job | the macOS runner, one job, then a Linux job that publishes |

The first draft of the PRD said the Delivery Queue's repository gate ran `make verify`. The Idea stage corrected it from the code and the configuration: the gate runs the configured command, and the Project Config of this repository sets it to `make verify-changed`. No log of that gate was found on the mini.

`make verify-changed` formats, vets and builds, then runs one of two fixed sets of packages, or both: every package outside the Baseline set, with `internal/cli` run apart, or the Baseline packages with the Baseline tests of `internal/cli`. A change to one Go file outside the Baseline set therefore runs every package outside it. It then runs the Repository Contract Tests that Contract Relevance selects.

## 4. Field measurements

### 4.1 Run: `make verify-changed` at settlement, on the mini

Source: the Run Database of the mini, table `run_events`, kind `daemon.verification`, the time between the event `started: make verify-changed` and the next `command passed: make verify-changed` of the same Run, for the Runs of this repository from 2026-11-18 01:22 UTC to 2026-11-25 00:50 UTC.

| Passing runs | Minimum | Median | p90 | p95 | Maximum | Sum |
| --- | --- | --- | --- | --- | --- | --- |
| 83 | 436 s | 607 s | 777 s | 882 s | 1,142 s | 14.4 h |

The same command failed 16 times in that window. The declared Verification that runs before it took a median of 15 s and 321 s at p90, over the 80 attempts that reached the command.

Counting the executions that failed as well, the same export holds 99 executions in 43 Runs of this repository, with a median of 583 s and 15.3 hours in all: 69 % of the settlement-check time of the machine with 41 % of its executions. The Runs of the other repositories on the mini ran their check 140 times, with a median of 39 s and 7.0 hours in all. The Idea stage made this reading.

Earlier records of the same command at settlement: ADR-0182, 2026-09-30, "measured at 206 to 223 seconds on this repository"; ADR-0257, 2026-10-08, "350 s on average".

`roundfix runs causes --since 2026-11-01`, read on 2026-11-25, lists 56 items from 125 terminal Spec Runs. Sixteen of them are failures of `make verify-changed` at settlement: nine carry the signature `lint-or-format` and seven `record-or-golden`, all in the class `repository-convention`. None carries `go-test-failure`. The signatures classify the diagnostic text.

### 4.1.1 The sixteen diagnostics

Read on 2026-11-25 for the rewrite of the PRD. Source: the Run Database of the mini, the events of kind `daemon.verification` whose summary ends `failed: make verify-changed` since 2026-11-01, sixteen of them between 2026-11-18 and 2026-11-24, and the diagnostic log each event names. What failed in each execution:

| Executions | What failed | Where the check lives today |
| --- | --- | --- |
| 1 | `gofmt`: one test file needed formatting | the format step |
| 2 | `TestTheSkillCoverageMapIsCurrent`: the record of command surfaces was stale | a Repository Contract Test |
| 2 | `TestEveryRealTimeBudgetIsWithinItsFileCeiling` once and `TestEveryTestInAParallelTestPackageRunsInParallel` once: a rule over the test sources | code tests of the test fixture package |
| 7 | the tests of the selector's package: six executions name `TestRepositoryContractTestsDeclareTheirRelevance`, the rule that every Repository Contract Test declares its Contract Relevance, with its pinned count of contract files, and one names no test in the lines read | code tests of the selector's package |
| 2 | `TestCoverageEquivalence`: the coverage record was stale; one of the two also failed the rule above | a code test of the Spec package |
| 2 | `TestBaselineUpdateSurfaceTranscriptThisRepositoryIsCurrent` and `TestBaselineUpdateSurfaceTranscriptKeepsTheClaudeInstructionLink`: this repository was not current with its Baseline record after a `Makefile` change | code tests of the CLI package |

Fifteen of the sixteen are a check of the repository against its own rules or records, and one names no test. None is a test of product behavior that the Task's change broke. Thirteen of them failed in the tests of a code package, not in a contract, so a settlement that ran only format, vet, build and the contracts would have missed them.

### 4.2 Pull request, push to main and release, on the runners

Source: `gh run view <id> --json jobs`, step start and end times, for the 67 successful runs of the workflow `CI | Verify` among its last 80, read on 2026-11-25. Every run before 2026-11-25 01:00 UTC ran on the Linux runner.

| Event, runner | Runs | Step | Minimum | Median | p90 | Maximum |
| --- | --- | --- | --- | --- | --- | --- |
| pull request, Linux | 35 | `Verify changed` | 250 s | 410 s | 631 s | 701 s |
| pull request, Linux | 35 | `Verify docs` | 144 s | 214 s | 278 s | 279 s |
| pull request, Linux | 35 | whole job | 462 s | 656 s | 860 s | 988 s |
| pull request, macOS | 1 | `Verify changed` | 1,671 s | | | |
| pull request, macOS | 1 | `Verify docs` | 232 s | | | |
| pull request, macOS | 1 | whole job | 1,965 s | | | |
| push to main, Linux | 31 | `Verify` | 219 s | 310 s | 387 s | 406 s |
| push to main, Linux | 31 | `Verify docs` | 181 s | 250 s | 308 s | 325 s |
| push to main, Linux | 31 | `Verify every contract` | 169 s | 239 s | 279 s | 289 s |
| push to main, Linux | 31 | whole job | 603 s | 840 s | 996 s | 1,022 s |

The one pull request on the macOS runner is run `45529968311`, the delivery of Spec 0420, whose change to the `Makefile` selected both sets and every contract. The one push to main on the macOS runner is run `45533963674`, commit `33ece0cf`: the step `Verify` failed after 1,017 s, because `go test` killed the test binary of `internal/cli` at its ten-minute limit (`panic: test timed out after 10m0s`, `FAIL roundfix/internal/cli 603.085s`), and the documentation, contract and budget steps were skipped.

Release, Linux runner, two runs: `45423344461` (v0.153.2) took 464 s in `Verify gate` and 350 s in `Verify every contract`; `45398713637` (v0.153.1) took 466 s and 347 s. No release has run on the macOS runner.

### 4.3 Packages on each runner

Seconds each test binary ran, from the log of the step that ran the whole suite. Linux: push `45510841870`, commit `4a53079d`, 2026-11-24, suite 355 s of the 420 s budget. macOS: push `45533963674`, the failed run above. The third column is the pull request `45529968311` on the macOS runner, where `internal/cli` ran alone after the other packages and without its Baseline tests.

| Package | Linux runner, push | macOS runner, push | macOS runner, pull request |
| --- | --- | --- | --- |
| `internal/cli` | 295 s | 603 s, killed | 562 s, plus 68 s for its Baseline tests |
| `internal/spec` | 157 s | 499 s | 227 s |
| `internal/daemon` | 123 s | 304 s | 175 s |
| `internal/store` | 120 s | 212 s | 102 s |
| `internal/baseline` | 62 s | 189 s | 111 s |
| `internal/speccheck` | 67 s | 117 s | 98 s |
| `internal/worktree` | 57 s | 117 s | 100 s |
| `internal/verifyselect` | 57 s | 53 s | 187 s |
| `internal/delivery` | 22 s | 52 s | 14 s |
| every package, summed | about 1,090 s | more than 2,230 s | |

When the test binary of `internal/cli` was killed, two of its tests had been running for minutes: `TestAgentSelectionProfilesMacro/mixed_profiles_configure_validate_fallback_persist_and_stream` for 4 min 20 s, which builds the product binary and a stub inside the test, and `TestPhaseModelKeepsTheFallbackForItsSpecAndModel` for 5 min 17 s, which was waiting for its parallel subtests to be given a turn.

Repository Contract Tests on the macOS runner, pull request `45529968311`: in `make verify-changed`, `internal/baseline` 372 s, `internal/config` 341 s, `internal/verifyselect` 201 s, `skills` 66 s, in one `go test` invocation; then again in `make verify-docs`, `internal/baseline` 214 s. On the Linux runner, release `45423344461`: `internal/baseline` 308 s, `internal/config` 283 s, `internal/verifyselect` 142 s, `skills` 137 s.

## 5. Lab measurements on the mini, 2026-11-25

Planning worktree at `33ece0cf`. Build cache warm. Commands run from the repository root.

| Measurement | Command | Result |
| --- | --- | --- |
| Whole suite, default cores | `go test -count=1 -json -parallel 16 ./...` | 391 s; 6,131 top-level tests, 16,056 leaf tests, none failed, 3 skipped; load average 40 to 60 |
| Compile only, warm cache | `go test -count=1 -run '^$' ./...` | 15 s |
| `make verify-changed`, no change | `GOMAXPROCS=5 GOFLAGS=-count=1 make verify-changed` | 17 s |
| `make verify-changed`, one comment line added to one file of `internal/queueforecast` | the same | 251 s: the packages outside `internal/cli` ended within 80 s (`internal/daemon` 80 s, `internal/spec` 47 s, `internal/speccheck` 44 s, `internal/worktree` 39 s), then `internal/cli` alone 159 s |
| `make verify` | `GOMAXPROCS=5 GOFLAGS=-count=1 make verify` | 376 s; `internal/cli` 372 s, `internal/baseline` 197 s, `internal/daemon` 122 s, `internal/worktree` 78 s, `internal/spec` 71 s |
| `make verify-docs` | `GOMAXPROCS=5 make verify-docs` | 77 s; `internal/baseline` 70 s |
| `make verify-contracts` | `GOMAXPROCS=5 make verify-contracts` | 81 s; `internal/config` 78 s, `internal/baseline` 73 s |
| One Git process | 500 runs of `git --version` from a shell loop | 1.89 s, 3.8 ms each |
| Whole suite in an environment like the runner's | `env -i` with an empty `HOME`, a `PATH` of the Go tool, Homebrew and the system directories, no `CODEX_PATH`, `GOMAXPROCS=5`, then `go test -count=1 -json -parallel 16 ./...` | 391 s; every test passed; one subtest more was skipped than in the normal run, `TestGuidanceCompositionJourney/rust-cli`, which names `rustfmt` |
| Whole suite, eight tests at once in a package | `GOMAXPROCS=5 go test -count=1 -parallel 8 ./...` | 395 s, against 391 s at 16 |
| Git processes of one run | the same suite with a script named `git` first on `PATH` that appends its arguments to a log and runs the real `git` | 153,850 calls in 435 s, every test passing. 35,875 carry the identity arguments the test fixture passes; 117,975 do not. By subcommand: `rev-parse` 36,787, `ls-tree` 15,384, `cat-file` 12,974, `add` 9,908, `commit` 9,846, `show` 7,964, `symbolic-ref` 5,761, `init` 4,682. A count of calls that reached `git` through `PATH`, so a lower bound |
| A checkout path the build cache never saw | a copy of the tree at a new path, then `go vet ./...` and `go test -count=1 -run '^$' ./...` under `GOMAXPROCS=5` | vet 5 s, every test binary compiled and started in 11 s; at the known path 1 s and 4 s. A fresh Task Worktree therefore costs about 15 s of compilation, not minutes |

### 5.1 Where the time of the whole suite goes

From the `-json` output of the whole-suite run. "Test seconds" is the sum of the time each leaf test ran, which excludes the time a parallel test waits for its turn. The machine was saturated, so these seconds are larger than the same tests take alone; they rank the tests.

| Package | Top-level tests | Test seconds | Wall time of its binary |
| --- | --- | --- | --- |
| `internal/cli` | 2,102 | 6,091 | 390 s |
| `internal/baseline` | 605 | 2,830 | 197 s |
| `internal/daemon` | 493 | 2,271 | 179 s |
| `internal/worktree` | 208 | 1,861 | 152 s |
| `internal/speccheck` | 456 | 1,447 | 101 s |
| `internal/spec` | 417 | 1,377 | 168 s |
| `internal/store` | 256 | 295 | 111 s |
| every package | 6,131 | 17,008 | 391 s for the suite |

The time is spread: the 10 costliest top-level tests hold 8 % of the test seconds, the 100 costliest 30 %, the 400 costliest 55 %. 343 top-level tests run for 10 s or more and 4,062 for less than a second.

ADR-0259 measured 1,409 top-level tests in `internal/cli` on 2026-10-08, and a whole suite of 375 s to 384 s before its change and 227 s after it on the same machine. On 2026-11-25 the package has 2,102 and the suite takes 391 s.

### 5.2 The costliest tests by category

Test seconds from the same run.

The Go tool run over the module from inside a code test:

| Test | Package | Test seconds | What it runs |
| --- | --- | --- | --- |
| `TestCoverageCollectionMatchesTheToolchainOnThisPlatform` | `internal/spec` | 123 | `go test -list` over every package of the module with a contract tag, which compiles every test binary again |
| `TestJournalConsumerCorpusReplaysEveryConsumer` | `internal/store` | 103 | a nested `go test` of `internal/cli` and of `internal/tui`; it is 93 % of the wall time of its package |
| `TestPartitionFollowsTheMakefileRecipes` | `internal/verifyselect` | 23 | `go test -list` of `internal/cli`, which compiles its test binary |
| `TestCoverageEquivalence` | `internal/spec` | 13 | `go list` over the module for each release platform; it is also the test that compares the coverage record |
| `TestCoverageRecordCountsNoPackageUnderDocs` | `internal/spec` | 11 | the same collection again |
| `TestDerivedQAVerificationPassesTheChecker` | `internal/spec` | 11 | a nested `go test` of `internal/speccheck` |
| `TestPartitionCoversEveryTestExactlyOnce` | `internal/verifyselect` | 9 | `go list` over the module |

A binary built inside a test, once per test: `TestAgentSelectionProfilesMacro` (`internal/cli`, 55 s on the mini; one of its cases was at 4 min 20 s on the macOS runner when the binary was killed) builds the product binary and a stub. The product binary is built inside a test in six places of `internal/cli` and `internal/daemon`, and once per package in a seventh; a stub binary is built per call by a fixture that fifteen test files use.

The same archived corpus written to disk again by each test: thirteen tests of `internal/spec`, `internal/speccheck`, `internal/judge` and `internal/authorization` each extract the pinned history into their own temporary directory. The five of `internal/spec` take 310 test seconds: `TestQAReportReaderAgreesWithTheDerivedVerificationOnTheArchive` 86, `TestQAGateLegacyArchivedManifestsLoadUnchanged` 76, `TestRepositorySpecCorpusStillLoads` 64, `TestArchivedPassCorpusRemainsArchiveEligible` 44, `TestArchivedQAOverrideCorpusIncludesFailedSpec` 43.

The same expensive product operation repeated by each test: the eighteen asset-sync tests of `internal/baseline` take 619 test seconds, each running the whole sync two to four times over a source it builds itself (`TestAssetSyncStillRefusesAnInvalidCatalogWithoutDrift` 103, `TestAnAssetSyncWithNoSourceChangeLeavesTheComposedSetup` 88, `TestAnAssetSyncRewritesTheComposedSetupWithItsComponents` 86); `TestTransactionFailureMatrix` takes 388 over 61 cases and `TestBaselinePlansTheClaudeInstructionLink` 143 over 18.

Families of `internal/cli`, test seconds: review 1,052 over 308 tests; Baseline 1,037 over 145; reconcile 692 over 68; implement 614 over 122; release plan 468 over 30; archive 379 over 123.

Waiting for a time to pass: `TestRunDetachedCommandMonitorReportAllowsRunCreationPastLivenessDeadline` (`internal/cli`) waits the ten seconds of the handshake it tests, on every run; three tests of `internal/suiteguard` about helper residue take 10 s to 12 s each. The whole class is about 45 test seconds: the Real-Time Budget rule of ADR-0271 already removed the rest.

Not hermetic: none found. The whole suite passes in an environment with an empty home directory, without `CODEX_PATH` and without the user's tools on `PATH` (section 5), after Spec 0420 gave the three detach tests a fixture and made every External Test Tool skip by name. What remains is a test that skips: on a machine without an External Test Tool, the tests that need it do not run, and no gate says so beyond the skip line.

Serial: none worth a Task. Every one of the seven slowest packages is a Parallel Test Package, and their Sequential Tests are at most 40, held by ceilings that only fall.

Contracts run twice: `make verify-docs` runs the regeneration contracts of `internal/baseline` that `make verify-changed`, on a pull request that selects them, and `make verify-contracts`, on every push to main, also run in the same job. One run of them takes 214 s to 372 s on the macOS runner and 245 s to 308 s on the Linux runner.

### 5.3 Git processes by test

From the same log, by the temporary directory in the arguments, which names the test for 78 % of the calls: `TestAnAssetSyncWithNoSourceChangeLeavesTheComposedSetup` 1,562, `TestAnAssetSyncRewritesTheComposedSetupWithItsComponents` 1,177, `TestReconcileListsAndReleasesAnItemBranchOfAMergedSpec` 1,168, `TestMergedSpecRunWithLeftoversBeneathADeclaredDirectoryIsSuperseded` 935, `TestReconcileDryRunListsTheDirtyDiscardPathsAndDigest` 795. One asset sync of the Baseline package starts about 390 Git processes, and the tests of the reconcile command start 400 to 1,170 each. ADR-0259 counted at least 83,660 Git processes per run on 2026-10-08.

## 6. Published passages

Read on 2026-11-25.

- The Go command, `go help testflag`, Go 1.26.5: "If a test binary runs longer than duration d, panic. If d is 0, the timeout is disabled. The default is 10 minutes (10m)." The same text for `-parallel`: "Allow parallel execution of test functions that call t.Parallel". `go help build`, for `-p`: "the number of programs, such as build commands or test binaries, that can be run in parallel. The default is GOMAXPROCS, normally the number of CPUs available."
- The Go command, `go help test`: "In package list mode only, go test caches successful package test results to avoid unnecessary repeated running of tests."
- GitHub Docs, "GitHub-hosted runners reference", https://docs.github.com/en/actions/reference/runners/github-hosted-runners: "Use of the standard GitHub-hosted runners is free and unlimited on public repositories." The table of standard runners for public repositories lists macOS on arm64 with "3 (M1)" processors and 7 GB of memory under the labels `macos-latest`, `macos-15`, `macos-26` and `xcode-27`, and Linux with 4 processors and 16 GB under `ubuntu-latest`. The repository is public: the GitHub CLI answered `PUBLIC` on 2026-11-25.
- DORA, "Capabilities: Continuous integration", https://dora.dev/capabilities/continuous-integration/: "The tests should not take more than a few minutes to run, with an upper limit of about 10 minutes according to DORA’s research (PDF). If your build takes longer than this, you should improve the efficiency of your tests, add more compute resources so you can run them in parallel, or split out longer-running tests into a separate build using the deployment pipeline pattern."

- GitHub Docs, "Actions limits", https://docs.github.com/en/actions/reference/limits, opened again by the author of the Idea: standard runners allow 5 concurrent macOS jobs on the Free, Pro and Team plans. "During periods of high demand, jobs that run on macOS runners may take longer to start."

Found through Exa on 2026-11-25 for the rewrite of the PRD and read as the excerpts the search returned, not as whole pages; they are practitioner posts, not primary sources, and no acceptance row rests on them:

- "How We Stopped Shipping Performance Regressions", https://runcycles.io/blog/how-we-stopped-shipping-performance-regressions: a gate at release time that compares a median of three trials with a rolling median and blocks the release on a regression above 25 %. "A GitHub-hosted runner has a noisy 10 minutes and the benchmark run is slow. The gate should not fail."
- WarpBuild, "Setting Build Time Budgets for GitHub Actions", https://www.warpbuild.com/guides/ci-slo-and-build-time-budgets, a vendor guide: a budget is a p50 and a p95 for a class of workflows, "checked against measured percentiles on a fixed cadence, with a named action when a target breaks", and "Queue time gets its own budget line and stays out of the duration budget."
- WarpBuild, "Alerting on GitHub Actions Build Time Regressions", https://www.warpbuild.com/guides/alerting-on-ci-regressions, a vendor guide: "per-run alerting produces noise on both ends. It pages on the outlier and stays silent through the slow slide".

The same day the Indexus search `pesquisar_acervos` of the Secondbrain, restricted to wiki pages, returned no relevant page for what follows a spent budget.

Read by the Research Subagent on 2026-11-25 and not opened again by the author; they inform the plan and no acceptance row rests on them:

- Martin Fowler, "Continuous Integration", https://martinfowler.com/articles/continuousIntegration.html: the ten-minute commit build and a slower secondary build.
- golang/go issue 58571, https://github.com/golang/go/issues/58571: the test cache records the modification time of each file a test opens, so a fresh checkout misses it.
- golang/go issues 39054 and 26473: the Go project keeps its own short run at three to four minutes and trims it when it reaches five.
- Nicholas Nethercote, "Faster Rust builds on Mac", https://nnethercote.github.io/2025/09/04/faster-rust-builds-on-mac.html: macOS scans each new executable on its first launch, one at a time for the whole machine.
- jj pull request 7493, https://github.com/jj-vcs/jj/pull/7493: a process-heavy suite measured several times faster with its temporary directory on a memory disk on macOS.
- CockroachDB pull request 166889 and DigitalOcean's `gta`: selection of tests by the packages a change affects, each with a fallback to the whole suite.

## 7. The Secondbrain

Read on 2026-11-25: `wiki/index.md`; `wiki/concepts/testes-paralelos-em-go.md`, which records the measurement of Spec 0254 and the count of at least 83,660 Git processes per run; `wiki/concepts/flaky-tests-e-prazos-de-tempo-real.md`, which records why a test never waits on real time; `wiki/sources/digest-roundfix-toolchain-formatter-e-subprocessos-go-no-macos-2026-09-09.md`. The Indexus search `pesquisar_acervos` returned them for two queries about the time of the suite. They changed one thing in this plan: serial execution is not proposed as a lever, because Spec 0254 already removed it. The sweeps of the code and of the web were run by the Research Subagent, and the author read again every passage section 6 lists as read by the author. The research digest went to the Secondbrain's owner as ai-memory message `01a2124b-f1b9-78b0-bfd1-02abf74caf88`.

## 8. Governed Paths, measured for the rewrite of the PRD

Measured on 2026-11-25 at `33ece0cf` with `speccheck.GovernedPath`, in a `go test -overlay` probe that added one test file to the checker's package from outside the repository and wrote nothing to it. The default branch was two commits ahead, `9f85cc5e` and `3f621ac6`; neither changes the checker.

Governed, among the files the plan changes or may change: `Makefile`, `.github/workflows/ci-verify.yml`, `.github/workflows/release.yml`, `.roundfixrc.yml`, `docs/agents/setup-context.json`, `docs/agents/specific-repository.md`, `docs/references/coverage-record.json`, `internal/cli/baseline_plan_test.go`, `internal/spec/coverage_test.go`.

Governed and left unchanged by the plan: `go.mod`, `go.sum`, `.coderabbit.yaml`, `.github/workflows/ci-conventions.yml`, `.github/workflows/secondbrain-sync.yml`, `docs/agents/agent-instructions.md`, `docs/agents/autonomous-work.md`, `docs/agents/spec-routing.md`, every file under `.agents/skills/` and `internal/baseline/assets/`, `skills/roundfix/SKILL.md`, `internal/gittest/gittest.go`, `internal/cli/cli_test.go`, `internal/cli/baseline_release_gate_test.go`, `internal/spec/archive_test.go`, `internal/speccheck/governed_repocontract_test.go`, `internal/speccheck/mechanical_test.go`, `internal/baseline/plan_test.go`, `internal/baseline/repository_test.go`, `internal/docscontract/doc.go`, `internal/docscontract/testdata/corpus-golden.json`, `internal/suiteguardcontract/regeneration.go`. The version file `dist/npm/roundfix/package.json` is governed and changes only in a release.

Not governed: `CONTEXT.md`, `README.md`, every file under `docs/adr/`, `docs/references/` other than the coverage record, `docs/user-guide/`, `docs/agents/go.md`, `cmd/verify-select/`, `internal/verifyselect/`, `internal/testfixture/`, `internal/suiteguard/`, `internal/gittest/pinned_history.go`, a new test file under `internal/docscontract/`, and the test files that hold the tests section 5.2 names: `internal/spec/coverage_platform_test.go`, `internal/store/journal_consumer_corpus_test.go`, `internal/cli/implement_test.go`, `internal/cli/detach_unix_test.go`, `internal/baseline/setup_composition_test.go`, `internal/baseline/transaction_test.go`.

Two facts bound the list. The pinned history is read through one function of the Git test package, in a file that is not governed, so its callers need no change, the governed `internal/spec/archive_test.go` among them. The product binary is built inside a test in `internal/cli/baseline_plan_test.go`, which is governed, and once for its whole package, on purpose from an empty cache, in `internal/cli/baseline_release_gate_test.go`, which the plan leaves as it is.


## 9. The release, measured for the revision of 2026-11-25

The TechSpec author measured this section on 2026-11-25, after the maintainer's answer that continuous integration runs only for a pull request. `gh` ran with the operator's session on the mini, read-only.

### 9.1 What protects main

`gh api repos/{owner}/{repo}/branches/main/protection`, `gh api repos/{owner}/{repo}/rulesets` and `gh api repos/{owner}/{repo}`:

- Required checks: `Verification gate` and `Validate PR title`, both from GitHub Actions.
- `required_status_checks.strict` is `false`: a branch need not hold main's tip to merge.
- `enforce_admins` is `false`: an administrator may merge or push past the checks.
- No ruleset. Squash is the only merge method. The repository is public.

### 9.2 The last fifteen release tags against the four conditions

For each tag, with `C` its commit and `H` the head of the pull request GitHub lists for `C`:

```bash
gh api "repos/{owner}/{repo}/commits/$C/pulls"
git fetch origin "refs/pull/$N/head"
git merge-base --is-ancestor "$C^" "$H"
test "$(git rev-parse "$C^{tree}")" = "$(git rev-parse "$H^{tree}")"
gh api "repos/{owner}/{repo}/commits/$H/check-runs?check_name=Verification%20gate&filter=latest"
git diff --name-only "$C^" "$C" | grep -c '^dist/npm/roundfix/package.json$'
```

| Tag | Commit | Pull request | Its merge commit | Parent in the head, same tree | `Verification gate` on the head | Changes the version file |
| --- | --- | --- | --- | --- | --- | --- |
| `v0.153.3` | `3f621ac6` | #960 | yes | yes | success | yes |
| `v0.153.2` | `cefd4e82` | #954 | yes | yes | success | yes |
| `v0.153.1` | `7aa773db` | #952 | yes | yes | success | yes |
| `v0.153.0` | `fe247d4c` | #945 | yes | yes | success | yes |
| `v0.152.0` | `851024fa` | #940 | yes | yes | success | yes |
| `v0.151.0` | `7807cfcf` | #937 | yes | yes | success | yes |
| `v0.150.0` | `416e16dd` | #934 | yes | yes | success | yes |
| `v0.149.1` | `4c095fe2` | #931 | yes | yes | success | yes |
| `v0.149.0` | `66e4190c` | #929 | yes | yes | success | yes |
| `v0.148.2` | `894710c3` | #924 | yes | yes | success | yes |
| `v0.148.1` | `145d6c76` | #920 | yes | no | success | yes |
| `v0.148.0` | `02c41e4b` | #912 | yes | yes | success | yes |
| `v0.147.1` | `bcd716df` | #893 | yes | yes | success | yes |
| `v0.147.0` | `1fc00c9a` | #888 | yes | yes | success | yes |
| `v0.146.0` | `8189c25e` | #884 | yes | yes | success | yes |

Fourteen of the fifteen meet every condition. The pull request of `v0.148.1` merged behind main: its parent is not in the head and the trees differ, so its check ran on a tree that was not the published one. Every tagged commit is the squash merge of the release's pull request and changes the version file, which lies outside every Go package.

### 9.3 Governed Paths, measured again

At `3f621ac6`, with the probe of section 8, over 122 paths: the 119 of the TechSpec's first measurement and three more, the test file of the Release Provenance, ADR-0253 and the authorization record. Fifteen are governed: the nine of section 8 and six the design leaves unchanged, `internal/cli/baseline_release_gate_test.go`, `internal/gittest/gittest.go`, `internal/baseline/plan_test.go`, `docs/agents/agent-instructions.md`, `docs/agents/autonomous-work.md` and `dist/npm/roundfix/package.json`. The TechSpec's first text said 120 paths; its probe lists 119.

### 9.4 Published passages

Read on 2026-11-25. The Research Subagent read the GitHub pages; the author read "About protected branches" through Exa.

- GitHub Docs, "About protected branches" (https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-protected-branches/about-protected-branches), on a required check that is "Loose": "The branch does not have to be up to date with the base branch before merging." and "Status checks may fail after you merge your branch if there are incompatible changes with the base branch."
- GitHub Docs, "REST API endpoints for commits" (https://docs.github.com/en/rest/commits/commits#list-pull-requests-associated-with-a-commit): "Lists the merged pull request that introduced the commit to the repository. If the commit is not present in the default branch, it will return merged and open pull requests".
- GitHub Docs, "REST API endpoints for check runs" (https://docs.github.com/en/rest/checks/runs#list-check-runs-for-a-git-reference): "Lists check runs for a commit ref. The ref can be a SHA, branch name, or a tag name." Its `check_name` parameter "Returns check runs with the specified name."
- GitHub Docs, "Permissions required for GitHub Apps" (https://docs.github.com/en/rest/authentication/permissions-required-for-github-apps): the first endpoint is listed under the repository permission "Pull requests" with read access, the second under "Checks" with read access. Whether either answers without a token on a public repository is not stated.
- GitHub Docs, "Workflow syntax for GitHub Actions" (https://docs.github.com/en/actions/reference/workflows-and-actions/workflow-syntax): "If you specify the access for any of these permissions, all of those that are not specified are set to `none`."
- GitHub Docs, "Events that trigger workflows" (https://docs.github.com/en/actions/reference/workflows-and-actions/events-that-trigger-workflows), on `pull_request`: "Note that `GITHUB_SHA` for this event is the last merge commit of the pull request merge branch."
- GitHub Docs, "Dependency caching reference" (https://docs.github.com/en/actions/reference/workflows-and-actions/dependency-caching), as the fetch returned it, shortened: "Workflow runs can restore caches created in either the current branch or the default branch (usually `main`)." A cache a pull request creates belongs to its merge ref and "cannot be restored by the base branch or other pull requests targeting that base branch." That a pull request starts with no cache when nothing runs on main is the author's inference.
- GitHub Docs, "Using GitHub CLI in workflows" (https://docs.github.com/en/actions/how-tos/write-workflows/choose-what-workflows-do/use-github-cli): "GitHub CLI is preinstalled on all GitHub-hosted runners."
- Two published examples, found through Exa and read as excerpts, drop the tests from a release because an earlier gate ran them: a commit of a Java client, "Since the main branch is already tested, there is no point in re-running the same tests." (https://github.com/jkvn/docker-java/commit/6a44f43963e4e96979539b99d6291779a6525d84), and a Go command's workflow that skips its tests on a release's pull request (https://github.com/proggarapsody/bitbottle/blob/main/.github/workflows/ci.yml). They are examples, not primary sources, and neither checks that the release's tree is the tested one.

### 9.5 The Secondbrain

Searched with the Indexus tool `pesquisar_acervos` on 2026-11-25. Two pages opened:

- `concepts/tags-protegidas-e-release-por-workflow-no-github.md`: a sibling repository moved its whole verification to the pull request on 2026-11-16, after one delivery and its version ran the same verification six times, and made main require an up-to-date branch.
- `projects/quaestor-legacy/mirror/docs/references/ci-repetido-no-release-2026-11-16.md`: the measurement behind that decision. It records one release whose pull request merged before its check ended, so the only complete verification before publication was the release's own, and it notes that this repository ran everything on main and at a release because its pull request gate is selective.

`wiki/index.md` did not open through the tool, which is a limit of this consultation.
