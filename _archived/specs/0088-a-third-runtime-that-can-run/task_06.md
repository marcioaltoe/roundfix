---
task: task_06
spec: 0088-a-third-runtime-that-can-run
status: completed
type: docs
complexity: low
---

# Task 06: Ship the roundfix skill with the CLI change

## Overview

This Spec changes observable CLI behavior twice: the Doctor Command's profile
readiness and adapter enumeration now cover every configured Agent Work Category,
and configuration refuses a non-empty reasoning effort on the `opencode` runtime.
The repository's HARD RULE is that a pull request changing CLI behavior ships the
roundfix skill update with it. This is authorized tooling work with an exact
bounded file.

## Requirements

1. MUST describe, in the section that already documents the Doctor Command, that
   Agent Selection Profile Readiness covers every configured Agent Work Category
   and that the adapter line enumerates the ACP Runtimes those tuples reference.
2. MUST describe that a non-empty reasoning effort on the `opencode` runtime is
   refused, and that the empty value means the Agent Model manages reasoning.
3. MUST NOT change the skill's fetch, resolve, watch, implement, settle,
   reconcile, archive, release, or baseline guidance.
4. MUST NOT change any agent bundle under the skill's `agents/` directory or any
   other skill in the repository.
5. MUST run the sanctioned synchronization command so the embedded copy matches
   the canonical skill, and MUST NOT hand-edit any generated copy or digest.

## Subtasks

- [x] Add the widened readiness description to the Doctor Command section.
- [x] Add the OpenCode model-managed reasoning refusal to the configuration
      guidance.
- [x] Run the sanctioned synchronization command.
- [x] Confirm the repository skill check passes.

## Acceptance Criteria

- [x] The canonical skill names the configured-category readiness scope.
- [x] The canonical skill names the OpenCode model-managed reasoning refusal.
- [x] The embedded copy matches the canonical skill.
- [x] No section of the skill outside the Doctor Command and configuration
      guidance differs.
- [x] No other skill in the repository differs.

## Context

- instruction: `docs/workflow/authorizations/2026-08-08-the-third-runtime-gets-a-profile-and-a-skill-entry.md`
- instruction: `docs/agents/skill-dispatch.md`
- instruction: `.agents/skills/roundfix/SKILL.md`

## Bounded scope

Authorized by
`docs/workflow/authorizations/2026-08-08-the-third-runtime-gets-a-profile-and-a-skill-entry.md`.
This Task may create or modify only:

- `.agents/skills/roundfix/SKILL.md`
- `docs/specs/0088-a-third-runtime-that-can-run/task_06.md`

Generated copies under `skills/` rewritten by `make skills-sync` are sanctioned
fallout under ADR-0081, not separate targets. Any other path is out of scope;
stop and fail the Task rather than widen it.

## Verification

- `grep -q 'opencode' .agents/skills/roundfix/SKILL.md` — expected: exits 0.
- `grep -qi 'model-managed' .agents/skills/roundfix/SKILL.md` — expected: exits 0, proving the reasoning refusal is documented.
- `make skills-sync` — expected: exits 0.
- `go run -buildvcs=false ./cmd/roundfix skills check` — expected: exits 0, proving the embedded copy matches the canonical skill.
- `git diff --name-only -- .agents/skills | grep -v '^\.agents/skills/roundfix/SKILL\.md$'` — expected: prints nothing, proving no other skill moved.

## References

- `_prd.md` → Project Constraints: Tooling authority.
- `_techspec.md` → Implementation Design: API Contracts; Build Order 6.
- `docs/agents/specific-repository.md` → the roundfix skill synchronization HARD
  RULE.
- ADR-0081, ADR-0106, ADR-0107.

## Result

The roundfix skill now describes both CLI behaviors this Spec changed.

**What changed.** The Doctor Command section says profile readiness covers every
Agent Work Category the effective configuration defines — the required five plus
each optional category a profile actually declares — that a category resolving
only by inheriting `general` adds no tuple, that the `adapter:` line follows the
same scope and names a runtime only an optional category selects, and that a
failing configured profile now fails the `profiles:` line instead of leaving it
`ok`. A second paragraph documents the `opencode` reasoning contract:
`reasoning_effort: ""` is required, any other value is refused at configuration
load and again by runtime validation so a `--reasoning-effort` override cannot
bypass it, the measured reason, and that such a selection proves with encoding
`runtime_managed` rather than `model_managed`. One word changed in the opening
sentence: the checks list now reads "configured Agent Selection Profiles"
instead of "required".

**Commands and outcomes.**

- `make skills-sync` — exit 0; rewrote `skills/roundfix/SKILL.md` as sanctioned
  fallout under ADR-0081.
- `roundfix skills check` — exit 0 over all fourteen owned skills.
- `grep -q 'opencode' .agents/skills/roundfix/SKILL.md` — exit 0.
- `grep -qi 'model-managed' .agents/skills/roundfix/SKILL.md` — exit 0.
- `git status --porcelain -- .agents/skills` — only
  `.agents/skills/roundfix/SKILL.md`.
- `make verify` — exit 0 on a genuinely cold cache, zero `(cached)` lines.

**Evidence per acceptance criterion.**

- Configured-category readiness scope named: the added paragraph enumerates the
  optional categories and states the inheritance exclusion.
- OpenCode refusal named: the second paragraph states the required value, both
  enforcement points, and the encoding.
- Embedded copy matches: `roundfix skills check` passed, and
  `make verify` runs `skills-sync-check` before `skills-check`.
- No other section differs: the diff touches two contiguous regions of the
  Doctor Command section and one word in its opening sentence.
- No other skill differs: the scoped `git status` printed nothing else.

**Bounded scope honored.** Two paths changed:
`.agents/skills/roundfix/SKILL.md`, which the authorization names, and
`skills/roundfix/SKILL.md`, which `make skills-sync` regenerates and ADR-0081
classifies as sanctioned fallout rather than a separate target. No digest was
hand-edited.
