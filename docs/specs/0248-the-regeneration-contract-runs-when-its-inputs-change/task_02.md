---
task: task_02
spec: 0248-the-regeneration-contract-runs-when-its-inputs-change
status: completed
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


## Result

Implemented the Task 02 slice for Daemon Verification. The regeneration
contract now carries the exact relevant directive from the TechSpec. Its
build constraint, body, assertions and all other bytes are unchanged.
The repository inventory checks its ordered patterns, rejects every boundary
contract, expects four contracts by change for an owned skill and exactly
`TestRegenerationIsDeclared` by change for `.roundfixrc.yml`. The existing
user-guide scenario still expects no contracts by change.

`TestRegenerationContractRunsWhenADerivedInputChanges` reads the real project
config through `config.ResolveConfigProposal(nil, contents)` and selects
against the real `DiscoverContracts` result. It checks all 15 path entries
across six declarations, including `lines.paths`, directory samples and
wildcard samples, plus `.roundfixrc.yml` and all seven generator files.
Failures identify the declaration, scope, original entry and sampled path.

### Focused checks and acceptance evidence

- Before changing the directive,
  `rtk proxy go test -count=1 -run '^TestRegenerationContractRunsWhenADerivedInputChanges$' ./internal/verifyselect`
  exited 1: every positive input was unselected under the original boundary
  directive, with each uncovered entry named. The first sandbox attempt was
  blocked by Go cache permissions; the executed check used approved access.
- After changing the directive,
  `rtk proxy go test -race -count=1 ./internal/verifyselect` exited 0.
  This exercises both repository tests and the rest of the selector package.
  Acceptance criterion 1: the project config, every declared sample and every
  generator selects the regeneration contract. Acceptance criterion 2:
  `docs/user-guide/example.md` and `internal/daemon/daemon.go` do not select it.
  The inventory scenarios also cover the selection counts in Surface
  Transcripts 1 and 2; their summary formatting belongs to Task 01.
- Acceptance criterion 3: compiled the package with
  `rtk proxy go test -c -o /tmp/roundfix-task02-verifyselect.test ./internal/verifyselect`,
  then used a Python temporary-repository probe to run that binary with
  `-test.run '^TestRegenerationContractRunsWhenADerivedInputChanges$'`.
  The copies retained real tagged test sources and project config; the real
  `.roundfixrc.yml` was never changed. Adding `uncovered/sample.json` to a
  declaration exited 1 and named that entry and sample as unselected.
  Reducing the copied config to three declarations exited 1 with
  `.roundfixrc.yml derived_paths has 3 declarations, want at least four`.
  Both expected failures were checked by the probe, which exited 0 and
  removed its temporary repositories.
- `rtk make verify-incremental GOCACHE=/Users/marcio/Library/Caches/go-build`
  exited 0 with approved cache and process access: formatting, Go analysis,
  repository package tests, skill checks and CLI build passed.
- Git diff inspection confirmed the regeneration directive is the only
  changed line in that contract file. A byte comparison confirmed the
  existing Makefile wiring test is unchanged. `git diff --check` exited 0.

Only the two Task 02 test files and this Result section were edited. The
initial worktree already contained the Daemon's `status: in_progress` change;
that status and the declared Verification section were preserved. No authored
Verification command, commit, push or pull request was run. No follow-up work
was identified. Task settlement remains with the Daemon.
