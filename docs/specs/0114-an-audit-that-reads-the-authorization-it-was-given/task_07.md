---
task: task_07
spec: 0114-an-audit-that-reads-the-authorization-it-was-given
status: pending # pending | in_progress | completed | failed — only implement-task changes this
type: test
complexity: low
---

# Task 07: Teach the Daemon's fixture what governed now means

## Overview

The Daemon's integration fixture folds `outside.txt` into a Task commit and
expects one `QA-AUTH-PATHS` finding. Under ADR-0130 that path is ordinary, so the
corrected audit returns nothing and the assertion fails — the last consumer still
asserting the behaviour this Spec deliberately changed. The fixture is updated to
assert the new contract in both directions rather than the old one in one.

## Requirements

1. MUST make the fixture's refusal case fold a governed path, so it still proves
   that an ungranted governed path blocks.
2. MUST add the case this Spec exists to allow: an ordinary path folded into the
   same commit produces no finding.
3. MUST leave `internal/daemon` and the repository Verification green.
4. MUST NOT weaken the assertion to make the failure disappear; the refusal case
   keeps asserting a blocking finding that names its path.

## Subtasks

- [ ] Fold a governed path in the refusal case.
- [ ] Add the ordinary-path case that must not block.
- [ ] Prove both through the Daemon's own mechanical request.

## Acceptance Criteria

- [ ] A Task commit folding an ungranted governed path still blocks, and the
      finding names that path.
- [ ] A Task commit folding an ordinary path produces no finding.
- [ ] `internal/daemon` passes.
- [ ] No assertion was removed or loosened; the refusal case still requires a
      blocking result with a named path.

## Verification

- `go test -count=1 ./internal/daemon -run 'TestQAMechanicalRequest' -v > /tmp/0114-t07.log 2>&1; s=$?; grep -q '^--- PASS: TestQAMechanicalRequest' /tmp/0114-t07.log || { cat /tmp/0114-t07.log; exit 1; }; exit $s` — expected: exits 0 and the log names the passing fixture; fails today, where the assertion expects a finding the corrected audit no longer produces.
- `test -s /tmp/0114-t07.log || { echo 'the suite produced no output'; exit 1; }; grep -q '^--- PASS: TestQAMechanicalRequest' /tmp/0114-t07.log || { echo 'the fixture did not run'; exit 1; }; grep -c '^    --- PASS\|^--- PASS' /tmp/0114-t07.log > /tmp/0114-t07-n.txt; test "$(cat /tmp/0114-t07-n.txt)" -ge 2 || { echo "expected the governed refusal and the ordinary-path pass as their own cases, got $(cat /tmp/0114-t07-n.txt)"; cat /tmp/0114-t07.log; exit 1; }` — expected: exits 0, proving both directions run rather than one assertion being flipped.
- `grep -n 'result.Blocking' internal/daemon/task_engine_test.go | head -1 > /tmp/0114-t07-assert.txt; test -s /tmp/0114-t07-assert.txt || { echo 'the blocking assertion was removed rather than retargeted'; exit 1; }; grep -q 'outside.txt' internal/daemon/task_engine_test.go && { echo 'the ordinary path is still used as the refusal case'; exit 1; }; exit 0` — expected: exits 0, proving the refusal assertion survives and no longer rests on an ungoverned path. Fails today, where `outside.txt` is still the refusal case.
- `go test -count=1 ./internal/daemon ./internal/speccheck > /tmp/0114-t07-regress.log 2>&1; s=$?; grep -q 'FAIL' /tmp/0114-t07-regress.log && { echo 'a consumer of the audit is still red:'; grep -B 3 -A 8 'FAIL' /tmp/0114-t07-regress.log | head -30; exit 1; }; exit $s` — expected: exits 0, proving the audit and the Daemon that calls it agree. Fails today.

## Context

- interface: `internal/daemon/task_engine_test.go`

## References

`_techspec.md` → Risks & Considerations, the direction that loses safety.
`_prd.md` → Core Feature 2; Goal 2; User Story 2. ADR-0130.
Evidence: this Spec's QA report `qa/qa-report-2026-08-15-01.md`, finding F-001.
