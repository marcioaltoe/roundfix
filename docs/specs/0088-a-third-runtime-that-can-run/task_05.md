---
task: task_05
spec: 0088-a-third-runtime-that-can-run
status: pending
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

- [ ] Choose the Agent Work Category that gets the OpenCode route.
- [ ] Write the profile with an empty reasoning effort and a distinct fallback.
- [ ] Write the comment recording why the effort is empty.
- [ ] Confirm configuration loads and the new tuple is proven.

## Acceptance Criteria

- [ ] The configuration defines at least one profile whose preferred or fallback
      selection has `runtime: opencode` and `reasoning_effort: ""`.
- [ ] Configuration loading succeeds.
- [ ] The Doctor Command's adapter line names `opencode`.
- [ ] The Doctor Command's profiles check passes with the new tuple counted.
- [ ] No key outside `profiles` differs from its committed value.

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
- `git diff -- .roundfixrc.yml | grep '^[-+]' | grep -v '^[-+][-+]' | grep -v -e 'profiles' -e 'preferred' -e 'fallbacks' -e 'runtime:' -e 'model:' -e 'reasoning_effort:' -e '^[-+] *#' -e '^[-+] *$' -e '^[-+] *- '` — expected: prints nothing, proving no key outside the profiles section moved.

## References

- `_prd.md` → Goal 2; Core Features: a reachable route in this repository;
  Project Constraints: Tooling authority.
- `_techspec.md` → Build Order 5; Project Constraints: Tooling authority.
- `references/2026-08-08-what-the-opencode-adapter-answers-before-its-first-prompt.md`
  → the eighteen subscribed `opencode-go` models and their default efforts.
- ADR-0106.
