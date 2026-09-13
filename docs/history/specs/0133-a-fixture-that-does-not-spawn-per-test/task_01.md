---
task: task_01
spec: 0133-a-fixture-that-does-not-spawn-per-test
status: completed
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

## Result

The canonical clause now says that a non-operative record grants no governed
mutation and does not by itself refuse implementation, commit or push. The
public Baseline update rendered that clause into the agent guide and refreshed
the Setup Manifest. `make baseline-digests` produced the sanctioned pin updates;
a repeat reported `changed: false`.

The public Implement preflight no longer asks every Spec for `implement`.
Before settlement, the Daemon now derives the Task's actual changed paths from
the unfiltered worktree snapshots. It asks the unchanged operation gate for
`implement` and `commit` only when that delta contains a Governed Path. Settle
uses the same boundary for recovered commits, and Spec Run push authority is
asked only when the resolved Spec authorization is operative. The existing
changed-path audit in `internal/speccheck` was not changed.

Focused checks and observations:

- Before the production edit, the combined public-command regression run
  reached both new tests and failed with exit 1: the ordinary Spec stopped in
  Preflight on `authorization operation "implement" is not permitted`, and the
  governed test recorded zero Agent calls. This reproduced the dispatch defect.
- `rtk proxy env GOCACHE=/tmp/roundfix-go-build-0133-task-01 go test
  ./internal/cli -run
  'Test(ImplementDispatchesWithoutRecordWhenNoGovernedMutation|GovernedChangeStillRefusesWithoutRecord|SettleRefusesMissingCommitAuthority|SettleCommitsOrdinaryWorkWithoutRecord)'`
  passed after the implementation.
- `rtk proxy env GOCACHE=/tmp/roundfix-go-build-0133-task-01 go test
  ./internal/authorization ./internal/spec ./internal/daemon` passed. The first
  unchanged attempt reached a sandboxed GPG signing error in a real-repository
  fixture; the permitted unchanged retry with access to the existing signing
  environment exited 0.
- `rtk proxy env GOCACHE=/tmp/roundfix-go-build-0133-task-01 go test
  ./internal/cli ./internal/daemon` passed (`internal/cli` 71.763s,
  `internal/daemon` 11.240s).
- The final targeted regression run covered the authorization refusal identity,
  empty-path validation, public Implement paths, Settle paths, Daemon
  implement/commit/push separation and unfiltered Governed Path detection
  across `internal/authorization`, `internal/spec`, `internal/cli` and
  `internal/daemon`; all four packages exited 0.
- The managed-refresh preview after regeneration returned `state: current` and
  `fileChanges: []`. A second `rtk make baseline-digests` exited 0 and reported
  that the derived artifacts already matched their canonical sources.

Acceptance evidence:

1. `TestImplementDispatchesWithoutRecordWhenNoGovernedMutation` removes the
   record from committed provenance, invokes the public Implement command and
   observes one Agent call plus a Clean command exit.
2. `TestGovernedChangeStillRefusesWithoutRecord` changes `Makefile` through the
   same public command, observes the Agent run, no commit, and the unchanged
   `authorization operation "implement" is not permitted` condition naming the
   canonical record path. `TestGovernedMutationRefusesMissingImplementAuthority`
   locks the same boundary in the Daemon. A source diff confirms the mechanical
   exact-path audit is untouched.
3. `TestAuthorizationReaderRefusesEmptyPaths` observes a typed `paths` refusal
   for an otherwise approved record whose `paths` list is empty.
4. The canonical module and both rendered agent guides contain `grants no
   governed mutation`; the public managed refresh converged to no file changes.
5. `TestRequireGovernedOperationKeepsRefusalIdentity` compares the governed
   refusal byte-for-byte with the existing `RequireOperation` error. The closed
   operation constants and their existing vocabulary test were unchanged and
   passed in the focused package checks.

The commands under `## Verification` remain unrun for Daemon execution.
