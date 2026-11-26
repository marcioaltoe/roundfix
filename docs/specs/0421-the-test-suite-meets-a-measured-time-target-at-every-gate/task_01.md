---
task: task_01
spec: 0421-the-test-suite-meets-a-measured-time-target-at-every-gate
status: completed
type: test
complexity: high
---

# Task 01: Code tests stop running the Go tool over other packages and building per test

## Overview

This Task delivers `_techspec.md` → Build Order 1. Three of the four code tests that start a nested `go test` stop doing so in the default suite: two become Repository Contract Tests under the `docscontract` tag, and one is merged into the test it starts. The two tests that collect the coverage record share one collection, the product binary and each stub binary are built once for a test binary, and two new repository rules, the Nested Go Tool rule and the One Build rule, keep both out of the code tests. The fourth nested run, in `internal/verifyselect`, changes in task_03 with the sets it tests. It is verifiable on its own: the moved tests are discovered as contracts, the rules pass on the tree and fail on their seeded cases, and the coverage record loses exactly one name.

## Requirements

1. MUST move `TestCoverageCollectionMatchesTheToolchainOnThisPlatform`, unchanged in name and assertions, from `internal/spec/coverage_platform_test.go` into a new file `internal/spec/coverage_platform_docscontract_test.go` that carries the `docscontract` build tag and the directive `//verify:relevant dist/npm/platforms.json go.mod go.sum`. The contract also runs when a file of its own package changes, which holds the collector it compares.
2. MUST move `TestJournalConsumerCorpusReplaysEveryConsumer`, unchanged in name and assertions, from `internal/store/journal_consumer_corpus_test.go` into a new file `internal/store/journal_consumer_corpus_docscontract_test.go` that carries the `docscontract` build tag and the directive `//verify:relevant internal/store/testdata/ internal/runevent/ internal/tui/ internal/cli/attach.go internal/cli/events.go internal/cli/gc*.go internal/cli/reconcile.go`, the inputs its two replays read. Its helpers may stay where the untagged and the tagged files both compile them.
3. MUST merge `TestDerivedQAVerificationPassesTheChecker` (`internal/spec/task_test.go`) into `TestDerivedQAVerificationIsAcceptedByTaskStageChecker` (`internal/speccheck/coherence_test.go`), the test it starts in a nested `go test`: the first is removed and the second keeps every assertion it has, so the checker-side contract runs once, in its own process.
4. MUST make `TestCoverageEquivalence` and `TestCoverageRecordCountsNoPackageUnderDocs` share one collection of the coverage record for each run of `internal/spec`, within the bound `_authorization.md` → What the grant covers gives `internal/spec/coverage_test.go`: only how the collection runs and is shared changes, never what the record holds, how it is recorded, or the comparison, and the sanctioned record command keeps working unchanged.
5. MUST add `testfixture.ProductBinary`, which builds `./cmd/roundfix` once for a test binary, and make `testfixture.FixtureBinary` build each stub once for a test binary and each pair of name and source, keeping its signature and its callers. Each executable is written once to a new path and is never rewritten in place; a test that changes an executable copies it first (ADR-0281). The helpers live in `internal/testfixture/fixture.go`.
6. MUST make the tests that build a binary inside a test take those helpers: `TestBaselinePlanPreflightRealCLI` in `internal/cli/baseline_plan_test.go`, only for its build, with no assertion changed, as the grant bounds that Governed Path; `historyRemovalBinary` in `internal/cli/history_removal_test.go`; `buildRoundfixBinaryForMacro` and `buildMacroCodexExecutable` in `internal/cli/implement_test.go`; `internal/cli/reopen_test.go`; `internal/cli/reopen_late_dependency_test.go`; the build in `internal/daemon/run_disposition_characterization_test.go`; and `buildRealACPXEchoAgent` in `internal/agent/acpx_integration_unix_test.go`. `TestDoctorAdapterNamesTheCursorLogin` in `internal/cli/doctor_cursor_test.go` stops compiling the test binary of `internal/agent` with `go test -c` and proves the same login diagnostics in its own process or through a binary the helpers build, keeping both of its cases.
7. MUST add the Nested Go Tool rule in `internal/testfixture/nested_go_tool_test.go`: `TestNoCodeTestStartsAGoTestOrGoVetRun` walks the Go test sources of the module and fails, naming the file, when a file without a `docscontract` or `repocontract` build constraint starts `go test` or `go vet`, outside a frozen inventory that gives each entry its reason and may only shrink; `TestTheNestedGoToolRuleReportsSeededViolations` proves on seeded sources that it reports such a file and passes a contract file. The inventory holds `internal/verifyselect/verifyselect_test.go`, for `TestPartitionFollowsTheMakefileRecipes`, with the reason that task_03 changes it, and `internal/spec/coverage_test.go`, if its toolchain listing stays there, with the reason that only the contract of requirement 1 calls it.
8. MUST add the One Build rule in `internal/testfixture/one_build_rule_test.go`: `TestEveryTestBuildsABinaryOnlyThroughItsPackagesOnePlace` fails, naming the file, when a Go test source outside `internal/testfixture` starts `go build`, `go run` or `go test -c`, outside a frozen inventory that gives each entry its reason and may only shrink; `TestTheOneBuildRuleReportsSeededViolations` proves it on seeded sources. The inventory holds `internal/cli/baseline_release_gate_test.go` with the reason that it is a Governed Path outside this Spec's grant, and that file MUST stay byte-identical.
9. MUST add `internal/testfixture/fixture_test.go` with `TestProductBinaryIsBuiltOncePerTestBinary` and `TestFixtureBinaryIsBuiltOncePerNameAndSource`, each proving that two calls return the same executable and that a different source returns another. Because these tests start the Go tool, the package becomes a spawning test package: add `internal/testfixture/main_test.go` installing `suiteguard.Main`, `internal/testfixture/suiteguard_repocontract_test.go` carrying `TestEverySpawningPackageInstallsTheSuiteGuard`, and `internal/testfixture` in `guardedSpawningPackages` of `internal/suiteguardcontract/contract.go`, as the repository's rule on spawning test packages requires.
10. MUST add to the `relevant` map of `TestRepositoryContractTestsDeclareTheirRelevance` in `internal/verifyselect/contracts_repository_test.go` the two contracts of requirements 1 and 2 with exactly the paths their directives declare. The two numbers of `docscontract` files that test pins do not change, because neither new file is in `internal/docscontract`, `internal/cli` or `internal/baseline`; nothing else in that file changes.
11. MUST write the two rules in `docs/agents/specific-repository.md`, under `## Repository rules`, in the form of the rules already there: a code test starts neither `go test` nor `go vet`, and the Nested Go Tool rule in `internal/testfixture` enforces it; a test builds a binary only through `testfixture.ProductBinary` or `testfixture.FixtureBinary`, once for its test binary, and the One Build rule in `internal/testfixture` enforces it. No other sentence of the guide changes in this Task.
12. MUST re-record the coverage record last with `go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1`, the command `_authorization.md` → Sanctioned regeneration declares, and hand-edit nothing in it. Against the tree before this Task the record loses exactly `TestDerivedQAVerificationPassesTheChecker` of `roundfix/internal/spec`, merged by requirement 3, and gains the six tests of requirements 7 to 9 in `roundfix/internal/testfixture`. The two moved tests keep their names in the record, because it is collected with the `docscontract` tag.
13. MUST NOT change any assertion, test case or platform of a moved or repaired test, skip a test on any platform, add a `-timeout` to any invocation, or change `internal/cli/baseline_release_gate_test.go`, `internal/gittest/gittest.go`, `internal/baseline/plan_test.go` or any file `_authorization.md` → Limits names.

