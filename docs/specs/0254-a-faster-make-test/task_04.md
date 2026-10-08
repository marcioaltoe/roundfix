---
task: task_04
spec: 0254-a-faster-make-test
status: completed
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


## Result

Implemented the Task slice for Daemon Verification:

- `ContractInvocations` returns no line for an empty selection and one line
  otherwise, with sorted unique comma-joined tags, sorted unique quoted test
  names in an anchored alternation, and sorted unique package arguments.
  The same construction preserves the existing single-tag form.
- Updated only the output expectations of the three declared tests. Contract
  discovery, Contract Relevance, summary diagnostics, the Makefile, CI, and
  the existing Makefile wiring tests remain unchanged.
- Added the ordinary, untagged `TestNoContractNameCrossesABuildTagClass`. Its
  repository scan parses top-level testing functions across packages and
  reports cross-class names with both files and their classes. Seeded cases
  reject a docscontract/untagged collision and a docscontract/repocontract
  collision, and allow a repeated name within one class.

### Focused evidence by acceptance criterion

1. Selector output: built the current entry point with
   `GOCACHE="$PWD/.gocache" go build -buildvcs=false -o /private/tmp/task04-verify-select ./cmd/verify-select`
   (exit 0), then captured `/private/tmp/task04-verify-select -contracts -all`
   with Python `subprocess.run`. It exited 0 and printed exactly one stdout
   line beginning `docscontract,repocontract`, with 29 package arguments.
   Stderr still reported `75 selected (every Repository Contract Test), 0 not selected`.
2. Full Contract Run: the stable after run of `make verify-contracts` exited
   0. The selector's one line and the unchanged Makefile loop yield one
   `go test` invocation. The unchanged
   `TestVerifyContractsRunsEveryContract` and
   `TestVerifyChangedRunsTheSelectedContracts` both passed in the focused
   check below.
3. Seeded collision: the new test passed on the real repository and asserted
   that its seeded docscontract/untagged tree produces this rejection:
   `TestShared crosses build-tag classes: docs/contract_test.go (docscontract) and other/shared_test.go (neither)`.
   The other-tag seeded case also produced a collision, while the same-class
   case produced none.

Focused command:

```sh
GOCACHE="$PWD/.gocache" go test -count=1 -v -run 'ContractInvocations|Prints.*Contract|NoContractName|VerifyContractsRuns|VerifyChangedRuns' ./internal/verifyselect
```

Final outcome: exit 0, package duration 4.755 s; all three declared output
checks, the new collision test and its seeded cases, and both unchanged
Makefile wiring tests passed. The first focused attempt exposed an incorrect
expected order for quoted names; correcting the expectation to
`TestA|TestZ|Test\.\+` preserved the existing sort of quoted names.

### Full Contract Run wall times

All measurements used this Task worktree, the same machine, installed Go
and Make, and the worktree's `.gocache`. Exit statuses are those of Make.

| Measurement | Command | Wall time | Exit status |
| --- | --- | --- | --- |
| Initial before, original source | `make verify-contracts` | 104.16 s | 0 |
| Initial after | `make verify-contracts` | 87.30 s | 2 |
| Stable before, original selector through Go overlay | `make verify-contracts GOFLAGS=-overlay=/private/tmp/task04-before-overlay.json` | 168.87 s | 0 |
| Stable after, merged selector | `make verify-contracts` | 209.70 s | 0 |

The initial pair used `/usr/bin/time -p`. Its after measurement is invalid
for comparison: the Agent corrected the quoted-name expectation while it
was running, and the baseline package's suite guard detected
`modified: internal/verifyselect/contracts_test.go`. This was an Agent source
mutation during measurement; the guard was retained.

The stable pair ran back to back in one Python process, timing each Make
subprocess with `time.monotonic()` and preserving each return code. No
repository file was edited during either run. The temporary overlay maps
only `internal/verifyselect/contracts.go` to its original `HEAD` bytes, so
both runs use the same final repository tests while comparing the original
two-line selector with the merged selector. Both selected all 75 contracts.

This pair did not demonstrate a speedup. Other Go test processes were
observed on the host during the pair, cache state was reused, and the before
run used an overlay. These measurements establish the observed wall times
and successful exits, without isolating the merge's performance effect.
Logs are `/private/tmp/task04-contracts-before.log`,
`/private/tmp/task04-contracts-after.log`,
`/private/tmp/task04-contracts-before-overlay.log`, and
`/private/tmp/task04-contracts-after-stable.log`; stable measurements are in
`/private/tmp/task04-contracts-pair.json`.

Task status remains Daemon-owned. Authored Verification commands, selected
repository Verification, incremental Verification, and the full ordinary
test suite were not run. No commit, push, or pull request was made.

## Carry-forward provenance

- Source Run: `run_20261008T215927Z_c26d3579de5392c8`
- Source commit: `7342ccf6c1456595d0b5deb7ec9d983ef68c855e`
