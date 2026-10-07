---
task: task_02
spec: 0247-a-run-gate-that-runs-the-repository-contracts
status: completed
type: backend
complexity: medium
---

# Task 02: make verify-changed runs the selected contracts, and every contract declares its relevance

## Overview

Wires task_01's selector into the Run gate. `make verify-changed` ends by
running the Repository Contract Tests that the change makes relevant, and each
existing contract test file declares its Contract Relevance as ADR-0252
decides. This Task answers the suiteguard installation failure recorded in
entries 237 and 238 of the operator's queue log on 2026-10-07. It is
verifiable on its own: the repository inventory test, the Makefile test, a
real run of the new target and `make verify-docs`.

## Requirements

1. MUST add a header directive to each file listed in `_techspec.md` →
   Directive assignments, exactly as listed. Each directive goes in the
   comment block before the `package` clause, and a file's existing
   `// Suite:` and `//go:build` lines stay. Where a header sentence says the
   test runs only at the pull request boundary, that sentence changes to say
   that the selective gate runs it by its Contract Relevance and that
   `make verify-docs` still runs it. `internal/config/regeneration_declared_test.go`
   gets `//verify:boundary` with its reason. No test body, assertion or
   other comment changes.
2. MUST NOT add a directive to any per-package
   `suiteguard_repocontract_test.go`. Those files stay class `package`.
3. MUST add to the `Makefile`:
   - a `verify-changed-contracts` target that implements `_techspec.md` →
     Invariant 12 through `$(VERIFY_SELECT) -contracts -base "$(VERIFY_BASE)"`;
   - its name in `.PHONY`;
   - an invocation from the `verify-changed` recipe after the set loop that
     runs on every run, including a run whose set list is empty.

   The new recipe MUST NOT contain the text `-run "$$pattern"`, because
   `TestPartitionFollowsTheMakefileRecipes` requires that text exactly once.
   Name its shell variable differently, for example `$$tests`.
4. MUST add a Makefile comment above `verify-changed` saying three things:
   - the gate runs the Repository Contract Tests that their Contract
     Relevance selects;
   - `make verify-docs` still runs all of them;
   - platform-only failures stay a known limit, because CI is the Linux gate.
     The example is `go list` printing `go: downloading` on stderr for
     modules the host never fetched (ADR-0252).
5. MUST leave the `verify`, `verify-docs`, `docs-test`, `repo-test`,
   `spec-budget`, `verify-changed-core` and `verify-changed-baseline`
   recipes byte-identical. It MUST NOT change `REPO_CONTRACT_TESTS`, the CI
   workflows or `.roundfixrc.yml`.
6. MUST add `internal/verifyselect/contracts_repository_test.go`, untagged,
   with these tests:
   - `TestRepositoryContractTestsDeclareTheirRelevance`: `DiscoverContracts`
     on the real repository returns no error and has these classes:
     - `TestEverySpawningPackageInstallsTheSuiteGuard` in
       `internal/suiteguard`, each `docscontract` test in
       `internal/docscontract`, and the four governed-set tests in
       `internal/speccheck` are `always`;
     - `TestMeasuredSanctionedOwnershipMatchesRecords`,
       `TestDeclaredStepRegenerationAndFrozenBoundaries`,
       `TestOwnedSkillEditLeavesDerivedArtifactsByteIdentical` and
       `TestRepositoryGateRunsTheAnalyzer` are `relevant`, with the paths of
       the Directive assignments;
     - `TestRegenerationIsDeclared` is `boundary`;
     - the `internal/config` suiteguard check is `package`;
     - every test named in the Makefile's `REPO_CONTRACT_TESTS` is
       discovered.
   - `TestVerifyChangedRunsTheSelectedContracts`: it runs
     `make -f <repository Makefile> verify-changed-contracts` in a temporary
     directory, with `VERIFY_SELECT` and `GO` set to stub scripts in that
     directory. The stub selector prints two invocation lines, and the stub
     `go` records its arguments. The test asserts the two recorded
     `test -count=1 -tags <tag> -run <pattern> <package>...` invocations in
     order. A stub `go` that fails on the first line makes `make` exit
     non-zero without a second invocation. A failing stub selector makes
     `make` exit non-zero without any invocation. An empty selector output
     runs no `go test`. The test also asserts that the `verify-changed`
     recipe invokes `verify-changed-contracts` after its set loop.
