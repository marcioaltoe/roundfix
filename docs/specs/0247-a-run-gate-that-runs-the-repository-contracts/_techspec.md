---
spec: 0247-a-run-gate-that-runs-the-repository-contracts
prd: _prd.md
created: 2026-10-07
---

# A Run gate that runs the repository contracts — Technical Spec

## Executive Summary

`make verify-changed` (`Makefile`, target `verify-changed`) asks
`cmd/verify-select` which test sets a change selects, `core` and `baseline`,
and runs `go test` without build tags. No `docscontract` or `repocontract`
test file is ever compiled there. Those Repository Contract Tests run only in
`make verify-docs`, through `docs-test` (`./internal/docscontract` only) and
`repo-test` (`REPO_CONTRACT_TESTS` over `./...`). So a Run and its QA gate
pass while CI fails, as happened twice on 2026-10-07 (`_prd.md`).

The fix gives each Repository Contract Test a Contract Relevance, read from a
`//verify:` directive in its source. `cmd/verify-select -contracts` then
selects the contracts a change makes relevant and prints one `go test`
invocation per build tag, and a new `verify-changed-contracts` target ends
`verify-changed` by running them. The trade-off accepted: a change to Baseline
assets or an owned skill now pays about 45 s for the regeneration contracts in
every Run gate. Those are the contracts such a change can break, and a
typical core change pays about 3 s (ADR-0252).

## Project Constraints

- Identifier strategy: not applicable. No identifier scheme changes. The
  directive classes are lowercase words, and the selector prints the existing
  build-tag names and Go package paths. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable. No request, credential or network
  call is added. The selector reads Go sources and local Git, and its tests
  use temporary repositories, a stub `go` and a stub selector. Source:
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0252 (this Spec) governs Contract
  Relevance, the always set, the budget and the platform limit. ADR-0182 runs
  "the configured repository Verification, as the last Verification command"
  of each Settlement Check, so every selected contract also gates each Task's
  settlement. ADR-0184: "A TechSpec states a command surface as a transcript",
  answered in Surface Transcripts. `docs/agents/specific-repository.md`:
  "repository contracts validate at the pull request boundary". That stays
  true: `make verify-docs` and `repo-test` are unchanged. Source:
  `docs/agents/domain.md`, `docs/agents/specific-repository.md`.
