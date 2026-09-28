---
status: approved
granted: 2026-09-28
action: add a Delivery Retry that carries a parked item's settled Tasks back and returns it to the queue through a live or new owner, stage carry-forward in the Run's integration order, commit Project Config only on a named grant and report every exclusion, and retry a timed-out profile proof once
consuming: 0173-a-delivery-queue-that-recovers
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

# Approved authority for Spec 0173

The maintainer approved Onda 2 of the efficiency sequence in chat on 2026-09-28
("Siga para onda 2 como sugerido até o release"); this Spec is part of it. The
Go test file rides the standing grant of 2026-09-21 for governed source; the
skill files ride the standing grant of 2026-09-18 for keeping the shipped
skills true to the CLI. The set was measured with `GovernedPath`.

## Why each governed path is unavoidable

- `internal/cli/cli_test.go` — `TestRunCommandHelp` pins the `deliver` help
  text, and task_05 adds `roundfix deliver retry <slug>` to the strings it
  expects. No top-level test is added, renamed or removed there, so
  `docs/references/coverage-record.json` does not change.
- `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md` — the Roundfix
  skill states the Daemon commit rule for Project Config (task_02), the Delivery
  Retry (task_05) and the profile proof retry (task_06); `.agents/skills/` is
  canonical and `skills/` its mirror.

## What is not governed

`internal/cli/carryforward.go`, `internal/cli/deliver.go`,
`internal/cli/deliver_workflow.go`, `internal/cli/cli.go`,
`internal/cli/profiles_validate.go`, `internal/daemon/engine.go`,
`internal/daemon/task_engine.go`, `internal/delivery/engine.go`,
`internal/store/delivery.go`, the new test files
`internal/cli/carryforward_integration_order_test.go`,
`internal/cli/deliver_recovery_test.go`, `internal/cli/deliver_retry_test.go`,
`internal/cli/profile_proof_retry_test.go`,
`internal/daemon/project_config_commit_test.go`,
`internal/delivery/retry_test.go` and `internal/store/delivery_retry_test.go`,
`docs/user-guide/commands.md` and `CONTEXT.md` are ordinary.

## Sanctioned regeneration

```yaml
command: make skills-sync
```

## Limits

- No Run Database schema change, no change to archived Specs, and no retry of a
  recorded push, pull request creation or merge.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
