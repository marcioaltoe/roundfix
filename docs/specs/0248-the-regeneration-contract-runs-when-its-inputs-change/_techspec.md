---
spec: 0248-the-regeneration-contract-runs-when-its-inputs-change
prd: _prd.md
created: 2026-10-07
---

# The regeneration contract runs when its inputs change — Technical Spec

## Executive Summary

`TestRegenerationIsDeclared` (`internal/config/regeneration_declared_test.go`,
build tag `docscontract`) is declared `//verify:boundary`, and no target runs
it: `docs-test` compiles only `./internal/docscontract`, and `repo-test` runs
the `repocontract` names in `REPO_CONTRACT_TESTS`. This Spec makes the test
`relevant` to the inputs its regeneration reads, adds a Full Contract Run
that runs every discovered contract (`verify-select -contracts -all`, `make
verify-contracts`), wires that run into the push-to-main and release
workflows, and makes the selective summary line name the expensive contracts
it leaves out. The trade-off accepted: a pull request that touches the
regeneration inputs (29 of the last 100 commits on main) pays about 52 s more
in its Run gate, and a selection miss is caught on the push to main after the
merge rather than before it. `make verify-docs` keeps its scope and its 45.5 s
warm wall time (ADR-0253).

## Project Constraints

- Identifier strategy: not applicable. No identifier scheme changes. The
  `-all` flag and the `verify-contracts` target are lowercase command words,
  and the selector prints existing build-tag names and Go package paths.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable. No request, credential or network
  call is added. The selector reads Go sources and local Git, and its tests
  use temporary repositories, a stub `go` and a stub selector. No Task or QA
  row triggers a GitHub workflow. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0253 (this Spec) governs the
  derived relevance, the Full Contract Run and the named exclusions.
  ADR-0252: "A later change can make an expensive contract relevant by giving
  it narrow inputs"; this Spec is that change. ADR-0182 runs "the configured repository Verification,
  as the last Verification command" of each Settlement Check, so every
  selected contract gates each Task's settlement. ADR-0184: "A TechSpec
  states a command surface as a transcript", answered in Surface Transcripts.
  `docs/agents/specific-repository.md`: "repository contracts validate at the
  pull request boundary". That stays true; task_04 corrects the sentence that
  says `make verify-docs` runs every contract. Source:
  `docs/agents/domain.md`, `docs/agents/specific-repository.md`.
- Tooling authority: applicable. task_03 changes the `Makefile` and both CI
  workflows, and task_04 changes `docs/agents/specific-repository.md`.
  Express maintainer authorization on 2026-10-07: "Rodar só quando relevante
  (Recommended)", the grant of the Governed Paths this Spec declares including
  the `Makefile`, the explicit authorization of CI workflow changes, and the
  standing "Concedo". Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`. Spec-contained authorization record:
  `docs/specs/0248-the-regeneration-contract-runs-when-its-inputs-change/_authorization.md`.
  Bounded files: `.github/workflows/ci-verify.yml`,
  `.github/workflows/release.yml`, `Makefile`,
  `docs/agents/specific-repository.md`.

## Current behavior

Measured on 2026-10-07 at `ae56aba0`, on the maintainer's 10-core machine,
with the worktree's own `GOCACHE`:

- `make verify-docs` exited 0 in 53.6 s on its first run and 45.5 s warm. In
  the warm run `repo-test` spent 42.7 s in `internal/baseline` and 14.8 s in
  `./skills`, in parallel; `docs-test` took 0.35 s.
- `go test -count=1 -tags docscontract -run '^TestRegenerationIsDeclared$' ./internal/config`
  passed in 52.8 s of wall time (51.6 s in the test).
- The same test run in one `go test` beside the three regeneration contracts
  of ADR-0252, with both tags, took 68 s of wall time, against about 95 s when
  the two tags run one after the other.
- `make test` with `GOFLAGS=-count=1` exited 0 in 310.9 s. `internal/cli`
  took 304 s, `internal/daemon` 133 s, `internal/spec` 116 s and
  `internal/baseline` 105 s; other processes shared the machine during part
  of that run.
- No Makefile target and no workflow names `TestRegenerationIsDeclared`, and
  `.github/workflows/release.yml` runs only `make verify` before publishing.
- Of the last 100 commits on main, 29 touched an input of the relevance
  below, 26 touched an input of ADR-0252's regeneration contracts, 13 touched
  `.roundfixrc.yml` or `internal/baseline/`, and 2 touched `Makefile`,
  `go.mod` or `go.sum`.

## System Architecture

No new package. The work extends three existing seams:

- `internal/verifyselect` (`Run`, `ContractInvocations`, `DiscoverContracts`)
  gains the `-all` mode and the named exclusions in the summary line.
- The header directive of `internal/config/regeneration_declared_test.go`
  changes from `boundary` to `relevant`.
- The `Makefile` gains `verify-contracts`, and the two workflows call it.

```text
push to main / release
  make verify-contracts
    verify-select -contracts -all      DiscoverContracts -> every contract, any class
    go test -count=1 -tags <tag> -run <names> <packages>   (per tag, stop at first failure)
