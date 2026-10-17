---
task: task_02
spec: 0277-a-run-never-reaches-the-judge-and-an-amended-plan-keeps-its-settled-tasks
status: completed
type: test
complexity: high
---

# Task 02: The bootstrap test reads only its own release, and every test compiles for Windows

## Overview

The bootstrap serialization test shares one release FIFO among its siblings,
so a release writer that stays open past its sibling's exit gives the next
sibling an end-of-file, which failed the v0.86.0 release gate; and seven test
files call Unix-only symbols without a build constraint, so the tests of three
packages do not compile for Windows (`_techspec.md` → Current behavior). This
Task gives each sibling its own release, moves every Unix-only test into a
`unix`-constrained file, re-records the coverage record, and adds a Repository
Contract Test that compiles every test for Windows. It changes no production
code.

## Requirements

1. MUST answer the operator observation of 2026-10-17 cited in `_prd.md` →
   Sources 2 (GitHub Actions run 39044083565) by implementing API Contract 3
   of `_techspec.md` in `internal/worktree/worktree_test.go`: one release FIFO
   per sibling, each sibling's command reading only its own, and the release
   helper writing only to the FIFO of the sibling whose start it observed. The
   lock collision check, the six rounds, the four siblings and every failure
   message keep their text, and no sleep, retry or timing allowance is added
   (ADR-0271).
2. MUST add `TestBootstrapReleaseWriterOutlivingItsSiblingReachesNoOtherSibling`
   in `internal/worktree/worktree_test.go`, calling `t.Parallel()` first: one
   round in which each release writer, after writing, stays open until the
   next sibling has ended and is closed at the round's end, asserting that
   every sibling succeeds and the shared lock is removed. It waits only on the
   events its assertions read or the test deadline.
3. MUST answer the operator observation of 2026-10-17 cited in `_prd.md` →
   Sources 3 by implementing API Contract 4: every top-level test and helper of
   `internal/agent/acpx_integration_test.go`,
   `internal/baseline/repository_test.go`,
   `internal/cli/deliver_item_probe_report_test.go`,
   `internal/cli/detach_test.go` and `internal/cli/implement_test.go` that uses
   a symbol absent from Windows moves, unchanged and under its name, into the
   new `unix`-constrained files named in Context;
   `internal/cli/implement_detach_teardown_test.go` gains `//go:build unix` and
   `internal/cli/detach_fixture_group_other_test.go` gains `unix` in its
   constraint. A test or helper that another moved test needs moves with it.
   MUST NOT add a Windows stub that skips or fakes a Unix behavior, rename or
   remove any test, or weaken any assertion. Each new file carries the
   `// Suite:` header its neighbours use.
4. MUST re-record `docs/references/coverage-record.json` with
   `go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1`,
   because the moved tests leave the Windows platform set; the record then
   lists each moved test for `darwin` and `linux` only. No record is
   hand-edited.
5. MUST implement API Contract 5: `TestEveryTestCompilesForWindows`, calling
   `t.Parallel()` first, in the new
   `internal/verifyselect/windows_test_compile_repocontract_test.go` with the
   `repocontract` build constraint and no relevance directive, runs
   `go vet -tags 'repocontract docscontract' ./...` from the repository root
   with `GOOS=windows` added to its environment and fails with the command's
   combined output when it exits non-zero (ADR-0285). It reads nothing outside
   the repository and the Go toolchain and writes nothing into the checkout.
6. MUST leave `GOOS=windows go vet -tags 'repocontract docscontract' ./...`
   exiting 0 and `GOOS=windows go build -buildvcs=false ./...` exiting 0.

## Subtasks

- [ ] Give each bootstrap sibling its own release FIFO and add the held-writer test.
- [ ] Move the Unix-only tests and helpers into `unix`-constrained files.
- [ ] Re-record the coverage record with its sanctioned command.
- [ ] Add the Windows compile contract.

## Acceptance Criteria

- [ ] The bootstrap serialization test passes 20 consecutive runs, and the held-writer test passes.
- [ ] Every test of the repository compiles for Windows with both contract tags.
- [ ] The coverage record lists every moved test for Darwin and Linux only, and the equivalence test passes.
- [ ] The Windows compile contract passes and names the failing file when an unconstrained test uses a Unix-only symbol.

## Consumers

- The coverage record's platform lists — consumers: `./internal/spec`
- The Repository Contract Test set — consumers: `./internal/verifyselect`

## Context

