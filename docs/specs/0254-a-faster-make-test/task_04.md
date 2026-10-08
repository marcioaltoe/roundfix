---
task: task_04
spec: 0254-a-faster-make-test
status: pending
type: infra
complexity: medium
---

# Task 04: One go test runs every selected Repository Contract Test

## Overview

`cmd/verify-select -contracts` prints one `go test` invocation per build tag,
and `make verify-contracts` and `make verify-changed-contracts` run them one
after the other (ADR-0253). On 2026-10-08 the Full Contract Run's two
invocations took 109 s and 91 s; one `go test` with both tags and the union of
names and packages took 62 s and 63 s in the same rounds, because Go runs the
two sets' packages side by side (`_techspec.md` → Current behavior). This Task
makes `ContractInvocations` print one invocation and adds a repository test
that keeps the merge exact (ADR-0259). It is verifiable on its own: the
selector prints one line and the Makefile loop runs it unchanged.

## Requirements

1. MUST make `verifyselect.ContractInvocations` follow `_techspec.md` → API
   Contract 3: no line for an empty selection, otherwise one line of the
   sorted comma-joined tags, the anchored alternation of the sorted unique
   quoted names, and the sorted unique `./<package>` arguments. A selection of
   one tag MUST print the line it prints today.
2. MUST update `TestContractInvocationsSortAndDeduplicate`,
   `TestRunPrintsContractInvocations` and `TestRunPrintsEveryContractWithAll`
   to the one-line form and to nothing else. These are the declared breaks;
   the Makefile wiring tests in `full_contract_run_test.go` and
   `contracts_repository_test.go` MUST pass unchanged.
3. MUST add `TestNoContractNameCrossesABuildTagClass` to
   `internal/verifyselect/contracts_repository_test.go`, following API
   Contract 4. It reads the real repository: it collects every top-level test
   name with its build-tag class (`docscontract`, `repocontract`, or neither)
   and fails, naming both files, when a Repository Contract Test's name is
   also a top-level test name in another class. It MUST also run on a seeded
   tree where a `docscontract` test and an untagged test share a name, and
   fail there. Like the other tests in that file it is an ordinary test with
   no build tag, so `make test` runs it on every change.
4. MUST NOT change the Makefile, the CI workflows, Contract Relevance,
   discovery or the selector's summary line (Invariant 5).
5. MUST record in the Result the wall time of `make verify-contracts` before
   and after the change, on the same machine, back to back, and both exit
   statuses.

## Subtasks

- [ ] Merge the invocations in `ContractInvocations`.
- [ ] Update the three declared tests to the one-line form.
- [ ] Add the name-uniqueness contract with its seeded case.
- [ ] Time the Full Contract Run before and after.

## Acceptance Criteria

- [ ] `go run -buildvcs=false ./cmd/verify-select -contracts -all` prints one
      line naming both tags.
- [ ] `make verify-contracts` passes and runs one `go test`.
- [ ] A seeded name that crosses a class fails the new contract.

## Context

- instruction: `docs/adr/0259-a-test-runs-in-parallel-unless-it-names-why-it-cannot.md`
- instruction: `docs/adr/0253-every-repository-contract-test-runs-on-main-and-before-a-release.md`
- instruction: `docs/adr/0252-the-selective-gate-runs-the-repository-contracts-a-change-makes-relevant.md`
- instruction: `Makefile`
- interface: `internal/verifyselect/contracts.go`
- interface: `internal/verifyselect/contracts_test.go`
- interface: `internal/verifyselect/contracts_repository_test.go`

## Verification

- `out="$(go run -buildvcs=false ./cmd/verify-select -contracts -all 2>/dev/null)" || { printf '%s\n' "$out"; exit 1; }; test "$(printf '%s\n' "$out" | grep -c .)" = 1 || { printf 'the selector printed more than one invocation\n' >&2; exit 1; }; printf '%s\n' "$out" | grep -q '^docscontract,repocontract ' || { printf 'the invocation does not name both tags\n' >&2; exit 1; }; out="$(go test -count=1 -v -run '^(TestContractInvocationsSortAndDeduplicate|TestRunPrintsContractInvocations|TestRunPrintsEveryContractWithAll)$' ./internal/verifyselect 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestContractInvocationsSortAndDeduplicate TestRunPrintsContractInvocations TestRunPrintsEveryContractWithAll; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the selector prints one line per tag, so the line count check fails.
- `out="$(go test -count=1 -v -run '^TestNoContractNameCrossesABuildTagClass$' ./internal/verifyselect 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestNoContractNameCrossesABuildTagClass ' || { printf 'the name-uniqueness contract did not run\n' >&2; exit 1; }; go test -count=1 ./internal/verifyselect || exit 1` — expected: exit 0; before this Task the contract does not exist, so no PASS line is printed and the command fails.

## References

- [_prd.md](_prd.md) — Goal 5; User Story 5; Core Feature 5; Success Metric 5
- [_techspec.md](_techspec.md) — Current behavior; API Contract 3; API Contract 4; Invariant 5; Testing Approach; Build Order 4
- ADR-0259; ADR-0252; ADR-0253
