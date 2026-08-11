---
granted: 2026-08-10
action: reroute-claude-profile-entries
paths:
  - .roundfixrc.yml
consuming: 0091-a-proof-that-can-refuse
---

# Tooling authorization — a broken tuple that blocks its own repair (2026-08-10)

Spec 0091's QA gate reported F-002: `profiles validate` rejects the configured
`claude/opus/high` with `effective_selection_mismatch` while its own diagnostic
prints requested and observed as the same `opus/high`. Measured on 2026-08-10,
the same binary built from `main` validates all five configured tuples, so this
is a regression introduced by this Spec's own catalogue and membership work.

The regression blocks its own repair. Preflight proves every configured tuple
and substitutes none, and the broken tuple sits in the fallback chains of five
categories, so no Run can be created at all — including the Run that would
execute the corrective Task. A complete one-Run override does not help: it
preserves the configured Fallback Chains, so the broken tuple is still proved.

Shown three ways out, the maintainer chose to reroute the profile rather than
have the Supervisor write the fix by hand or revert the Spec:

> Vamos pela opção 2, mude para opencode deekseek v4 flash 0731

## What this covers

The `claude / opus / high` entries in `.roundfixrc.yml` become
`opencode / openrouter/deepseek/deepseek-v4-flash-0731 / max`, which the same
binary already proves. This is the runtime the file's own comments record as
the measured OpenRouter choice, so no new model enters the configuration.

The reroute exists to let Task 07 run and fix F-002. It is not a judgement
about Claude's fitness for the `frontend` category, and the file's comments
recording why Claude leads there are preserved.

## Authorized paths

- `.roundfixrc.yml`, limited to the Agent Selection entries naming
  `claude`/`opus`/`high`.

## Bounded by purpose

The purpose is unblocking the repair of a regression this Spec introduced. It
does not authorize changing any other selection, adding a model the file does
not already carry, or altering the reasoning-effort policy of an untouched
category. Restoring the Claude entries once F-002 is fixed is the expected end
state and needs no separate grant.

## Commit choreography

This record lands as its own commit, before the configuration change.
