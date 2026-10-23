---
task: task_03
spec: 0302-a-subsumed-test-selector-is-flagged-and-a-qa-refusal-names-its-commit
status: completed
type: backend
complexity: medium
---

# Task 03: A QA Report the Daemon closes without an Agent records the commit it audited

## Overview

The Daemon's precondition refusal report and its blocking mechanical report
record no `build`, so a refusal an earlier Run left on the branch reads like
the current result: oraculum's Spec 0087 on 2026-10-21, the adopted Backlog
Entry `references/2026-10-22-a-qa-precondition-refusal-report-records-no-audited-commit.md`.
This Task resolves the Run Worktree's `HEAD` with the auditor evidence and
writes it as `build`, quoted, right after `verdict` in both reports, while the
pending seed an Agent finishes and every report already written keep their
bytes (`_prd.md` → Core Feature 3).

## Requirements

1. MUST add the field `AuditedCommit string` to `spec.AuditorEvidence`
   (`internal/spec/auditor_evidence.go`) with a comment saying it is the
   commit the gate audited, written as a final report's `build` key, and
   empty when not resolved. `ResolveAuditorEvidence` does not set it.
2. MUST make `spec.WritePreconditionRefusalReport` (`internal/spec/qa.go`)
   write `build: <strconv.Quote(commit)>` on the line right after `verdict:`
   when `strings.TrimSpace(evidence.AuditedCommit)` is not empty, using that
   trimmed value, and no `build` key otherwise; every other byte it writes is
   unchanged (`_techspec.md` → API Contract 6).
3. MUST make `resolveQAAuditorEvidence` (`internal/daemon/task_engine.go`)
   set `AuditedCommit` to the result of `qaAuditedHead(ctx, plan.WorkDir)`,
   leaving it empty and still returning the evidence when that call fails, and
   make `mechanicalQAReportContent` write the same `build` line right after
   `verdict:` in the report it writes for a blocking result without a
   precondition refusal, only when `result.Blocking` is true and the trimmed
   commit is not empty. The non-blocking seed keeps its bytes (`_techspec.md`
   → API Contracts 5 and 7).
4. MUST create `internal/spec/qa_build_test.go` with
   `TestPreconditionRefusalReportRecordsItsAuditedCommit` and
   `internal/daemon/qa_report_build_test.go` with
   `TestMechanicalQAReportRecordsTheAuditedCommitWhenItIsFinal`, exactly as
   `_techspec.md` → Testing Approach describes them: the commit
   `1234567e89012345678901234567890123456789`, written with surrounding
   whitespace and expected trimmed and quoted; the prefix `---`, `verdict:
   fail`, `build`, `auditing_binary:`; one `build` key; readability through
   `ReadQAReport` with the refusal unchanged; no key for a blank commit; the
   refusal, the blocking result and the seed through
   `mechanicalQAReportContent`; and the written refusal on the task-cycle
   fixture naming the `HEAD` that `runGitForTest(t, plan.WorkDir,
   "rev-parse", "HEAD")` returns. Every top-level test calls `t.Parallel()`
   first and passes no real-time duration.
5. MUST update, in `internal/daemon/task_engine_test.go`, only
   `TestWriteMechanicalQAReportWritesThePreconditionRefusal`: the expected
   refusal frontmatter of its subtest "the refusal is recorded as one terminal
   row and its frontmatter" and the expected report of its subtest "a gate
   that refused nothing writes the report it wrote before" gain `build:
   <strconv.Quote(head)>` right after `verdict: fail`, where `head` is the
   fixture's `HEAD` read with `runGitForTest`. This is a declared break of
   pinned bytes; no other assertion or test in the file changes.
6. MUST NOT change `spec.QAReport`, `readQAReport`, `PreconditionRefusal`,
   the report shape detector in `internal/speccheck`, the full report
   template, or which report the mechanical stage or `NewestQAReport` reads.

## Subtasks

- [ ] Add the field, resolve it in the Daemon, and write `build` in both final reports.
- [ ] Write the two tests and update the pinned refusal frontmatter.

## Acceptance Criteria

- [ ] The precondition refusal and the blocking mechanical report name the audited commit, quoted, right after `verdict`.
- [ ] The pending seed, a refusal without a resolved commit, and every reader keep their behavior.

## Consumers

- `spec.AuditorEvidence.AuditedCommit` and the `build` line of the refusal report — consumers: `internal/daemon`, `internal/spec`

## Context

- interface: `internal/spec/auditor_evidence.go`
- interface: `internal/spec/qa.go`
- interface: `internal/daemon/task_engine.go`
- interface: `internal/daemon/task_engine_test.go`
- creates: `internal/spec/qa_build_test.go`
- creates: `internal/daemon/qa_report_build_test.go`
- instruction: `internal/spec/qa_test.go`
- instruction: `internal/daemon/qa_hollow_report_test.go`
- instruction: `internal/daemon/qa_report_sequence_test.go`
- instruction: `docs/adr/0132-a-refused-gate-records-the-refusal-as-its-terminal-row.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestPreconditionRefusalReportRecordsItsAuditedCommit|TestMechanicalQAReportRecordsTheAuditedCommitWhenItIsFinal|TestWriteMechanicalQAReportWritesThePreconditionRefusal|TestWritePreconditionRefusalReportWritesOneTerminalRow|TestWritePreconditionRefusalReportIsReadableAsARefusal|TestPreconditionRefusalReportNamesItsAuditor|TestMechanicalQAReportSeedStartsPending|TestWriteMechanicalQAReportRecordsTheRefusal|TestMechanicalQAReportStartsUnsuffixedOnAFreshDate)$' ./internal/spec ./internal/daemon 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestPreconditionRefusalReportRecordsItsAuditedCommit TestMechanicalQAReportRecordsTheAuditedCommitWhenItIsFinal TestWriteMechanicalQAReportWritesThePreconditionRefusal TestWritePreconditionRefusalReportWritesOneTerminalRow TestWritePreconditionRefusalReportIsReadableAsARefusal TestPreconditionRefusalReportNamesItsAuditor TestMechanicalQAReportSeedStartsPending TestWriteMechanicalQAReportRecordsTheRefusal TestMechanicalQAReportStartsUnsuffixedOnAFreshDate; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; GOOS=windows go vet ./internal/spec ./internal/daemon` — expected: exit 0. Before this Task the two new tests do not exist, so the command fails at its first pass check. Afterwards both new tests pass, the updated refusal test passes with the `build` line, the refusal, seed and sequence tests that pin other bytes pass unchanged, and both packages' tests compile for Windows.

