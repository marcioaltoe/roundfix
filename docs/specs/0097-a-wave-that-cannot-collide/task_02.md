---
status: pending
type: backend
---

# Task: The checker reports a collision at authoring

A collision is produced when the graph is authored, and ADR-0117 places a check
with the stage that can produce its defect. Reporting it at authoring costs one
edit; reporting it at integration costs a Run.

## Work

- The Spec Consistency Check reports a graph whose same-wave Tasks share a path,
  naming both Tasks, each shared path, and the source that produced it.
- It is a finding rather than a gap: both sides are written down in the Spec's
  own artifacts, so the checker is not guessing.
- The message names the fix — give one Task a `needs` edge on the other —
  because a reader who learns only that two Tasks collide still has to work out
  what to do.
- Register it against the stage where a Task Graph exists, so a Spec checked
  before its graph is written is not reported for a graph it does not have.
- Cover a colliding graph, a serialized one, and a Spec with no graph yet.

## References

- `_prd.md` → Goal 1, Core Feature 1
- `_techspec.md` → Build Order 2; API Contracts
- ADR-0117 places a check with the stage that produces its defect; ADR-0148 is
  the one-rule-two-callers shape this Task is the first half of

## Verification
- `grep -q "CodeWaveCollision" internal/speccheck/coherence.go || exit 1; grep -q "TestWaveCollisionIsReportedAtAuthoring" internal/speccheck/coherence_test.go || exit 1; go test -count=1 ./internal/speccheck ./internal/spec`
