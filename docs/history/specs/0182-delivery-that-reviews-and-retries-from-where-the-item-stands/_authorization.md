---
status: approved
granted: 2026-09-29
action: make the pre-PR review diff the candidate from its merge base, record a Task already completed on its carry-forward target as nothing to carry, and make a Delivery Retry carry forward from every terminal Implement Run of the item's Spec newest first
consuming: 0182-delivery-that-reviews-and-retries-from-where-the-item-stands
paths:
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0182

The maintainer approved the Onda 4 plan in chat on 2026-09-29. Answering a
structured question with "Duas filas", they chose to deliver Specs 0181 and
0182 through one `roundfix deliver start`, and a second queue afterwards. This
Spec is part of that wave. The skill files ride the standing grant of
2026-09-18 for keeping the shipped skills true to the CLI. The set was measured
with `GovernedPath` on the authoring branch at `a626494c`.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md` — the Roundfix
  skill describes the Pre-PR review base and record (task_01), the Reconcile
  Command's carry-forward and the implement Preflight (task_02), and the
  Delivery Retry (task_03). `.agents/skills/` is canonical and `skills/` its
  mirror.

## What is not governed

`internal/cli/review.go`, `internal/cli/carryforward.go`,
`internal/cli/reconcile.go`, `internal/cli/implement.go`,
`internal/cli/deliver.go`, `internal/cli/deliver_workflow.go`,
`internal/delivery/engine.go`, the existing tests
`internal/cli/deliver_recovery_test.go` and `internal/delivery/retry_test.go`,
the new test files `internal/cli/review_merge_base_test.go`,
`internal/cli/carryforward_completed_target_test.go` and
`internal/cli/deliver_retry_runs_test.go`, `docs/user-guide/commands.md` and
`CONTEXT.md` are ordinary. `internal/cli/cli_test.go` and
`docs/references/coverage-record.json` are governed and are not touched: no
help text changes and no existing top-level test is renamed or removed.

## Sanctioned regeneration

```yaml
command: make skills-sync
```

## Limits

- No change to archived Specs, existing review records, the reviewer prompt
  wording or the Run Database schema.
- The live Run Database under `~/.roundfix` is never opened for writing by a
  Task or the gate.
- No paid API use, provider call from a test, release, tag, deployment or
  branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
