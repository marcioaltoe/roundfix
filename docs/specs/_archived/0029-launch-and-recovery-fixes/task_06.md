---
task: task_06
spec: 0029-launch-and-recovery-fixes
status: completed
type: docs
complexity: medium
---

# Task 06: Sync the Roundfix Skill and docs to the shipped behavior

## Overview

Skill-sync hard rule: the Roundfix Skill's Detached Run, Doctor, and Settle sections are updated to describe the shipped two-phase handshake semantics and failure diagnostics, the `model:` doctor line, the `Settle surface:` reporting and failed-surface preference, and the actionable Batch-time model-rejection reason — in the canonical skill source with the embedded mirror regenerated and the skill-check anchors kept green.

## Requirements

1. MUST update the canonical Roundfix Skill's Detached Run section: startup failures now print explicit diagnostics (phase, exit code or signal, output presence), and a slow Preflight Validation no longer fails a detach start.
2. MUST update the Doctor section for the `model:` check line and what its failure content means (advertised models, next action).
3. MUST update the Settle section for the `Settle surface:` line, the failed-surface preference, and the per-candidate refusal.
4. MUST mention, wherever the skill explains Batch failures, that an Agent Model rejection reports the model and the advertised list rather than a generic protocol error.
5. MUST regenerate the embedded skill mirror and update skill-check anchors together with any anchored text changed, keeping the check green.
6. MUST NOT alter upstream-managed skills (skill-governance ownership split).

## Subtasks

- [x] Detached Run, Doctor, Settle, and Batch-failure sections updated in the canonical skill source
- [x] Mirror regenerated via the repository's skills-sync target
- [x] Skill-check anchors updated with the text they anchor
- [x] Cross-read the PRD's Core Features and confirm each shipped behavior is described accurately

## Acceptance Criteria

- [x] The skills check passes with the updated anchors
- [x] The skill documents the detach failure diagnostics and no longer implies a silent failure mode
- [x] The skill documents the doctor `model:` line and the settle surface reporting
- [x] The full test suite passes

## Context

- instruction: `docs/agents/skill-governance.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/skills.go`

## Verification

- `go run -buildvcs=false ./cmd/roundfix skills check` — expected: skill check passes
- `grep -q "Settle surface:" skills/roundfix/SKILL.md` — expected: exit 0 (settle contract updated)
- `grep -qi "liveness" skills/roundfix/SKILL.md` — expected: exit 0 (detach semantics documented)
- `go test ./...` — expected: all tests pass
- `go build -buildvcs=false ./...` — expected: clean build

## References

`_prd.md` → Core Feature 5; `_techspec.md` → Build Order 6, System Architecture (skills/roundfix); `docs/agents/skill-governance.md`.

## Result

- Updated the canonical Roundfix Skill to document the two-phase Detached Run handshake, liveness marker, separate Run-creation ceiling, explicit timeout and child-exit diagnostics, and the fact that slow live Preflight Validation no longer fails detach startup only because it exceeds 10 seconds.
- Updated the Doctor section to include the `model:` line, the passing `model: ok (<model>)` shape, and rejected-model failure content with advertised models plus `next:` recovery guidance.
- Updated the Settle section to describe failed-first surface selection, the `Settle surface: <path>` stderr line, and per-candidate refusal messages when no surface has the Task `failed`.
- Updated Batch-failure documentation to state that Agent Model not-advertised failures report the rejected model, runtime, and advertised list instead of a generic `agent/protocol error`.
- Regenerated the embedded `skills/roundfix` mirror with `rtk make skills-sync` and added skill-check anchors for `model: ok`, liveness diagnostics, `Settle surface: <path>`, and not-advertised Batch reasons.
- Evidence: `rtk go run -buildvcs=false ./cmd/roundfix skills check` passed; `rtk proxy grep -q "Settle surface:" skills/roundfix/SKILL.md` passed; `rtk proxy grep -qi "liveness" skills/roundfix/SKILL.md` passed; `rtk go test ./...` passed with 1249 tests across 19 packages; `rtk go build -buildvcs=false ./...` passed.
- Verification Feedback attempt 1: inspected the diagnostic artifact without copying its body; the failure was an intermittent `internal/agent` cancellation-session test outside the task_06 docs/skill-sync slice. No non-task_06 files were changed; a fresh `rtk go test ./...` rerun passed with 1249 tests across 19 packages.
