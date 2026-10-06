---
task: task_03
spec: 0234-an-openrouter-key-per-stage
status: completed
type: backend
complexity: low
---

# Task 03: Implementation on an open model reads its own OpenRouter key first and records the variable

## Overview

Spec 0233 adds one helper that chooses the OpenRouter key for implementation
on an open model through OpenCode, a spend record for that stage and the
`openrouter.implement_monthly_ceiling_usd` ceiling. This Task makes that
helper read `ROUNDFIX_OPENROUTER_IMPLEMENT_API_KEY` first and the shared key
second through the stage key list task_01 added, and makes each
implementation spend record name the variable it used. Spec 0233 must be
merged before this Task runs; its helper's file is located in the tree, not
named here.

## Requirements

1. MUST make Spec 0233's implementation-key helper return the variable chosen
   by `openrouterkey.Select(environ, openrouterkey.StageImplement)` per
   Invariant 12, called in that form, and make every place that hands
   OpenCode the key, or checks `openrouter.implement_monthly_ceiling_usd`,
   use that one selected variable. No reader may name
   `ROUNDFIX_OPENROUTER_API_KEY` itself for this stage.
2. MUST add `key_variable`, the selected variable's name, to each
   implementation spend record Spec 0233 writes, per Invariant 13; a record
   never contains a key value, and records written before this change still
   parse and still count toward the ceiling.
