---
status: approved
granted: 2026-09-30
action: refresh the default branch before a merged item's cleanup, name amending commits and the recovery order in a Delivery Retry refusal, warn when a queue owner is older than its item's starting main, and authorize a Task commit by the grant it ran under when the delivery target holds that grant
consuming: 0187-a-queue-that-recovers-without-a-supervisor
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

# Approved authority for Spec 0187

The maintainer asked on 2026-09-30 to "Continue até o final para o release de
todas as implementações e ajustes": to deliver and release every remaining
implementation and adjustment of this cycle. This Spec delivers the queue
residuals recorded in
`docs/findings/2026-09-29-the-first-full-queue-trial-needed-manual-recovery.md`.
The skill files ride the standing grant of 2026-09-18 for keeping the shipped
skills true to the CLI. The set was measured with `GovernedPath` on the
authoring branch at `379caa31`, through a `go test -overlay` probe that wrote
nothing to the repository.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md` — the Roundfix
  skill's Delivery queue section must state the cleanup refresh (task_01), the
  retry recovery commands (task_02) and the `owner-older-than-main` warning
  (task_03), and its QA section must state the authorizing revision rule
  (task_04). `.agents/skills/` is
  canonical and `skills/` its mirror.

## What is not governed

`internal/cli/deliver_workflow.go`, `internal/cli/carryforward.go`,
`internal/cli/deliver_revalidate.go`, `internal/delivery/engine.go`,
`internal/speccheck/mechanical.go`, the existing test
`internal/cli/deliver_release_evidence_test.go`, the new test files
`internal/cli/deliver_release_refresh_test.go`,
`internal/cli/deliver_retry_amendment_test.go`,
`internal/cli/deliver_owner_staleness_test.go`,
`internal/delivery/owner_staleness_test.go` and
`internal/speccheck/mechanical_grant_ran_under_test.go`, and
`docs/user-guide/commands.md` are ordinary. `internal/cli/cli_test.go`,
`internal/speccheck/mechanical_test.go` and
`docs/references/coverage-record.json` are governed and are not touched: no
help text changes and no existing top-level test is renamed or removed.

## Sanctioned regeneration

```yaml
command: make skills-sync
```

## Limits

- No change to archived Specs, existing QA Reports, the Run Database schema,
  the retry limit or any command's flags.
- Roundfix prints the retry recovery commands and never runs them.
- The live Run Database under `~/.roundfix` is never opened for writing by a
  Task or the gate, and no test reaches GitHub or a provider.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
