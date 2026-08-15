---
task: task_09
spec: 0113-a-gate-report-that-does-not-block-its-successor
status: pending # pending | in_progress | completed | failed — only implement-task changes this
type: backend
complexity: high
---

# Task 09: A refusal row records a refusal, and the precondition breaks no journey

## Overview

Two regressions this Spec introduced, both found by its own gate.

The refusal row is written whenever the result carries no rows, without requiring
that the result blocks. A non-blocking stage that starts the Agent therefore
seeds a report claiming `QA-PRECONDITION | fail | mechanical refusal` for a
refusal that did not happen — the writer's own words contradicting its verdict.
ADR-0132 says a refused gate records its refusal; absence of rows is not refusal.

And routing the precondition through `speccheck.Check` in production made the
public Implement journeys red: the CLI's Spec fixtures were never
consistency-clean, and now the gate they never had to satisfy refuses them.

## Requirements

1. MUST write the terminal refusal row only when the mechanical result actually
   blocks, never merely because it carries no rows.
2. MUST leave a non-blocking zero-row seed truthful: no refusal row, no `fail`
   verdict, and nothing claiming a cause that did not occur.
3. MUST NOT convert an optional detector skip into a refusal cause.
4. MUST make the CLI's on-disk Spec fixtures satisfy the precondition the public
   journey now enforces, rather than removing the precondition.
5. MUST leave `make verify` green, including the Implement journeys that passed
   before this Spec.

## Subtasks

- [ ] Gate the refusal row on an actual blocking result.
- [ ] Keep a non-blocking zero-row seed truthful.
- [ ] Stop skips becoming refusal causes.
- [ ] Repair the CLI Spec fixtures.

## Acceptance Criteria

- [ ] A blocking result with no rows still writes the terminal refusal row.
- [ ] A non-blocking result with no rows writes no refusal row and no `fail`
      verdict.
- [ ] A result whose only absences are optional detector skips is not reported as
      a refusal.
- [ ] The Implement journeys that passed at `66c60d4b` pass again.
- [ ] `make verify` exits 0.

## Verification

- `go test -count=1 ./internal/daemon -run 'TestWriteMechanicalQAReport' -v > /tmp/0113-t09.log 2>&1; s=$?; grep -q '^--- PASS: TestWriteMechanicalQAReport' /tmp/0113-t09.log || { cat /tmp/0113-t09.log; exit 1; }; grep -qE '^(    )?--- PASS: TestWriteMechanicalQAReport.*NonBlocking' /tmp/0113-t09.log || { echo 'the zero-row non-blocking case does not exist; the current tests cover only a carried row'; cat /tmp/0113-t09.log; exit 1; }; exit $s` — expected: exits 0, proving the writer covers the zero-row non-blocking case the current tests miss. The existing writer tests pass on an untouched tree, so the assertion is on the named new case rather than on a count they already satisfy.
- `go test -count=1 ./internal/cli -run '^TestRunImplementQAPromptStatesSpecTargetBranchFromRunRecord$' -v > /tmp/0113-t09-cli.log 2>&1; s=$?; grep -q '^--- PASS: TestRunImplementQAPromptStatesSpecTargetBranchFromRunRecord' /tmp/0113-t09-cli.log || { cat /tmp/0113-t09-cli.log; exit 1; }; exit $s` — expected: exits 0, proving the public Implement journey works again. Fails today, where the precondition refuses the CLI fixture.
- `grep -n 'len(result.Carried)+len(result.Blocked) == 0' internal/daemon/task_engine.go && { echo 'the refusal row is still written on row count alone'; exit 1; }; grep -q 'Blocking' internal/daemon/task_engine.go || { echo 'the writer no longer consults the blocking result'; exit 1; }; exit 0` — expected: exits 0, proving the refusal row is gated on refusal rather than on emptiness. It prints the offending line on failure. Fails today.
- `make verify > /tmp/0113-t09-verify.log 2>&1; s=$?; grep -q 'FAIL' /tmp/0113-t09-verify.log && { echo 'the authoritative gate is not green:'; grep -B 3 -A 8 'FAIL' /tmp/0113-t09-verify.log | head -40; exit 1; }; exit $s` — expected: exits 0, proving the repository gate recovers. Fails today, where the CLI journeys are red.

## Context

- interface: `internal/daemon/task_engine.go`
- interface: `internal/cli/implement_test.go`

## References

`_techspec.md` → System Architecture, the report writer; Risks & Considerations,
a refused report looks structurally like a passing one. `_prd.md` → Core
Feature 1; Goal 1. ADR-0132.
Evidence: this Spec's QA report `qa/qa-report-2026-08-15-01.md`, findings F-001
and F-002.
