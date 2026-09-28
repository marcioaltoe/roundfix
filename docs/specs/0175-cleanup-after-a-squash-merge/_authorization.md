---
status: approved
granted: 2026-09-28
action: prove a squash-merged Spec's Runs against the merged head, release them after a delivery merge and through reconcile --apply, key legacy Runs from their Run Worktree, own and sweep carry-forward staging worktrees, refuse to delete a nested registered worktree, and route every production git worktree administration through the worktree administration lock
consuming: 0175-cleanup-after-a-squash-merge
paths:
  - internal/cli/cli_test.go
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0175

The maintainer approved Onda 2 of the efficiency sequence in chat on 2026-09-28
("Siga para onda 2 como sugerido até o release"); this Spec is part of it. The
Go test file rides the standing grant of 2026-09-21 for governed source; the
skill files ride the standing grant of 2026-09-18 for keeping the shipped
skills true to the CLI. The set was measured with `GovernedPath`.

## Why each governed path is unavoidable

- `internal/cli/cli_test.go` — `TestRunReconcileJSONMatchesTextFields` pins the
  exact top-level field names of the `roundfix-reconcile/v1` report, and
  task_05 adds `stagingCandidates`, so that list gains one name. The test keeps
  its name, so `docs/references/coverage-record.json` does not change.
- `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md` — the Roundfix
  skill describes the Reconcile Command (task_03, task_05) and the Delivery
  queue (task_04); `.agents/skills/` is canonical, `skills/` its mirror.

## What is not governed

`internal/worktree/` (`worktree.go`, `adminlock.go`, the new `merged_head.go`,
`staging.go`, `disposal.go` and every test file), `internal/store/store.go` and
its new test file, `internal/cli/reconcile.go`, `internal/cli/carryforward.go`,
`internal/cli/deliver_workflow.go`, the new `internal/cli/merged_release.go`
and `internal/cli/deliver_release.go`,
the other `internal/cli` test files this Spec names,
`internal/delivery/engine.go` and its tests, `internal/daemon/reconcile.go`,
`internal/speccheck/checkout.go`, the new ADR under `docs/adr/`,
`docs/user-guide/commands.md` and `CONTEXT.md` are ordinary.

## Sanctioned regeneration

```yaml
command: make skills-sync
```

## Limits

- No change to archived Specs, to the Branch Disposition record or to the Run
  Database schema version.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
