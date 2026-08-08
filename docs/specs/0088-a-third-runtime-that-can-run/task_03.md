---
task: task_03
spec: 0088-a-third-runtime-that-can-run
status: pending
type: backend
complexity: medium
---

# Task 03: Make OpenCode reasoning effort model-managed

## Overview

Roundfix applies a non-empty reasoning effort as a separate ACP config option
after ensuring the Agent Session, and on `opencode` that call cannot succeed
before the Run's first prompt — the live agent process is still on the runtime
default and advertises no `effort` option, so the adapter answers ACP `-32602`.
This Task makes `opencode` a model-managed reasoning runtime: configuration
refuses a non-empty effort with the repair named, and runtime validation refuses
it again so no invocation override slips past.

## Requirements

1. MUST stop mapping `opencode` to a reasoning-effort config key, so runtime
   validation rejects a non-empty effort on that runtime for every entry point,
   including a `--reasoning-effort` invocation override.
2. MUST refuse a non-empty `reasoning_effort` on `runtime: opencode` during Agent
   Selection normalization, so a maintainer sees it when the configuration loads
   rather than when a Run starts.
3. MUST name the repair in the refusal text: the empty value, and why — OpenCode
   advertises reasoning effort per model and only after an Agent Session's first
   prompt, so a token-free Exact Agent Selection Proof cannot apply one.
4. MUST apply the same refusal to the legacy runtime-defaults path, so a legacy
   `runtimes.opencode.reasoning_effort` cannot produce a profile that
   normalization would reject.
5. MUST keep an `opencode` Agent Selection with an empty `reasoning_effort` fully
   valid, and MUST NOT let it reach the disposable-effort application path.
6. MUST NOT change how any Codex or Claude selection is validated, mapped, or
   applied.
7. MUST re-record the coverage record in this Task's own commit if any test is
   renamed or removed.

## Subtasks

- [ ] Remove `opencode` from the reasoning-effort config key mapping.
- [ ] Add the normalization refusal with its repair text.
- [ ] Cover the legacy runtime-defaults path with the same refusal.
- [ ] Confirm an empty-effort `opencode` selection skips effort application on
      both the disposable and the fallback paths.
- [ ] Edit the break-half characterization test that pinned today's acceptance,
      and declare the break in this Task's Result.

## Acceptance Criteria

- [ ] Decoding a profile with `runtime: opencode` and a non-empty
      `reasoning_effort` fails, and the message names the empty value as the
      repair and states why.
- [ ] Decoding the same profile with `reasoning_effort: ""` succeeds.
- [ ] Runtime validation rejects an `opencode` runtime carrying a non-empty
      reasoning effort, independently of configuration decoding.
- [ ] A legacy `runtimes.opencode.reasoning_effort` that is non-empty is refused
      with the same repair.
- [ ] Codex and Claude selections with a non-empty effort still validate and
      still map to their existing config keys.
- [ ] An `opencode` selection with an empty effort issues no reasoning config set
      on the disposable-session path.

## Context

- interface: `internal/agent/acpx_runner.go`
- interface: `internal/agent/fallback.go`
- interface: `internal/config/profiles.go`
- interface: `internal/config/config.go`

## Bounded scope

This Task may create or modify only:

- `internal/agent/acpx_runner.go`
- `internal/agent/acpx_runner_test.go`
- `internal/agent/fallback.go`
- `internal/config/profiles.go`
- `internal/config/config.go`
- `internal/config/config_test.go`
- `internal/config/profiles_characterization_test.go`
- `docs/references/coverage-record.json`
- `docs/specs/0088-a-third-runtime-that-can-run/task_03.md`

## Verification

- `go build -buildvcs=false ./...` — expected: exits 0.
- `go test ./internal/agent ./internal/config -count=1` — expected: exits 0.
- `go test ./internal/config -run 'OpenCodeReasoning' -count=1 -v` — expected: exits 0 and names at least one test; `no tests to run` fails this Task.
- `go test ./internal/agent -run 'ReasoningEffortConfigKey' -count=1 -v` — expected: exits 0 and names at least one test.
- `go test ./internal/spec -run '^TestCoverageEquivalence$' -count=1` — expected: exits 0.
- `grep -q 'model-managed' internal/config/profiles.go` — expected: exits 0, proving the repair is named in the refusal text rather than only in the Spec.

## References

- `_prd.md` → Goal 4; Core Features: model-managed reasoning for OpenCode.
- `_techspec.md` → Implementation Design: API Contracts; Build Order 3.
- `references/2026-08-08-what-the-opencode-adapter-answers-before-its-first-prompt.md`
  → the nine-step sequence proving the effort cannot be applied before the first
  prompt.
- ADR-0106.
