---
status: approved
granted: 2026-09-25
action: seed QA Reports as pending and refuse a hollow pass, allocate reports where the reader looks, keep a deterministic failure across a temporary retry, classify an unobserved Verification, and let the event stream survive a record it cannot project
consuming: 0172-a-qa-gate-that-tells-the-truth
paths:
  - internal/cli/cli_test.go
  - internal/spec/archive_test.go
  - .agents/skills/qa-gate/SKILL.md
  - skills/qa-gate/SKILL.md
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0172

The maintainer approved the 2026-09-25 efficiency sequence (Ondas 0 and 1) in
chat on 2026-09-25; this Spec is part of it. The Go test files ride the standing
grant of 2026-09-21 for governed source; the skill files ride the standing grant
of 2026-09-18 for keeping the shipped skills true to the CLI. The set was
measured with `GovernedPath`.

## Why each governed path is unavoidable

- `internal/cli/cli_test.go` — the `roundfix events` tests live there, and the
  test that pins an abort on a malformed record is replaced by one that pins
  the warning (task_05).
- `internal/spec/archive_test.go` — task_01's Verification runs its archived
  `pass` corpus test as the preservation gate for the hollow-report rule; the
  file is expected to stay byte-identical.
- `.agents/skills/qa-gate/SKILL.md`, `skills/qa-gate/SKILL.md` — the qa-gate
  skill states the seeded verdict, the hollow-report refusal (task_01) and the
  report naming rule (task_02); `.agents/skills/` is canonical, `skills/` its
  mirror.
- `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md` — the Roundfix
  skill describes the Supervisor Run Event Stream's classified records
  (task_04) and its exit behavior (task_05).

## What is not governed

`internal/spec/qa.go`, `internal/spec/qa_test.go`, `internal/daemon/`
(`task_engine.go`, `engine.go`, their tests and the new test files),
`internal/worktree/worktree.go` and its test, `internal/runevent/stream.go` and
its test, `internal/cli/events.go`, `internal/cli/qa_report_test.go`,
`internal/cli/settle_test.go`, `internal/cli/archive_test.go`,
`docs/user-guide/commands.md` and `CONTEXT.md` are ordinary.

## Sanctioned regeneration

```yaml
command: make skills-sync
```

## Limits

- No change to the derived QA Verification command, to archived Specs or to
  archived QA Reports.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
