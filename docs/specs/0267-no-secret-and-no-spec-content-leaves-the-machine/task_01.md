---
task: task_01
spec: 0267-no-secret-and-no-spec-content-leaves-the-machine
status: completed
type: backend
complexity: medium
---

# Task 01: A QA Report that carries a secret is never committed

## Overview

Three QA Runs of Spec 0259 committed `make -p` output that held an API key,
and nothing between the QA Report commit and the push looked at it. This Task
adds `internal/secretscan`, the one reviewed table of Secret Patterns and the
functions that turn content or a pushed history into Secret Matches, and makes
the Daemon's QA step refuse to commit a QA Report whose files carry a Secret
Match. The hazard: QA evidence is written by the QA Agent from any command
output, and its destination is the Spec branch, which the Delivery Queue
pushes to GitHub. It is verifiable on its own through the package tests and
the two QA step tests.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-13,
   [QA evidence records environment values, secrets included](references/2026-10-13-qa-evidence-records-environment-values.md),
   for its mechanical half: a QA Report commit never carries a Secret Match.
2. MUST create `internal/secretscan` with the API of `_techspec.md` →
   Interfaces, importing only the standard library. `patterns.go` holds the
   table of `_techspec.md` → Data Models in that order and is the only place a
   pattern is defined; `scan.go` holds `ScanContent`, `PushedHistoryArgs`,
   `ScanPushedHistory`, `Match` and `MatchesError` under API Contracts 1 and 3
   and Invariants 1, 2 and 4.
3. MUST make the QA step compute the QA Report commit's path set once and run
   `ScanContent` over each of those files in the Run Worktree after the QA
   Agent hands back and before the QA Task settles, as `_techspec.md` → API
   Contract 2 states. On a Secret Match it settles the QA Task `failed` with
   the reason of API Contract 2, publishes the one `daemon.qa` event with
   phase `secret_scan`, and skips the QA Report commit; the report, the
   evidence and the Task file stay in the Run Worktree. On no Match nothing
   observable changes (Invariant 3): no new event and the same commit.
4. MUST add `internal/secretscan/scan_test.go` with
   `TestSecretPatternsMatchTheirKindAndNeverEchoTheValue`,
   `TestScanContentReportsPathLineAndKind`,
   `TestScanSkipsBinaryAndPlaceholders` and
   `TestScanPushedHistoryReportsCommitPathAndLine`, as `_techspec.md` →
   Testing Approach states. Every sample MUST be built at runtime by
   concatenation, so this file holds no string a pattern matches.
   `TestScanPushedHistoryReportsCommitPathAndLine` MUST read a synthetic log
   in the shape API Contract 3's arguments produce, written in the test, so
   the package starts no process and needs no suite guard; task_03 proves the
   shape against real Git. Each test calls `t.Parallel()` first (ADR-0259).
5. MUST add `internal/daemon/qa_secret_scan_test.go` with
   `TestQAReportWithASecretMatchSettlesFailedWithoutACommit` and
   `TestQAReportWithoutASecretMatchKeepsItsCommit`, through the task engine's
   existing QA fixtures: a fake QA Agent writes a passing report and one
   evidence file whose second line carries a runtime-built `openrouter-key`
   sample. The first test requires the QA Task `failed`, the reason naming
   the evidence path, line 2 and `openrouter-key`, the `secret_scan` event,
   no QA Report commit in the Run Worktree's history, the evidence file still
   present, and the sample absent from the event, the progress output, the
   Task file and the reason. The second requires today's QA Report commit.
   Each follows ADR-0259 and ADR-0271.
6. MUST NOT change any existing QA step test, the QA Report commit message,
   the QA format command, or what a clean QA step publishes.

## Subtasks

- [ ] Write the Secret Pattern table and the scan functions with their tests.
- [ ] Extract the QA Report commit's path set and scan it before settlement.
- [ ] Add the refused and the clean QA step tests.

## Acceptance Criteria

- [ ] Every kind in the table matches its runtime-built sample, and no Match,
      error or event carries a sample byte run.
- [ ] A QA step whose evidence carries a sample settles `failed` and leaves no
      QA Report commit; a clean QA step commits as before.
- [ ] No file this Task adds holds a string a Secret Pattern matches.

## Consumers

- The QA Report commit refusal — consumers: `./internal/daemon`
- checkpoint: task_03

## Context

- instruction: `docs/adr/0273-no-secret-and-no-spec-content-leaves-the-machine-from-a-run.md`
- instruction: `docs/adr/0259-a-test-runs-in-parallel-unless-it-names-why-it-cannot.md`
- instruction: `docs/adr/0271-a-test-bounds-by-real-time-only-what-it-asserts.md`
- creates: `internal/secretscan/patterns.go`
- creates: `internal/secretscan/scan.go`
- creates: `internal/secretscan/scan_test.go`
- interface: `internal/daemon/task_engine.go`
- creates: `internal/daemon/qa_secret_scan_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestSecretPatternsMatchTheirKindAndNeverEchoTheValue|TestScanContentReportsPathLineAndKind|TestScanSkipsBinaryAndPlaceholders|TestScanPushedHistoryReportsCommitPathAndLine)$' ./internal/secretscan 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestSecretPatternsMatchTheirKindAndNeverEchoTheValue TestScanContentReportsPathLineAndKind TestScanSkipsBinaryAndPlaceholders TestScanPushedHistoryReportsCommitPathAndLine; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the package does not exist.
- `out="$(go test -count=1 -v -run '^(TestQAReportWithASecretMatchSettlesFailedWithoutACommit|TestQAReportWithoutASecretMatchKeepsItsCommit)$' ./internal/daemon 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestQAReportWithASecretMatchSettlesFailedWithoutACommit TestQAReportWithoutASecretMatchKeepsItsCommit; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task neither test exists.

