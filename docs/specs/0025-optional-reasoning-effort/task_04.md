---
task: task_04
spec: 0025-optional-reasoning-effort
status: pending
type: docs
complexity: low
---

# Task 04: Ship gpt-5.6-sol config and optional-effort guidance

## Overview

Move this repository's Project Config to the gpt-5.6-sol Agent Model with
model-managed reasoning, and align every guidance surface with the optional
Default Reasoning Effort semantics: the glossary, the README, and both
Roundfix Skill copies. The slice is verifiable through the Skill sync and
Skill validation gates plus the updated config content.

## Requirements

1. MUST set the repository Project Config codex runtime to the `gpt-5.6-sol`
   Agent Model with an explicitly empty reasoning effort, with a comment
   explaining that the empty value means the model manages reasoning and
   overrides any User Config value.
2. MUST update the CONTEXT.md Default Reasoning Effort entry so an empty
   value means the Agent Model manages reasoning, keeping the
   never-inherits-runtime-configuration stance for non-empty values.
3. MUST document the optional reasoning semantics — empty config value,
   explicit empty flag, model-managed header line, the empty claude built-in
   default, and the preflight failure for rejected non-empty values — in the
   README selection and configuration sections.
4. MUST update the canonical Roundfix Skill and its embedded copy together so
   the Skill documents the same optional-effort behavior with zero sync
   drift.

## Subtasks

- [ ] Flip the Project Config codex runtime to gpt-5.6-sol with an explicitly
      empty reasoning effort and the explanatory comment.
- [ ] Update the CONTEXT.md glossary entry.
- [ ] Update the README selection and configuration guidance.
- [ ] Update both Roundfix Skill copies with the optional-effort contract.

## Acceptance Criteria

- [ ] The Project Config selects gpt-5.6-sol with an explicitly empty
      reasoning effort and passes config validation.
- [ ] CONTEXT.md defines the empty Default Reasoning Effort as model-managed.
- [ ] The README documents the optional-effort semantics and the gpt-5.6
      model-managed behavior.
- [ ] The canonical and embedded Roundfix Skill copies document the
      optional-effort contract with zero drift.

## Verification

- `rtk make skills-sync-check` - expected: canonical and embedded Skill
  bundles have zero drift.
- `rtk go run -buildvcs=false ./cmd/roundfix skills check` - expected: every
  shipped Skill passes validation.
- `rtk go test ./internal/config` - expected: config loading accepts the
  shipped Project Config shape.
- `rtk make verify` - expected: formatting, all tests, Skill sync checks,
  Skill checks, and build pass.

## References

- `_prd.md` → Goals; Core Feature 4.
- `_techspec.md` → Coverage Map; Build Order 4; Risks & Considerations.
- ADR-0040; `docs/agents/skill-governance.md`.
