---
task: task_06
spec: 0174-operator-surfaces-that-tell-the-truth
status: pending
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
