---
task: task_02
spec: 0420-roundfix-stops-verifying-on-linux
status: pending
type: infra
complexity: medium
---

# Task 02: Pull requests, main and releases are verified on macOS, and Linux is compiled

## Overview

This Task delivers `_techspec.md` → Build Order 2, on the maintainer's decision of 2026-11-24, "Tirar de vez", and its context, "Talvez só deixar algum tipo de teste no CI antes do merge na main". The job `Verification gate` of `CI | Verify` runs on the standard macOS runner `macos-26` with every step it has; the release verifies on that runner in a job of its own that the publishing job needs; the Repository Contract Test that compiled every test for Windows compiles them for Linux as well; and the Makefile, the README and the repository's own guide say what is true after the change. It is verifiable on its own through three new workflow contracts, the moved assertions of an existing one, and the renamed compile contract with a subtest over a temporary module.

## Requirements

1. MUST change in `.github/workflows/ci-verify.yml` only the runner of the job `verify` to `runs-on: macos-26`; the workflow's name, triggers, permissions and concurrency, the job's name `Verification gate` and its timeout, and each of its nine steps with its name, identifier, condition, environment and command keep their bytes and their order (API Contract 6, invariant 6). This is the Delegated Decision of 2026-11-24 of the `techspec` phase, "Which label names the runner, the moving `macos-latest` or a versioned one?".
2. MUST give `.github/workflows/release.yml` the two jobs of `_techspec.md` → Interfaces, as API Contract 7 states: a job `verify` on `macos-26` with `contents: read`, whose steps are the checkout with the full history and no persisted credentials, the Go setup with the options the release uses today, `Verify gate` and `Verify every contract` with today's names, commands and comments and no condition, and which carries the checkout comment that explains the full history; and the job `release` with `needs: verify`, without those two steps, whose runner, permissions, environment and every other step keep their bytes and their order, and whose checkout keeps its options while its comment says that no step of the job reads the history any more and that the option stays as it was. The triggers, the input and the workflow's name do not change, and no condition lets `release` run when `verify` did not succeed. This is the Delegated Decision of 2026-11-24 of the `techspec` phase on how the release keeps its verification before the publication.
3. MUST delete `internal/verifyselect/windows_test_compile_repocontract_test.go` and create `internal/verifyselect/platform_test_compile_repocontract_test.go`, with the `repocontract` build tag, no Contract Relevance declaration (as the file it replaces) and a `Suite`, `Invariant` and `Boundary` header, holding `compiledPlatforms`, `compileEveryTest` and `TestEveryTestCompilesForLinuxAndWindows` with the signatures of `_techspec.md` → Interfaces and the behavior of API Contract 8. The test has the subtests `linux` and `windows`, each running `go vet -tags "repocontract docscontract" ./...` at the repository root with `GOOS` set and `GOARCH` inherited and failing with the platform's name and the compiler's output, and `a_file_only_linux_fails_to_compile`, which runs `compileEveryTest` over a temporary module with no dependency whose one test file builds only for Linux and does not compile, and asserts that the Linux run fails with an output naming that file and the Windows run passes.
4. MUST create `internal/verifyselect/verifying_runner_test.go`, untagged like the tests that already read the workflows, holding `verifyingRunner` as `_techspec.md` → Interfaces gives it and:
   - `TestEveryVerifyingJobRunsOnTheStandardMacOSRunner`, which reads every file of this repository's `.github/workflows` and asserts invariant 7 in both directions: every job whose steps run `make verify`, `make verify-changed`, `make verify-docs`, `make verify-contracts`, `go test` or `go vet` names `runs-on: macos-26`, and no job on another runner runs one of them;
   - `TestTheReleaseVerifiesBeforeItBuildsOrPublishes`, which asserts API Contract 7 and invariant 8: the job `verify` exists on `macos-26` with `contents: read` and the two steps, the job `release` names `needs: verify`, holds neither step, keeps its steps in today's order from the checkout to `GitHub Release`, and no job carries a condition that runs `release` after a failed `verify`;
   - `TestTheVerificationGateKeepsItsNamesAndItsSteps`, which asserts invariant 6: the workflow `CI | Verify`, its job `Verification gate`, its timeout of 60 minutes and the names of its nine steps in today's order.
   None of them reads a Baseline asset; each top-level test calls `t.Parallel()` first.
