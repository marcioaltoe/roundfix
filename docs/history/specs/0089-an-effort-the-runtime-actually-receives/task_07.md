---
task: task_07
spec: 0089-an-effort-the-runtime-actually-receives
status: completed
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
- `! grep -iq 'model-managed reasoning runtime' docs/references/model-selection.md` — expected: exits 0, proving the superseded claim is gone.
- `! grep -q 'docs/specs/' docs/references/model-selection.md` — expected: exits 0, proving durable knowledge does not depend on a Spec.
- `go run -buildvcs=false ./cmd/roundfix spec check 0089-an-effort-the-runtime-actually-receives` — expected: exits 0.

## References

- `_prd.md` → Core Features.
- `_techspec.md` → Build Order 7.
- `docs/agents/specific-repository.md` → the durable-knowledge-flows-upstream HARD RULE.
- ADR-0108.

## Result

The model selection reference now records the four OpenRouter candidates'
advertised effort ranges and measured opening values. It identifies the three
models that open at the floor, dates the snapshot, names the measured runtime
and adapter versions, and replaces the superseded refusal with Roundfix's
session warm-up contract and its ADR-0108 citation.

Pre-change signal:

- `rtk rg -n 'OpenCode reasoning effort|Roundfix therefore treats|ADR-0108|2026-08-09' docs/references/model-selection.md` — returned only the superseded heading and the statement that Roundfix treated OpenCode as model-managed; the durable reference contained neither the 2026-08-09 snapshot nor an ADR-0108 citation.

Focused checks after the documentation edit:

- `rtk rg -n '2026-08-09|OpenCode 1\.18\.15|acpx 0\.13\.0|openrouter/x-ai/grok-4\.5|openrouter/moonshotai/kimi-k3|openrouter/deepseek/deepseek-v4-flash-0731|openrouter/deepseek/deepseek-v4-pro|Three of four candidates|ADR-0108' docs/references/model-selection.md` — exited 0 and returned the dated/versioned snapshot, all four candidate rows, the three-of-four observation, and the ADR citation.
- `rtk rg -n 'model-managed reasoning runtime|refuses any non-empty|docs/specs/' docs/references/model-selection.md` — exited 1 with no matches, confirming that neither the superseded refusal nor a Spec-directory path remains.
- `rtk git -c core.fsmonitor=false diff --check` — exited 0.

Acceptance evidence:

- Criterion 1: the focused positive search returned all four candidate rows; their recorded defaults are `low` for `grok-4.5`, `kimi-k3`, and `deepseek-v4-flash-0731`, and `high` for `deepseek-v4-pro`.
- Criterion 2: the same search returned the statement that three of four candidates open at the floor of their advertised range.
- Criterion 3: the same search returned the session warm-up contract's ADR-0108 citation.
- Criterion 4: the focused absence search found neither the model-managed claim nor the non-empty-effort refusal.
- Criterion 5: the focused absence search found no `docs/specs/` path in the durable reference.

No follow-up work was found inside this Task's slice. The commands authored
under `## Verification` were not run; Daemon Verification remains the
settlement boundary.