- instruction: `docs/adr/0285-every-test-file-compiles-for-windows-and-a-repository-contract-proves-it.md`
- instruction: `docs/adr/0271-a-test-bounds-by-real-time-only-what-it-asserts.md`
- instruction: `docs/adr/0250-a-module-version-is-chosen-when-recorded-and-the-coverage-record-lists-every-platform.md`
- instruction: `internal/worktree/worktree.go`
- instruction: `internal/testwait/testwait.go`
- interface: `internal/worktree/worktree_test.go`
- interface: `internal/agent/acpx_integration_test.go`
- creates: `internal/agent/acpx_integration_unix_test.go`
- interface: `internal/baseline/repository_test.go`
- creates: `internal/baseline/repository_unix_test.go`
- interface: `internal/cli/deliver_item_probe_report_test.go`
- creates: `internal/cli/deliver_item_probe_report_unix_test.go`
- interface: `internal/cli/detach_test.go`
- creates: `internal/cli/detach_unix_test.go`
- interface: `internal/cli/implement_detach_teardown_test.go`
- interface: `internal/cli/detach_fixture_group_other_test.go`
- interface: `internal/cli/implement_test.go`
- creates: `internal/cli/implement_unix_test.go`
- creates: `internal/verifyselect/windows_test_compile_repocontract_test.go`
- interface: `docs/references/coverage-record.json`

## Verification

- `out="$(go test -count=1 -v -run '^(TestBootstrapSerializesAcrossSiblings|TestBootstrapReleaseWriterOutlivingItsSiblingReachesNoOtherSibling)$' ./internal/worktree 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestBootstrapSerializesAcrossSiblings TestBootstrapReleaseWriterOutlivingItsSiblingReachesNoOtherSibling; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; go test -count=20 -run '^TestBootstrapSerializesAcrossSiblings$' ./internal/worktree || exit 1; out="$(GOOS=windows go vet -tags 'repocontract docscontract' ./... 2>&1)" || { printf '%s\n' "$out"; exit 1; }; GOOS=windows go build -buildvcs=false ./... || exit 1; out="$(go test -count=1 -tags repocontract -v -run '^TestEveryTestCompilesForWindows$' ./internal/verifyselect 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- "--- PASS: TestEveryTestCompilesForWindows " || { printf 'missing pass: TestEveryTestCompilesForWindows\n' >&2; exit 1; }; go test -count=1 -run '^TestCoverageEquivalence$' ./internal/spec || exit 1; go test -count=1 ./internal/worktree || exit 1; go vet ./internal/agent ./internal/baseline ./internal/cli ./internal/verifyselect` — expected: exit 0. Before this Task the held-writer test does not exist and the Windows vet fails in three packages, so the command fails; afterwards the bootstrap tests, the Windows compile, the contract, the coverage record and the moved tests pass.

## References

- `_prd.md` → Sources; Goals; User Story 2; User Story 3; Core Feature 4; Core Feature 5; Core Feature 6; Success Metric 2; Success Metric 3
- `_techspec.md` → API Contract 3; API Contract 4; API Contract 5; API Contract 9; Invariants 4 and 5; Current behavior; Testing Approach; Build Order 2
- ADR-0285; ADR-0271; ADR-0266; ADR-0250; ADR-0252; ADR-0253; ADR-0259

## Result

Implemented the test changes for Daemon Verification. Status and authored
checkboxes remain unchanged. No production code changed.

### Implementation

Each bootstrap sibling owns a release FIFO selected by its observed start.
Six rounds, four siblings, collision detection and existing failure messages
remain unchanged. The new parallel regression keeps every writer open until
the round ends, verifies that a later reader reaches no earlier FIFO, and
asserts sibling success and lock removal. Cancellation joins the workers.

Unix-only tests and their necessary callers moved under their original names
into the five new Unix files. The teardown and fallback group-probe files now
exclude Windows. The parallel Windows contract retains default relevance and
reports combined Windows vet output with both contract tags.

Dependency changes beyond Context are necessary for this slice:
`internal/baseline/release_gate_test.go` moves its unchanged journey that calls
the unsafe-alias test; `internal/cli/deliver_item_binary_test.go` and
`internal/cli/deliver_item_binary_rebuild_test.go` move their probe-report
callers unchanged into the Unix report file, removing the emptied rebuild
file; `internal/testfixture/real_time_budget_test.go` transfers the existing
ceiling of two and its reason to the moved detach tests without increasing it.

### Focused evidence

Commands ran through `rtk proxy`.

