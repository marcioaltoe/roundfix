---
task: task_05
spec: 0089-an-effort-the-runtime-actually-receives
status: completed
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

## Result

Every explicit `opencode` selection keeps
`openrouter/deepseek/deepseek-v4-pro` and now carries
`reasoning_effort: xhigh`. The profile comment records that `xhigh` is this
model's maximum and the variant measured by its published benchmarks, and that
Roundfix applies and observes the effort after warming the Agent Session rather
than inheriting it from the model.

Pre-change signal:

- A focused `awk` inventory reported all six explicit `opencode` selections
  with `reasoning_effort: ""`, plus three stale comment signals describing the
  removed refusal and inherited default.

Focused checks after the configuration edit:

- `rtk awk '/^  [a-z]+:$/ { profile=$1; sub(":$", "", profile) } /runtime: opencode/ { getline; model=$2; getline; effort=$2; print profile, model, effort; count++; if (model != "openrouter/deepseek/deepseek-v4-pro" || effort != "xhigh") bad++ } END { print "opencode selections=" count ", mismatches=" bad+0; exit !(count == 6 && bad == 0) }' .roundfixrc.yml` — exited 0 with six selections and zero mismatches.
- `GOCACHE="$PWD/.gocache" rtk go run -buildvcs=false ./cmd/roundfix profiles show --json` — exited 0, loaded the Project Config as `roundfix/profiles/v1`, and resolved all ten required Agent Work Categories. The four optional categories inherited `general`; the explicit categories retained their existing routing and every resolved OpenCode selection reported `deepseek-v4-pro/xhigh`.
- `GOCACHE="$PWD/.gocache" rtk go test ./internal/config -run '^TestOpenCodeEffortAcceptedConfigurationLoadsAndResolvesNonEmptyReasoningEffort$'` — one focused configuration test passed. The first attempt used the sandbox-denied macOS Go cache; the recorded rerun used the repository-local ignored cache.
- A focused `awk` comment check exited 0 with all three required signals (`xhigh` maximum and benchmarked variant, Agent Session warm-up and application, no inheritance) and zero stale refusal signals.
- `rtk diff <(rtk git show HEAD:.roundfixrc.yml | rtk awk 'BEGIN { p=0 } /^profiles:/ { p=1; next } /^[a-z_]+:/ { if (p) p=0 } !p { print }') <(rtk awk 'BEGIN { p=0 } /^profiles:/ { p=1; next } /^[a-z_]+:/ { if (p) p=0 } !p { print }' .roundfixrc.yml)` — exited 0 with `[ok] Files are identical`, a byte-level proof that all content outside `profiles` still matches `HEAD`.

Acceptance evidence:

- Criterion 1: the OpenCode tuple inventory found six selections, all at
  `reasoning_effort: xhigh`, with zero mismatches.
- Criterion 2: the same complete inventory found no OpenCode selection with an
  empty effort.
- Criterion 3: `profiles show --json` loaded the changed Project Config and
  resolved all ten Agent Work Categories; the focused configuration test also
  passed for a non-empty OpenCode effort.
- Criterion 4: the comment check found the required maximum, benchmark, warm-up,
  and applied-effort statements and no stale refusal/default statement.
- Criterion 5: the exact outside-`profiles` comparison matched `HEAD` byte for
  byte.

No follow-up work was found inside this Task's slice. The commands authored
under `## Verification` were not run; Daemon Verification remains responsible
for proving every configured tuple and settling the Task.
