---
task: task_07
spec: 0089-an-effort-the-runtime-actually-receives
status: pending
type: docs
complexity: low
---

# Task 07: Record that the runtime hands back the floor

## Overview

The durable half of the adopted measurement — that three of four candidate
models open at the bottom of their advertised range, and what Roundfix now does
about it — belongs in the repository's model selection reference, where the
configuration already points. An archived Spec may be deleted; this reference
may not.

## Requirements

1. MUST record the measured default effort each candidate model opens at, and
   that three of four open at the floor of their own range.
2. MUST record that Roundfix now applies a requested effort after a session
   warm-up rather than inheriting the default, citing the decision by its ADR
   identifier rather than restating its reasoning.
3. MUST correct the section that currently states OpenCode reasoning is
   model-managed and refused, which this Spec makes false.
4. MUST date the snapshot and name the runtime and adapter versions.
5. MUST NOT reference any path under the Spec directory.
6. MUST NOT change `CONTEXT.md`; a glossary gap is raised by the closing gate.

## Subtasks

- [ ] Record the per-model defaults and the floor observation.
- [ ] Replace the model-managed paragraph with the warm-up contract.
- [ ] Date the snapshot with runtime and adapter versions.
- [ ] Confirm no Spec path is referenced.

## Acceptance Criteria

- [ ] The reference records each candidate model's default effort.
- [ ] It states that three of four open at the floor of their range.
- [ ] It cites ADR-0108 for what Roundfix does about it.
- [ ] It no longer claims OpenCode reasoning effort is refused.
- [ ] It contains no path under the Spec directory.

## Bounded scope

This Task may create or modify only:

- `docs/references/model-selection.md`
- `docs/specs/0089-an-effort-the-runtime-actually-receives/task_07.md`

## Verification

- `grep -q 'ADR-0108' docs/references/model-selection.md` — expected: exits 0.
- `grep -q '2026-08-09' docs/references/model-selection.md` — expected: exits 0.
- `grep -i 'model-managed reasoning runtime' docs/references/model-selection.md` — expected: exits non-zero with no output, proving the superseded claim is gone.
- `grep 'docs/specs/' docs/references/model-selection.md` — expected: exits non-zero with no output.
- `go run -buildvcs=false ./cmd/roundfix spec check 0089-an-effort-the-runtime-actually-receives` — expected: exits 0.

## References

- `_prd.md` → Core Features.
- `_techspec.md` → Build Order 7.
- `docs/agents/specific-repository.md` → the durable-knowledge-flows-upstream HARD RULE.
- ADR-0108.