- Tooling authority: applicable. task_02 changes the `Makefile` and the
  headers of four governed contract test files, and task_03 changes
  `docs/agents/specific-repository.md`. Express maintainer authorization on
  2026-10-07: "Contratos sempre no gate (Recommended)", the grant of the
  Governed Paths this Spec declares including the Makefile, and the standing
  "Concedo". Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`. Spec-contained authorization record:
  `docs/specs/0247-a-run-gate-that-runs-the-repository-contracts/_authorization.md`.
  Bounded files: `Makefile`, `docs/agents/specific-repository.md`,
  `internal/baseline/derived_regeneration_repocontract_test.go`,
  `internal/docscontract/publicdocs_test.go`,
  `internal/speccheck/governed_repocontract_test.go`,
  `skills/owned_skill_edit_repocontract_test.go`.

## Current behavior

Measured on 2026-10-07 at `59548f13`, on the maintainer's 10-core machine:

- `make verify-changed VERIFY_BASE=origin/main` in a temporary clone took
  247 s and exited 0. The branch had one commit adding a comment line to
  `internal/config/config.go`, and the run used `GOFLAGS=-count=1` and the
  worktree's warm `GOCACHE`. It selected `core` and compiled no tagged file.
- In that clone, deleting `"internal/config",` from `guardedSpawningPackages`
  in `internal/suiteguardcontract/contract.go` let
  `go test -count=1 ./internal/suiteguardcontract ./internal/suiteguard ./internal/config`
  pass (606 tests). `go test -count=1 -tags repocontract -run '^TestEverySpawningPackageInstallsTheSuiteGuard$' ./internal/suiteguard`
  failed with "internal/config spawns a process but is not enumerated as
  guarded". This is the CI failure of 2026-10-07, and the selective gate
  could not see it.
- `make repo-test` took 43 s warm and 51 s on its first run.
  `TestMeasuredSanctionedOwnershipMatchesRecords` and
  `TestDeclaredStepRegenerationAndFrozenBoundaries` took 42.6 s and 39.2 s in
  parallel. `TestOwnedSkillEditLeavesDerivedArtifactsByteIdentical` took
  15.3 s, `TestCleanupHistoricalGrantEvidence` 1.2 s, and each other test
  0.5 s or less.
- `go test -count=1 -tags docscontract ./internal/docscontract ./internal/config`
  took 44 s, of which `TestRegenerationIsDeclared` took 42.9 s. The
  `internal/docscontract` tests took 0.34 s together. No Makefile target or
  workflow runs `TestRegenerationIsDeclared` today.
- Here is the extra cost of the always set over packages already compiled
  without tags. Each figure is warm, then with an empty `GOCACHE`. The
  suiteguard audit cost 0.7 s, then 1.0 s. The `internal/docscontract` tests
  cost 0.8 s, then 7.4 s. The four governed-set contracts in
  `internal/speccheck` cost 1.1 s, then 2.4 s.

## System Architecture

`internal/verifyselect` gains contract discovery and selection, and the
existing `Run` gains the `-contracts` mode. It reuses `ChangedPaths` for the
change list and the `go/parser` approach already used by
`BaselineCLITestNames`. `cmd/verify-select` stays a thin `main`. The Makefile
gains one target, which `verify-changed` invokes after the set loop. No new
package is created, and no test body changes.

```text
make verify-changed
  fmt-check, vet, build
  verify-changed-core / verify-changed-baseline   (unchanged)
  verify-changed-contracts
    verify-select -contracts -base $(VERIFY_BASE)
      DiscoverContracts(repo) -> []ContractTest   (build tags + //verify: directives)
      ChangedPaths(repo, base) -> []string        (fail safe: every non-boundary)
      SelectContractPaths -> selected
      ContractInvocations -> "<tag> <pattern> <package>..." per tag
    go test -count=1 -tags <tag> -run <pattern> <package>...   (per line, stop at first failure)
```

## Implementation Design

### Interfaces

```go
type ContractClass string

const (
	ContractAlways   ContractClass = "always"
	ContractPackage  ContractClass = "package"  // no directive
	ContractRelevant ContractClass = "relevant"
	ContractBoundary ContractClass = "boundary"
)

type ContractTest struct {
	Name, Package, Tag, File string // Package and File are repository-relative, slash-separated
	Class                    ContractClass
	Paths                    []string // relevant only
	Reason                   string   // boundary only
}

func DiscoverContracts(repoRoot string) ([]ContractTest, error)
func SelectContractPaths(contracts []ContractTest, changed []string) []ContractTest
func SelectContracts(ctx context.Context, repoRoot, baseRef string) ([]ContractTest, []ContractTest, error) // all, selected
func ContractInvocations(selected []ContractTest) []string
```

### Data Models

None. There is no persisted state. A directive is a source line comment.

### Directive grammar

```text
//verify:always
//verify:relevant <pattern> [<pattern>...]
//verify:boundary <reason>
```

- A directive in a contract test file's header, which is any comment before
  the `package` clause, applies to every test function in the file. A
  directive in a test function's doc comment overrides it for that test.
- A pattern is repository-relative and slash-separated. A pattern ending in
  `/` matches every path under that directory. Any other pattern matches when
  `path.Match(pattern, changed)` is true. A pattern must not start with `/` or
  `./`, contain a `..` segment or `**`, or fail `path.Match`'s syntax check.

### Invariants

1. A Repository Contract Test is a top-level `func TestXxx(t *testing.T)` in
   a `_test.go` file whose build constraint is true when exactly one of
   `docscontract` and `repocontract` is set, and false when neither is set.
   That tag is its `Tag`. A constraint that names one of the tags and fits
   neither case is a discovery error. Directories named `.git`, `testdata`,
   `vendor` and `node_modules`, and any directory whose name starts with `.`,
   are not walked.
2. A test with no directive in either scope has class `package`.
3. Each of the following is a discovery error that names the file, and the
   test when there is one: an unknown class, `relevant` without a pattern,
   `boundary` without a reason, two directives in one scope, and an invalid
   pattern. `verify-select -contracts` then exits 1, so `make verify-changed`
   fails. A malformed declaration never narrows the gate.
4. `always` is selected on every change, including an empty change list.
5. `boundary` is never selected by `-contracts`.
6. `package` is selected when a changed path's `path.Dir` equals its package
   directory. `relevant` is selected in that case too, and also when a
   changed path matches one of its patterns. A file in a subdirectory of the
   package does not select it.
7. A changed `go.mod`, `go.sum` or `Makefile` selects every contract that is
   not `boundary`.
8. When `ChangedPaths` fails, `SelectContracts` selects every contract that
   is not `boundary` and returns the error. `Run` prints the error to stderr
   followed by `; selecting every contract that is not boundary`, and exits 0
   with the invocations. This matches today's fail-safe for sets.
9. `ContractInvocations` prints one line per tag that has a selected test,
   ordered by tag name: `<tag> ^(<names>)$ <package>...`. The names are the
   unique, sorted, `regexp.QuoteMeta`-quoted selected names of that tag. The
   packages are the unique, sorted `./<dir>` paths of that tag's selected
   tests. A test outside the selection that shares a selected name in a
   selected package also runs. In this repository that is only each
   package's own suiteguard check, which is intended.
10. `Run -contracts` prints one stderr summary line:
    `verify-select: contracts: <s> selected (<a> always, <c> by change), <n> not selected, <b> boundary`.
    `<c>` counts selected `package` and `relevant` tests.
11. `-contracts` combined with `-packages` or `-baseline-cli-pattern` exits 2.
12. `verify-changed-contracts` runs `$(GO) test -count=1 -tags <tag> -run <pattern> <package>...`
    for each printed line, in order. It exits with the first non-zero status,
    and exits non-zero when the selector fails. It never runs `go test` with
    an empty package list. `verify-changed` invokes it after the set loop on
    every run, including a run whose set list is empty.

### API Contracts

1. API Contract: `go run -buildvcs=false ./cmd/verify-select -contracts [-base <ref>] [-repo <dir>]`.
   Input: a repository root and a base ref, with `main` and `.` as the
   defaults. Output: Invariant 9 on stdout and Invariant 10 on stderr.
   Failures: exit 1 on a discovery error (Invariant 3), exit 2 on a usage
   error (Invariant 11), and exit 0 with a diagnostic when no change list can
   be produced (Invariant 8).
2. API Contract: `make verify-changed-contracts` and `make verify-changed`.
   They take the existing `VERIFY_BASE`, `GO` and `VERIFY_SELECT` variables.
   The behavior is Invariant 12. `make verify`, `make verify-docs`, `docs-test`
   and `repo-test` keep their recipes.

### Directive assignments

| File | Directive |
| --- | --- |
| `internal/suiteguard/installation_repocontract_test.go` | `//verify:always` |
| `internal/speccheck/governed_repocontract_test.go` | `//verify:always` |
| the eight `docscontract` files in `internal/docscontract` | `//verify:always` |
| `internal/baseline/repository_gate_repocontract_test.go` | `//verify:relevant internal/baseline/analyzer/` |
| `internal/baseline/derived_regeneration_repocontract_test.go` | `//verify:relevant internal/baseline/assets/ internal/baseline/testdata/ .agents/skills/ skills/` |
| `skills/owned_skill_edit_repocontract_test.go` | `//verify:relevant internal/baseline/assets/ internal/baseline/testdata/ .agents/skills/ skills/` |
| `internal/config/regeneration_declared_test.go` | `//verify:boundary` with a reason saying it regenerates Baseline derived artifacts in a repository copy (about 43 s) and is not yet selected by its inputs |
| every `suiteguard_repocontract_test.go` per-package check | none (`package`) |

The header comments that say "make verify-docs runs it at the pull request
boundary" change to say that the selective gate runs the test by its Contract
Relevance and that `make verify-docs` still runs it.

### Surface Transcripts

1. Surface Transcript: in a temporary clone of the built tree, a committed
   change to one `docs/user-guide` file selects only the always set.

   ```transcript
   $ go run -buildvcs=false ./cmd/verify-select -contracts -base <base>
   stdout:
   docscontract ^(<names>)$ ./internal/docscontract
   repocontract ^(<names>)$ ./internal/speccheck ./internal/suiteguard
   stderr:
   verify-select: contracts: <s> selected (<a> always, 0 by change), <n> not selected, 1 boundary
   exit: 0
   ```

2. Surface Transcript: in that clone, a directive `//verify:sometimes` added
   to one contract test file refuses.

   ```transcript
   $ go run -buildvcs=false ./cmd/verify-select -contracts -base <base>
   stdout:
   stderr:
   verify-select: <file>: unknown Contract Relevance class "sometimes"
   exit: 1
   ```

3. Surface Transcript: in that clone, a committed change to
   `.agents/skills/roundfix/SKILL.md` adds the regeneration contracts.

   ```transcript
   $ go run -buildvcs=false ./cmd/verify-select -contracts -base <base>
   stdout:
   docscontract ^(<names>)$ ./internal/docscontract
   repocontract ^(<names>)$ ./internal/baseline ./internal/speccheck ./internal/suiteguard ./skills
   stderr:
   verify-select: contracts: <s> selected (<a> always, 3 by change), <n> not selected, 1 boundary
   exit: 0
   ```

## Coverage Map

- Goal "cheap global contracts on every change" → Directive assignments
  (always), Invariant 4, `verify-changed-contracts`.
- Goal "other contracts by package or declared inputs" → `SelectContractPaths`,
  Invariants 6 and 7, Directive assignments.
- Goal "a contract declares its own inputs" → Directive grammar,
  `DiscoverContracts`.
- Goal "budget of about 20 s" → the always set chosen by measured cost in
  Current behavior; Success Metric 2 is measured at QA.
- Goal "never narrows" → Invariants 3 and 8.
- Core Feature 1 → `DiscoverContracts`, `SelectContracts`,
  `ContractInvocations`, API Contract 1.
- Core Feature 2 → the Makefile target, API Contract 2, Directive assignments.
- Core Feature 3 → `CONTEXT.md` and `docs/agents/specific-repository.md`
  (Build Order 3).
- Success Metric 1 → Invariant 4 (suiteguard audit always) and the QA replay.
- Success Metric 2 → QA measurement against the 247 s baseline.
- Success Metric 3 → Surface Transcripts 1 and 3.
- Success Metric 4 → Invariants 3 and 8, Surface Transcript 2.
- Success Metric 5 → task_01 and task_02 tests, and QA.

## Integration Points

- Git, through the existing `runGit` and `ChangedPaths`.
- The Go toolchain, through `go test` in the Makefile recipe. The selector
  itself parses sources with `go/parser` and `go/build/constraint` and never
  runs `go list`. Spec 0245 showed that `go list` stderr can differ on Linux.
- GitHub CI is unchanged. Its "Verify changed" step now also runs the
  selected contracts, and "Verify docs" still runs all of them.

## Testing Approach

- `internal/verifyselect/contracts_test.go` (task_01) uses fixture
  repositories in `t.TempDir()`, with `git init` where a change list is
  needed, as the existing `TestSelect*` tests do. The package already installs
  `suiteguard.Main` and is in the guarded list.
- `internal/verifyselect/contracts_repository_test.go` (task_02) reads the
  real repository with `DiscoverContracts` and asserts the Directive
  assignments by name. It also runs the real `Makefile` target
  `verify-changed-contracts` with `make -f <repo>/Makefile` in a temporary
  directory, with `GO` and `VERIFY_SELECT` pointing at stub scripts that
  record their arguments. It asserts the invocations and the propagation of
  a failing stub, and it asserts that the `verify-changed` recipe invokes the
  target. These tests are untagged, so the `core` set runs them.
- No contract test body changes. The proof that the gate now sees the class
  of failure is the QA replay of Success Metric 1.

## Build Order

1. Contract discovery, directive parsing, selection, invocations and the
   `-contracts` mode of `verify-select`, with fixture tests (`internal/verifyselect`).
2. Directives on the existing contract test files, the
   `verify-changed-contracts` target and its invocation from `verify-changed`,
   the Makefile comment stating the known platform limit, and the
   repository inventory and Makefile tests (depends on: 1).
3. `CONTEXT.md` terms **Repository Contract Test** and
   **Contract Relevance**, and the hard rule in
   `docs/agents/specific-repository.md` (depends on: none; it is independent
   of 1 and 2 and describes ADR-0252).
4. QA gate (depends on: 1, 2, 3).

## Risks & Considerations

- A change in the Baseline assets or an owned skill pays about 45 s in every
  Settlement Check after it lands in the Run. The ADR records this as
  deliberate. Narrowing those inputs further is a later change.
- In CI the selected contracts run twice, once in "Verify changed" and once
  in "Verify docs". This costs seconds on a typical change, and the CI
  workflow stays unchanged.
- Spec 0248 will change `TestRegenerationIsDeclared` from `boundary` to
  `relevant`. Its header directive is the only line it needs.
- A contract added without a directive defaults to `package`, which is never
  narrower than the gate is today.

## Glossary

- adds: **Repository Contract Test**
- adds: **Contract Relevance**

## Decisions

- Contract Relevance lives in `//verify:` directives next to the tests; see
  ADR-0252.
- The always set is the suiteguard audit, the `internal/docscontract` tests
  and the governed-set contracts, chosen by measured cost; see ADR-0252.
- One `go test` per build tag with a union pattern, so the selected packages
  of a tag run in parallel.
- The CI workflows stay unchanged, and `make verify-docs` keeps running every
  contract.
- Platform-only failures are a documented known limit; see ADR-0252.
