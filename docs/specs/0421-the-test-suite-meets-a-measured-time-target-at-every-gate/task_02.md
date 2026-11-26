---
task: task_02
spec: 0421-the-test-suite-meets-a-measured-time-target-at-every-gate
status: completed
type: test
complexity: high
---

# Task 02: Tests share what they only read, and wait on events

## Overview

This Task delivers `_techspec.md` → Build Order 2. The pinned history is extracted once for each set of paths in a test binary and shared read-only; the Baseline asset-sync tests that assert a later step, and each fault point of `TestTransactionFailureMatrix`, start from a prepared copy instead of running the whole sync again; and the four tests that wait about ten seconds for an allowance to pass end on the event they read. No assertion, test case or platform is removed. It is verifiable on its own: a new test proves the shared extraction is made once and cannot be written, another proves a prepared copy equals a fresh sync, and the four waiting tests pass without their wait.

## Requirements

1. MUST make `gittest.PinnedHistory` extract each set of paths once for a test binary into a directory its callers can read and cannot write, keeping its signature and every caller unchanged. A test that writes into the shared directory fails. Add `TestPinnedHistoryExtractsEachPathSetOnceAndReadOnly` to `internal/gittest/pinned_history_test.go`: two calls with the same paths return the same directory, a call with other paths returns another, and a write into the directory fails.
2. MUST make the asset-sync tests of `internal/baseline` whose assertion is about a step after the sync start from a prepared copy of the synced state, made once for the test binary and copied for each test, and make `TestTransactionFailureMatrix` copy one prepared repository for each fault point. A test that asserts the sync itself keeps running it. Add `TestAPreparedAssetSyncMatchesAFreshSync` to `internal/baseline/assets_sync_test.go`, proving that the prepared copy holds the same tree as a sync run for that test. Every existing test keeps its name, its cases and its assertions.
3. MUST make each of these four tests end on the event its assertion reads, or on an allowance the test gives the code under test, instead of waiting about ten seconds for real time to pass (`_prd.md` → Core Feature 12; ADR-0271): `TestRunDetachedCommandMonitorReportAllowsRunCreationPastLivenessDeadline` (`internal/cli/detach_unix_test.go`), `TestGuardReportsAndKillsHelperResidue` (`internal/suiteguard/suiteguard_test.go`), `TestHelperResidueSkipsAnUnreadableProcessAndKeepsScanning` and `TestHelperResidueSkipsUnreadableProcessesWithoutResidue` (`internal/suiteguard/residue_test.go`). Each keeps what it proves. No new Real-Time Budget is added; a file's ceiling in `internal/testfixture/real_time_budget_test.go` may only fall.
4. MAY give the code under test the seam its test needs, in `internal/cli/detach.go`, `internal/cli/detach_unix.go`, `internal/suiteguard/residue.go` or `internal/suiteguard/suiteguard.go`, with no change to any flag, output, exit code or detached behavior. If `internal/cli/detach.go` changes, it is a source of the Behavior Surfaces `command resolve` and `command watch`, and this Task MUST then record their Coverage Review and the record: prepend to the `review` value of those two entries of `docs/references/skill-coverage.json` one sentence, dated the day this Task writes it, that begins `<YYYY-MM-DD> — Spec 0421 task_02` and says that only a test's allowance for the detach handshake changed, with no flag, output, exit code or detached behavior changed, so the review-runs reference needs no text change, keeping the existing text after it; then run the `-record-skill-coverage` step, `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`, the conditional command `_authorization.md` → Sanctioned regeneration declares (`_prd.md` → Skills). If `internal/cli/detach.go` does not change, neither file changes in this Task.
5. MUST re-record the coverage record last with `go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1` and hand-edit nothing in it. Against the tree task_01 leaves, the record loses no name and gains `TestPinnedHistoryExtractsEachPathSetOnceAndReadOnly` of `roundfix/internal/gittest` and `TestAPreparedAssetSyncMatchesAFreshSync` of `roundfix/internal/baseline`.
6. MUST NOT change any assertion, test case or platform, skip a test, add a `-timeout` to any invocation, change any Baseline module, template, profile, Setup Snapshot or plan golden, or change `internal/gittest/gittest.go`, `internal/baseline/plan_test.go` or any file `_authorization.md` → Limits names.

## Subtasks

- [ ] Extract each set of pinned paths once, read-only, and prove it with a test.
- [ ] Prepare the synced Baseline source once and copy it into the tests that assert a later step and into each fault point.
- [ ] Make the four waiting tests end on their event.
- [ ] Record the Coverage Review and the Behavior Surface Record only if the detach source changed.
- [ ] Re-record the coverage record.

## Acceptance Criteria

- [ ] Two calls of `gittest.PinnedHistory` with the same paths in one test binary share one extraction that cannot be written.
- [ ] A prepared asset-sync copy equals a fresh sync, and every asset-sync and transaction test passes with its assertions unchanged.
- [ ] Each of the four waiting tests passes in well under its former ten seconds.
- [ ] The coverage record gains the two new tests and loses none; the Skill Coverage Map and the Behavior Surface Record are current.

