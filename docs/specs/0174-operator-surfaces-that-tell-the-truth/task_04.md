---
task: task_04
spec: 0174-operator-surfaces-that-tell-the-truth
status: pending
type: backend
complexity: medium
---

# Task 04: The governed set covers every bounded file and its contract runs in a gate

## Overview

`GovernedPath` in `internal/speccheck/governed.go` does not match `skills/_ownership.yml`, the code-generator declaration archived Spec 0163 bounds, so `go test -tags repocontract -run TestEveryBoundedPathIsGoverned ./internal/speccheck` fails on `0160f70a`. Nothing notices, because the Makefile's `repo-test` target runs only four named `repocontract` tests in `./internal/baseline` and `./skills`, and `make verify` passes no build tag. The same gap hides `TestEverySpawningPackageInstallsTheSuiteGuard` in `./internal/suiteguard`, which fails because `internal/authorization` and `internal/verifyselect` spawn processes in their tests without installing `suiteguard.Main` (ADR-0126). The governed set is read by the changed-path audit and by `spec check`; the gates are read by every pull request through `make verify-docs`, which CI already runs.

## Requirements

1. MUST add `skills/_ownership.yml` to the historical bounded-path literal list in `internal/speccheck/governed.go`, so `GovernedPath` matches it under the clause ADR-0130 names, and add it to the newly governed list of `TestGovernedSetOnlyGrows` in `internal/speccheck/governed_repocontract_test.go`.
2. MUST add `TestGovernedPathMatchesTheSkillOwnershipDeclaration` and `TestGovernedPathKeepsAnotherSkillsRootFileOrdinary` to `internal/speccheck/governed_test.go`, asserting `GovernedPath("skills/_ownership.yml")` is true and `GovernedPath("skills/_notes.yml")` is false.
3. MUST install the suite guard in `internal/authorization` and `internal/verifyselect`: each gains a `main_test.go` whose `TestMain` calls `suiteguard.Main` with the repository root, and a `suiteguard_repocontract_test.go` that calls `suiteguardcontract.CheckCurrentPackage`, and both packages join `guardedSpawningPackages` in `internal/suiteguardcontract/contract.go`. A test that writes into the repository is fixed at the test, never by weakening the guard.
4. MUST make the Makefile's `REPO_CONTRACT_TESTS` list every top-level test declared in a `repocontract` file, adding `TestEverySpawningPackageInstallsTheSuiteGuard`, `TestGovernedSetCoversOwnedShippedTemplates`, `TestGovernedSetOnlyGrows`, `TestCleanupHistoricalGrantEvidence` and `TestEveryBoundedPathIsGoverned` to the four it names, and make `repo-test` run them with `-count=1 -tags repocontract` across `./...`; `verify-docs` keeps depending on `repo-test`, and `repo-test`'s help text says it runs the repository-contract tests.
5. MUST add `internal/suiteguardcontract/repository_gate_test.go`, with no build tag so `make verify` runs it. `TestEveryRepositoryContractTestRunsInTheRepositoryGate` walks the module as the Go tool does (skipping directories named `testdata` or starting with `.` or `_`), collects every top-level `Test` function in a `_test.go` file whose `//go:build` line requires `repocontract`, and fails naming each one missing from `REPO_CONTRACT_TESTS`; it also fails when `repo-test` does not pass `-tags repocontract` and `./...`, or `verify-docs` does not depend on `repo-test`. `TestRepositoryGateAuditNamesAMissingContractTest` and `TestRepositoryGateAuditNamesADroppedRepositoryWideRun` run the same audit over fixture modules in a temporary directory and assert the finding names the missing test and the missing `./...` respectively.
6. MUST NOT change `skills/_ownership.yml`, the CI workflow, or any other target's recipe.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] `GovernedPath` matches `skills/_ownership.yml` and not an unbounded sibling, and every authorization record's bounded paths are governed.
- [ ] Every internal package whose tests spawn a process installs the suite guard.
- [ ] `make repo-test` runs and passes every `repocontract` test in the repository, and leaving one out of its list fails a test `make verify` runs.

## Context

- interface: `internal/speccheck/governed.go`
- interface: `internal/speccheck/governed_repocontract_test.go`
- interface: `internal/speccheck/governed_test.go`
- interface: `Makefile`
- interface: `internal/suiteguardcontract/contract.go`
- creates: `internal/suiteguardcontract/repository_gate_test.go`
- creates: `internal/authorization/main_test.go`
- creates: `internal/authorization/suiteguard_repocontract_test.go`
- creates: `internal/verifyselect/main_test.go`
- creates: `internal/verifyselect/suiteguard_repocontract_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestGovernedPathMatchesTheSkillOwnershipDeclaration|TestGovernedPathKeepsAnotherSkillsRootFileOrdinary|TestEveryRepositoryContractTestRunsInTheRepositoryGate|TestRepositoryGateAuditNamesAMissingContractTest|TestRepositoryGateAuditNamesADroppedRepositoryWideRun)$" ./internal/speccheck ./internal/suiteguardcontract 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestGovernedPathMatchesTheSkillOwnershipDeclaration TestGovernedPathKeepsAnotherSkillsRootFileOrdinary TestEveryRepositoryContractTestRunsInTheRepositoryGate TestRepositoryGateAuditNamesAMissingContractTest TestRepositoryGateAuditNamesADroppedRepositoryWideRun; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done && contract="$(go test -count=1 -v -tags repocontract -run "^(TestGovernedSetOnlyGrows|TestEveryBoundedPathIsGoverned|TestEverySpawningPackageInstallsTheSuiteGuard)$" ./internal/speccheck ./internal/suiteguard ./internal/authorization ./internal/verifyselect 2>&1)" || { printf "%s\\n" "$contract"; exit 1; }; for name in TestGovernedSetOnlyGrows TestEveryBoundedPathIsGoverned TestEverySpawningPackageInstallsTheSuiteGuard; do printf "%s\\n" "$contract" | grep -q -- "--- PASS: $name" || exit 1; done && make repo-test` — expected: exit 0; before this Task the new named tests do not exist, `TestEveryBoundedPathIsGoverned` fails on `skills/_ownership.yml` and `TestEverySpawningPackageInstallsTheSuiteGuard` fails on two unguarded packages, so the command fails.

## References

- [_techspec.md](_techspec.md) — The governed set and the repository gate
- `_prd.md` → Goal 5; Core Feature 4; Success Metric 4
- `_techspec.md` → API Contracts 7-8; Testing Approach 4
