---
task: task_01
spec: 0133-a-fixture-that-does-not-spawn-per-test
status: pending
type: backend
complexity: high
---

# Task 01: State the two authorities apart

## Overview

The canonical clause says an absent authorization record "grants nothing", and
the operation gate obeyed it by requiring the `implement` operation from every
Spec. Since a record must name at least one path to be valid, a Spec touching no
Governed Path can neither write a record nor be implemented. This slice narrows
the clause to governed mutation and moves the gate to match. It is verifiable
alone: a Spec with no protected mutation dispatches, and a governed change
without a record still refuses.

This Task's governed scope is exactly the three paths the approved grant bounds:
`internal/baseline/assets/modules/core.json`, `docs/agents/agent-instructions.md`
and `docs/agents/setup-context.json`, plus the derived pins that
`make baseline-digests` rewrites as sanctioned regeneration. It also changes
ordinary source in `internal/authorization`, `internal/cli` and
`internal/daemon`, which is not governed and stays within this Spec's declared
behavior. Stop before any other mutation.

## Requirements

1. MUST narrow the canonical clause so an absent, proposed, contradictory or
   withdrawn record withholds governed mutation and does not by itself refuse
   implementation, commit or push.
2. MUST regenerate the rendered guides and the manifest from the module through
   the public Baseline update, then regenerate the sanctioned derived pins, and
   MUST NOT hand-edit a derived value.
3. MUST move the operation gate to ask for an operation where a governed
   mutation is at stake rather than at every dispatch, so a Spec that declares
   no protected tooling mutation dispatches with no record present.
4. MUST keep every governed-path refusal exactly as it is: a change to a
   Governed Path still requires an operative record naming that exact path, an
   absent record still refuses it, and the changed-path audit is untouched.
5. MUST keep a record that exists required to name at least one path. The
   repair is that a Spec with nothing to bound needs no record, not that it
   writes an empty one.
6. MUST keep the `operations` vocabulary and every refusal message identity
   unchanged, so a caller matching on them keeps working.

## Subtasks

- [ ] Narrow the clause in the canonical module.
- [ ] Regenerate the rendered guides and the sanctioned pins.
- [ ] Move the gate to the governed-mutation boundary.
- [ ] Prove every governed-path refusal is unchanged.

## Acceptance Criteria

- [ ] A Spec declaring no protected tooling mutation, with no record present,
      dispatches through the public Implement command; it refuses today with
      `operation "implement" is not permitted`.
- [ ] A change to a Governed Path with no operative record still refuses, and
      the refusal names the same condition it names today.
- [ ] A record that exists with an empty `paths` list is still invalid.
- [ ] The rendered guides state the division between governed mutation and the
      work itself, and a second managed refresh reports no further change.
- [ ] The `operations` vocabulary and every refusal message identity are
      unchanged.

## Context

- interface: `internal/cli/implement.go`
- interface: `internal/authorization/authorization.go`
- instruction: `docs/agents/agent-instructions.md`

## Verification

- `grep -q 'grants no governed mutation' internal/baseline/assets/modules/core.json && grep -q 'grants no governed mutation' docs/agents/agent-instructions.md` — the clause states the division and the rendered guide carries it; neither says this today.
- `grep -q 'func TestImplementDispatchesWithoutRecordWhenNoGovernedMutation' internal/cli/implement_test.go && go test -count=1 ./internal/cli -run '^TestImplementDispatchesWithoutRecordWhenNoGovernedMutation$'` — a Spec with no protected mutation dispatches with no record; this fails today.
- `grep -q 'func TestGovernedChangeStillRefusesWithoutRecord' internal/cli/implement_test.go && go test -count=1 ./internal/cli -run '^TestGovernedChangeStillRefusesWithoutRecord$'` — a governed change with no record still refuses, so narrowing cost no authority.
- `grep -q 'grants no governed mutation' internal/baseline/assets/modules/core.json || exit 1; go test -count=1 ./internal/authorization ./internal/spec ./internal/speccheck` — every reader and audit still passes with the narrowed rule.
- `grep -q 'grants no governed mutation' internal/baseline/assets/modules/core.json || exit 1; raw="$(mktemp)"; before="$(mktemp)"; after="$(mktemp)"; go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text > /dev/null || exit 1; find docs/agents internal/baseline -type f -exec shasum {} + > "$raw" || exit 1; sort "$raw" > "$before" || exit 1; go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text > /dev/null || exit 1; find docs/agents internal/baseline -type f -exec shasum {} + > "$raw" || exit 1; sort "$raw" > "$after" || exit 1; diff "$before" "$after"` — the managed refresh converges: a second run reproduces the first byte for byte.

## References

- `_prd.md` → Goals 1, 3; User Stories 1-2; Core Features 1-3; Decisions: Declared intentional breaks 1.
- `_techspec.md` → Implementation Design: Two authorities, stated apart; Build Order 1.
- `_authorization.md` → approved bounded paths and sanctioned regeneration.
