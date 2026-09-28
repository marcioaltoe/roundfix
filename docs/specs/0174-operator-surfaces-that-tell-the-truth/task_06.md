---
task: task_06
spec: 0174-operator-surfaces-that-tell-the-truth
status: completed
type: backend
complexity: low
---

# Task 06: The qa-report accept refusal test pins the new front-matter message

## Overview

Corrective Task from the first QA gate of 2026-09-28, whose repository Verification precondition failed: task_03 made an unparseable QA Report refuse with `QA Report front matter must open with a "---" first line`, but `TestRunQAReportAcceptCommandFailsClosed/unparseable_report` in `internal/cli/qa_report_test.go` still expects the substring `frontmatter`, so `internal/cli` fails.

## Requirements

1. MUST make the `unparseable report` case assert the exact refusal cause task_03 introduced (`front matter must open with a "---" first line`), not a looser word.
2. MUST search `internal/**` tests for other assertions of the old front-matter wording and update them the same way.
3. MUST NOT change the production message or any other case of the test.

## Subtasks

- [ ] Implement the requirements above.

## Acceptance Criteria

- [ ] `TestRunQAReportAcceptCommandFailsClosed` passes with every case, and the whole `internal/cli` and `internal/spec` packages pass.

## Context

- interface: `internal/cli/qa_report_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^TestRunQAReportAcceptCommandFailsClosed$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; printf "%s\\n" "$out" | grep -q -- "--- PASS: TestRunQAReportAcceptCommandFailsClosed" && grep -qF -- 'front matter must open with a' internal/cli/qa_report_test.go && go test -count=1 ./internal/cli ./internal/spec >/dev/null` — expected: exit 0; before this Task the unparseable case expects `frontmatter` and fails, so the command fails.

## References

- [_techspec.md](_techspec.md) — Build Order

## Result

- Updated the `unparseable report` CLI case to require `front matter must open
  with a "---" first line`, the refusal cause introduced by Task 03. Production
  wording and the other table cases are unchanged.
- Searched `internal/**` tests for stale front-matter assertions. The only other
  assertion of this QA Report failure was
  `TestQAVerdictReportsUnreadableReports/no_frontmatter`; it now requires the
  same cause. Its former `frontmatter` check passed accidentally because the
  temporary path contained the subtest name.
- Pre-change reproduction:
  `GOCACHE=/private/tmp/roundfix-task06-gocache go test -count=1 -run
  '^TestRunQAReportAcceptCommandFailsClosed$/^unparseable_report$'
  ./internal/cli` failed because stderr contained the new `front matter` cause
  while the test expected `frontmatter`.
- Focused evidence:
  `GOCACHE=/private/tmp/roundfix-task06-gocache go test -count=1 -run
  '^TestRunQAReportAcceptCommandFailsClosed$' ./internal/cli` passed all cases;
  `GOCACHE=/private/tmp/roundfix-task06-gocache go test -count=1 -run
  '^TestQAVerdictReportsUnreadableReports$/^no_frontmatter$' ./internal/spec`
  passed with the tightened assertion.
- Acceptance evidence: `GOCACHE=/private/tmp/roundfix-task06-gocache make
  verify-incremental` passed with process-table access, including the complete
  `internal/cli` and `internal/spec` package tests, `go vet`, skill checks, and
  the build. The sandboxed attempt reached the suite but its two force-stop
  integration tests were blocked by `operation not permitted`; rerunning with
  the required OS access resolved that environment-only block.
- The Task's declared `## Verification` command was not run; Daemon
  Verification remains the settlement authority.