## References

- [_prd.md](_prd.md) — Goals; User Stories 1-2; Core Features 2-3; Success Metric 2
- [_techspec.md](_techspec.md) — API Contract 1; API Contract 2; API Contract 3; Invariants 1-4; Data Models; Testing Approach; Build Order 1
- [2026-10-13-qa-evidence-records-environment-values.md](references/2026-10-13-qa-evidence-records-environment-values.md)
- ADR-0273; ADR-0015; ADR-0259; ADR-0271

## Result

Implemented this Task's slice for Daemon Verification. Task status and authored
Verification remain Daemon-owned; no declared Verification command, repository
Verification, commit, push or pull request was run by this Agent.

### Implementation

- Added the standard-library-only `internal/secretscan` API and the ordered,
  compiled Secret Pattern table in `patterns.go`. Matches hold only path, line,
  kind and optional commit; errors list every location without storing values.
  Content scanning applies the 8000-byte binary check and the first-kind rule.
- Added pushed-history arguments and a process-free patch parser that preserves
  commit, new path and new-file line numbers. It handles multiple hunks, quoted
  paths, additions resembling patch headers, removed lines, deleted files,
  binary-file notices and merge commits without patches.
- Extracted the existing QA Report staging calculation. The QA step computes
  that filtered path set once and scans it after evidence snapshot recording,
  before Task settlement. A Match changes settlement to `failed`, records the
  contracted location-only reason and one `daemon.qa` event with
  `phase: secret_scan`, and skips the report commit. The report, evidence and
  Task file remain in the Run Worktree. Clean reports retain their commit
  message, formatter invocation and existing event phases.

### Focused checks and acceptance evidence

Before the guard was added,
`go test -count=1 -run '^TestQAReportWithASecretMatch' ./internal/daemon`
failed because QA accepted the sample-bearing evidence. After implementation:

- `go test -count=1 ./internal/secretscan` — exit 0.
- `go test -race -shuffle=on -count=1 -run 'TestSecretPatterns|TestScan|TestQAReport(With|Without)ASecretMatch|TestQASecretScan|TestTaskCycleQAVerdictMatrix|TestTaskCycleQAOnlyRun|TestQAFormat(Commits|Formats|Reverts|Leaves|Empty)|TestQAGateCommitsTheEvidenceSnapshot|TestQAGateCarriesARow' ./internal/secretscan ./internal/daemon`
  — exit 0 for both packages after the last code edit.
- `go test -count=1 -overlay /var/folders/_7/68y3l_1s55jcsdmmcmmm4dhh0000gn/T/roundfix-task01-source-audit-cdwnrbwe/overlay.json -run '^TestTask01AddedSourcesContainNoSecretMatches$' ./internal/secretscan`
  — exit 0 after the last code edit. This temporary overlay reads all four added
  Go files and runs `ScanContent` against them; it adds no repository file.
- `git diff --check` — exit 0.

| Acceptance criterion | Focused evidence |
| --- | --- |
| Every table kind matches its runtime sample; diagnostics never echo values | `TestSecretPatternsMatchTheirKindAndNeverEchoTheValue` covers every kind, prefix/operator variants, table order and copied kind names. It checks Match JSON, `String()` and `MatchesError.Error()` for sample byte runs of eight or more. The QA refusal test checks every captured event, progress, settlement reason and Task file for the same leak. |
| Evidence with a Match settles failed without a QA Report commit; clean evidence commits as before | Both required QA tests pass through `TaskCycle`, the existing fake QA Agent and real temporary Git repositories. The refusal names the evidence path, line 2 and `openrouter-key`, emits exactly one refusal event and preserves the evidence bytes and report. The clean test checks the unchanged commit message, committed evidence and existing QA event phases. Existing verdict, formatter, prior-pass carry and snapshot checks also pass. |
| No added file holds a Secret Pattern match | The temporary source audit reports zero Matches across `patterns.go`, `scan.go`, `scan_test.go` and `qa_secret_scan_test.go`. Every synthetic sample is constructed by concatenation at runtime. |

Additional focused coverage in `TestQASecretScanUsesTheReportCommitPathSet`
proves forced report/Task paths, imported evidence and ordinary QA changes are
scanned; Project Config, executable files, pre-existing unrelated dirt and
repository-Verification output are excluded. Binary content and tracked
file deletions contribute no Match. Read errors and cancellation refuse the
scan rather than accepting unread content. Every added top-level test calls
`t.Parallel()` first and introduces no short real-time budget.

### Follow-up boundary

The assigned Task explicitly requires a synthetic pushed-history log and no
scanner-package process. Its test therefore uses synthetic additions followed
by removal, rather than the real-Git example in the TechSpec's Testing Approach.
Task 03 owns proving that shape against real Git and integrating both push
refusals. Judge environment changes and QA guidance remain with their assigned
Tasks.
