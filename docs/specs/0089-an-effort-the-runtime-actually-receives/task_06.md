---
task: task_06
spec: 0089-an-effort-the-runtime-actually-receives
status: pending
type: docs
complexity: low
---

# Task 06: Ship the roundfix skill with the removed refusal

## Overview

The skill currently tells an Agent that `opencode` requires
`reasoning_effort: ""` and that any other value is refused at configuration load
and again at runtime validation. After this Spec that text is false. The
repository's HARD RULE is that a pull request changing CLI behavior ships the
skill update with it. Authorized tooling work with exact bounded files.

## Requirements

1. MUST replace the OpenCode model-managed refusal with the contract this Spec
   delivers: a non-empty effort is accepted, applied after a session warm-up,
   and observed before any work turn.
2. MUST name the `runtime_deferred` encoding and state what Preflight proves
   against what the Run proves.
3. MUST NOT change the skill's fetch, resolve, watch, implement, settle,
   reconcile, archive, release, or baseline guidance.
4. MUST NOT change any agent bundle under the skill's `agents/` directory or any
   other skill.
5. MUST run the sanctioned synchronization command so the embedded copy matches,
   and MUST NOT hand-edit any generated copy or digest.

## Subtasks

- [ ] Replace the refusal paragraph with the new contract.
- [ ] Name the encoding and the split proof.
- [ ] Run the sanctioned synchronization command.
- [ ] Confirm the repository skill check passes.

## Acceptance Criteria

- [ ] The canonical skill no longer claims a non-empty OpenCode effort is refused.
- [ ] The canonical skill names `runtime_deferred` and the session warm-up.
- [ ] The embedded copy matches the canonical skill.
- [ ] No section outside the OpenCode reasoning guidance differs.
- [ ] No other skill in the repository differs.

## Context

- instruction: `docs/workflow/authorizations/2026-08-09-an-effort-the-opencode-runtime-actually-receives.md`
- instruction: `.agents/skills/roundfix/SKILL.md`

## Bounded scope

Authorized by
`docs/workflow/authorizations/2026-08-09-an-effort-the-opencode-runtime-actually-receives.md`.
This Task may create or modify only:

- `.agents/skills/roundfix/SKILL.md`
- `docs/specs/0089-an-effort-the-runtime-actually-receives/task_06.md`

Generated copies under `skills/` rewritten by `make skills-sync` are sanctioned
fallout under ADR-0081, not separate targets.

## Verification

- `grep -q 'runtime_deferred' .agents/skills/roundfix/SKILL.md` — expected: exits 0.
- `! grep -iq 'must be empty' .agents/skills/roundfix/SKILL.md` — expected: exits 0, proving the refusal text is gone.
- `make skills-sync` — expected: exits 0.
- `go run -buildvcs=false ./cmd/roundfix skills check` — expected: exits 0.
- `test -z "$(git diff --name-only -- .agents/skills | grep -v '^\.agents/skills/roundfix/SKILL\.md$')"` — expected: exits 0, proving no other skill moved.

## References

- `_prd.md` → Project Constraints: Tooling authority.
- `_techspec.md` → Implementation Design: API Contracts; Build Order 6.
- ADR-0081, ADR-0108.
