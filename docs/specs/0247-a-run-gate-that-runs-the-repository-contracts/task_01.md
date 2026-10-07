---
task: task_01
spec: 0247-a-run-gate-that-runs-the-repository-contracts
status: completed
type: backend
complexity: medium
---

# Task 01: verify-select discovers Repository Contract Tests and selects them by Contract Relevance

## Overview

Teaches `cmd/verify-select` to find every Repository Contract Test, read its
Contract Relevance directive, and print the `go test` invocations that a
change makes relevant. This Task answers the two CI failures of 2026-10-07
recorded in the operator's queue log, entries 237 and 238. On that day the
Run's `make verify-changed` passed and CI then failed on contract tests it had
never compiled. Nothing calls the new mode yet, which is task_02's work. This
Task is verifiable on its own through fixture tests.

## Requirements

1. MUST add `internal/verifyselect/contracts.go` with the types and functions
   of `_techspec.md` → Interfaces: `ContractClass`, `ContractTest`,
   `DiscoverContracts`, `SelectContractPaths`, `SelectContracts` and
   `ContractInvocations`.
2. MUST discover contracts as `_techspec.md` → Invariant 1 states. It uses
   `go/build/constraint` for the tag and `go/parser` for the functions and
   comments, and it never runs `go list`. Results are sorted by package, then
   name.
3. MUST parse directives by `_techspec.md` → Directive grammar and
   Invariants 2 and 3. A header directive is any `//verify:` line comment
   before the `package` clause, and a test's own doc-comment directive
   overrides it. The error messages are:
   - `<file>: unknown Contract Relevance class "<class>"`
   - `<file>: <test>: unknown Contract Relevance class "<class>"`
   - `<file>[: <test>]: relevant needs at least one path`
   - `<file>[: <test>]: boundary needs a reason`
   - `<file>[: <test>]: more than one Contract Relevance directive`
   - `<file>[: <test>]: invalid relevant path "<pattern>"`
4. MUST select by `_techspec.md` → Invariants 4 to 8. `SelectContracts`
   reuses `ChangedPaths`. When `ChangedPaths` fails, it returns every contract
   that is not `boundary` as selected, together with the wrapped error.
5. MUST format invocations and the summary by `_techspec.md` → Invariants 9
   and 10.
6. MUST add the `-contracts` flag to `Run` in
   `internal/verifyselect/verifyselect.go`. It follows `_techspec.md` → API
   Contract 1 and Invariant 11: exit 1 on a discovery error, exit 2 on a
   flag conflict, and exit 0 with the fail-safe diagnostic. It MUST NOT
   change the output of the existing modes.
7. MUST add `internal/verifyselect/contracts_test.go`. Its fixtures are
   temporary repositories, never the real tree, and it has these tests:
   - `TestContractDiscoveryReadsTagsAndDirectives`: `docscontract` and
     `repocontract` files with header and per-test directives, a test
     without a directive (class `package`), an untagged test file (ignored)
     and a `testdata` directory (not walked).
   - `TestContractDirectiveRefusesMalformedDeclarations`: one case per
     message of requirement 3, plus the patterns `/abs`, `./x`, `a/../b`,
     `a/**` and `[`. Each case asserts the exact message.
   - `TestContractSelectionFollowsChangedPaths`: `always` with an empty
     change list, `package` selected by a direct file and not by a file in a
     subdirectory, `relevant` selected by a directory pattern and by a glob,
     `boundary` never selected, and each of `go.mod`, `go.sum` and `Makefile`
     selecting every contract that is not `boundary`.
   - `TestContractSelectionFailsSafeWithoutAChangeList`: an unknown base ref
     selects every contract that is not `boundary` and returns an error.
   - `TestRunPrintsContractInvocations`: in a fixture Git repository with a
     committed change, `Run` with `-contracts` prints the per-tag lines and
     the summary line exactly. A malformed directive exits 1 with its
     message, and `-contracts -packages core` exits 2. The refusal matches
     `_techspec.md` → Surface Transcript 2.
8. MUST NOT change the `core` and `baseline` set logic, the Makefile, or any
   contract test file. Those are task_02's.

## Subtasks

- [ ] Discover contract tests from build constraints and test declarations.
- [ ] Parse and validate header and per-test directives.
- [ ] Select from changed paths with the module-wide and fail-safe rules.
- [ ] Print the per-tag invocations and the summary, and add `-contracts`.
- [ ] Write the fixture tests.

## Acceptance Criteria

- [ ] `verify-select -contracts` prints one invocation per tag with the
      selected tests.
- [ ] A malformed directive exits 1 and names its file.
- [ ] A missing change list selects every contract that is not `boundary`.
- [ ] Every existing `./internal/verifyselect` test passes unchanged.

## Context

