---
status: approved
granted: 2026-09-25
action: stamp the QA Task status an override waives, align the Roundfix and archive-spec skills with the command, and make ordinal claims unique within a Spec
consuming: 0169-override-and-ordinal-follow-ups
paths:
  - internal/spec/archive.go
  - internal/spec/archive_test.go
  - .agents/skills/archive-spec/SKILL.md
  - skills/archive-spec/SKILL.md
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
  - .agents/skills/qa-gate/SKILL.md
  - skills/qa-gate/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0169

Minted on 2026-09-25 as the follow-up Spec 0159 recorded. The skill files ride
the maintainer's express authorization of 2026-09-24 for the archive override
guidance and the standing grant of 2026-09-18 for keeping the shipped skill true
to the CLI; the Go sources ride the standing grant of 2026-09-21. The set was
measured with `GovernedPath`.

## Why each governed path is unavoidable

- `internal/spec/archive.go`, `internal/spec/archive_test.go` — the override
  stamp is written there.
- The four skill files — the Roundfix and archive-spec skills describe the
  override; `.agents/skills/` is canonical, `skills/` its mirror.

## What is not governed

`internal/speccheck/ordinal.go`, `internal/speccheck/citations.go`, their tests
and `docs/user-guide/commands.md` are ordinary.

## Sanctioned regeneration

```yaml
command: make skills-sync
```

## Added after the pre-PR review

`.agents/skills/qa-gate/SKILL.md` and its mirror carry the override field list;
the maintainer's authorization of 2026-09-24 for the archive override guidance
covers them.

## Limits

- No change to when the override is accepted or refused.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
