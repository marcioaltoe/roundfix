---
task: task_02
spec: 0248-the-regeneration-contract-runs-when-its-inputs-change
status: pending
type: backend
complexity: medium
---

# Task 02: TestRegenerationIsDeclared runs when an input its regeneration reads changes

## Overview

Moves `TestRegenerationIsDeclared` from `boundary`, where no target ran it,
to `relevant` over the inputs its regeneration reads, so `make
verify-changed` runs it when a change touches them and not otherwise. A
repository test derives the sample paths from the `derived_paths`
declarations in `.roundfixrc.yml`, so a new declaration that the directive
does not cover fails here. This answers the maintainer decision of
2026-10-07 recorded in `_authorization.md`.

## Requirements

1. MUST replace the `//verify:boundary` line in the header of
   `internal/config/regeneration_declared_test.go` with the exact directive
   in `_techspec.md` → Directive assignment. No other line of the file
   changes: not the build constraint, not the test body, not an assertion.
2. MUST update `TestRepositoryContractTestsDeclareTheirRelevance` in
   `internal/verifyselect/contracts_repository_test.go`:
   - `TestRegenerationIsDeclared` is `relevant`, and its paths equal the
     directive's patterns in order;
   - the check that a boundary contract exists is removed, and the test
     asserts instead that no discovered contract is `boundary`;
   - the scenario "owned skill selects regeneration" expects four contracts
     by change: the three with ADR-0252's regeneration paths and
     `TestRegenerationIsDeclared` with its own;
   - a new scenario for `.roundfixrc.yml` expects exactly one contract by
     change, `TestRegenerationIsDeclared`;
   - the scenario for `docs/user-guide/example.md` still expects none.

   Every other assertion of that test stays as it is.
3. MUST add `TestRegenerationContractRunsWhenADerivedInputChanges` to
   `internal/verifyselect/contracts_repository_test.go`. It reads the real
   `.roundfixrc.yml` through `config.ResolveConfigProposal` and, with
   `SelectContractPaths` over the real `DiscoverContracts` result, asserts
   that `TestRegenerationIsDeclared` is selected by:
   - `.roundfixrc.yml`;
   - a sample of every entry of every declaration's `paths` and
     `lines.paths`. An entry ending in `/` samples a file inside it, such as
     `<entry>sample.json`. An entry with `*` samples it with `*` replaced by
     `sample`. Any other entry samples itself.
   - one file of each generator: `internal/baseline/plan.go`,
     `skills/skills.go`, `internal/spec/archive.go`,
     `internal/cli/baseline_update.go`, `cmd/roundfix/main.go`,
     `internal/suiteguard/suiteguard.go` and
     `internal/suiteguardcontract/regeneration.go`.

   It asserts that `TestRegenerationIsDeclared` is not selected by
   `docs/user-guide/example.md` or `internal/daemon/daemon.go`. It fails, and
   names the uncovered entry, when the declarations have fewer than four
   entries or when any sample does not select the test.
4. MUST NOT change `.roundfixrc.yml`, any other contract test's directive or
   body, or the selector's code.
5. MUST keep `_techspec.md` → Surface Transcript 1 and Surface Transcript 2
   true of the real tree.

## Subtasks

- [ ] Change the regeneration test's directive.
- [ ] Update the repository inventory test and its scenarios.
- [ ] Add the declaration-derived selection test.
- [ ] Run the package's tests on the real tree.

## Acceptance Criteria

- [ ] A change to `.roundfixrc.yml` or to any declared path selects `TestRegenerationIsDeclared`.
- [ ] A change only to `docs/user-guide` or `internal/daemon` does not.
- [ ] A declaration whose path the directive does not cover fails `TestRegenerationContractRunsWhenADerivedInputChanges`.

## Context

- interface: `internal/config/regeneration_declared_test.go`
- interface: `internal/verifyselect/contracts_repository_test.go`
- instruction: `internal/config/config.go`
- instruction: `.roundfixrc.yml`
- instruction: `internal/verifyselect/contracts.go`
- instruction: `docs/adr/0253-every-repository-contract-test-runs-on-main-and-before-a-release.md`

## Verification

- `grep -q '^//verify:relevant .roundfixrc.yml internal/baseline/ .agents/skills/ skills/ docs/agents/ docs/references/coverage-record.json internal/spec/ internal/cli/baseline_. cmd/roundfix/ internal/suiteguard/ internal/suiteguardcontract/$' internal/config/regeneration_declared_test.go || { printf 'the regeneration test lacks its relevant directive\n' >&2; exit 1; }; ! grep -q '^//verify:boundary' internal/config/regeneration_declared_test.go || { printf 'the regeneration test is still boundary\n' >&2; exit 1; }; out="$(go test -count=1 -v -run '^(TestRepositoryContractTestsDeclareTheirRelevance|TestRegenerationContractRunsWhenADerivedInputChanges)$' ./internal/verifyselect 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for t in TestRepositoryContractTestsDeclareTheirRelevance TestRegenerationContractRunsWhenADerivedInputChanges; do printf '%s\n' "$out" | grep -q -- "--- PASS: $t " || { printf 'missing passing test %s\n' "$t" >&2; exit 1; }; done; go test -count=1 ./internal/verifyselect` — expected: exit 0. Before this Task the header still says `boundary` and the derived-selection test does not exist, so the first check fails. After it, the directive is in place, both repository tests pass on the real tree, and every `./internal/verifyselect` test passes.

## References

- `_prd.md` → Core Feature 3; Goals; Success Metric 1; Success Metric 2; Success Metric 6
- `_techspec.md` → Directive assignment; Invariant 6; Invariant 7; Surface Transcript 1; Surface Transcript 2; Testing Approach; Build Order 2
- ADR-0253; ADR-0252