- creates: `internal/verifyselect/contracts.go`
- creates: `internal/verifyselect/contracts_test.go`
- interface: `internal/verifyselect/verifyselect.go`
- instruction: `internal/verifyselect/verifyselect_test.go`
- instruction: `internal/suiteguardcontract/repository_gate_test.go`
- instruction: `docs/adr/0252-the-selective-gate-runs-the-repository-contracts-a-change-makes-relevant.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestContractDiscoveryReadsTagsAndDirectives|TestContractDirectiveRefusesMalformedDeclarations|TestContractSelectionFollowsChangedPaths|TestContractSelectionFailsSafeWithoutAChangeList|TestRunPrintsContractInvocations)$' ./internal/verifyselect 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for t in TestContractDiscoveryReadsTagsAndDirectives TestContractDirectiveRefusesMalformedDeclarations TestContractSelectionFollowsChangedPaths TestContractSelectionFailsSafeWithoutAChangeList TestRunPrintsContractInvocations; do printf '%s\n' "$out" | grep -q -- "--- PASS: $t " || { printf 'missing passing test %s\n' "$t" >&2; exit 1; }; done; go test -count=1 ./internal/verifyselect` — expected: exit 0. Before this Task none of the five tests exists, so the first check fails. After it, discovery, directive refusals, selection, the fail-safe and the CLI mode pass, and so does every existing selector test.

## References

- `_prd.md` → Core Feature 1; Goals; Success Metric 4; Success Metric 5
- `_techspec.md` → API Contract 1; Surface Transcript 2; Invariant 1; Invariant 2; Invariant 3; Invariant 4; Invariant 5; Invariant 6; Invariant 7; Invariant 8; Invariant 9; Invariant 10; Invariant 11; Build Order 1
- ADR-0252

## Result

Implemented the Task 01 selector slice. Discovery walks temporary or supplied
repository sources without `go list`, evaluates modern and legacy build
constraints, recognizes top-level Go test signatures, and sorts contracts by
package and name. Header relevance applies across detached header comments;
per-test doc directives override it. Invalid directives return the required
file- and test-scoped messages.

Selection implements the always, package, relevant and boundary rules,
module-wide changes, and a wrapped Git-error fail-safe. `Run -contracts`
prints sorted, deduplicated per-tag invocations and the stderr summary;
malformed discovery exits 1 and conflicting output flags exit 2. The summary's
not-selected count includes boundary tests, which are also counted separately.
An empty discovered inventory still exits 0 when Git cannot list changes.

Acceptance evidence:

- Per-tag invocations: `TestRunPrintsContractInvocations` asserts exact stdout
  and summary for a committed package change in a fixture Git repository.
  `TestContractInvocationsSortAndDeduplicate` covers tag/package/name ordering,
  duplicate names, regular-expression escaping and an empty selection.
- Malformed directive exits 1 and names its file:
  `TestRunPrintsContractInvocations/malformed_directive` asserts the exact
  Surface Transcript 2 refusal. `TestContractDirectiveRefusesMalformedDeclarations`
  checks every required message in both header and test scopes, including all
  required invalid patterns.
- Missing change list selects every non-boundary contract:
  `TestContractSelectionFailsSafeWithoutAChangeList` exercises an unknown base,
  checks both returned selections and the wrapped error, and asserts CLI exit
  0 with the fail-safe diagnostic. Its empty-repository case covers an empty
  inventory with the same diagnostic behavior.
- Existing tests remain unchanged. Focused regression checks exercised the
  existing classification, Git selection, package partition, Baseline
  dependency/pattern and command-output tests. The two existing complete-test
  partition/Makefile-recipe tests and the full package run remain for Daemon
  Verification; no full-package passing claim is made here.

Focused checks after the final implementation edit:

- `GOCACHE=/private/tmp/roundfix-task01-gocache rtk proxy go test ./internal/verifyselect -run 'TestContract|TestRunPrints|TestClassify|TestDocumentation|TestRoot|TestFixture|TestUnknown|TestSelect|TestPackages|TestBaseline' -count=1 -timeout=120s`
  — exit 0 (`ok roundfix/internal/verifyselect`, 4.782 s). This includes all
  five authored fixture tests plus additional build-constraint and invocation
  edge cases, and the focused existing regressions listed above.
- `rtk proxy git -c core.fsmonitor=false diff --check` — exit 0.

Starting evidence: neither new contract file nor the five fixture tests
existed, and `Run` had no `-contracts` flag. The Task file was already modified
by the Daemon; its status and authored sections were preserved. Newly changed
paths are only the two new contract files, `internal/verifyselect/verifyselect.go`
and this Result section. Task 02 owns directive assignments and Makefile
integration. Declared Verification, settlement and commits remain Daemon-owned.

## Carry-forward provenance

- Source Run: `run_20261007T202708Z_1511c9319dd8dbf6`
- Source commit: `7743dcf928c0b1c6d39d98fdde6ae6eee78ee93e`