pull request / Run gate (unchanged path)
  make verify-changed -> verify-changed-contracts -> verify-select -contracts -base <base>
    stderr names each relevant contract not selected and each boundary contract
```

## Implementation Design

### Interfaces

```go
// Run gains one flag; no exported function changes signature.
//   -all   with -contracts: print the invocations of every discovered contract.
func Run(ctx context.Context, args []string, stdout, stderr io.Writer) int

// Unchanged and reused by -all: ContractInvocations(all) with the full discovery.
func DiscoverContracts(repoRoot string) ([]ContractTest, error)
func ContractInvocations(selected []ContractTest) []string
```

### Data Models

None. `ContractTest` and `ContractClass` keep their fields and values.

### Directive assignment

`internal/config/regeneration_declared_test.go` replaces its `//verify:boundary`
line with exactly:

```text
//verify:relevant .roundfixrc.yml internal/baseline/ .agents/skills/ skills/ docs/agents/ docs/references/coverage-record.json internal/spec/ internal/cli/baseline_* cmd/roundfix/ internal/suiteguard/ internal/suiteguardcontract/
```

Each pattern answers one input of the regeneration:

- `.roundfixrc.yml` holds the declarations.
- `internal/baseline/` holds the Baseline modules, profiles, formatter
  fixtures, setups, source Baselines, Baseline testdata, the module version
  record and the code of `-record-module-versions` and `make baseline-digests`.
- `.agents/skills/` and `skills/` hold the owned skills whose version lines a
  declaration scopes, `skills/testdata/owned-skill-versions.json`, and the
  code of `-record-skill-versions` and `TestAuthorialSkillSync`.
- `docs/agents/` holds `setup-context.json` and the managed guides that
  `roundfix baseline update` rewrites.
- `docs/references/coverage-record.json` and `internal/spec/` are the coverage
  record and the code that writes it.
- `internal/cli/baseline_*` and `cmd/roundfix/` are the code of
  `roundfix baseline update`.
- `internal/suiteguard/` and `internal/suiteguardcontract/` hold the sanctioned
  regeneration guard the record commands run under.

The package class still adds `internal/config` itself. `Makefile`, `go.mod`
and `go.sum` select every contract, as ADR-0252 decides. The header comment
keeps its other lines; no test body changes.

### API Contracts