## Subtasks

- [ ] Move the two nested-run tests into relevant `docscontract` files and register them in the relevance map.
- [ ] Merge the derived QA Verification wrapper into the checker-side test.
- [ ] Share one coverage collection between the two record tests.
- [ ] Build the product binary and each stub once, and move every build site onto the helpers.
- [ ] Add the Nested Go Tool rule and the One Build rule with their inventories, seeded cases and guide sentences.
- [ ] Install the suite guard in `internal/testfixture` and re-record the coverage record.

## Acceptance Criteria

- [ ] Contract discovery lists `TestCoverageCollectionMatchesTheToolchainOnThisPlatform` and `TestJournalConsumerCorpusReplaysEveryConsumer` under the `docscontract` tag as `relevant`, and both pass when run with the tag.
- [ ] No code test outside the two inventories starts `go test`, `go vet`, `go build`, `go run` or `go test -c`; each rule fails on its seeded case.
- [ ] The product binary and each stub binary are built once for a test binary; every test that took the helpers passes with its assertions unchanged.
- [ ] The coverage record differs from the Delivery Base by the one merged name and the six new names, and `TestCoverageEquivalence` passes.

## Context

- instruction: `internal/testfixture/real_time_budget_test.go`
- instruction: `internal/testfixture/hard_link_rule_test.go`
- instruction: `internal/suiteguard/installation_repocontract_test.go`
- interface: `internal/spec/coverage_platform_test.go`
- creates: `internal/spec/coverage_platform_docscontract_test.go`
- interface: `internal/spec/coverage_test.go`
- interface: `internal/spec/task_test.go`
- interface: `internal/speccheck/coherence_test.go`
- interface: `internal/store/journal_consumer_corpus_test.go`
- creates: `internal/store/journal_consumer_corpus_docscontract_test.go`
- interface: `internal/testfixture/fixture.go`
- creates: `internal/testfixture/fixture_test.go`
- creates: `internal/testfixture/nested_go_tool_test.go`
- creates: `internal/testfixture/one_build_rule_test.go`
- creates: `internal/testfixture/main_test.go`
- creates: `internal/testfixture/suiteguard_repocontract_test.go`
- interface: `internal/suiteguardcontract/contract.go`
- interface: `internal/cli/baseline_plan_test.go`
- interface: `internal/cli/history_removal_test.go`
- interface: `internal/cli/implement_test.go`
- interface: `internal/cli/reopen_test.go`
- interface: `internal/cli/reopen_late_dependency_test.go`
- interface: `internal/cli/deliver_quaestor_built_binary_test.go`
- interface: `internal/cli/version_freshness_isolation_test.go`
- interface: `internal/cli/doctor_cursor_test.go`
- interface: `internal/daemon/run_disposition_characterization_test.go`
- interface: `internal/agent/acpx_integration_unix_test.go`
- interface: `internal/verifyselect/contracts_repository_test.go`
- interface: `docs/agents/specific-repository.md`
- interface: `docs/references/coverage-record.json`