## Context

- instruction: `internal/gittest/suiteguard_repocontract_test.go`
- interface: `internal/gittest/pinned_history.go`
- interface: `internal/gittest/pinned_history_test.go`
- interface: `internal/baseline/assets_sync_test.go`
- interface: `internal/baseline/setup_composition_test.go`
- interface: `internal/baseline/assets_sync_owned_membership_test.go`
- interface: `internal/baseline/assets_sync_retired_test.go`
- interface: `internal/baseline/assets_sync_upstream_rename_test.go`
- interface: `internal/baseline/transaction_test.go`
- interface: `internal/baseline/claude_link_plan_test.go`
- interface: `internal/cli/detach_unix_test.go`
- interface: `internal/cli/detach_test.go`
- interface: `internal/cli/detach.go`
- interface: `internal/cli/detach_unix.go`
- interface: `internal/suiteguard/residue.go`
- interface: `internal/suiteguard/residue_test.go`
- interface: `internal/suiteguard/suiteguard.go`
- interface: `internal/suiteguard/suiteguard_test.go`
- interface: `internal/testfixture/real_time_budget_test.go`
- interface: `docs/references/skill-coverage.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `docs/references/coverage-record.json`

## Consumers

- `gittest.PinnedHistory`, now one shared read-only extraction for each set of paths — consumers: `./internal/spec`, `./internal/speccheck`, `./internal/judge`, `./internal/authorization`

## Verification

- `grep -q '^func TestPinnedHistoryExtractsEachPathSetOnceAndReadOnly(' internal/gittest/pinned_history_test.go || exit 1; out="$(go test -count=1 -v -run '^TestPinnedHistoryExtractsEachPathSetOnceAndReadOnly$' ./internal/gittest 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestPinnedHistoryExtractsEachPathSetOnceAndReadOnly ' || exit 1; go test -count=1 ./internal/spec ./internal/speccheck ./internal/judge ./internal/authorization` — expected: exit 0, the pinned history is extracted once and read-only, and every package that reads it passes.
- `grep -q '^func TestAPreparedAssetSyncMatchesAFreshSync(' internal/baseline/assets_sync_test.go || exit 1; out="$(go test -count=1 -v -run '^(TestAPreparedAssetSyncMatchesAFreshSync|TestTransactionFailureMatrix|TestAssetSyncStillRefusesAnInvalidCatalogWithoutDrift|TestAnAssetSyncWithNoSourceChangeLeavesTheComposedSetup|TestAnAssetSyncRewritesTheComposedSetupWithItsComponents)$' ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestAPreparedAssetSyncMatchesAFreshSync TestTransactionFailureMatrix TestAssetSyncStillRefusesAnInvalidCatalogWithoutDrift TestAnAssetSyncWithNoSourceChangeLeavesTheComposedSetup TestAnAssetSyncRewritesTheComposedSetupWithItsComponents; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done` — expected: exit 0, a prepared copy equals a fresh sync and the costliest asset-sync and transaction tests pass on it.
- `out="$(go test -count=1 -v -run '^TestRunDetachedCommandMonitorReportAllowsRunCreationPastLivenessDeadline$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | awk '$1 == "---" && $2 == "PASS:" && $3 == "TestRunDetachedCommandMonitorReportAllowsRunCreationPastLivenessDeadline" { seconds = $4; gsub(/[()s]/, "", seconds); found = 1; if (seconds + 0 >= 5) slow = 1 } END { exit !(found && !slow) }' || { printf '%s\n' "$out" | grep -- '--- ' >&2; echo "the detach test still waits" >&2; exit 1; }; out="$(go test -count=1 -v -run '^(TestGuardReportsAndKillsHelperResidue|TestHelperResidueSkipsAnUnreadableProcessAndKeepsScanning|TestHelperResidueSkipsUnreadableProcessesWithoutResidue)$' ./internal/suiteguard 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestGuardReportsAndKillsHelperResidue TestHelperResidueSkipsAnUnreadableProcessAndKeepsScanning TestHelperResidueSkipsUnreadableProcessesWithoutResidue; do printf '%s\n' "$out" | awk -v n="$name" '$1 == "---" && $2 == "PASS:" && $3 == n { seconds = $4; gsub(/[()s]/, "", seconds); found = 1; if (seconds + 0 >= 5) slow = 1 } END { exit !(found && !slow) }' || { printf '%s\n' "$out" | grep -- '--- ' >&2; echo "still waits: $name" >&2; exit 1; }; done` — expected: exit 0, each of the four tests passes in under five seconds instead of about ten.
- `test -f docs/references/coverage-record.json && for name in TestPinnedHistoryExtractsEachPathSetOnceAndReadOnly TestAPreparedAssetSyncMatchesAFreshSync; do grep -q "\"$name\"" docs/references/coverage-record.json || { echo "record lacks $name" >&2; exit 1; }; done && go test -count=1 -run '^TestCoverageEquivalence$' ./internal/spec && go test -count=1 -timeout 20m -tags docscontract -run '^TestTheSkillCoverageMapIsCurrent$' ./internal/docscontract` — expected: exit 0, the coverage record holds the two new tests and passes, and the Skill Coverage Map and the Behavior Surface Record are current.

