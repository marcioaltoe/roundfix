---
status: approved
granted: 2026-09-24
action: make the pre-PR review classify verdicts by substance and keep the answer, run the claude provider, and review a Spec's delivery against its decisions, and keep the shipped skill true to it
consuming: 0160-a-review-that-reaches-a-verdict
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

# Approved authority for Spec 0160

The skill files ride the standing authorization of 2026-09-18, valid across the
remaining slices of this queue for one purpose: keeping the shipped skill true
to the CLI behavior the slice itself delivers. The bounded set was measured with
`GovernedPath` over every path this Spec changes.

## Why each governed path is unavoidable

The Roundfix skill states that `claude` is refused and that only an exact
`No findings` passes; both statements become false with this Spec, and the
repository's hard rule requires the skill update to ship with the CLI change.

## What is not governed

`internal/cli/review.go`, its tests and `docs/user-guide/commands.md` are
ordinary source.

## Sanctioned regeneration

```yaml
command: make skills-sync
```

## Limits

- No action, operation or path beyond those above.
- No change to the exit-code contract of `roundfix review`.
- No paid API use: every test stubs `agent.Runner`.
- No release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
