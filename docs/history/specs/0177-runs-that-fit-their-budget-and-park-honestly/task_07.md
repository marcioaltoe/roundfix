---
task: task_07
spec: 0177-runs-that-fit-their-budget-and-park-honestly
status: completed
type: backend
complexity: low
---

# Task 07: The QA settlement renews the budget before its report commit, and wrap-fix remediations are shell-safe

## Overview

Corrective Task from the pre-PR review of 2026-09-28. In `internal/daemon/task_engine.go` the QA Task's budget renewal happens only after `runQAGate` returns, but the QA Task is settled before the QA Report commit; if that post-settlement commit crosses the previous deadline, the watchdog cancels the gate and the timely settlement never renews the allowance. And the `SC-VERIFY-WRAP-FRAGILE` remediation built in `internal/speccheck/verification.go` quotes the phrase with `strconv.Quote` (Go syntax, not shell quoting) and interpolates the file operand raw, so a phrase or path containing `$()`, backticks, `;` or spaces yields an unsafe or malformed command.

## Requirements

1. MUST renew the Implement Run budget at the QA Task's settlement, before the QA Report commit, exactly as every other Task settlement renews it, so a report commit that starts before the renewed deadline is never cancelled by the old one.
2. MUST quote both the phrase and the file operand of the suggested remediation with POSIX single-quote shell quoting (a `'` inside becomes `'\''`), so the printed command is safe to paste for any phrase or path.
3. MUST keep every existing test and message of this Spec unchanged, and put the new tests in files of their own.

## Subtasks

- [x] Implement the requirements above.
- [x] Add a named test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [x] A QA settlement near the old deadline renews the budget, and the report commit completes without cancellation.
- [x] A phrase containing `$(`, a backtick, `;`, a space and a `'`, and a path with a space, produce a remediation that a shell runs as a literal `grep -qF` of exactly that phrase against exactly that file.

## Context

- interface: `internal/daemon/task_engine.go`
- interface: `internal/speccheck/verification.go`
- creates: `internal/daemon/qa_budget_renewal_test.go`
- creates: `internal/speccheck/wrap_fix_quoting_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestQASettlementRenewsTheBudgetBeforeTheReportCommit|TestWrapFragileRemediationQuotesAHostilePhrase|TestWrapFragileRemediationQuotesAPathWithASpace)$" ./internal/daemon ./internal/speccheck 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestQASettlementRenewsTheBudgetBeforeTheReportCommit TestWrapFragileRemediationQuotesAHostilePhrase TestWrapFragileRemediationQuotesAPathWithASpace; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task the three new tests do not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — Build Order

## Result

The QA path now renews the shared Implement Run Budget immediately after the Daemon settles the QA Task and before it prepares or creates the QA Report commit. Wrap-fragile remediations preserve the existing rendering for the established shell-safe phrase and path while POSIX single-quoting hostile phrases and paths that require quoting; embedded single quotes use the `'''` shell sequence without evaluating command substitutions, backticks, or separators.

Acceptance evidence:

- `TestQASettlementRenewsTheBudgetBeforeTheReportCommit` failed before the implementation with `context deadline exceeded`, then passed after the renewal moved to the settlement boundary. Its fake clock advances past the prior deadline during the QA Report commit and observes the QA settlement's renewed deadline afterwards.
- `TestWrapFragileRemediationQuotesAHostilePhrase` failed before the implementation after executing the embedded `$()` and backtick substitutions, then passed after the fix. It executes the suggested remediation through `sh`, proves neither marker command ran, and proves a near-match does not satisfy `grep -qF`.
- `TestWrapFragileRemediationQuotesAPathWithASpace` failed before the implementation because the path was raw, then passed after the fix. It executes the remediation through `sh` against the exact spaced Markdown path.

Focused checks:

- `GOCACHE=/tmp/roundfix-task07-gocache go test -count=1 ./internal/daemon ./internal/speccheck` — passed after the final production changes; existing affected-package tests and messages remain unchanged.
- `make verify-incremental` — the sandboxed attempt reached the full suite but could not read the process table in two force-stop integration tests; the rerun with process-table access passed `go vet`, the full Go suite, skill checks, and the build.
- Task 07's authored `## Verification` command was not run; the Daemon owns that check and terminal settlement.