## References

- `_prd.md` → Core Feature 3; Goals; Success Metrics 3-4
- `_techspec.md` → API Contracts 5-7; Testing Approach; Build Order 3
- ADR-0132


## Result

Implemented this Task's report provenance slice. `AuditorEvidence.AuditedCommit`
is documented and remains unset by `ResolveAuditorEvidence`. The Daemon resolves
it from the Run Worktree's `HEAD`; an unsuccessful lookup leaves it empty and
returns the other auditor evidence. Both final report writers emit the trimmed,
quoted `build` immediately after `verdict`. The mechanical seed does not emit it.
Only the two declared expectations and their local fixture plumbing changed in
`TestWriteMechanicalQAReportWritesThePreconditionRefusal`.

Acceptance evidence:

- Final refusal and blocking reports: the two new named tests assert the exact
  frontmatter prefix, whitespace trimming, quoting of
  `1234567e89012345678901234567890123456789`, and exactly one `build` key.
  The Daemon test writes a refusal through the task-cycle fixture and compares
  its `build` with `runGitForTest(t, plan.WorkDir, "rev-parse", "HEAD")`.
- Compatibility: new tests assert no key for blank/unresolved commits, unchanged
  seed bytes, unchanged report bytes apart from the added line, unchanged
  refusal readback through `ReadQAReport`, and retained auditor evidence when
  HEAD cannot resolve. Both top-level tests start with `t.Parallel()` and add
  no real-time duration. Readers, shape detection, templates, and report
  selection were not modified.

Focused checks:

- Before the writer changes, the focused run of the two new tests exited 1:
  final report prefixes lacked `build` (captured in
  `/tmp/roundfix-task03-red.log`).
- After the final code/test edits, this focused command exited 0; all seven
  selected top-level tests and their subtests passed:

  ```sh
  go test -count=1 -v -run '^(TestPreconditionRefusalReportRecordsItsAuditedCommit|TestMechanicalQAReportRecordsTheAuditedCommitWhenItIsFinal|TestWriteMechanicalQAReportWritesThePreconditionRefusal|TestWritePreconditionRefusalReportKeepsTheRefusalOnOneLine|TestWritePreconditionRefusalReportRecordsAnUnnamedRefusal|TestMechanicalQAReportSeedStartsPending|TestMechanicalQAReportStartsUnsuffixedOnAFreshDate)$' ./internal/spec ./internal/daemon
  ```

  Invoked through `rtk proxy`; output captured in
  `/tmp/roundfix-task03-focused.log`.
- `git diff --check`: exit 0.

Environment: the sandbox rejected the symlinked writable root; operations used
its resolved `/Volumes/Brain-SSD/cache/.roundfix/worktrees/` checkout with
escalation. The pre-existing `status: in_progress` change was preserved.

Declared Task Verification, Windows vet, and repository Verification remain
Daemon-owned and were not run. No follow-up work identified; no commit, push,
or Pull Request was made.

## Carry-forward provenance

- Source Run: `run_20261023T070317Z_62af7ae2274ed43c`
- Source commit: `820a82968405e1fc8bc11a8f4e0451311c3a5ed5`
