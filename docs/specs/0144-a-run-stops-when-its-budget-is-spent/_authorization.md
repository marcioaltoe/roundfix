---
status: approved
granted: 2026-09-18
action: describe the budget-bounded Run outcome and the widened carry-forward acceptance in the Roundfix skill that ships with the CLI
consuming: 0144-a-run-stops-when-its-budget-is-spent
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

# Approved authority for Spec 0144

The Spec bounds an Implement Run by its configured maximum duration, settles the
bounded Run as `BudgetExceeded`, and widens Task Carry-Forward to accept that
outcome. Both are public CLI behavior, and the repository's hard rule requires a
pull request that changes CLI behavior to ship the Roundfix skill update with
it. The maintainer authorized both copies on 2026-09-18.

## Why a governed path is unavoidable

The Roundfix skill is how an operating Agent learns what the CLI does, and it
currently states that carry-forward "accepts only Runs whose outcome is
`Stopped` or `Unresolved`". That sentence becomes false with this change. The
skill is a Roundfix-owned Skill, so changing its text needs an express grant,
and both the canonical copy and the embedded mirror are bounded here.

The Daemon, the CLI, the Run Database and the user guide are ordinary source
that no authorization has bounded.

## Approved bounded mutation

Describe, in the Roundfix skill and its mirror:

- that an Implement Run is bounded by the configured maximum Run duration when
  the budget is enabled, and settles `BudgetExceeded` with a reason naming the
  maximum and the elapsed time;
- that the bounded Run preserves its Run Worktree and Run Branch;
- that Task Carry-Forward accepts `BudgetExceeded` beside `Stopped` and
  `Unresolved`, under its existing proof requirements.

After the canonical edit, regenerate the distributed mirror with
`make skills-sync`. If any derived pin changes as a result, rewrite it only with
`make baseline-digests`.

## Sanctioned regeneration

```yaml
command: make skills-sync
```

```yaml
command: make baseline-digests
```

## Limits

- No action, operation or path beyond those above.
- No change to the skill's version, its owned-skill contract, or any behavior it
  documents beyond the bound and the widened acceptance.
- No Baseline module, authoring skill, qa-gate skill, linter or Verification
  configuration edit.
- No edit to this repository's own `.roundfixrc.yml`, and no change to the
  budget's default value or default enablement.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