- `go test ./internal/worktree -race -run '^Test(BootstrapSerializesAcrossSiblings|BootstrapReleaseWriterOutlivingItsSiblingReachesNoOtherSibling|TryWriteBootstrapReleaseReturnsWithoutAReader)$' -count=3`: exit 0. This covers three held-writer runs and eighteen serialization rounds. The authored twenty-run check remains Daemon-owned.
- A temporary overlay restoring the shared FIFO makes the held-writer test exit 1 with a reader reaching a completed sibling's FIFO. No checkout file was changed by that probe.
- Windows test compilation with both contract tags and `go test -c` passed separately for agent, baseline and CLI. Before the moves, agent compilation failed at `acpx_integration_test.go:375`, undefined `syscall.Kill`.
- `go test -tags repocontract ./internal/verifyselect -run 'TestEveryTestCompilesForWindows|TestClassifyPath' -count=1`: exit 0 on the final source tree. This exercises repository-wide Windows vet with both tags.
- The same focused contract invocation with a temporary GOFLAGS overlay adding an unconstrained `syscall.Kill` exits 1 and names `/tmp/roundfix-task02-windows-broken_test.go:26:36` and the undefined symbol.
- Focused moved-flow selections passed in agent, baseline and CLI, including BaselineReleaseGate, delivery report and rebuild flows, detach handshakes and ownership teardown, Implement process-group survival, Verification Capacity and queued cancellation. The real acpx integration remains opt-in; its body is unchanged.
- `go test ./internal/testfixture -run '^TestEvery(RealTimeBudgetIsWithinItsFileCeiling|TestInAParallelTestPackageRunsInParallel)$' -count=1`: exit 0.
- `go test ./internal/spec -run '^Test(CoverageEquivalence|CompareCoverageRecordsReportsMissingTest)$' -count=1`: exit 0 on the final source tree.
- A Go-parser audit retained all 312 original test/helper bodies in the affected agent, baseline and CLI files byte-for-byte. All 20 tests in the new Unix files and constrained teardown file are recorded for exactly darwin and linux. The coverage record was generated exclusively by its sanctioned command; its additional entries also collect pre-existing tests missing from the old record.

These checks supply focused evidence for every acceptance criterion. The
Daemon retains the entire authored Verification, including twenty consecutive
serialization runs and the Windows production build.

### Regeneration limitation and follow-up

The required coverage-record command wrote the record and its test passed,
but its initial process exited 1: suiteguard rejected the changed record.
The generator lacks the in-process DeclareSanctionedRegeneration call that
other generators use, although this Spec declares its command and output.
An earlier concurrent CLI check also observed that record mutation; its final
focused run passed after generation ended.

The declaration belongs in internal/spec/coverage_test.go, a historically
governed path absent from this Spec's grant. An attempted ancillary change
was removed; that file has no final diff. A separately authorized follow-up
must add the declaration before a record-changing regeneration can exit 0 on
its first invocation. Final read-only coverage equivalence passes with the
existing generator source and the generated record.

No commit, push, Pull Request, other Task file or Task Graph edit was made.


### Verification feedback, attempt 1

Inspected the Daemon diagnostic at
`/Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20261017T210659Z_bf9b102645f6f0d3/verification/batch-002-attempt-1.log`.
Its repository-gate completeness audit required the new Windows contract in
the legacy manual list, despite ADR-0253 making discovery the source of every
contract in the Full Contract Run. The existing Makefile already runs that
complete discovered set on main and before release.

Updated `internal/suiteguardcontract/repository_gate_test.go`, a necessary
dependency of installing this Task's new contract. Its audit now recognizes
unlisted contracts only when the Full Contract Run invokes the repository
selector with both discovery flags, pipes its output into the execution loop,
forwards tags, test names and packages, executes nonempty package selections,
and propagates selector and test failures. Legacy recipe checks remain.
The Makefile and governed configuration have no changes.

Expanded the existing missing-contract test with one valid discovery fixture
and ten broken-wiring cases. The valid fixture first failed against the old
audit; all cases pass after the repair. No top-level test name was added, so
the generated coverage inventory is unaffected.

Focused checks, through `rtk proxy`:

- `go test ./internal/suiteguardcontract -run '^TestRepositoryGateAuditNamesAMissingContractTest$' -count=1`: initially exit 1 on the valid full-discovery fixture, establishing the regression.
- `go test ./internal/suiteguardcontract -run '^Test(EveryRepositoryContractTestRunsInTheRepositoryGate|RepositoryGateAuditNamesAMissingContractTest|RepositoryGateAuditNamesADroppedRepositoryWideRun)$' -count=1`: exit 0 after the final repair, including the original failing audit and the valid and broken discovery fixtures.
- `go test ./internal/verifyselect -run '^Test(VerifyContractsRunsEveryContract|MainAndReleaseRunEveryContract|RepositoryContractTestsDeclareTheirRelevance)$' -count=1`: exit 0. These existing checks exercise the actual Makefile loop with command stubs and validate the workflow wiring and contract discovery.
- `git diff --check`: exit 0.

The configured and authored Verification sequences were not rerun. The
Daemon retains their rerun and settlement; status remains unchanged. The
coverage-generator authorization follow-up above remains separate.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `internal/baseline/release_gate_test.go`
- `internal/cli/deliver_item_binary_rebuild_test.go`
- `internal/cli/deliver_item_binary_test.go`
- `internal/suiteguardcontract/repository_gate_test.go`
- `internal/testfixture/real_time_budget_test.go`
