---
task: task_07
spec: 0088-a-third-runtime-that-can-run
status: pending
type: docs
complexity: low
---

# Task 07: Move the measured OpenCode facts upstream

## Overview

The adopted measurement lives inside this Spec, and an archived Spec may be
deleted at any time. The durable half of what it recorded — which models the
subscription grants, how the effort vocabulary varies by model, and why Roundfix
treats OpenCode reasoning as model-managed — belongs to the repository's model
selection reference, which `.roundfixrc.yml` already cites.

## Requirements

1. MUST record, in the repository's model selection reference, the OpenCode model
   identifiers Roundfix can select, distinguishing the subscribed `opencode-go`
   tier from the pay-per-use and aggregator tiers.
2. MUST record that the reasoning-effort vocabulary is advertised per model and
   varies, with the measured examples, and that Roundfix therefore treats
   OpenCode reasoning as model-managed.
3. MUST date the snapshot and name the runtime and adapter versions it was taken
   against, matching the existing snapshot convention in that document.
4. MUST reference the decision rather than restate its reasoning, so the two
   documents cannot drift into disagreement.
5. MUST NOT reference any path under the Spec directory, because durable
   knowledge never depends on a Spec that may be deleted.
6. MUST NOT change `CONTEXT.md`; a glossary gap, if any, is raised by the closing
   gate.

## Subtasks

- [ ] Add the OpenCode section to the model selection reference.
- [ ] Record the subscribed tier and the measured per-model effort variation.
- [ ] Date the snapshot with its runtime and adapter versions.
- [ ] Confirm no Spec path is referenced.

## Acceptance Criteria

- [ ] The model selection reference names the `opencode-go` subscribed models.
- [ ] It records that the effort vocabulary is per model, with the measured
      examples.
- [ ] It carries a dated snapshot line naming the runtime and adapter versions.
- [ ] It cites the model-managed reasoning decision by its ADR identifier.
- [ ] It contains no path under the Spec directory.

## Bounded scope

This Task may create or modify only:

- `docs/references/model-selection.md`
- `docs/specs/0088-a-third-runtime-that-can-run/task_07.md`

## Verification

- `grep -q 'opencode-go' docs/references/model-selection.md` — expected: exits 0.
- `grep -q 'ADR-0106' docs/references/model-selection.md` — expected: exits 0, proving the decision is cited rather than restated.
- `grep -q '2026-08-08' docs/references/model-selection.md` — expected: exits 0, proving the snapshot is dated.
- `grep 'docs/specs/' docs/references/model-selection.md` — expected: exits non-zero with no output, proving durable knowledge does not depend on a Spec.
- `go run -buildvcs=false ./cmd/roundfix spec check 0088-a-third-runtime-that-can-run` — expected: exits 0.

## References

- `_prd.md` → Core Features.
- `_techspec.md` → Build Order 7.
- `docs/agents/specific-repository.md` → the durable-knowledge-flows-upstream
  HARD RULE.
- ADR-0106.