7. MUST leave the real tree so that `_techspec.md` → Surface Transcript 1
   and Surface Transcript 3 hold. A change only under `docs/user-guide`
   selects the always set, and a change to `.agents/skills/roundfix/SKILL.md`
   adds the three regeneration contracts.
8. MUST NOT write outside the test's temporary directories. The stubs are
   created and run there.

## Subtasks

- [ ] Add the directives to the contract test headers.
- [ ] Add `verify-changed-contracts` and invoke it from `verify-changed`.
- [ ] State the known platform limit in the Makefile comment.
- [ ] Write the repository inventory and Makefile tests.
- [ ] Run the new target on the real tree and `make verify-docs`.

## Acceptance Criteria

- [ ] `make verify-changed` ends with the selected contracts and fails when
      one fails.
- [ ] Every contract test file carries the directive ADR-0252 assigns.
- [ ] `make verify-docs` still passes and still runs every contract.

## Context

- interface: `Makefile`
- creates: `internal/verifyselect/contracts_repository_test.go`
- interface: `internal/suiteguard/installation_repocontract_test.go`
- interface: `internal/speccheck/governed_repocontract_test.go`
- interface: `internal/baseline/repository_gate_repocontract_test.go`
- interface: `internal/baseline/derived_regeneration_repocontract_test.go`
- interface: `skills/owned_skill_edit_repocontract_test.go`
- interface: `internal/config/regeneration_declared_test.go`
- interface: `internal/docscontract/command_documentation_test.go`
- interface: `internal/docscontract/commands_index_test.go`
- interface: `internal/docscontract/corpus_test.go`
- interface: `internal/docscontract/model_selection_test.go`
- interface: `internal/docscontract/publicdocs_test.go`
- interface: `internal/docscontract/release_step_test.go`
- interface: `internal/docscontract/secondbrain_export_test.go`
- interface: `internal/docscontract/user_guide_contract_test.go`
- instruction: `internal/verifyselect/verifyselect_test.go`
- instruction: `internal/suiteguardcontract/repository_gate_test.go`
- instruction: `docs/adr/0252-the-selective-gate-runs-the-repository-contracts-a-change-makes-relevant.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestRepositoryContractTestsDeclareTheirRelevance|TestVerifyChangedRunsTheSelectedContracts)$' ./internal/verifyselect 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for t in TestRepositoryContractTestsDeclareTheirRelevance TestVerifyChangedRunsTheSelectedContracts; do printf '%s\n' "$out" | grep -q -- "--- PASS: $t " || { printf 'missing passing test %s\n' "$t" >&2; exit 1; }; done; go test -count=1 ./internal/verifyselect && make verify-changed-contracts && make verify-docs` — expected: exit 0. Before this Task neither test exists and the Makefile has no `verify-changed-contracts` target, so the first check fails. After it, the real tree's directives parse, the new target runs the selected contracts, and `make verify-docs` still runs every contract.

## References

- `_prd.md` → Core Feature 2; Goals; Success Metric 1; Success Metric 3; Success Metric 4; Success Metric 5
- `_techspec.md` → API Contract 2; Surface Transcript 1; Surface Transcript 3; Directive assignments; Invariant 9; Invariant 12; Build Order 2
- ADR-0252

## Result

Implemented the Task 02 slice for Daemon Verification. `verify-changed` now
invokes `verify-changed-contracts` after its set loop, including an empty set
list. The new target captures selector output before running tests, runs one
uncached invocation per selected tag in order, skips empty package lists, and
propagates selector or test failures. Its recipe uses `tests`, preserving the
existing baseline recipe's sole `-run "$$pattern"` occurrence.

