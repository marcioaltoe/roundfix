# Grok through Cursor, measured on two replayed Tasks

Measured on 2026-10-03 on the maintainer's machine, with `NODE_OPTIONS`
unset, `bin/roundfix` built from `cddff2de` (tasks 01 to 03), acpx `0.19.4`,
`cursor-agent` `2026.07.09-a3815c0` and the maintainer's existing Cursor
login (`cursor-agent status` reported a login; nothing logged in, and no key
or token was read or passed). Prompts went through the maintainer's Cursor
plan, as authorized on 2026-10-01.

## Measured

### The recorded session

An ACP client that sent only `initialize` and `session/new` to
`cursor-agent acp`, in an empty scratch directory, received a session and its
`configOptions` without any `authenticate` call. `initialize` advertised one
auth method, `cursor_login`, and `loadSession: true`. The `configOptions`
carry no account field; they are committed unchanged as
`internal/agent/testdata/cursor-session-recorded.json`: a `mode` option and a
`model` option with 43 values, current value `default[]`.

Advertised Grok values, in the order the session lists them:

- `grok-4.7[context=256k,reasoning_effort=high,fast=true]`
- `grok-4.6[effort=high,fast=true]`
- `grok-4.5[effort=high,fast=true]`

No Grok value states `fast=false`, so the protocol's preference had no
variant to choose, and the first value was taken.

- **acpx opened the session alone: yes.** `bin/roundfix profiles validate
  --category docs`, in a scratch clone whose `.roundfixrc.yml` named
  `cursor / grok-4.7[context=256k,reasoning_effort=high,fast=true] / ""` as
  the `docs` preferred selection and `codex / gpt-5.6-luna / max` as its
  fallback, passed both tuples (29 s). The Grok tuple was proved
  `model-managed` through acpx's built-in `cursor` agent, with no
  `authenticate` from Roundfix.
- **`session/load` after `session/new`: did not work.** Loading the id just
  returned by `session/new`, on the same connection and again from a second
  `cursor-agent acp` process, answered `-32602 Invalid params` with
  `Session "<id>" not found`. No prompt had been sent in that session.
- **Billing mode the value states: `fast=true`.** Every Grok value states
  Cursor's fast variant; none states a price, a quota or a plan pool. The
  Cursor sessions reported no token usage to Roundfix, so the plan's charge
  for these Runs is not visible from the Run record.

### The replays

Each replay ran `bin/roundfix implement --spec <slug> --agent <runtime>
--model <model> --reasoning-effort <effort> --no-input` in a fresh scratch
clone of this repository, on a pre-state built as the protocol states: the
Spec's squash merge commit, the Task's declared files restored to its first
parent, the Spec directory moved back under `docs/specs/` with that Task
`pending` and its Result removed, every other Task removed from the graph,
and `qa: declined`. Both pre-states failed the Task's Verification before the
replay, so no reserve was needed. The two clones of each Task had identical
trees. The default is the built-in `docs` and `chore` recommendation,
`codex / gpt-5.6-luna / max`. Replays ran one at a time. Wall time is the
whole `implement` command, preflight proof included.

| Task | Selection | Outcome | Wall time | Prompts | Verification repairs | Person needed | Tokens |
| --- | --- | --- | --- | --- | --- | --- | --- |
| `0210-evidence-snapshots-that-stay-small/task_02` (docs) | `cursor / grok-4.7[context=256k,reasoning_effort=high,fast=true] / ""` | Clean, Verification passed on attempt 1 | 264 s | 1 | 0 | no | none reported |
| `0210-evidence-snapshots-that-stay-small/task_02` (docs) | `codex / gpt-5.6-luna / max` | Clean, Verification passed on attempt 1 | 167 s | 1 | 0 | no | 457,797 |
| `0195-owned-skills-and-a-release-step-that-follow-the-bundle/task_06` (chore) | `cursor / grok-4.7[context=256k,reasoning_effort=high,fast=true] / ""` | Clean, Verification passed on attempt 1 | 278 s | 1 | 0 | no | none reported |
| `0195-owned-skills-and-a-release-step-that-follow-the-bundle/task_06` (chore) | `codex / gpt-5.6-luna / max` | Clean, Verification passed on attempt 1 | 189 s | 1 | 0 | no | 910,160 |

Runs: `run_20261003T135501Z_04ace8f85ce01177`,
`run_20261003T135928Z_a1394e4dfe7bc1c0`,
`run_20261003T140232Z_1927170b378f8404` and
`run_20261003T140711Z_c40a89b26474e75e`, in that row order.

No replay showed a permission prompt, `cursor/ask_question` or
`cursor/create_plan` in its Agent console. On `0195/task_06`, both selections
committed the same skill version change and a byte-identical
`skills/testdata/owned-skill-versions.json`. On `0210/task_02`, both wrote the
two required phrases and passed the same Verification.

The `0195/task_06` pre-state is not the Task's original one: restoring the
declared files to the squash merge's first parent leaves the Roundfix skill
at `0.0.2` with no version record, while the Task text says `0.0.3`. Both
selections noticed the difference, raised the skill to `0.0.4` and recorded
the versions, so the comparison stays on identical inputs.

## Reading

Grok 4.7 through Cursor settled both replays with no person needed: one
prompt each, no Verification repair, the same outcome as the default. It was
slower, about 1.5 times the default's wall time on both Tasks (264 s against
167 s, and 278 s against 189 s), and Cursor reported no tokens, so its cost
cannot be compared from the Run record. Two low-complexity Tasks are a small
sample and say nothing about harder work, where the published benchmark gap
named in the TechSpec still stands.

The protocol's condition holds, so this Task mints
`docs/backlog/2026-10-03-a-grok-fallback-for-docs-and-chore.md`, proposing
Grok as a `docs` and `chore` fallback. It changes no profile and no default.
Two observations bound that proposal: every advertised Grok value is a
`fast=true` variant whose plan charge Roundfix cannot see, and `session/load`
does not find a session that has not been prompted.