5. MUST move the release assertions of `TestMainAndReleaseRunEveryContract` in `internal/verifyselect/full_contract_run_test.go` from the job `release` to the job `verify`, keeping its name and its assertions about `ci-verify.yml`, and asserting that the job `release`, which needs `verify`, holds `Publication preflight`.
6. MUST change in the `Makefile` only the three places `_authorization.md` → What the grant covers bounds (API Contract 9):
   - the comment above `verify-changed`, which no longer calls continuous integration the Linux gate and says that the selective gate and continuous integration run on macOS and that a failure only Linux or Windows produces is a known limit, because those platforms are only compiled (ADR-0252);
   - the comment of `SUITE_BUDGET_SECONDS`, which keeps the measurement it states and says that those runs were on the Linux runner the gate used until this change, that the gate now runs on the standard macOS runner, where the suite has not been measured, and that the first push to main there is the measurement the number waits for; the value stays `420` (the Delegated Decision of 2026-11-24 of the `techspec` phase on the suite-time budget);
   - `REPO_CONTRACT_TESTS`, which names `TestEveryTestCompilesForLinuxAndWindows` in the place of `TestEveryTestCompilesForWindows`.
7. MUST change in `docs/agents/specific-repository.md` only the sentence outside the Baseline's markers that calls continuous integration the Linux gate, so that it says what requirement 6's first comment says; and add to `README.md`, after the sentence that names the platform packages, the sentence "The test suite runs on macOS only; the Linux and Windows packages are compiled and not tested." (API Contract 8).
8. MUST keep passing unchanged `TestTheSuiteTimeBudgetIsRecordedThenJudged`, `TestTheCIPushJudgesTheSuiteTimeBudgetLast`, `TestTheCIPullRequestSkipsOnlyTheRepeatedDocumentationContracts`, `TestRepositoryContractTestsDeclareTheirRelevance` and the tests of `internal/verifyselect/release_workflow_test.go` and `internal/verifyselect/release_workflow_compare_test.go`.
9. MUST NOT change `.github/workflows/ci-conventions.yml`, `.github/workflows/secondbrain-sync.yml`, any Baseline asset (the starter workflow and `internal/baseline/starter_files_test.go` keep `ubuntu-latest`), any other target, recipe or variable of the Makefile, any other line of `docs/agents/specific-repository.md`, `docs/references/coverage-record.json`, `go.mod`, `go.sum`, or any file `_authorization.md` → Limits names. No test reaches GitHub or the network.

## Subtasks

- [ ] Move the job `Verification gate` to `macos-26`.
- [ ] Split the release into the job `verify` and the job `release` that needs it.
- [ ] Replace the Windows compile contract with the Linux and Windows one and its proof on a temporary module.
- [ ] Add the three workflow contracts and move the release assertions of the existing one.
- [ ] Correct the Makefile's two comments and contract name, the guide's sentence and the README.

## Acceptance Criteria

- [ ] Every job that runs a verification, a test, a vet or a contract runs on `macos-26`, and no other job runs one; the job `Verification gate` keeps its names, its timeout and its steps.
- [ ] The release runs `Verify gate` and `Verify every contract` on `macos-26` in a job the publishing job needs, and every publishing step keeps its bytes, order and runner.
- [ ] The compile contract fails on a test file that compiles for macOS and Windows and not for Linux, naming the file, and passes on this repository for both platforms.
- [ ] The Makefile, the guide's sentence and the README name macOS as the only tested platform, the budget stays 420 seconds with its comment, and the Makefile lists the renamed contract.

## Context

