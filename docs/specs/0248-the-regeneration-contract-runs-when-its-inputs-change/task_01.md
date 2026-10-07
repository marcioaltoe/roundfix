---
task: task_01
spec: 0248-the-regeneration-contract-runs-when-its-inputs-change
status: pending
type: backend
complexity: medium
---

# Task 01: verify-select lists every contract with -all and names the contracts it leaves out

## Overview

Adds the selection half of the Full Contract Run to `cmd/verify-select`: the
`-all` flag of the `-contracts` mode prints the invocations of every
discovered Repository Contract Test, whatever its Contract Relevance. It also
makes the selective mode's summary line name each `relevant` contract it left
out and each `boundary` contract, so no exclusion is silent. It is verifiable
on its own through fixture tests and a real run on the repository.

## Requirements

1. MUST add a boolean flag `-all` to `verifyselect.Run`, implementing
   `_techspec.md` → API Contract 1 exactly: stdout is
   `ContractInvocations` of every contract `DiscoverContracts` returns, every
   class included; stderr is the one line
   `verify-select: contracts: <n> selected (every Repository Contract Test), 0 not selected`;
   exit 0. It MUST NOT run Git or `go list` in this mode.
2. MUST refuse `-all` without `-contracts` with
   `verify-select: -all requires -contracts` and exit 2, and `-all` with an
   explicitly set `-base` with
   `verify-select: -all cannot be combined with -base` and exit 2, writing
   nothing to stdout. The existing refusal of `-contracts` with `-packages` or
   `-baseline-cli-pattern` keeps its message and exit 2, with or without
   `-all`.
3. MUST, in the selective `-contracts` mode, keep stdout and the existing
   summary prefix byte-identical and extend the line as `_techspec.md` → API
   Contract 2 states: `; relevant not selected: <Name> (<package>), ...` when
   a `relevant` contract was not selected, then `; boundary: <Name> (<package>), ...`
   when a `boundary` contract exists, in discovery order. The fail-safe path
   prints the same extended line after its existing diagnostic.
4. MUST NOT change which contracts the selective mode selects, the directive
   grammar, `SelectContractPaths`, `SelectContracts` or `ContractInvocations`.
5. MUST add these tests to `internal/verifyselect/contracts_test.go`, each on
   a fixture repository in `t.TempDir()`:
   - `TestRunPrintsEveryContractWithAll`: a fixture with an `always`, a
     `package`, a `relevant` and a `boundary` contract under two tags prints
     one line per tag naming all four, and the exact `-all` summary line. It
     also passes when the fixture has no Git repository, which proves the
     mode never reads Git, and exits 1 with the discovery error when one
     directive is malformed.
   - `TestRunRefusesAllOutsideTheContractMode`: `-all` alone, `-contracts -all -base main`,
     and `-contracts -all -packages core` each exit 2 with the stated message
     and empty stdout.
   - `TestContractSummaryNamesEveryExclusion`: a change that selects one of
     two `relevant` contracts prints the exact extended line naming only the
     other, with its package; a fixture without `relevant` or `boundary`
     contracts prints the unchanged prefix alone.
6. MUST update the exact expected summary in
   `TestRunPrintsContractInvocations`, whose fixture holds one boundary
   contract, to the extended line ending in `; boundary: TestBoundary (pkg)`.
   No other existing assertion changes.

## Subtasks

- [ ] Add `-all` and its two refusals to `Run`.
- [ ] Extend the selective summary line with the named exclusions.
- [ ] Write the three fixture tests and update the one exact summary assertion.
- [ ] Run the new mode on the real repository.

## Acceptance Criteria

- [ ] `verify-select -contracts -all` prints every discovered contract, boundary ones included, and never reads Git.
- [ ] `-all` outside `-contracts`, or with `-base`, exits 2 with its message.
- [ ] The selective summary names every relevant contract left out and every boundary contract, and selection is unchanged.

## Context

- interface: `internal/verifyselect/verifyselect.go`
- interface: `internal/verifyselect/contracts.go`
- interface: `internal/verifyselect/contracts_test.go`
- instruction: `docs/adr/0253-every-repository-contract-test-runs-on-main-and-before-a-release.md`
- instruction: `docs/adr/0252-the-selective-gate-runs-the-repository-contracts-a-change-makes-relevant.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestRunPrintsEveryContractWithAll|TestRunRefusesAllOutsideTheContractMode|TestContractSummaryNamesEveryExclusion|TestRunPrintsContractInvocations)$' ./internal/verifyselect 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for t in TestRunPrintsEveryContractWithAll TestRunRefusesAllOutsideTheContractMode TestContractSummaryNamesEveryExclusion TestRunPrintsContractInvocations; do printf '%s\n' "$out" | grep -q -- "--- PASS: $t " || { printf 'missing passing test %s\n' "$t" >&2; exit 1; }; done; sel="$(go run -buildvcs=false ./cmd/verify-select -contracts -all 2>&1)" || { printf '%s\n' "$sel"; exit 1; }; printf '%s\n' "$sel" | grep -q -- '^docscontract .*TestRegenerationIsDeclared.* ./internal/config ./internal/docscontract$' || { printf 'the -all selection does not list TestRegenerationIsDeclared\n%s\n' "$sel" >&2; exit 1; }; printf '%s\n' "$sel" | grep -q -- 'selected (every Repository Contract Test), 0 not selected$' || { printf 'missing the -all summary line\n%s\n' "$sel" >&2; exit 1; }; go test -count=1 ./internal/verifyselect` — expected: exit 0. Before this Task the three new tests do not exist and `-all` is an unknown flag, so the first check fails. After it, the fixture tests pass, the real repository's `-all` selection lists the still-boundary `TestRegenerationIsDeclared` in `./internal/config`, and every `./internal/verifyselect` test passes.

## References

- `_prd.md` → Core Feature 1; Core Feature 2; Goals; Success Metric 3; Success Metric 6
- `_techspec.md` → API Contract 1; API Contract 2; Invariants 1-5; Surface Transcript 3; Surface Transcript 4; Surface Transcript 5; Build Order 1
- ADR-0253; ADR-0252
