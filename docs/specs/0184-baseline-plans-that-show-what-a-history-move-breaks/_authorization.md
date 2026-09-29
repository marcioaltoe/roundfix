---
status: approved
granted: 2026-09-29
action: make every Baseline Plan that carries History Relocations report, as warnings the Plan Digest binds, each tracked file whose citations the relocations would break, and document it in the Roundfix skill
consuming: 0184-baseline-plans-that-show-what-a-history-move-breaks
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

# Approved authority for Spec 0184

The maintainer approved this plan in chat on 2026-09-29, answering structured
questions. With "Duas filas" they chose a second delivery queue, Onda 5, for
this Spec and Spec 0183. With "Mostrar o impacto" they chose, for Backlog Entry
B18, to show each History Relocation's citation impact as warnings inside the
same Plan Digest, over a deferral flag or a separate digest. The skill files
ride the standing grant of 2026-09-18 for keeping the shipped skills true to
the CLI. The set was measured with `GovernedPath` on the authoring branch at
`fe517d2b`.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md` — the Roundfix
  skill's Baseline section describes what `baseline update` and `baseline plan`
  report. Task 03 adds the Relocation Citation warnings, their codes and the
  digest binding. `.agents/skills/` is canonical and `skills/` is its mirror.

## What is not governed

These files are ordinary:

- `internal/baseline/history_layout.go` and `internal/baseline/plan.go`
- the new `internal/baseline/history_citations.go`
- the new test files `internal/baseline/history_citations_test.go` and
  `internal/cli/baseline_history_citation_test.go`
- `internal/cli/baseline_update_test.go`
- `docs/user-guide/context-driven-development.md` and `CONTEXT.md`

These governed files are not touched: `internal/baseline/plan_test.go`,
`internal/cli/baseline_plan_test.go`, `internal/cli/cli_test.go` and
`docs/references/coverage-record.json`. No help text changes, and no existing
top-level test is renamed or removed. The new CLI tests reuse the package's
existing Baseline Plan test helpers without editing them.

## Sanctioned regeneration

```yaml
command: make skills-sync
```

## Limits

- No change to the Baseline Plan or result schema versions, the Setup
  Manifest, the embedded catalog, `internal/baseline/assets/**`,
  `skills/setup-context-driven/**` or archived Specs.
- No new command or flag, and no rewrite of any citing file.
- Tests use disposable repositories only. The live Run Database under
  `~/.roundfix` is never opened for writing by a Task or the gate.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
