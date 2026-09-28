---
task: task_01
spec: 0177-runs-that-fit-their-budget-and-park-honestly
status: pending
type: backend
complexity: low
---

# Task 01: Instruction paths never make two Tasks collide

## Overview

`declaredTaskTouches` in `internal/spec/collision.go` adds every existing
`## Context` path of a Task to the Task's touch set, whatever its kind. Every
Task that names the shared `instruction:` file
`.agents/skills/implement-task/SKILL.md` therefore collides with every other
such Task. `roundfix spec check` then reports `SC-WAVE-COLLISION`, the Daemon
refuses the plan at Task Capacity above `1`, and the author has to chain Tasks
that share nothing they edit. The touch set is read from the Task files by
`spec.Collisions`, and two readers act on it: `detectWaveCollisions` in
`internal/speccheck/coherence.go` at authoring, and
`refuseTaskPlanWaveCollisions` in `internal/daemon/task_engine.go` before
dispatch. An `instruction:` path is read-only by the write-tasks contract, so it
is not evidence of a touch. A path that really is edited must keep colliding,
or two Tasks editing one file would run in one Wave from a stale base.

## Requirements

1. MUST make `declaredTaskTouches` skip Context references of kind
   `spec.ContextKindInstruction`, and keep adding `interface:` and `creates:`
   references, Verification operands from `TaskVerificationFiles` and prior-Run
   settlement paths from `priorRunTaskTouches` exactly as today.
2. MUST keep reporting a collision when a path one Task names as an
   `instruction:` is also a Verification operand or an `interface:` path of
   both Tasks; only the instruction source stops counting.
3. MUST NOT change `internal/speccheck/coherence.go` or
   `internal/daemon/task_engine.go`: both readers call `spec.Collisions`, and a
   test through `speccheck.Check` proves the authoring reader follows.
4. MUST state in the declared-path rules of
   `.agents/skills/write-tasks/SKILL.md` that an `instruction:` path is
   read-only and never makes two Tasks collide, using the phrase
   `never makes two Tasks collide`, and regenerate `skills/write-tasks/SKILL.md`
   with `make skills-sync`.
5. MUST put the new `internal/spec` tests in
   `internal/spec/collision_instruction_test.go` and the new `speccheck` test in
   `internal/speccheck/wave_instruction_test.go`, each negative case a test of
   its own, and MUST NOT rename or remove an existing top-level test.

## Subtasks

- [ ] Skip instruction references in the touch set.
- [ ] Add the positive and negative collision tests and the `speccheck` test.
- [ ] State the rule in the write-tasks skill and run `make skills-sync`.

## Acceptance Criteria

- [ ] Two unordered Tasks that share only an `instruction:` path produce no
      collision from `spec.Collisions` and no `SC-WAVE-COLLISION` from
      `speccheck.Check`.
- [ ] Two unordered Tasks that share an `interface:` path still collide.
- [ ] A path named as an instruction and also read by both Tasks'
      Verifications still collides.

## Context

- interface: `internal/spec/collision.go`
- creates: `internal/spec/collision_instruction_test.go`
- creates: `internal/speccheck/wave_instruction_test.go`
- interface: `.agents/skills/write-tasks/SKILL.md`
- interface: `skills/write-tasks/SKILL.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestCollisionsIgnoresASharedInstructionPath|TestCollisionsStillReportsASharedInterfacePath|TestCollisionsReportsAnInstructionPathReadByBothVerifications|TestCollisionsLearnsPathFromDeclaredContext|TestWaveCollisionCheckIgnoresASharedInstructionPath|TestWaveCollisionCheckStillReportsASharedInterfacePath)$" ./internal/spec ./internal/speccheck 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestCollisionsIgnoresASharedInstructionPath TestCollisionsStillReportsASharedInterfacePath TestCollisionsReportsAnInstructionPathReadByBothVerifications TestCollisionsLearnsPathFromDeclaredContext TestWaveCollisionCheckIgnoresASharedInstructionPath TestWaveCollisionCheckStillReportsASharedInterfacePath; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done && tr -s '[:space:]' ' ' < .agents/skills/write-tasks/SKILL.md | grep -qF -- "never makes two Tasks collide" && diff -r .agents/skills/write-tasks skills/write-tasks >/dev/null` — expected: exit 0; before this Task none of the new named tests exists and the write-tasks skill never says an instruction path never makes two Tasks collide, so the command fails.

## References

- `_prd.md` → Goal 1; Core Feature 1; Success Metric 1.
- `_techspec.md` → Instruction paths and Wave collisions; API Contract 1;
  Testing Approach 1; ADR-0025; ADR-0056; ADR-0093; ADR-0117.
