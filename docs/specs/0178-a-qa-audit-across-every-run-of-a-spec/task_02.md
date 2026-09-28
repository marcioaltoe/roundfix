---
task: task_02
spec: 0178-a-qa-audit-across-every-run-of-a-spec
status: pending
type: backend
complexity: medium
---

# Task 02: The audit table names each commit

## Overview

The seeded QA Report's `## Authorization audit inputs` table, written by
`WriteMechanicalResult` in `internal/speccheck/report.go`, lists Task, Outcome,
Record, Revision and Detail but no commit. After task_01 a Task can have several
audited commits, and the qa-gate skill still tells the QA Agent to audit every
Task commit by command because it cannot tell which commits the stage already
audited. The table is written by the Daemon into the Spec's `qa/` directory and
read by the QA Agent and by the maintainer reviewing the gate.

## Requirements

1. MUST add `Commit string` to `speccheck.MechanicalAuthorizationRead`, and
   MUST set it to the audited Task commit's SHA on every read
   `detectMechanicalAuthPaths` in `internal/speccheck/mechanical.go` appends:
   granted, refused and unresolved, including the read for an unavailable Task
   commit and for an unavailable delivery target.
2. MUST make `WriteMechanicalResult` write the table header
   `| Task | Commit | Outcome | Record | Revision | Detail |` and the full SHA
   in each row's Commit cell. The `None.` form for no reads stays unchanged.
3. MUST put new tests in `internal/speccheck/mechanical_commit_column_test.go`,
   using the package's existing Git fixture helpers without changing
   `internal/speccheck/mechanical_test.go`:
   - `TestMechanicalAuthorizationReadNamesTheAuditedCommit`: a granted read
     carries the Task commit's SHA;
   - `TestMechanicalUnresolvedAuthorizationReadNamesTheAuditedCommit`: the read
     for a Task commit absent from the repository carries that SHA;
   - `TestMechanicalReportAuditTableHasACommitColumn`: the written report
     contains the new header and the SHA in the row;
   - `TestMechanicalReportAuditTableListsEachCommitOfATask`: two commits of one
     Task produce two rows with their two SHAs.
4. MUST state in the commit-dependent tooling audit of
   `.agents/skills/qa-gate/SKILL.md` that the mechanical stage audits every Task
   commit between the delivery base and the audited head (that phrase), that its
   authorization audit table names each audited commit, and that the gate
   audits by command only a Task commit in that range the table does not list
   and that changes a repository-tooling path, or every Task commit when the
   report's mechanical skips name `Task commits of earlier Runs` (that phrase).
   MUST regenerate `skills/qa-gate/SKILL.md` with `make skills-sync`.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] Every authorization read, including an unresolved one, names its commit.
- [ ] The report's authorization audit table has a Commit column and one row
      per audited commit.
- [ ] The qa-gate skill audits by command only what the stage did not.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- creates: `internal/speccheck/mechanical_commit_column_test.go`
- interface: `internal/speccheck/mechanical.go`
- interface: `internal/speccheck/report.go`
- interface: `.agents/skills/qa-gate/SKILL.md`
- interface: `skills/qa-gate/SKILL.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestMechanicalAuthorizationReadNamesTheAuditedCommit|TestMechanicalUnresolvedAuthorizationReadNamesTheAuditedCommit|TestMechanicalReportAuditTableHasACommitColumn|TestMechanicalReportAuditTableListsEachCommitOfATask)$" ./internal/speccheck 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestMechanicalAuthorizationReadNamesTheAuditedCommit TestMechanicalUnresolvedAuthorizationReadNamesTheAuditedCommit TestMechanicalReportAuditTableHasACommitColumn TestMechanicalReportAuditTableListsEachCommitOfATask; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done && tr -s '[:space:]' ' ' < .agents/skills/qa-gate/SKILL.md | grep -qF -- "every Task commit between the delivery base and the audited head" && tr -s '[:space:]' ' ' < .agents/skills/qa-gate/SKILL.md | grep -qF -- "Task commits of earlier Runs" && diff -r .agents/skills/qa-gate skills/qa-gate >/dev/null` — expected: exit 0; before this Task none of the four new tests exists and the qa-gate skill names neither phrase, so the command fails.

## References

- [_techspec.md](_techspec.md) — The commit column