1. API Contract: `verify-select -contracts -all`. It discovers every
   Repository Contract Test with `DiscoverContracts` and prints
   `ContractInvocations` of all of them, every class included, on stdout. It
   never runs Git. It prints one stderr line,
   `verify-select: contracts: <n> selected (every Repository Contract Test), 0 not selected`,
   and exits 0. A discovery error prints `verify-select: <error>` and exits 1.
   `-all` without `-contracts` prints `verify-select: -all requires -contracts`
   and exits 2. `-all` with an explicitly set `-base` prints
   `verify-select: -all cannot be combined with -base` and exits 2. The
   existing conflict with `-packages` and `-baseline-cli-pattern` still exits 2
   with its existing message.
2. API Contract: the selective summary line. `verify-select -contracts` keeps
   its stdout and its prefix
   `verify-select: contracts: <s> selected (<a> always, <c> by change), <n> not selected, <b> boundary`.
   When a `relevant` contract was not selected, the line continues with
   `; relevant not selected: ` and the list of those contracts. When a
   `boundary` contract exists, it continues with `; boundary: ` and the list
   of those. Each list item is `<Name> (<package>)`, items are joined by
   `, `, and the order is discovery order: package, then name. The line ends
   with a newline. The fail-safe path, when Git cannot list changes, prints
   its existing diagnostic and then this same line.
3. API Contract: `make verify-contracts`. It runs
   `$(VERIFY_SELECT) -contracts -all` and then, for each printed line
   `<tag> <pattern> <package>...`, runs
   `$(GO) test -count=1 -tags <tag> -run <pattern> <package>...` in order. It
   exits with the first non-zero status, exits non-zero without running
   `go test` when the selector fails, and runs nothing for an empty
   selection. It is listed in `.PHONY` and in `make help`. The recipes of
   `verify`, `verify-changed`, `verify-changed-contracts`, `verify-docs`,
   `docs-test`, `repo-test` and `spec-budget`, and `REPO_CONTRACT_TESTS`, stay
   byte-identical. The new recipe must not contain the text `-run "$$pattern"`,
   which `TestPartitionFollowsTheMakefileRecipes` requires exactly once.
4. API Contract: the workflows. `ci-verify.yml` gains, after its "Verify docs"
   step, a step named `Verify every contract` with
   `if: github.event_name == 'push'` and `run: make verify-contracts`. The
   workflow's triggers, its pull request steps and their order are unchanged.
   `release.yml` gains a step named `Verify every contract` with
   `run: make verify-contracts` and no `if`, placed directly after its
   "Verify gate" step and before "Publication preflight". Each new step
   carries a comment naming the Full Contract Run and ADR-0253.

### Invariants

```text
1. The set -all prints is the set DiscoverContracts returns; no list in the Makefile or a workflow names a contract for it.
2. -all ignores Contract Relevance: always, package, relevant and boundary contracts are all printed.
3. -all never consults Git, so a missing base cannot narrow or fail it.
4. The selective mode selects exactly what it selected before; only its stderr line grows.
5. Every relevant contract the selective mode leaves out, and every boundary contract, is named on its stderr line.
6. A change to .roundfixrc.yml, or to any path a derived_paths declaration names in paths or lines.paths, selects TestRegenerationIsDeclared.
7. A change only to docs/user-guide or to internal/daemon does not select TestRegenerationIsDeclared.
8. The pull request path of ci-verify.yml runs the same steps as before.
```

### Surface Transcripts

1. Surface Transcript: in a temporary clone of the built tree, on a branch
   with one committed change to `docs/user-guide/commands/history.md`, against
   the clone's base.

   ```transcript
   $ go run -buildvcs=false ./cmd/verify-select -contracts -base <base>
   stdout:
   docscontract ^(<names>)$ ./internal/docscontract
   repocontract ^(<names>)$ ./internal/speccheck ./internal/suiteguard
   stderr:
   verify-select: contracts: <s> selected (<a> always, 0 by change), <n> not selected, 0 boundary; relevant not selected: TestDeclaredStepRegenerationAndFrozenBoundaries (internal/baseline), TestMeasuredSanctionedOwnershipMatchesRecords (internal/baseline), TestRepositoryGateRunsTheAnalyzer (internal/baseline), TestRegenerationIsDeclared (internal/config), TestOwnedSkillEditLeavesDerivedArtifactsByteIdentical (skills)
   exit: 0
   ```

