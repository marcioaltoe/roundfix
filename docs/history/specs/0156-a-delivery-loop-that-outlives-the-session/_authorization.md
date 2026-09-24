---
status: approved
granted: 2026-09-24
action: add a durable delivery queue that advances Specs from Run to merge, and keep the shipped skill true to it
consuming: 0156-a-delivery-loop-that-outlives-the-session
paths:
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
  - internal/cli/cli_test.go
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0156

The maintainer approved this delivery on 2026-09-24 as the first of the
restructured queue. Two standing grants are consumed.

The skill files ride the authorization of 2026-09-18, for keeping the shipped
skill true to the CLI behavior this slice delivers. `internal/cli/cli_test.go`
rides the authorization of 2026-09-21 for governed paths a slice needs, measured
against `internal/speccheck/governed.go`.

## Why each governed path is unavoidable

A new command family changes the public CLI surface, which the shipped skill
documents and the repository's hard rule requires to ship together.
`internal/cli/cli_test.go` holds the command surface's own tests.

## What is not governed

Measured by intersection: the new `internal/delivery` package, the store
additions in `internal/store`, `internal/cli/deliver.go` and its tests, and
`docs/user-guide/commands.md` are ordinary source.

## Approved bounded mutation

Register the `deliver` command family in `internal/cli/cli_test.go` alongside the
commands already covered there.

State, in the Roundfix skill and its mirror, that `roundfix deliver` advances a
queue of Specs from Run to merge in the order the maintainer set — Run, pre-PR
review, archive on the branch, repository gate, pull request, checks, squash
merge — that a blocker parks its item, that a resumed queue reconciles each
recorded action before retrying it, and that publication needs the Spec's
authorization to grant push, pull request and merge.

After the canonical edit, regenerate the distributed mirror with
`make skills-sync`.

## Sanctioned regeneration

```yaml
command: make skills-sync
```

## Limits

- No action, operation or path beyond those above.
- No release, tag or deployment is ever performed by the queue.
- No change to the Implement executor's own contract, to `roundfix review` or to
  `roundfix archive`.
- No paid API use in any test: the GitHub boundary is faked.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