## References

- `_prd.md` → Goal 4; User Story 6; Core Features 10, 11, 12 and 13; Success Metric 9; Skills; Categories of work; Sizing, Task 2
- `_techspec.md` → Build Order 2; Testing Approach; ADR Alternatives (ADR-0259, ADR-0271, ADR-0281); Risks & Considerations
- `_authorization.md` → What the grant covers; Sanctioned regeneration; Limits
- [`references/sources.md`](references/sources.md), sections 5.2 and 5.3
- ADR-0259, ADR-0271, ADR-0281


## Result

Skills loaded: implement-task, coding-guidelines, golang-testing, testing-boss, golang-concurrency, golang-cli, golang-context, golang-error-handling, tech-writer, systematic-debugging, no-workarounds, evidence-gate.

Implementation:

- `PinnedHistory` caches each normalized path set under the physical repository identity, serializes extraction, and returns the same directory to every caller. Files and directories are read-only throughout every caller's test lifetime. Reader accounting restores directory permissions after the last reader ends so the Go tool can clean its build directory; later callers reuse and protect the existing extraction. Its signature and consumers are unchanged. The new test checks shared identity, different path sets, failed directory/file writes, reordered/duplicate paths, and reuse after the first reader ends.
- Baseline tests that inspect a later sync step copy an immutable, once-prepared synced repository: composed-component rewriting, unchanged composition, and invalid-catalog refusal without drift. Tests asserting sync normalization, provenance, rollback, membership, retirement, or rename still exercise their original sync boundary. The new equivalence test compares both bytes and modes against a fresh sync. Every transaction fault point copies one prepared repository and reuses its portable plan; all existing cases and assertions remain.
- The detach fixture expires its old liveness timer before releasing the child's Run-creation report. Production timer defaults remain unchanged. Residue tests control retry permission or give the fixture guard zero cleanup grace; production still uses its original grace, scan, diagnostics and proven kill. No Real-Time Budget or ceiling was added.
- Prepended the required dated Coverage Review to `command resolve` and `command watch`. Re-recording the Behavior Surface Record produced identical bytes, consistent with unchanged command behavior. Re-recorded the coverage inventory last through its sanctioned recorder: exactly the two required new tests were added, no test was removed, and all three recorded platforms remain.

Focused checks (all commands used `rtk proxy`):

- `go test -count=1 -json -run 'Test.*(AssetsSync|AssetSync|Transaction|ComposedSetup|ComponentsThat)' ./internal/baseline`: exit 0; 35 top-level tests, including prepared/fresh equivalence, and all 61 transaction matrix subcases passed. Package elapsed 46.180 s. An earlier invocation's assertions passed but suiteguard rejected concurrent source edits; the reported invocation ran on a stable tree.
- `go test -race -count=1 -json -run 'TestPinnedHistory|TestRunDetachedCommand|TestHelperResidue|TestGuardReportsAndKillsHelperResidue' ./internal/gittest ./internal/cli ./internal/suiteguard`: exit 0 for all three packages. The four target waits took 0.07 s (detach), 0.14 s (guard residue), 0.06 s (unreadable beside residue), and 0.00 s (unreadable without residue). Shared-cache identity, protection and later-reader reuse passed with race detection.
- `go test -count=1 -json -run 'TestPinnedHistory|Test.*Grouping|Test.*Question' ./internal/gittest ./internal/judge`: exit 0; exercises unchanged pinned-history consumers, including a package without suiteguard lifecycle cleanup.
- `go test -count=1 -json ./internal/suiteguard`: exit 0; repository boundary, residue killing, unreadable-process handling and regeneration diagnostics passed.
- `go test -count=1 -json -run 'TestEveryRealTimeBudgetIsWithinItsFileCeiling|TestRealTimeBudgetRuleRejectsSeededViolations' ./internal/testfixture`: exit 0; existing ceilings and their negative probes passed unchanged.
- `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`: exit 0 after the final detach source edit; Skill Coverage Map and Behavior Surface Record are current.
- `go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1`: exit 0. A separate parsed comparison against `HEAD` confirmed exactly `roundfix/internal/gittest: TestPinnedHistoryExtractsEachPathSetOnceAndReadOnly` and `roundfix/internal/baseline: TestAPreparedAssetSyncMatchesAFreshSync` added, zero removed, and unchanged platforms. Neither generated record was hand-edited.

The declared Verification commands and repository Verification remain for the Daemon. Task status, the Task Graph and other Task files were not edited. No commit, push or pull request was made.

## Carry-forward provenance

- Source Run: `run_20261126T024343Z_4d6d0cf3419b8574`
- Source commit: `e89f17afd558553fe65a28717091040a983e835d`