2. Surface Transcript: in that clone, on a branch with one committed change
   to `.roundfixrc.yml`, against the clone's base.

   ```transcript
   $ go run -buildvcs=false ./cmd/verify-select -contracts -base <base>
   stdout:
   docscontract ^(<names>)$ ./internal/config ./internal/docscontract
   repocontract ^(<names>)$ ./internal/speccheck ./internal/suiteguard
   stderr:
   verify-select: contracts: <s> selected (<a> always, 1 by change), <n> not selected, 0 boundary; relevant not selected: TestDeclaredStepRegenerationAndFrozenBoundaries (internal/baseline), TestMeasuredSanctionedOwnershipMatchesRecords (internal/baseline), TestRepositoryGateRunsTheAnalyzer (internal/baseline), TestOwnedSkillEditLeavesDerivedArtifactsByteIdentical (skills)
   exit: 0
   ```

3. Surface Transcript: the Full Contract Run's selection at the audited head.

   ```transcript
   $ go run -buildvcs=false ./cmd/verify-select -contracts -all
   stdout:
   docscontract ^(<names>)$ ./internal/config ./internal/docscontract
   repocontract ^(<names>)$ <packages>
   stderr:
   verify-select: contracts: <n> selected (every Repository Contract Test), 0 not selected
   exit: 0
   ```

4. Surface Transcript: `-all` outside the contract mode refuses.

   ```transcript
   $ go run -buildvcs=false ./cmd/verify-select -all
   stdout:
   stderr:
   verify-select: -all requires -contracts
   exit status 2
   exit: 1
   ```

5. Surface Transcript: `-all` with a base refuses.

   ```transcript
   $ go run -buildvcs=false ./cmd/verify-select -contracts -all -base main
   stdout:
   stderr:
   verify-select: -all cannot be combined with -base
   exit status 2
   exit: 1
   ```

## Coverage Map

- Goal "runs when an input changes" → Directive assignment, Invariants 6
  and 7, Surface Transcripts 1 and 2.
- Goal "inputs derived from the declarations" → Directive assignment,
  `TestRegenerationContractRunsWhenADerivedInputChanges`, Invariant 6.
- Goal "every push to main and every release runs every contract" → API
  Contracts 1, 3 and 4, Invariants 1 to 3.
- Goal "no silent exclusion" → API Contract 2, Invariant 5.
- Goal "`make verify-docs` keeps its scope and stays under 60 s" → API
  Contract 3 (byte-identical recipes), Success Metric 5 measured at QA.
- Core Feature 1 → API Contracts 1, 3 and 4.
- Core Feature 2 → API Contract 2.
- Core Feature 3 → Directive assignment and its repository test.
- Core Feature 4 → `CONTEXT.md` and `docs/agents/specific-repository.md`
  (Build Order 4).
- Success Metric 1 → Directive assignment; QA replay.
- Success Metric 2 → Surface Transcripts 1 and 2, Invariants 6 and 7.
- Success Metric 3 → API Contracts 1 and 3, Surface Transcript 3.
- Success Metric 4 → API Contract 4 and its workflow test.
- Success Metric 5 → QA measurement against the 45.5 s baseline.
- Success Metric 6 → every task's `./internal/verifyselect` run and QA.

## Integration Points

- Git, through the existing `ChangedPaths`, only in the selective mode.
- The Go toolchain, through `go test` in the Makefile recipes. The selector
  parses sources with `go/parser`; `-all` adds no `go list` call.
- GitHub Actions. `ci-verify.yml` on push to main and `release.yml` on a tag
  or a manual dispatch call `make verify-contracts`. The release job checks
  out full history and restores no build cache, so its run compiles the
  tagged test binaries from nothing. No Task or QA row runs a workflow; the
  workflow test parses the YAML.

