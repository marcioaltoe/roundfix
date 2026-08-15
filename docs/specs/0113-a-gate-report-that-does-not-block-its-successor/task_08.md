---
task: task_08
spec: 0113-a-gate-report-that-does-not-block-its-successor
status: pending # pending | in_progress | completed | failed — only implement-task changes this
type: backend
complexity: high
---

# Task 08: Wire the gate's new inputs into the request it actually receives

## Overview

The repairs of task_05 and task_06 exist and are tested, and nothing in
production reaches them. `GatePrecondition` has no production caller; the Daemon's
`qaMechanicalRequest` supplies neither `TaskRepairPaths` nor `AssignedRepairs`;
and `spec.Task` parses no repair declaration, so a Task file cannot name a repair
even if the gate would perform it. Unit callers pass while the real request
carries two empty slices.

## Requirements

1. MUST let a Task file declare an assigned repair and its bounded paths, parsed
   into the Task the Daemon already loads.
2. MUST populate `TaskRepairPaths` and `AssignedRepairs` in the request the
   Daemon builds, from that declaration.
3. MUST route the gate's static precondition through `GatePrecondition`, so a
   Spec-owned vocabulary finding is classified as a repair input rather than a
   refusal.
4. MUST prove the wiring through the request the Daemon builds, not through a
   hand-made one.
5. MUST leave a Task that declares no repair producing exactly the request it
   produces today.

## Subtasks

- [ ] Parse the repair declaration into the Task.
- [ ] Populate both request fields from it.
- [ ] Route the precondition through the classifier.
- [ ] Cover the declared, undeclared, and vocabulary-input cases.

## Acceptance Criteria

- [ ] A Task declaring a repair produces a request whose `AssignedRepairs` and
      `TaskRepairPaths` carry it.
- [ ] A Task declaring none produces the same request as before this Task.
- [ ] The gate's precondition runs through `GatePrecondition` in production, so a
      Spec's own undeclared-term finding becomes a repair input.
- [ ] The proof runs against the Daemon's own request builder.

## Verification

- `go test -count=1 ./internal/daemon -run 'TestQAMechanicalRequestCarriesAssignedRepairs' -v > /tmp/0113-t08.log 2>&1; s=$?; grep -q '^--- PASS: TestQAMechanicalRequestCarriesAssignedRepairs' /tmp/0113-t08.log || { cat /tmp/0113-t08.log; exit 1; }; exit $s` — expected: exits 0 and the log names the passing test; fails today, where no such test exists.
- `test -s /tmp/0113-t08.log || { echo 'the suite produced no output'; exit 1; }; grep -q '^--- PASS: TestQAMechanicalRequestCarriesAssignedRepairs' /tmp/0113-t08.log || { echo 'the wiring test did not run'; exit 1; }; grep -c '^    --- PASS\|^--- PASS' /tmp/0113-t08.log > /tmp/0113-t08-n.txt; test "$(cat /tmp/0113-t08-n.txt)" -ge 3 || { echo "expected the declared, undeclared, and vocabulary-input cases, got $(cat /tmp/0113-t08-n.txt)"; cat /tmp/0113-t08.log; exit 1; }` — expected: exits 0, refusing a vacuous run and proving each direction runs.
- `grep -rq 'GatePrecondition' internal/daemon --include='*.go' || { echo 'the gate precondition still has no production caller'; exit 1; }; grep -rq 'AssignedRepairs' internal/daemon --include='*.go' || { echo 'the Daemon request still does not carry assigned repairs'; exit 1; }; grep -rq 'AssignedRepairs\|repair' internal/spec/spec.go || { echo 'a Task file still cannot declare a repair'; exit 1; }` — expected: exits 0, proving all three layers are connected rather than only the checker. Fails today on every clause.
- `go test -count=1 ./internal/daemon ./internal/speccheck ./internal/spec > /tmp/0113-t08-regress.log 2>&1; s=$?; grep -q 'FAIL' /tmp/0113-t08-regress.log && { echo 'a layer regressed:'; grep -B 3 -A 8 'FAIL' /tmp/0113-t08-regress.log | head -30; exit 1; }; grep -rq 'TestQAMechanicalRequestCarriesAssignedRepairs' internal/daemon || { echo 'the three packages pass, but the wiring case does not exist'; exit 1; }; exit $s` — expected: exits 0, proving every layer the request crosses still passes, anchored to the case this Task adds.

## Context

- interface: `internal/daemon/task_engine.go`
- interface: `internal/spec/spec.go`
- interface: `internal/speccheck/mechanical.go`

## References

`_techspec.md` → System Architecture, the gate's contract execution.
`_prd.md` → Core Features 6 and 7; Goals 4 and 5. ADR-0134.
Evidence: this Spec's QA report `qa/qa-report-2026-08-15.md`, finding F-001.
