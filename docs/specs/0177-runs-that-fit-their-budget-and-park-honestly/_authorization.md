---
status: approved
granted: 2026-09-28
action: stop counting instruction paths as Wave collisions, renew the Implement Run Budget at each Task settlement, park a budget-stopped delivery item as a Run outcome, stage carry-forward commits without repository hooks, and report line-bound phrase checks against Markdown at authoring
consuming: 0177-runs-that-fit-their-budget-and-park-honestly
paths:
  - .agents/skills/write-tasks/SKILL.md
  - skills/write-tasks/SKILL.md
  - .agents/skills/write-tasks/references/task-template.md
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
  - internal/speccheck/coherence.go
  - internal/docscontract/testdata/corpus-golden.json
  - internal/spec/archive_layout_characterization_test.go
  - skills/baseline_skill_contract_test.go
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0177

The maintainer authorized continuing with the next wave after the v0.18.0
release in chat on 2026-09-28 ("Após o release, pode continuar com as
implementações da onda seguinte"); this Spec is part of that wave. The Go
sources and test files ride the standing grant of 2026-09-21 for governed
source a slice genuinely needs; the skill and template files ride the standing
grant of 2026-09-18 for keeping the shipped skills true to the CLI. The set was
measured with `GovernedPath`.

## Why each governed path is unavoidable

- `.agents/skills/write-tasks/SKILL.md`, `skills/write-tasks/SKILL.md` — the
  write-tasks skill's declared-path rules state that an `instruction:` path
  never makes two Tasks collide (task_01); `.agents/skills/` is canonical and
  `skills/` its mirror.
- `.agents/skills/write-tasks/references/task-template.md` — the Task
  template's Verification guidance teaches the wrap-tolerant phrase form and
  names `SC-VERIFY-WRAP-FRAGILE` (task_05). Its mirror
  `skills/write-tasks/references/task-template.md` is not governed.
- `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md` — the Roundfix
  skill describes the Implement Run Budget (task_02), the Delivery Queue's
  `run-budget-exceeded` blocker (task_03) and carry-forward staging (task_04).
- `internal/speccheck/coherence.go` — `stagedDetectors` registers
  `SC-VERIFY-WRAP-FRAGILE` as a Task-stage detector (task_05).
- `internal/docscontract/testdata/corpus-golden.json`,
  `internal/spec/archive_layout_characterization_test.go` — the corpus golden
  counts every characterized code and its pin must match it; the new code joins
  with count `0` (task_05).
- `skills/baseline_skill_contract_test.go` — the contract test that pins the
  Task template's Verification guidance gains the wrap-tolerant form (task_05).

## What is not governed

`internal/spec/collision.go` and the new collision test,
`internal/daemon/task_engine.go` and the new budget test,
`internal/cli/implement.go`, `internal/cli/deliver_workflow.go`,
`internal/cli/carryforward.go` and their new tests,
`internal/config/config.go` and its new test, `internal/delivery/engine.go` and
its new test, `internal/speccheck/verification.go`,
`internal/speccheck/citations.go` and the new phrase test,
`internal/docscontract/corpus_test.go`, the new ADR under `docs/adr/`,
`docs/user-guide/commands.md`, `docs/user-guide/configuration.md`, `CONTEXT.md`
and the mirror `skills/write-tasks/references/task-template.md` are ordinary.
`.roundfixrc.yml` is governed and is not touched.

## Sanctioned regeneration

```yaml
command: make skills-sync
```

## Limits

- No change to `.roundfixrc.yml`, to `GovernedPath`, to the governed set or to
  the changed-path audit, and no rename or removal of an existing top-level
  test.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