3. MUST add, beside the helper and the record and with the fakes Spec 0233's
   tests use, `TestOpenModelImplementationPrefersTheImplementStageKey` (both
   variables set to different sentinels: the implementation variable is
   chosen and only its value reaches OpenCode's configuration),
   `TestOpenModelImplementationFallsBackToTheSharedKey` (only the shared key
   set: it is chosen; only the generic `OPENROUTER_API_KEY` set: no key and
   0233's existing no-key behavior), and
   `TestOpenModelImplementationSpendRecordNamesItsKeyVariable` (the record
   carries `key_variable` and neither sentinel).
4. MUST update only the existing light-tier tests whose expected output this
   Task changes, and record any file it changes beyond those it can name now
   through the Daemon's `## Recorded paths`.
5. MUST NOT start OpenCode against OpenRouter, read the process environment
   in a test, change which models a light-tier selection may name, change the
   ceiling's value or sum, or change `CONTEXT.md` or `CHANGELOG.md`.

## Subtasks

- [ ] Locate Spec 0233's helper, spend record and ceiling check.
- [ ] Route the helper through the implementation stage's list.
- [ ] Add `key_variable` to the spend record.
- [ ] Add the three tests beside them.

## Acceptance Criteria

- [ ] With both variables set, implementation uses the implementation key
      and its spend record names it.
- [ ] With only the shared key, implementation behaves as Spec 0233 shipped
      it and names the shared key.
- [ ] The generic key alone gives no key, and no record contains a value.

## Context

- instruction: `docs/adr/0239-each-openrouter-stage-reads-its-own-roundfix-key-first.md`
- instruction: `docs/adr/0235-openai-and-anthropic-models-run-only-through-the-codex-and-claude-subscriptions.md`

## Verification

- `for name in TestOpenModelImplementationPrefersTheImplementStageKey TestOpenModelImplementationFallsBackToTheSharedKey TestOpenModelImplementationSpendRecordNamesItsKeyVariable; do file="$(grep -rl --include='*_test.go' -e "func $name(" internal)" || { printf 'missing test: %s\n' "$name" >&2; exit 1; }; out="$(go test -count=1 -v -run "^$name\$" "./$(dirname "$file")" 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task none of the three tests exists, so the command fails on any tree.
- `matches="$(grep -rn --include='*.go' -e 'openrouterkey[.]Select(.*openrouterkey[.]StageImplement' internal)" || { printf 'no implementation reader selects the implementation stage key\n' >&2; exit 1; }; printf '%s\n' "$matches" | grep -v -e '_test[.]go:' -e '^internal/openrouterkey/' | grep -q . || { printf 'only tests or the key package select the implementation stage key\n' >&2; exit 1; }` — expected: exit 0; before this Task no non-test reader outside the key package calls Select for the implementation stage, so the command fails.

## References

- `_prd.md` → Prerequisites; Goals; User Stories 1, 2, 4; Core Features 2, 4; Success Metric 3; Success Metric 5
- `_techspec.md` → Invariants 12, 13; Data Models; API Contract 4; Testing Approach; Build Order 3; Risks & Considerations
- ADR-0239; ADR-0235

## Result

Implemented the Task 03 slice for Daemon Verification. Task status, the
authored Verification commands, and the Subtasks/Acceptance checkboxes remain
Daemon-owned. No commit, push or Pull Request was made.

Spec 0233's helper `OpenRouterImplementKey` in `internal/config/light_tier.go`
now takes the command environ and returns exactly
`openrouterkey.Select(environ, openrouterkey.StageImplement)`, so the
implementation stage reads `ROUNDFIX_OPENROUTER_IMPLEMENT_API_KEY` first and
the shared `ROUNDFIX_OPENROUTER_API_KEY` second. The helper no longer names the
shared variable itself. `implementLightTierPlan` passes `environment.environ`
straight to it; when neither variable is set it keeps the stage's preferred
variable name so the existing no-key skip still names a real variable. Every
place that hands OpenCode the key (`Plan.KeyVariable` ->
`RuntimeSpec.OpenRouterKeyVariable` -> `lightSessionEnvironment`) and the
ceiling check in `applyLightTier` use that one selected variable.

`lighttier.SpendLine` gained `key_variable`, and `recordLightSpend` writes
`owner.lightPlan.KeyVariable` into each record. No record carries a key value,
and records written before this change still parse and still sum toward the
ceiling.

Acceptance-criterion evidence:

- **Both variables set: implementation key used and recorded:**
  `TestOpenModelImplementationPrefersTheImplementStageKey` (both OpenRouter
  variables and the generic one set to different sentinels) selects
  `ROUNDFIX_OPENROUTER_IMPLEMENT_API_KEY` and builds OpenCode's inline
  configuration with only `{env:ROUNDFIX_OPENROUTER_IMPLEMENT_API_KEY}`; the
  generic `OPENROUTER_API_KEY` is removed and neither the shared variable nor
  a sentinel value appears.
- **Shared fallback and no key:** `TestOpenModelImplementationFallsBackToTheSharedKey`
  selects and configures the shared variable when it is the only Roundfix key,
  and selects no variable when only the generic `OPENROUTER_API_KEY` is set.
  `TestOpenModelImplementationSpendRecordNamesItsKeyVariable` writes a record
  with `key_variable` and neither sentinel, and proves a pre-change record
  without the field still parses and still counts.
- **Every reader selects through the stage list:** a non-test file outside the
  key package calls `openrouterkey.Select(..., openrouterkey.StageImplement)`
  (`internal/config/light_tier.go`); no non-test reader names the shared key for
  this stage.

Focused checks (all with `GOCACHE=/private/tmp/roundfix-0234-task03-gocache`),
run directly rather than through the declared Verification commands:

- `go build ./...`: exit 0.
- `go test -count=1 -v -run '^TestOpenModelImplementation' ./internal/agent
  ./internal/lighttier`: the three named tests each printed
  `--- PASS: ...`; `ok` for both packages.
- `go test -count=1 ./internal/config`: exit 0, 558 passed.
- `go test -count=1 ./internal/agent ./internal/lighttier`: exit 0, 553
  passed.
- `go test -count=1 -run 'TestImplement' ./internal/cli`: 59 passed, including
  the updated `TestImplementBuildsTheLightTierPlan` (its no-key warning now
  names the implementation variable).
- `go test -count=1 ./internal/daemon`: 567 passed, including
  `TestLightTierRecordsSpend` and the skip-warning cases.
- `gofmt -l` on every changed Go file: no output.

The declared Verification commands were not run. Files changed beyond the
helper and record include `internal/config/light_tier_test.go`,
`internal/agent/light_session_test.go`, `internal/cli/implement_light_tier_test.go`
and `internal/daemon/light_tier_test.go`, all adapting to the helper's new
environ parameter or adding the required tests. No `CONTEXT.md`, `CHANGELOG.md`,
model selection, ceiling value or sum changed.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `internal/agent/light_session_test.go`
- `internal/cli/implement.go`
- `internal/cli/implement_light_tier_test.go`
- `internal/config/light_tier.go`
- `internal/config/light_tier_test.go`
- `internal/daemon/agent_session_owner.go`
- `internal/daemon/light_tier_test.go`
- `internal/lighttier/spend_log.go`
- `internal/lighttier/spend_log_test.go`
