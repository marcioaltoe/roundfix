---
task: task_06
spec: 0029-launch-and-recovery-fixes
status: pending
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

- [ ] Detached Run, Doctor, Settle, and Batch-failure sections updated in the canonical skill source
- [ ] Mirror regenerated via the repository's skills-sync target
- [ ] Skill-check anchors updated with the text they anchor
- [ ] Cross-read the PRD's Core Features and confirm each shipped behavior is described accurately

## Acceptance Criteria

- [ ] The skills check passes with the updated anchors
- [ ] The skill documents the detach failure diagnostics and no longer implies a silent failure mode
- [ ] The skill documents the doctor `model:` line and the settle surface reporting
- [ ] The full test suite passes

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
