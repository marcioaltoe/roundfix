---
status: completed
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

## Result

Implemented the Task Graph-stage `SC-WAVE-COLLISION` finding by calling the
shared `spec.Collisions` rule from the Spec Consistency Check. Each finding is
an error, names both Tasks, renders every shared repository path with its touch
source, and tells the author to give one Task a `needs` edge on the other.

Acceptance evidence:

- `TestWaveCollisionIsReportedAtAuthoring/colliding_graph` exercises two
  independent Tasks that name two shared files. It asserts one error finding,
  both Task IDs, both paths, the `verification command` source, and the
  `needs`-edge remedy.
- `TestWaveCollisionIsReportedAtAuthoring/serialized_graph` gives `task_02` a
  `needs: [task_01]` edge and asserts that no collision finding is emitted.
- `TestWaveCollisionIsReportedAtAuthoring/Spec_with_no_graph_yet` checks a Spec
  without `_tasks.md`, asserts no collision finding, and requires the detector
  to record the missing Task Graph as a skip.
- The stage-scope tests require `SC-WAVE-COLLISION` to be skipped at PRD and
  TechSpec authoring and available at the Task Graph stage.

Focused checks:

- Red: `rtk go test -run '^TestWaveCollisionIsReportedAtAuthoring$' ./internal/speccheck`
  failed to compile because `speccheck.CodeWaveCollision` did not exist.
- Green: `rtk go test -count=1 -run '^(TestWaveCollisionIsReportedAtAuthoring|TestStageScope)' ./internal/speccheck`
  passed 10 tests.
- `rtk git diff --check` exited 0 with no diagnostics.

The Daemon-owned Verification command was not run.
