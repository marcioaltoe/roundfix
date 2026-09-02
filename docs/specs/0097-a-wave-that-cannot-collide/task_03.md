---
status: pending
type: backend
---

# Task: The Run refuses before it dispatches

The checker catches a collision while the graph is being written. This is the
backstop for a graph authored before the checker carried the rule, or edited
after it passed — and it is where the Success Metric is measured: the collision
is reported before any Agent Session opens.

## Work

- Before dispatching, the Run asks the same rule the checker asks. Same
  function, not a second implementation: ADR-0148 exists because a checker that
  approves what the Run later refuses is one defect wearing two faces.
- A collision refuses. No Agent Session opens, no Task Worktree is created, and
  the refusal names both Tasks, each shared path, its source, and the `needs`
  edge that settles it.
- Never reorder the plan. Silently serializing would replace a visible failure
  with an invisible one, and the Supervisor authored that plan.
- Cover a refusal asserted against the Run Database rather than stdout alone —
  no Agent Session, no worktree — and a serialized graph that dispatches
  normally.

## References

- `_prd.md` → Goal 1, User Story 1; Success Metrics, the collision reported
  before any Agent Session opens; Open Questions, settled as refuse
- `_techspec.md` → Build Order 3
- ADR-0148 is why this calls the rule rather than restating it

## Verification
- `grep -q "Collisions" internal/daemon/task_engine.go || exit 1; grep -q "TestTaskCycleRefusesAWaveCollision" internal/daemon/task_engine_test.go || exit 1; go test -count=1 ./internal/daemon`
