---
task: task_05
spec: 0089-an-effort-the-runtime-actually-receives
status: pending
type: infra
complexity: low
---

# Task 05: Route this repository to deepseek-v4-pro at xhigh

## Overview

The capability is only demonstrated by a profile that uses it. The maintainer
chose `openrouter/deepseek/deepseek-v4-pro` at `xhigh` — the model's own maximum
and the variant the published benchmarks measure — and asked for it as soon as
the effort capability exists. This is authorized tooling work with exact bounded
files.

## Requirements

1. MUST set every `opencode` selection to
   `openrouter/deepseek/deepseek-v4-pro` with `reasoning_effort: xhigh`.
2. MUST record in a comment that `xhigh` is this model's maximum and the
   benchmarked variant, and that the effort is now applied after a session
   warm-up rather than inherited.
3. MUST NOT change any key of the configuration outside the `profiles` section.
4. MUST leave every required Agent Work Category resolvable, so no existing
   routing silently changes beyond the effort.
5. MUST remove the comment lines that describe the removed refusal, so the file
   stops documenting a rule that no longer exists.

## Subtasks

- [ ] Set the effort on every `opencode` selection.
- [ ] Rewrite the comment that explained the empty effort.
- [ ] Confirm configuration loads and every tuple is proven.

## Acceptance Criteria

- [ ] Every `opencode` selection carries `reasoning_effort: xhigh`.
- [ ] No `opencode` selection carries an empty `reasoning_effort`.
- [ ] Configuration loading succeeds.
- [ ] The comment no longer claims a non-empty effort is refused.
- [ ] No key outside `profiles` differs from its committed value.

## Context

- instruction: `docs/workflow/authorizations/2026-08-09-an-effort-the-opencode-runtime-actually-receives.md`

## Bounded scope

Authorized by
`docs/workflow/authorizations/2026-08-09-an-effort-the-opencode-runtime-actually-receives.md`.
This Task may create or modify only:

- `.roundfixrc.yml`, limited to the `profiles` section and its comments
- `docs/specs/0089-an-effort-the-runtime-actually-receives/task_05.md`

Any other path is out of scope; stop and fail the Task rather than widen it.

## Verification

- `grep -q 'reasoning_effort: xhigh' .roundfixrc.yml` — expected: exits 0.
- `! grep -A 2 'runtime: opencode' .roundfixrc.yml | grep -q 'reasoning_effort: ""'` — expected: exits 0, proving no OpenCode selection kept the empty effort.
- `go run -buildvcs=false ./cmd/roundfix profiles validate` — expected: exits 0, proving every configured tuple including the deferred-effort one passes Exact Agent Selection Proof.
- `test -z "$(git diff -- .roundfixrc.yml | grep '^[-+]' | grep -v '^[-+][-+]' | grep -vE '^[+-] *(#|$)' | grep -vE '^[+-] *-? *(preferred|fallbacks|runtime|model|reasoning_effort|general|backend|frontend|data|infra|docs|test|chore|qa|review):' | grep -v '^[+-]profiles:')"` — expected: exits 0, proving no key outside the profiles section moved.

## References

- `_prd.md` → Goal 1; Project Constraints: Tooling authority.
- `_techspec.md` → Build Order 5.
- ADR-0108.
