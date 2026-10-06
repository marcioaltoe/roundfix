---
task: task_03
spec: 0234-an-openrouter-key-per-stage
status: pending
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