## Testing Approach

- `internal/verifyselect/contracts_test.go` (task_01) builds fixture
  repositories in `t.TempDir()`, as the existing contract tests do, and calls
  `Run` for `-all`, its refusals and the named exclusions. The existing exact
  summary assertion in `TestRunPrintsContractInvocations` changes, because
  its fixture has a boundary contract that the line now names.
- `internal/verifyselect/contracts_repository_test.go` (task_02) reads the
  real repository with `DiscoverContracts`. It asserts the new class and
  paths of `TestRegenerationIsDeclared`. A new test reads `.roundfixrc.yml`
  through `config.ResolveConfigProposal` and selects with
  `SelectContractPaths` for a sample path of every declared entry. A
  directory entry ending in `/` samples a file inside it, and a `*` in a
  pattern samples a literal name. The changed scenario "owned skill selects
  regeneration" expects four contracts by change.
- `internal/verifyselect/full_contract_run_test.go` (task_03, new) runs the
  real `make -f <repository Makefile> verify-contracts` in a temporary
  directory with stub `VERIFY_SELECT` and `GO` scripts, as
  `TestVerifyChangedRunsTheSelectedContracts` does. It also parses both
  workflow files with `gopkg.in/yaml.v3`, already a module dependency.
- The package already installs `suiteguard.Main` and is in the guarded list,
  so the new spawning tests need no new wiring.

## Build Order

1. `-all` mode, its refusals and the named exclusions in the summary line,
   with fixture tests (`internal/verifyselect`).
2. The relevant directive on `TestRegenerationIsDeclared` and the repository
   tests that pin it to the declarations (depends on: none; it shares no file
   with step 1).
3. `make verify-contracts`, the two workflow steps, the Makefile comment and
   the Makefile and workflow tests (depends on: 1, because the target calls
   `-all`).
4. `CONTEXT.md` terms and the hard rule in
   `docs/agents/specific-repository.md` (depends on: none; it describes
   ADR-0253).
5. QA gate (depends on: 1, 2, 3, 4).

## Risks & Considerations

- A change to any regeneration input pays about 52 s in each Settlement
  Check, and about 95 s when ADR-0252's regeneration contracts are selected
  too, because the two tags run in series. Running both tags in one
  `go test` measured 68 s and is a next contributor, not part of this Spec.
- A selection miss, such as a regression in a package the regeneration
  imports but the directive does not name, fails only on the push to main.
  The directive names the generators the declared commands run, not their
  transitive imports; the Full Contract Run is the safety net.
- The release workflow restores no build cache, so `make verify-contracts`
  compiles every tagged test binary there. That adds minutes to a release,
  not to a pull request.
- The test's own package class still selects it for any change to a file in
  `internal/config`, as ADR-0252 decides for every contract. Such a change pays
  the 52 s too; `internal/config` holds the declaration matching the test
  exercises.
- `make verify-docs` and the pull request CI path are unchanged. Their
  measured wall time is the baseline QA compares against.
- The selective summary line grows by the names of five relevant contracts on
  a typical change. That is the visibility the maintainer asked for.

## Glossary

- adds: **Full Contract Run**
- changes: **Repository Contract Test**
- changes: **Contract Relevance**

## Decisions

- `TestRegenerationIsDeclared` is `relevant` to inputs derived from the
  declarations and to the generators they run; see ADR-0253.
- The Full Contract Run is discovery-driven, through `-all`; see ADR-0253.
- CI runs the Full Contract Run as its own step on push to main and in the
  release workflow, not through a switch inside `make verify-docs`; see
  ADR-0253.
- The selective summary line names relevant contracts left out and every
  boundary contract; package-class contracts stay counted only, because each
  is that package's own check.
- `verify-changed-contracts` stays byte-identical; `verify-contracts` repeats
  its loop with `-all`.
