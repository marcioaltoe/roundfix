---
task: task_03
spec: 0248-the-regeneration-contract-runs-when-its-inputs-change
status: pending
type: infra
complexity: medium
---

# Task 03: make verify-contracts runs every contract, on every push to main and before every release

## Overview

Adds the Full Contract Run as a Makefile target and wires it into CI. `make
verify-contracts` runs every Repository Contract Test that task_01's
`verify-select -contracts -all` lists. `ci-verify.yml` runs it on every push
to main, and `release.yml` runs it after its `make verify` gate and before
any build or publication. Before this Task, `TestRegenerationIsDeclared` ran
in no workflow and the release ran no contract at all. The pull request path
is unchanged.

## Requirements

1. MUST add a `verify-contracts` target to the `Makefile` that implements
   `_techspec.md` → API Contract 3, list it in `.PHONY`, and give it a `##`
   help text that names the Full Contract Run. The recipe captures the
   selector's output before running any test, runs the printed invocations
   in order with `$(GO) test -count=1 -tags <tag> -run <pattern> <package>...`,
   stops at the first failure, and skips a line with no packages. Its shell
   variable for the test pattern MUST NOT be named so that the text
   `-run "$$pattern"` appears, because `TestPartitionFollowsTheMakefileRecipes`
   requires that text exactly once.
2. MUST add a Makefile comment above `verify-contracts` saying that it is the
   Full Contract Run, that discovery rather than a list decides the set, and
   that CI runs it on every push to main and before every release
   (ADR-0253).
3. MUST leave the `verify`, `verify-changed`, `verify-changed-contracts`,
   `verify-docs`, `docs-test`, `repo-test`, `spec-budget`, `verify-changed-core`
   and `verify-changed-baseline` recipes and `REPO_CONTRACT_TESTS`
   byte-identical.
4. MUST add the two workflow steps of `_techspec.md` → API Contract 4. In
   `.github/workflows/ci-verify.yml` the step `Verify every contract` follows
   "Verify docs" and runs only when `github.event_name == 'push'`. In
   `.github/workflows/release.yml` the step `Verify every contract` runs
   unconditionally, directly after "Verify gate" and before "Publication
   preflight". Each step carries a comment naming the Full Contract Run and
   ADR-0253. No other line of either workflow changes, and the pull request
   steps of `ci-verify.yml` keep their order and content.
5. MUST create `internal/verifyselect/full_contract_run_test.go` with:
   - `TestVerifyContractsRunsEveryContract`: it runs
     `make -f <repository Makefile> verify-contracts` in a temporary
     directory, with `VERIFY_SELECT` and `GO` set to stub scripts there. The
     stub selector exits 92 unless its arguments are exactly
     `-contracts -all`, and the stub `go` records its arguments. Four cases:
     - two printed lines run two invocations in order;
     - a failing first invocation stops before the second and fails `make`;
     - a failing selector runs no invocation and fails `make`;
     - an empty selection runs nothing and succeeds.
   - `TestMainAndReleaseRunEveryContract`: it parses both workflow files
     with `gopkg.in/yaml.v3` and asserts these facts:
     - `ci-verify.yml` is triggered by `push` to `main`;
     - its job holds a step whose `run` is `make verify-contracts` with
       `if: github.event_name == 'push'`, after the step named "Verify docs";
     - no step that runs on `pull_request` runs `make verify-contracts`;
     - `release.yml` holds a step whose `run` is `make verify-contracts` with
       no `if`, after the step whose `run` is `make verify` and before the
       step named "Publication preflight".
6. MUST NOT write outside the tests' temporary directories, and MUST NOT
   start a workflow, push, or reach the network.

## Subtasks

- [ ] Add the `verify-contracts` target, its help text and its comment.
- [ ] Add the push-to-main step and the release step.
- [ ] Write the Makefile and workflow tests.
- [ ] Run the target on the real tree and `make verify-docs`.

## Acceptance Criteria

- [ ] `make verify-contracts` runs every discovered contract and fails on the first failure.
- [ ] A push to main and a release each run `make verify-contracts`; a pull request does not run it.
- [ ] The listed recipes and `REPO_CONTRACT_TESTS` are byte-identical to `ae56aba0`.

## Context

- interface: `Makefile`
- interface: `.github/workflows/ci-verify.yml`
- interface: `.github/workflows/release.yml`
- creates: `internal/verifyselect/full_contract_run_test.go`
- instruction: `internal/verifyselect/contracts_repository_test.go`
- instruction: `internal/verifyselect/verifyselect_test.go`
- instruction: `docs/adr/0253-every-repository-contract-test-runs-on-main-and-before-a-release.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestVerifyContractsRunsEveryContract|TestMainAndReleaseRunEveryContract|TestVerifyChangedRunsTheSelectedContracts|TestPartitionFollowsTheMakefileRecipes)$' ./internal/verifyselect 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for t in TestVerifyContractsRunsEveryContract TestMainAndReleaseRunEveryContract TestVerifyChangedRunsTheSelectedContracts TestPartitionFollowsTheMakefileRecipes; do printf '%s\n' "$out" | grep -q -- "--- PASS: $t " || { printf 'missing passing test %s\n' "$t" >&2; exit 1; }; done; make verify-contracts && make verify-docs` — expected: exit 0. Before this Task neither new test exists and the `Makefile` has no `verify-contracts` target, so the first check fails. After it, the stubbed target and the workflows hold their contracts, the real Full Contract Run passes on the tree, and `make verify-docs` still passes.

## References

- `_prd.md` → Core Feature 1; Goals; Success Metric 3; Success Metric 4; Success Metric 5
- `_techspec.md` → API Contract 3; API Contract 4; Invariant 1; Invariant 8; Integration Points; Build Order 3
- ADR-0253
