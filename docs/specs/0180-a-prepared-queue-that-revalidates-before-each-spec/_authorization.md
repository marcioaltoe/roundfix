---
status: approved
granted: 2026-09-28
action: add a read-only Delivery Plan and refuse a queue without delivery authority, revalidate each queued Spec on its starting main before its first Run, record and enforce queue limits, present one Pending Question, and make the implement-spec skill hand implementation to Roundfix
consuming: 0180-a-prepared-queue-that-revalidates-before-each-spec
paths:
  - internal/cli/cli_test.go
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
  - .agents/skills/implement-spec/SKILL.md
  - skills/implement-spec/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0180

The maintainer authorized continuing with the next wave after the v0.18.0
release in chat on 2026-09-28 ("Após o release, pode continuar com as
implementações da onda seguinte"); this Spec is part of that wave. The Go test
file rides the standing grant of 2026-09-21 for governed source a slice
genuinely needs. The skill files ride the standing grant of 2026-09-18 for
keeping the shipped skills true to the CLI. The set was measured with
`GovernedPath` on `b92aefda`.

## Why each governed path is unavoidable

- `internal/cli/cli_test.go` — `TestRunCommandHelp` pins the `deliver` help
  text, and task_03 adds `roundfix deliver plan [--json] [<slug>...]` to the
  strings it expects. No top-level test is added, renamed or removed there, so
  `docs/references/coverage-record.json` does not change.
- `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md` — the
  Delivery queue section documents the revalidation and its blockers (task_01),
  the Delivery Plan and the start refusal (task_03), and the limit flags and
  the Pending Question (task_04); `.agents/skills/` is canonical and `skills/`
  its mirror.
- `.agents/skills/implement-spec/SKILL.md`, `skills/implement-spec/SKILL.md` —
  the owned entry point stops carrying its own Task loop and hands
  implementation to `roundfix implement` or `roundfix deliver` (task_05).

## What is not governed

`internal/delivery/engine.go` and the new `internal/delivery/question.go`;
`internal/store/store.go` and `internal/store/delivery.go`;
`internal/cli/deliver.go`, `internal/cli/deliver_workflow.go`,
`internal/cli/cli.go` and the new `internal/cli/deliver_plan.go` and
`internal/cli/deliver_revalidate.go`; the existing tests
`internal/delivery/engine_test.go`, `internal/delivery/retry_test.go`,
`internal/cli/deliver_test.go`, `internal/cli/deliver_recovery_test.go` and
`internal/cli/deliver_retry_test.go`; the new tests
`internal/delivery/revalidate_test.go`, `internal/delivery/limits_test.go`,
`internal/delivery/question_test.go`, `internal/store/delivery_limits_test.go`,
`internal/cli/deliver_revalidate_test.go`, `internal/cli/deliver_plan_test.go`
and `internal/cli/deliver_limits_test.go`; `docs/user-guide/commands.md` and
`CONTEXT.md` are ordinary.

## Sanctioned regeneration

The repository-owned commands resolve their generated outputs; these
declarations record the regeneration that follows the approved skill edits and
add no source paths.

```yaml
command: make skills-sync
```

```yaml
command: make baseline-digests
```

## Limits

- No change to an exported function signature, to archived Specs, to
  `skills/_ownership.yml` or to the CI workflow, and no renamed or removed
  top-level test.
- The live Run Database under `~/.roundfix` is never opened by a Task or the
  gate for writing.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