- instruction: `internal/verifyselect/suite_budget_test.go`
- instruction: `internal/verifyselect/release_workflow_test.go`
- instruction: `internal/verifyselect/release_workflow_compare_test.go`
- instruction: `internal/verifyselect/contracts_repository_test.go`
- interface: `.github/workflows/ci-verify.yml`
- interface: `.github/workflows/release.yml`
- interface: `Makefile`
- interface: `README.md`
- interface: `docs/agents/specific-repository.md`
- interface: `internal/verifyselect/full_contract_run_test.go`
- creates: `internal/verifyselect/platform_test_compile_repocontract_test.go`
- creates: `internal/verifyselect/verifying_runner_test.go`
- deletes: `internal/verifyselect/windows_test_compile_repocontract_test.go`

## Consumers

- the two workflows and the Makefile's `REPO_CONTRACT_TESTS` — consumers: `./internal/verifyselect`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestEveryVerifyingJobRunsOnTheStandardMacOSRunner|TestTheReleaseVerifiesBeforeItBuildsOrPublishes|TestTheVerificationGateKeepsItsNamesAndItsSteps|TestMainAndReleaseRunEveryContract)$' ./internal/verifyselect 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestEveryVerifyingJobRunsOnTheStandardMacOSRunner TestTheReleaseVerifiesBeforeItBuildsOrPublishes TestTheVerificationGateKeepsItsNamesAndItsSteps TestMainAndReleaseRunEveryContract; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done` — expected: exit 0, every verifying job runs on macos-26, the release verifies in the job its publishing job needs, and the gate keeps its names and steps.
- `out="$(go test -count=1 -timeout 30m -tags repocontract -v -run '^TestEveryTestCompilesForLinuxAndWindows$' ./internal/verifyselect 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestEveryTestCompilesForLinuxAndWindows/linux TestEveryTestCompilesForLinuxAndWindows/windows TestEveryTestCompilesForLinuxAndWindows/a_file_only_linux_fails_to_compile TestEveryTestCompilesForLinuxAndWindows; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done` — expected: exit 0, every test compiles for Linux and Windows and the contract names a file only Linux fails to compile.
- `grep -qx '    runs-on: macos-26' .github/workflows/ci-verify.yml && ! grep -q 'ubuntu' .github/workflows/ci-verify.yml && grep -qx '    needs: verify' .github/workflows/release.yml && grep -qx '    runs-on: macos-26' .github/workflows/release.yml && grep -qx 'SUITE_BUDGET_SECONDS ?= 420' Makefile && grep -q 'TestEveryTestCompilesForLinuxAndWindows' Makefile && ! grep -q 'TestEveryTestCompilesForWindows' Makefile && ! grep -q 'Linux gate' Makefile && grep -q 'macOS' Makefile && if tr -s '[:space:]' ' ' < docs/agents/specific-repository.md | grep -qF -- 'Linux gate'; then echo 'the guide still calls CI the Linux gate' >&2; exit 1; fi && { tr -s '[:space:]' ' ' < README.md | grep -qF -- 'The test suite runs on macOS only; the Linux and Windows packages are compiled and not tested.' || { printf 'missing phrase in %s: %s\n' README.md 'The test suite runs on macOS only' >&2; exit 1; }; } && go test -count=1 -timeout 20m ./internal/verifyselect` — expected: exit 0, both workflows name macos-26, the release job needs verify, the Makefile, the guide and the README say macOS is the only tested platform with the budget unchanged, and the workflow and Makefile contracts pass whole.

## References

- `_prd.md` → Goals 3 and 4; User Stories 3 and 5; Core Features 6, 7, 8, 9 and 11; Success Metric 6; Success Metric 7; Success Metric 9; Acceptance evidence; Unreachable Acceptance; Decisions
- `_techspec.md` → The two workflows and the compile contract; Interfaces; API Contract 6; API Contract 7; API Contract 8; API Contract 9; Invariants; Integration Points; Testing Approach; Build Order 2; Risks & Considerations; Decisions (the Delegated Decisions of 2026-11-24 on the suite-time budget, the runner's label and the release's jobs)
- `_authorization.md` → What the grant covers; Limits
- [`references/2026-11-24-techspec-measurements.md`](references/2026-11-24-techspec-measurements.md), sections 1, 2, 7 and 8
- ADR-0252, ADR-0253, ADR-0285, ADR-0236, ADR-0387, ADR-0031, ADR-0250