All assigned contract headers now declare their Contract Relevance. The
regeneration boundary includes its cost and reason. Existing header sentences
about the pull request boundary describe selective execution and retain
`make verify-docs`. The Makefile comment records the Linux CI platform limit
and the `go list` downloading-stderr example.

Focused evidence by acceptance criterion:

- Selected contracts end the gate and failures stop it:
  `TestVerifyChangedRunsTheSelectedContracts` passed all seven subtests. The
  real Makefile ran with temporary selector and Go scripts: ordered tag and
  package arguments, first-test failure, selector failure, empty selection,
  empty package list, parent execution with empty sets, and propagation of a
  contract failure through the parent recipe. All stub writes stayed under
  each subtest's temporary directory.
- Assigned directives are present:
  `TestRepositoryContractTestsDeclareTheirRelevance` passed against the real
  repository. It checks all eight docscontract files, the suiteguard audit,
  the four governed-set contracts, all four relevant contracts and their
  exact patterns, the regeneration boundary, package-class suiteguard checks,
  and discovery of every name in `REPO_CONTRACT_TESTS`. Its two selection
  subtests prove that a user-guide path selects only the always set and an
  owned-skill path adds exactly the three regeneration contracts.
- The all-contract gate retains its coverage:
  comparison against `HEAD` confirmed byte-identical `verify`, `verify-docs`,
  `docs-test`, `repo-test`, `spec-budget`, `verify-changed-core`, and
  `verify-changed-baseline` recipes and `REPO_CONTRACT_TESTS`.
  `TestEveryRepositoryContractTestRunsInTheRepositoryGate` passed. Four
  focused documentation checks passed, covering active corpus validity,
  user-guide command coverage and links, and roundfix skill command coverage.
  Full `make verify-docs` success remains for Daemon Verification.

Focused commands (each exited 0 with the final implementation):

- `GOCACHE=/private/tmp/roundfix-task02-gocache rtk proxy go test -run 'TestRepositoryContractTestsDeclareTheirRelevance|TestVerifyChangedRunsTheSelectedContracts' -v ./internal/verifyselect`
- `GOCACHE=/private/tmp/roundfix-task02-gocache rtk proxy go test -run '^TestPartitionFollowsTheMakefileRecipes$' ./internal/verifyselect`
- `GOCACHE=/private/tmp/roundfix-task02-gocache rtk proxy go test -run '^TestEveryRepositoryContractTestRunsInTheRepositoryGate$' ./internal/suiteguardcontract`
- `GOCACHE=/private/tmp/roundfix-task02-gocache rtk proxy go test -tags docscontract -run '^(TestCheckActiveCorpusHasNoErrors|TestEveryCommandIsNamedInTheUserGuide|TestUserGuideLinksResolve|TestEveryCommandIsNamedInTheRoundfixSkill)$' ./internal/docscontract`
- `GOCACHE=/private/tmp/roundfix-task02-gocache rtk proxy go test -tags repocontract -run '^TestEverySpawningPackageInstallsTheSuiteGuard$' ./internal/suiteguard`
- `rtk proxy git -c core.fsmonitor=false diff --check`

Before wiring the gate and headers, the two new focused tests failed on the
missing parent invocation and the missing relevance declarations. The first
attempt using the host Go cache was blocked by sandbox permissions; the
Task-scoped cache resolved that environment limitation.

Postflight inspection confirmed unchanged existing test bodies, assertions,
and comments after every package clause; no per-package suiteguard header,
CI workflow, or `.roundfixrc.yml` changed. All new changed paths belong to
this Task's declared slice. The Task Graph and other Task files were untouched.

The authored Verification commands, including real `make
verify-changed-contracts` and `make verify-docs`, were not run during this
Daemon-assigned turn. Status and checkboxes remain for the Daemon to settle.
No commit, push, or pull request was made. No follow-up was identified.

## Carry-forward provenance

- Source Run: `run_20261007T202708Z_1511c9319dd8dbf6`
- Source commit: `186248f2f7e2763bab28670a9ce8b427e30c32cf`
