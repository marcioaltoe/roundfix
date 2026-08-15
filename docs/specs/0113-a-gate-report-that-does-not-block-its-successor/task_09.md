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

And routing the precondition through `speccheck.Check` in the Daemon's request
builder broke the public Implement journeys — because it contradicts a contract
this repository states by name. `TestRunImplementHasNoSpecCheckPrecondition`
builds a deliberately inconsistent Spec, proves the check would report errors,
and asserts that Implement runs anyway. The Implement command has no
Spec-consistency precondition, on purpose. The gate's classification belongs
where the gate already runs its strict check, not in the path that starts a Run.

## Requirements

1. MUST write the terminal refusal row only when the mechanical result actually
   blocks, never merely because it carries no rows.
2. MUST leave a non-blocking zero-row seed truthful: no refusal row, no `fail`
   verdict, and nothing claiming a cause that did not occur.
3. MUST NOT convert an optional detector skip into a refusal cause.
4. MUST remove the Spec-consistency check from the Implement path, leaving
   `TestRunImplementHasNoSpecCheckPrecondition` passing unchanged. The Implement
   command has no such precondition by design, and the CLI's deliberately
   inconsistent fixtures are that contract's evidence, not a defect to repair.
5. MUST apply the gate's precondition classification where the gate already runs
   its strict check, so a Spec's own declared term stops refusing it there.
6. MUST leave `make verify` green, including the Implement journeys that passed
   before this Spec.

## Subtasks

- [ ] Gate the refusal row on an actual blocking result.
- [ ] Keep a non-blocking zero-row seed truthful.
- [ ] Stop skips becoming refusal causes.
- [ ] Take the Spec check back out of the Implement path.
- [ ] Move the precondition classification to where the gate checks.

## Acceptance Criteria

- [ ] A blocking result with no rows still writes the terminal refusal row.
- [ ] A non-blocking result with no rows writes no refusal row and no `fail`
      verdict.
- [ ] A result whose only absences are optional detector skips is not reported as
      a refusal.
- [ ] The Implement journeys that passed at `66c60d4b` pass again.
- [ ] `TestRunImplementHasNoSpecCheckPrecondition` passes with its fixture
      unchanged, and the Implement path runs no Spec-consistency check.
- [ ] A Spec's own declared term still stops refusing its gate, proven where the
      gate runs its strict check.
- [ ] `make verify` exits 0.

## Verification

- `go test -count=1 ./internal/daemon -run 'TestWriteMechanicalQAReport' -v > /tmp/0113-t09.log 2>&1; s=$?; grep -q '^--- PASS: TestWriteMechanicalQAReport' /tmp/0113-t09.log || { cat /tmp/0113-t09.log; exit 1; }; grep -qE '^(    )?--- PASS: TestWriteMechanicalQAReport.*NonBlocking' /tmp/0113-t09.log || { echo 'the zero-row non-blocking case does not exist; the current tests cover only a carried row'; cat /tmp/0113-t09.log; exit 1; }; exit $s` — expected: exits 0, proving the writer covers the zero-row non-blocking case the current tests miss. The existing writer tests pass on an untouched tree, so the assertion is on the named new case rather than on a count they already satisfy.
- `go test -count=1 ./internal/cli -run '^TestRunImplement(QAPromptStatesSpecTargetBranchFromRunRecord|HasNoSpecCheckPrecondition)$' -v > /tmp/0113-t09-cli.log 2>&1; s=$?; for t in TestRunImplementQAPromptStatesSpecTargetBranchFromRunRecord TestRunImplementHasNoSpecCheckPrecondition; do grep -q "^--- PASS: $t" /tmp/0113-t09-cli.log || { echo "FAIL: $t"; cat /tmp/0113-t09-cli.log; exit 1; }; done; exit $s` — expected: exits 0, proving the public Implement journey works again *and* that the contract forbidding a Spec-consistency precondition there still holds. Both fail today.
- `grep -n 'len(result.Carried)+len(result.Blocked) == 0' internal/daemon/task_engine.go > /tmp/0113-t09-cond.txt; test -s /tmp/0113-t09-cond.txt || { echo 'the row-count condition vanished entirely; the refusal row must still be written for a blocking result'; exit 1; }; grep -qE 'result\.Blocking &&.*len\(result\.Carried\)' /tmp/0113-t09-cond.txt || { echo 'the refusal row is still written on row count alone:'; cat /tmp/0113-t09-cond.txt; exit 1; }` — expected: exits 0, proving the row count is still consulted and is now guarded by an actual blocking result. It prints the offending line on failure. Fails today, where the condition carries no blocking guard.
- `make verify > /tmp/0113-t09-verify.log 2>&1; s=$?; grep -q 'FAIL' /tmp/0113-t09-verify.log && { echo 'the authoritative gate is not green:'; grep -B 3 -A 8 'FAIL' /tmp/0113-t09-verify.log | head -40; exit 1; }; exit $s` — expected: exits 0, proving the repository gate recovers. Fails today, where the CLI journeys are red.

## Context

- interface: `internal/daemon/task_engine.go`
- interface: `internal/cli/implement_test.go`

## References

`_techspec.md` → System Architecture, the report writer; Risks & Considerations,
a refused report looks structurally like a passing one. `_prd.md` → Core
Feature 1; Goal 1; Non-Goals, removing or weakening the mechanical stage.
ADR-0132, ADR-0134.
Evidence: this Spec's QA report `qa/qa-report-2026-08-15-01.md`, findings F-001
and F-002.
