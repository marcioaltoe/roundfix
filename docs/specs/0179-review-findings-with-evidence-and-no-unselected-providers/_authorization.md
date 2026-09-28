---
status: approved
granted: 2026-09-28
action: scope every CodeRabbit request to the repository that selected it, give each pre-PR review finding an evidence-backed disposition tied to its head, let a findings verdict stand until its candidate changes, and park a blocking review after archive as needing a corrective Spec
consuming: 0179-review-findings-with-evidence-and-no-unselected-providers
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

# Approved authority for Spec 0179

The maintainer authorized continuing with the next wave after the v0.18.0
release in chat on 2026-09-28 ("Após o release, pode continuar com as
implementações da onda seguinte"); this Spec is part of that wave. The skill
files ride the standing grant of 2026-09-18 for keeping the shipped skills true
to the CLI. The set was measured with `GovernedPath`.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md` — the Roundfix
  skill tells agents when to request a CodeRabbit review (task_01) and
  describes the Pre-PR review command (task_02, task_03, task_04) and the
  Delivery queue (task_03, task_04); `.agents/skills/` is canonical, `skills/`
  its mirror.

## What is not governed

`internal/cli/review.go`, `internal/cli/cli.go`,
`internal/cli/deliver_workflow.go`, `internal/delivery/engine.go`,
`internal/config/config.go`, the new test files this Spec names under
`internal/cli/`, `internal/config/` and `internal/delivery/`,
`docs/user-guide/commands.md`, `docs/user-guide/usage.md`,
`docs/user-guide/configuration.md` and `CONTEXT.md` are ordinary.
`internal/cli/cli_test.go` and `docs/references/coverage-record.json` are
governed and are not touched: no existing top-level test is replaced or
renamed, and the public help test matches the review usage line, which stays.

## Sanctioned regeneration

```yaml
command: make skills-sync
```

## Limits

- No change to archived Specs, existing review records, the legacy `watch`,
  `resolve` and `fetch` behavior, or the Run Database schema.
- No paid API use, provider call from a test, release, tag, deployment or
  branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
