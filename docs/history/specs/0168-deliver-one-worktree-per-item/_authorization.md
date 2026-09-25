---
status: approved
granted: 2026-09-24
action: run every roundfix deliver item in its own linked worktree so the queue never switches, resets or cleans the user's checkout, and keep the shipped skill true to it
consuming: 0168-deliver-one-worktree-per-item
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

# Approved authority for Spec 0168

Minted on 2026-09-24 as the follow-up that Spec 0161 recorded, under the
maintainer's standing rules: at the corrective ceiling, deliver with limits
recorded and carry them to a Spec of their own.

The skill files ride the standing grant of 2026-09-18, for keeping the shipped
skill true to the CLI behaviour this slice changes. The set was measured with
`GovernedPath` itself over every path this Spec changes.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/SKILL.md` — the canonical skill's delivery queue
  section describes what `roundfix deliver` does to the checkout and what
  `deliver status` prints; both change here.
- `skills/roundfix/SKILL.md` — the distributed mirror must stay identical to the
  canonical skill, and is regenerated, never edited.

## What is not governed

Measured with `GovernedPath`: `internal/delivery`, `internal/cli/deliver.go`,
`internal/cli/deliver_workflow.go`, `internal/store/delivery.go`,
`internal/store/store.go`, `internal/worktree/worktree.go`, their tests and
`docs/user-guide/commands.md` are ordinary source. `internal/cli/cli_test.go`
is governed and is not changed: the `deliver` help text it pins stays the same.

## Approved bounded mutation

State, in the Roundfix skill and its mirror, that each queued Spec runs in its
own linked worktree under `worktree.location`, created from the refreshed
default branch; that `roundfix deliver` never switches, resets or cleans the
user's checkout and does not need it clean; that a parked item keeps its
worktree, which `deliver status` prints; and that a merged item's worktree and
local branch are removed.

After the canonical edit, regenerate the distributed mirror with
`make skills-sync`.

## Sanctioned regeneration

```yaml
command: make skills-sync
```

## Limits

- No action, operation or path beyond those above.
- No new stage, no new external action and no new configuration key.
- No change to `roundfix implement`, `roundfix review` or `roundfix archive`.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