## Consumers

None. `testfixture.FixtureBinary` keeps its signature and returns an executable path as before, `testfixture.ProductBinary` is new, and the coverage record keeps its format; no route, persisted format or report that code outside this slice reads changes.

## Verification

- `for pair in internal/spec/coverage_platform:TestCoverageCollectionMatchesTheToolchainOnThisPlatform internal/store/journal_consumer_corpus:TestJournalConsumerCorpusReplaysEveryConsumer; do base="${pair%%:*}"; name="${pair#*:}"; tagged="${base}_docscontract_test.go"; test -f "$tagged" || { echo "missing $tagged" >&2; exit 1; }; grep -q '^//go:build docscontract' "$tagged" && grep -q '^//verify:relevant ' "$tagged" && grep -q "^func $name(" "$tagged" || { echo "$name is not a relevant docscontract contract in $tagged" >&2; exit 1; }; if grep -q "^func $name(" "${base}_test.go"; then echo "$name still in ${base}_test.go" >&2; exit 1; fi; grep -q "\"$name\"" internal/verifyselect/contracts_repository_test.go || { echo "relevance map lacks $name" >&2; exit 1; }; done; out="$(go test -count=1 -v -run '^TestRepositoryContractTestsDeclareTheirRelevance$' ./internal/verifyselect 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestRepositoryContractTestsDeclareTheirRelevance '` — expected: exit 0, both nested-run tests live only in relevant docscontract files and the relevance map declares them.
- `test -f internal/spec/coverage_platform_docscontract_test.go && test -f internal/store/journal_consumer_corpus_docscontract_test.go || exit 1; out="$(go test -count=1 -timeout 20m -tags docscontract -v -run '^TestCoverageCollectionMatchesTheToolchainOnThisPlatform$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestCoverageCollectionMatchesTheToolchainOnThisPlatform ' || exit 1; out="$(go test -count=1 -timeout 20m -tags docscontract -v -run '^TestJournalConsumerCorpusReplaysEveryConsumer$' ./internal/store 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestJournalConsumerCorpusReplaysEveryConsumer '` — expected: exit 0, both contracts pass under the docscontract tag.
- `out="$(go test -count=1 -v -run '^(TestNoCodeTestStartsAGoTestOrGoVetRun|TestTheNestedGoToolRuleReportsSeededViolations|TestEveryTestBuildsABinaryOnlyThroughItsPackagesOnePlace|TestTheOneBuildRuleReportsSeededViolations|TestProductBinaryIsBuiltOncePerTestBinary|TestFixtureBinaryIsBuiltOncePerNameAndSource)$' ./internal/testfixture 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestNoCodeTestStartsAGoTestOrGoVetRun TestTheNestedGoToolRuleReportsSeededViolations TestEveryTestBuildsABinaryOnlyThroughItsPackagesOnePlace TestTheOneBuildRuleReportsSeededViolations TestProductBinaryIsBuiltOncePerTestBinary TestFixtureBinaryIsBuiltOncePerNameAndSource; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; out="$(go test -count=1 -tags repocontract -v -run '^TestEverySpawningPackageInstallsTheSuiteGuard$' ./internal/testfixture 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestEverySpawningPackageInstallsTheSuiteGuard '` — expected: exit 0, both rules pass on the tree and report their seeded cases, both helpers build once, and the package installs the suite guard.
- `test -f internal/cli/doctor_cursor_test.go && ! grep -qE '"go", "(build|test)"' internal/cli/baseline_plan_test.go internal/cli/history_removal_test.go internal/cli/implement_test.go internal/cli/reopen_test.go internal/cli/reopen_late_dependency_test.go internal/cli/doctor_cursor_test.go internal/daemon/run_disposition_characterization_test.go internal/agent/acpx_integration_unix_test.go && ! grep -q 'func TestDerivedQAVerificationPassesTheChecker' internal/spec/task_test.go && out="$(go test -count=1 -v -run '^TestDerivedQAVerificationIsAcceptedByTaskStageChecker$' ./internal/speccheck 2>&1)" && printf '%s\n' "$out" | grep -q -- '--- PASS: TestDerivedQAVerificationIsAcceptedByTaskStageChecker ' && out="$(go test -count=1 -v -run '^(TestBaselinePlanPreflightRealCLI|TestReopenThroughTheBuiltBinary|TestReopenReopensALateDependencyThroughTheBuiltBinary|TestDoctorAdapterNamesTheCursorLogin)$' ./internal/cli 2>&1)" && for name in TestBaselinePlanPreflightRealCLI TestReopenThroughTheBuiltBinary TestReopenReopensALateDependencyThroughTheBuiltBinary TestDoctorAdapterNamesTheCursorLogin; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf '%s\n' "$out"; echo "missing pass: $name" >&2; exit 1; }; done` — expected: exit 0, no repaired test builds or tests another package itself, the merged checker test passes, and the tests that took the helpers pass.
- `test -f docs/references/coverage-record.json && ! grep -q '"TestDerivedQAVerificationPassesTheChecker"' docs/references/coverage-record.json && for name in TestNoCodeTestStartsAGoTestOrGoVetRun TestTheNestedGoToolRuleReportsSeededViolations TestEveryTestBuildsABinaryOnlyThroughItsPackagesOnePlace TestTheOneBuildRuleReportsSeededViolations TestProductBinaryIsBuiltOncePerTestBinary TestFixtureBinaryIsBuiltOncePerNameAndSource TestCoverageCollectionMatchesTheToolchainOnThisPlatform TestJournalConsumerCorpusReplaysEveryConsumer; do grep -q "\"$name\"" docs/references/coverage-record.json || { echo "record lacks $name" >&2; exit 1; }; done && out="$(go test -count=1 -v -run '^(TestCoverageEquivalence|TestCoverageRecordCountsNoPackageUnderDocs)$' ./internal/spec 2>&1)" && printf '%s\n' "$out" | grep -q -- '--- PASS: TestCoverageEquivalence ' && printf '%s\n' "$out" | grep -q -- '--- PASS: TestCoverageRecordCountsNoPackageUnderDocs '` — expected: exit 0, the record lost the merged name, holds the moved and the new names, and both record tests pass on one collection.
- `tr -s '[:space:]' ' ' < docs/agents/specific-repository.md | grep -qF -- "Nested Go Tool rule" || { printf 'missing phrase in %s: %s\n' docs/agents/specific-repository.md "Nested Go Tool rule" >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/agents/specific-repository.md | grep -qF -- "One Build rule" || { printf 'missing phrase in %s: %s\n' docs/agents/specific-repository.md "One Build rule" >&2; exit 1; }` — expected: exit 0, the repository's guide states both rules.

