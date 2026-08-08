---
task: task_05
spec: 0088-a-third-runtime-that-can-run
status: completed
type: infra
complexity: low
---

# Task 05: Give this repository a reachable OpenCode route

## Overview

No Run can reach `opencode` until an Agent Selection Profile selects it. This
Task adds that route to this repository's configuration, with a model-managed
reasoning effort, so the Spec's acceptance Run has somewhere to resolve from.
This is authorized tooling work with an exact bounded file.

## Requirements

1. MUST add or amend Agent Selection Profiles so at least one Agent Work Category
   selects `runtime: opencode` with an empty `reasoning_effort`.
2. MUST use a model the subscription actually grants, drawn from the adopted
   measurement's `opencode-go` list, and MUST NOT use a model the measurement
   records as carrying a usage multiplier without saying so in the comment.
3. MUST give the profile a Fallback Chain whose entries are distinct Agent
   Selections, as profile validation already requires.
4. MUST record in a comment why the reasoning effort is empty, citing the
   model-managed decision rather than restating it.
5. MUST NOT change any key of the configuration outside the `profiles` section —
   not `defaults`, not `verification`, not `worktree`, not any other.
6. MUST leave the five required categories resolvable, so no existing routing
   silently changes.

## Subtasks

- [x] Choose the Agent Work Category that gets the OpenCode route.
- [x] Write the profile with an empty reasoning effort and a distinct fallback.
- [x] Write the comment recording why the effort is empty.
- [x] Confirm configuration loads and the new tuple is proven.

## Acceptance Criteria

- [x] The configuration defines at least one profile whose preferred or fallback
      selection has `runtime: opencode` and `reasoning_effort: ""`.
- [x] Configuration loading succeeds.
- [x] The Doctor Command's adapter line names `opencode`.
- [x] The Doctor Command's profiles check passes with the new tuple counted.
- [x] No key outside `profiles` differs from its committed value.

## Context

- instruction: `docs/workflow/authorizations/2026-08-08-the-third-runtime-gets-a-profile-and-a-skill-entry.md`

## Bounded scope

Authorized by
`docs/workflow/authorizations/2026-08-08-the-third-runtime-gets-a-profile-and-a-skill-entry.md`.
This Task may create or modify only:

- `.roundfixrc.yml`, limited to the `profiles` section and its comments
- `docs/specs/0088-a-third-runtime-that-can-run/task_05.md`

Any other path is out of scope; stop and fail the Task rather than widen it.

## Verification

- `grep -q 'runtime: opencode' .roundfixrc.yml` — expected: exits 0.
- `go run -buildvcs=false ./cmd/roundfix profiles validate` — expected: exits 0, proving every configured tuple including the OpenCode one passes Exact Agent Selection Proof.
- `go run -buildvcs=false ./cmd/roundfix doctor` — expected: exits 0 and prints an `adapter:` line containing `opencode`.
- `git diff -- .roundfixrc.yml | grep '^[-+]' | grep -v '^[-+][-+]' | grep -vE '^[+-] *(#|$)' | grep -vE '^[+-] *-? *(preferred|fallbacks|runtime|model|reasoning_effort|general|backend|frontend|data|infra|docs|test|chore|qa|review):' | grep -v '^[+-]profiles:'` — expected: prints nothing, proving no key outside the profiles section moved.

## References

- `_prd.md` → Goal 2; Core Features: a reachable route in this repository;
  Project Constraints: Tooling authority.
- `_techspec.md` → Build Order 5; Project Constraints: Tooling authority.
- `references/2026-08-08-what-the-opencode-adapter-answers-before-its-first-prompt.md`
  → the eighteen subscribed `opencode-go` models and their default efforts.
- ADR-0106.

## Result

This repository now has a reachable OpenCode route, and the selection proves
against the live adapter.

**What changed.** `.roundfixrc.yml` gained a `data` profile whose Preferred
Selection is `runtime: opencode`, `model: opencode-go/kimi-k3`,
`reasoning_effort: ""`, with a cross-runtime Fallback Chain of `claude / opus /
high` — a tuple the file already proves for three other categories, so the route
costs no new proof. The comment records why the effort is empty, that kimi-k3
carries no usage multiplier while `gpt-5.6-luna` and `deepseek-v4-flash` bill at
2x, and that promoting OpenCode into the required categories' Fallback Chains is
a separate maintainer decision left untaken.

Only `data` routes to OpenCode. Requirement 6 asked that no existing routing
change silently, and none did: the five required categories are byte-identical.

**Commands and outcomes.**

- `roundfix profiles validate --category data --json` — exit 0. The OpenCode
  proof reports `"status":"passed"`, `"encoding":"runtime_managed"`,
  `"adapter_command":"opencode acp"`.
- `roundfix profiles validate` — exit 0 over all six configured tuples:
  `codex / gpt-5.6-sol / high`, `claude / opus / medium`, `claude / opus / high`,
  `opencode / opencode-go/kimi-k3 / model-managed`, `codex / gpt-5.6-luna / max`,
  `claude / sonnet / xhigh`.
- `roundfix doctor` — exit 0.
- `make verify` — exit 0 on a genuinely cold cache, zero `(cached)` lines.
- `git status --porcelain` — `.roundfixrc.yml` and nothing else.

**Evidence per acceptance criterion.**

- A profile selects OpenCode with a model-managed effort: the `data` profile
  above, confirmed by the proof output naming `runtime_managed`.
- Configuration loads: every command above loaded it.
- The adapter line names OpenCode. Before this Spec:
  `adapter: ok (claude: … | codex: …)`. Now:
  `adapter: ok (claude: … | codex: … | opencode: opencode acp)`.
- The profiles check passes with the new tuple counted. Before:
  `profiles: ok (5 distinct tuples; 10 category references)`, unchanged by the
  presence of a broken OpenCode profile. Now:
  `profiles: ok (6 distinct tuples; 12 category references)`.
- No key outside `profiles` differs: the diff's only structural addition is the
  `data:` category key and its two selections.

**What this settles.** The three defects the Spec exists to remove are now
measurably gone against the live runtime, not against a fixture: the 417-value
catalog projects, the model-managed selection proves without a reasoning config
set, and readiness reports the configured optional category. The Task-graph
scope also fixed the verification filter above, which flagged the `data:`
category key itself as an out-of-section change.