## References

- `_prd.md` → Goal 4; User Story 6; Core Features 9, 10 and 13; Success Metric 5; Success Metric 9; Categories of work; Sizing, Task 1
- `_techspec.md` → Build Order 1; What a settlement runs; Testing Approach; Coverage Map; Deviations from the PRD; Risks & Considerations
- `_authorization.md` → What the grant covers; Sanctioned regeneration; Limits
- [`references/sources.md`](references/sources.md), section 5.2
- ADR-0252, ADR-0259, ADR-0271, ADR-0281, ADR-0326

## Result

Skills loaded: implement-task, coding-guidelines, golang-testing, testing-boss, golang-concurrency, golang-cli, golang-error-handling, golang-structs-interfaces, systematic-debugging, tech-writer, evidence-gate.

Implemented only this Task's slice. The pre-existing change at entry was this Task's Daemon-owned status; it was preserved. No Task Graph or other Task file was changed.

- Contract discovery: moved both nested-run tests into their declared relevant docscontract files and registered exactly their directive inputs. A comparison against HEAD confirmed both moved function bodies are byte-identical. Focused execution of both tagged contracts passed: `go test -tags docscontract ./internal/spec ./internal/store -run 'Test(CoverageCollectionMatches|JournalConsumerCorpusReplays)' -count=1`. The focused `TestRepositoryContractTestsDeclareTheirRelevance` run also passed.
- Nested runs and builds: removed the derived-QA subprocess wrapper while preserving the existing checker-side test and every assertion. Added the Nested Go Tool and One Build AST rules, frozen reasoned inventories, and seeded positive/negative cases, including variadic arguments and contract constraints. Both rules pass on the repository and their seeded cases. Added only the two requested guide rules. `go test ./internal/testfixture -tags repocontract -count=1` passed, including the existing executable-write, hard-link, timing and parallelism rules and the new suite guard installation contract.
- Shared executables: added process-wide, concurrency-safe ProductBinary and name/source-keyed FixtureBinary caches. Each output has a fresh path in the Go test build directory; individual test cleanup cannot remove it. Windows executable suffixes are preserved. The macro Codex fixture copies the cached executable through the existing fork-safe executable writer before use. The six new tests passed; focused cache and seeded-rule tests also passed with `-race`.
- Migrated consumers: focused checks passed for the Baseline real CLI, both built-binary reopen tests, both Doctor login cases, version-cache preservation, both history transcripts, the agent-selection macro, measured Quaestor queue, daemon unintegrated-branch preflight and derived-QA checker. The actual local echo-agent integration passed with `ROUNDFIX_REAL_ACPX=1 go test ./internal/agent -run '^TestRealACPXCommandOverrideRoundTripCancelCrashResume$' -count=1 -v`; it exercised cancellation, crash recovery and teardown. Non-test packages compiled for Linux (`internal/agent` and `internal/testfixture`) and Windows (`internal/testfixture`).
- Coverage: after the final source edit, ran the sanctioned command unchanged: `go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1`. It passed. A JSON comparison against HEAD proved exactly TestDerivedQAVerificationPassesTheChecker was removed from roundfix/internal/spec, exactly the six requested tests were added to roundfix/internal/testfixture, and every other record value is unchanged. Both moved names remain recorded. `go test ./internal/spec -race -run 'TestCoverage(Equivalence|RecordCountsNoPackageUnderDocs)$' -count=1` passed.
- Scope checks: the three expressly frozen files (baseline_release_gate_test.go, internal/gittest/gittest.go and internal/baseline/plan_test.go) are byte-identical to HEAD. Formatting and `git diff --check` checks passed.

One intermediate macro check had passing assertions but was rejected by suiteguard because source edits overlapped its run. It was rerun with the tree held still and passed. Move-related unused imports found by the initial compile check were removed.

Declared Verification and Task settlement remain with the Daemon. No commit, push or pull request was made.

## Carry-forward provenance

- Source Run: `run_20261126T024343Z_4d6d0cf3419b8574`
- Source commit: `3de59336950189bd0f40f23ce165661244ec2e7b`
